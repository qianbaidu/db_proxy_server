package service

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/models"
)

func RecordEvent(e *models.EventLog) error {
	if _, err := e.InsertData(); err != nil {
		l, _ := json.Marshal(e)
		log.Errorf("insert event log error %s , log : %s ", err, l)
		return err
	}
	return nil
}
