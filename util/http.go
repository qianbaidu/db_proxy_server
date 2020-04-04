package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"io"
	"time"
	"github.com/prometheus/common/log"
	"fmt"
)

const (
	HttpTimeout  = 10
	HttpTryTimes = 3
)

func httpPost(postUrl string, postData interface{}) (res []byte, err error) {
	jsonStr, err := json.Marshal(postData)
	if err != nil {
		return res, err
	}
	log.Debug(fmt.Sprintf("url:%s data: %s", postUrl, jsonStr))

	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("request error: ", err, " ,resp: ", resp)
		return res, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	return body, nil
}

func HttpPost(postUrl string, postData interface{}) (res []byte, err error) {
	for i := 1; i <= HttpTryTimes; i++ {
		res, err = httpPost(postUrl, postData)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5 * time.Duration(i))
	}

	return
}

func httpGet(url string) (res string, err error) {
	client := http.Client{Timeout: HttpTimeout * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return res, err
		}
	}

	res = result.String()
	return
}

func HttpGet(url string) (res string, err error) {
	for i := 1; i <= HttpTryTimes; i++ {
		res, err = httpGet(url)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5 * time.Duration(i))
	}

	return
}
