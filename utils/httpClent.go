package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ApiResponseData struct {
	Rtn  int                    `json:"rtn"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type ApiResponseMsg struct {
	Rtn int    `json:"rtn"`
	Msg string `json:"msg"`
}

// HttpGet 根据struct返回不同形式定义的接口数据
func (s *ApiResponseData) HttpGet(url string) (resp ApiResponseData, err error) {
	apiResponse, err := HttpGetBody(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(apiResponse, &resp)
	if err != nil {
		err = errors.New("Unmarshal failed " + err.Error())
		return
	}

	return resp, err
}

// HttpGet 根据struct返回不同形式定义的接口数据
func (s *ApiResponseMsg) HttpGet(url string) (resp ApiResponseMsg, err error) {
	apiResponse, err := HttpGetBody(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(apiResponse, &resp)
	if err != nil {
		err = errors.New("Unmarshal failed " + err.Error())
		return
	}

	return resp, err
}

// HttpGetBody 通用的获取http的Body
func HttpGetBody(url string) (body []byte, err error) {
	res, _ := http.Get(url)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = errors.New("Api接口调用错误:" + res.Status)
		return
	}

	body, err = ioutil.ReadAll(res.Body)

	return body, err
}

// HttpPostForm 通用的Post请求(Form)
func HttpPostForm(link string, form url.Values) ([]byte, error) {
	resp, err := http.PostForm(link, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// HttpPostJson 通用的Post请求(Json)
func HttpPostJson(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
