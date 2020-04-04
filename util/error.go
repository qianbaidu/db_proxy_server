package util

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ERR_DATABASES_NOT_EXISTS = errors.New("数据库信息不存在")
)

func IsFatalError(errorMsg string, err error) {
	if err != nil {
		if len(errorMsg) < 1 {
			errorMsg = err.Error()
		}
		logrus.Fatal(errorMsg, err)
	}
}

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type RestFulResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	//Result  string `json:"result"`
}

func CheckResponseStatus(resp string) error {
	respData := Response{}
	err := json.Unmarshal([]byte(resp), &respData)
	IsFatalError("Failed to parse request return data，error : ", err)
	if respData.Status != API_RESPONSE_STATUS_SUCCESS {
		return errors.New(fmt.Sprintf("request error : %s", err))
	}
	return nil
}

func CheckRestFulResponseStatus(resp string) error {
	respData := RestFulResp{}
	err := json.Unmarshal([]byte(resp), &respData)
	IsFatalError("Failed to parse request return data， error : ", err)
	if respData.Code != REST_API_RESPONSE_STATSUS_SUCCESS {
		return errors.New(fmt.Sprintf("%s", respData.Message))
	}
	return nil
}
