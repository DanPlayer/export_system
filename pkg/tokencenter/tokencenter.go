package tokencenter

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Options 令牌中心初始化参数
type Options struct {
	UniqueCode      string          // 唯一标识码(不能包含_字符)
	Cache           MemoryCache     // 缓存
	RefreshHandle   RefreshHandle   // 令牌刷新 Handle
	RefreshSW       bool            // 是否开启自定义刷新，默认：false 刷新策略
	RefreshTime     time.Duration   `default:"7200"` // 令牌刷新时间
	RefreshStrategy []time.Duration // 自定义刷新策略
}

type MemoryCache interface {
	Set(k, v string, expires time.Duration) error
	Get(ctx context.Context, k string) (string, error)
	Del(ctx context.Context, k string) error
	Scan(cursor uint64, match string, count int64) (keys []string, newCursor uint64, err error)
	Subscribe(k string, pb func(message string))
	SubscribeAllEvents(pb func(message string))
}

type RefreshHandle func(token string) (string, error)
type ExpiredHandle func(key string)

type TokenCenter struct {
	uniqueCode        string          // 唯一标识码
	mutex             sync.Mutex      // 互斥锁
	cache             MemoryCache     // 缓存
	refreshHandle     RefreshHandle   // 令牌刷新
	refreshSW         bool            // 是否开启自定义刷新，默认：false 刷新策略
	refreshTime       time.Duration   // 令牌刷新时间
	refreshStrategy   []time.Duration // 自定义刷新策略
	expiredHandleList []ExpiredHandle // 令牌过期处理
}

var ctx = context.Background()

func New(options *Options) *TokenCenter {
	if options.RefreshHandle == nil {
		panic("TokenCenter: refresh handle not be nil")
	}

	if options.RefreshSW && len(options.RefreshStrategy) == 0 {
		panic("TokenCenter: refresh strategy not be null if RefreshSW is true")
	}

	if !options.RefreshSW && options.RefreshTime == 0 {
		field, _ := reflect.TypeOf(options).FieldByName("RefreshTime")
		refreshTime, _ := strconv.Atoi(field.Tag.Get("default"))
		options.RefreshTime = time.Duration(refreshTime)
	}

	t := TokenCenter{
		mutex:             sync.Mutex{},
		uniqueCode:        options.UniqueCode,
		cache:             options.Cache,
		refreshHandle:     options.RefreshHandle,
		refreshSW:         options.RefreshSW,
		refreshTime:       options.RefreshTime,
		refreshStrategy:   options.RefreshStrategy,
		expiredHandleList: make([]ExpiredHandle, 0),
	}

	// 处理全局过期事件
	t.cache.SubscribeAllEvents(func(message string) {
		// 过滤无效事件
		if !t.isUsableKey(message) {
			return
		}
		contentKey := t.transformKey2Content(message)
		content, _ := t.cache.Get(ctx, contentKey)
		oldToken, count := t.parseContent(content)
		originKey := t.getOriginKey(message)

		// 重复定时刷新
		if !t.refreshSW && t.refreshTime != 0 {
			curToken, err := t.refreshHandle(oldToken)
			if err != nil {
				_ = t.cache.Del(ctx, contentKey)
				for _, handle := range t.expiredHandleList {
					handle(originKey)
				}
				return
			}
			// fmt.Printf("%s refresh success: %s \n", originKey, curToken)
			_ = t.cache.Set(contentKey, t.mergeContent(curToken, count), -1)
			_ = t.cache.Set(message, curToken, t.refreshTime)
			return
		}

		// 开启自定义刷新策略
		if t.refreshSW && len(t.refreshStrategy) > 0 {
			if count < len(t.refreshStrategy)-1 {
				count += 1
				curToken, err := t.refreshHandle(oldToken)
				if err != nil {
					_ = t.cache.Del(ctx, contentKey)
					for _, handle := range t.expiredHandleList {
						handle(originKey)
					}
					return
				}
				// fmt.Printf("%s refresh success: %s, count: %v \n", originKey, curToken, count)
				_ = t.cache.Set(contentKey, t.mergeContent(curToken, count), -1)
				_ = t.cache.Set(message, curToken, t.refreshStrategy[count])
				return
			}

			for _, handle := range t.expiredHandleList {
				_ = t.cache.Del(ctx, contentKey)
				handle(t.getOriginKey(message))
			}
		}
	})

	// 自动恢复当前实例历史未完成任务
	taskListKey, _, _ := t.cache.Scan(0, t.getContentKeyPrefix()+"*", 100000)
	// 批量恢复任务
	for _, taskKey := range taskListKey {
		contentStr, err := t.cache.Get(ctx, taskKey)
		if err != nil {
			continue
		}
		content, count := t.parseContent(contentStr)
		if content == "" {
			continue
		}
		curKey := t.transformContent2Key(taskKey)
		curKeyContent, err := t.Get(curKey)

		if err == nil && curKeyContent == "" {
			if t.refreshSW {
				_ = t.cache.Set(curKey, content, t.refreshStrategy[count])
			} else {
				_ = t.cache.Set(curKey, content, t.refreshTime)
			}
		}
	}

	return &t
}

func (t *TokenCenter) Set(key, token string) (err error) {
	refreshTime := time.Duration(0)
	if !t.refreshSW && t.refreshTime != 0 {
		refreshTime = t.refreshTime
	}
	if t.refreshSW && len(t.refreshStrategy) > 0 {
		refreshTime = t.refreshStrategy[0]
	}

	curKey := t.generateKey(key)
	curContent := t.generateContentKey(key)

	if err := t.cache.Set(curContent, t.mergeContent(token, 0), -1); err != nil {
		return err
	}
	return t.cache.Set(curKey, token, refreshTime)
}

func (t *TokenCenter) Get(key string) (token string, err error) {
	ctx := context.Background()
	contentKey := t.generateContentKey(key)
	content, err := t.cache.Get(ctx, contentKey)
	if err != nil {
		return "", err
	}
	token, _ = t.parseContent(content)
	return token, nil
}

func (t *TokenCenter) Del(key string) error {
	ctx := context.Background()
	contentKey := t.generateContentKey(key)
	err := t.cache.Del(ctx, contentKey)
	if err != nil {
		return err
	}
	return nil
}

// SubscribeExpiredEvent 订阅令牌过期事件
func (t *TokenCenter) SubscribeExpiredEvent(handle ExpiredHandle) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.expiredHandleList = append(t.expiredHandleList, handle)
}

func (t *TokenCenter) generateKey(key string) string {
	return "token-center-" + t.uniqueCode + "-key_" + key
}

func (t *TokenCenter) generateContentKey(key string) string {
	return "token-center-" + t.uniqueCode + "-content_" + key
}

func (t *TokenCenter) isUsableKey(key string) bool {
	return strings.HasPrefix(key, "token-center-"+t.uniqueCode+"-key_")
}

func (t *TokenCenter) getContentKeyPrefix() string {
	return "token-center-" + t.uniqueCode + "-content_"
}

func (t *TokenCenter) getOriginKey(key string) string {
	keyList := strings.Split(key, "_")
	if len(keyList) > 1 {
		return strings.Join(keyList[1:], "_")
	}
	return ""
}

func (t *TokenCenter) transformKey2Content(key string) string {
	originKey := t.getOriginKey(key)
	return t.generateContentKey(originKey)
}

func (t *TokenCenter) transformContent2Key(content string) string {
	originKey := t.getOriginKey(content)
	return t.generateKey(originKey)
}

func (t *TokenCenter) mergeContent(token string, count int) (con string) {
	return strconv.Itoa(count) + "_" + token
}

func (t *TokenCenter) parseContent(con string) (token string, count int) {
	conList := strings.Split(con, "_")
	if len(conList) > 1 {
		token = strings.Join(conList[1:], "_")
		count, _ = strconv.Atoi(conList[0])
	}
	return
}
