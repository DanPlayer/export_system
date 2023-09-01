package qiniu

import (
	"bytes"
	"context"
	"export_system/internal/config"
	"export_system/utils"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"math/rand"
	"sync"
	"time"
)

type Qiniu struct {
	Bucket string
	Mac    *qbox.Mac
}

var lock sync.RWMutex
var qn *Qiniu = nil

type UploadParam struct {
	Token  string `json:"token"`  // 上传token
	Domain string `json:"domain"` //  域名 支持https
	Bucket string `json:"bucket"` // bucket
	Expire string `json:"expire"` // token到期时间
	Prefix string `json:"prefix"` // 上传路径前缀 最终拼接  你的文件名 + prefix
}

type RespUploadCertificate struct {
	Rtn  int         `json:"rtn"`
	Msg  string      `json:"msg"`
	Data UploadParam `json:"data"`
}

func InitQiniu() error {
	lock.RLock()
	if qn != nil {
		lock.RUnlock()
		return nil
	}
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()
	qn = &Qiniu{
		Bucket: config.Config.Qiniu.Bucket,
		Mac:    qbox.NewMac(config.Config.Qiniu.AccessKey, config.Config.Qiniu.SecretKey),
	}
	return nil
}

func UploadFile2Qiniu(key, local string) error {
	putPolicy := storage.PutPolicy{
		Scope: qn.Bucket,
	}
	putPolicy.Expires = 7200
	upToken := putPolicy.UploadToken(qn.Mac)

	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}
	return formUploader.PutFile(context.Background(), &ret, upToken, key, local, nil)
}

func UploadData2Qiniu(key string, data []byte, size int64) error {
	putPolicy := storage.PutPolicy{
		Scope: qn.Bucket,
	}
	putPolicy.Expires = 7200
	upToken := putPolicy.UploadToken(qn.Mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExtra := storage.PutExtra{}
	return formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), size, &putExtra)
}

func UploadStreamQiniu(key string, reader io.Reader, size int64) error {
	putPolicy := storage.PutPolicy{
		Scope: qn.Bucket,
	}
	putPolicy.Expires = 7200
	upToken := putPolicy.UploadToken(qn.Mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExtra := storage.PutExtra{}
	return formUploader.Put(context.Background(), &ret, upToken, key, reader, size, &putExtra)
}

func MakeUpToken(id string) UploadParam {
	putPolicy := storage.PutPolicy{
		Scope:      qn.Bucket,
		SaveKey:    "upload_$(etag)$(ext)",
		ReturnBody: `{"key":"$(key)","ext":"$(ext)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)"}`,
	}
	putPolicy.Expires = 7200
	upToken := putPolicy.UploadToken(qn.Mac)

	now := utils.Now()
	prefix := "c/" + id + fmt.Sprintf("c/%02d/%02d/%02d/%02d/%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Unix()+rand.Int63())
	up := UploadParam{
		Token:  upToken,
		Domain: config.Config.Qiniu.Domain,
		Prefix: prefix,
		Bucket: qn.Bucket,
		Expire: now.Add(time.Second * 7000).Format(utils.DATETIME_FORMAT),
	}
	return up
}
