package service

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/qianbaidu/db-proxy/models"
	"time"
)

func RecordEvent(e *models.EventLog) error {
	e.UpdateDatatime = time.Now().Format("2006-01-02 15:04:05")
	if _, err := e.InsertData(); err != nil {
		l, _ := json.Marshal(e)
		log.Errorf("insert event log error %s , log : %s ", err, l)
		return err
	}
	return nil
}
