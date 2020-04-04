package util

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"github.com/siddontang/go/hack"
	"strconv"
	"strings"
	"time"
)

func FormatValue(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case []byte:
		return v, nil
	case string:
		return hack.Slice(v), nil
	case nil:
		return hack.Slice(""), nil
	default:
		log.Error(fmt.Sprintf("FormatValue : invalid type %T", value))
		return nil, fmt.Errorf("invalid type %T", value)
	}
}

func FormatDefaultValue(typeId uint8) ([]byte, error) {
	switch typeId {
	case 8:
		//log.Info("8 bigint")
		return strconv.AppendInt(nil, int64(0), 10), nil
	case 1, 3, 9, 16, 11:
		//log.Info("1 :tinyint, 3 int , 9 mediumint, 16 bit ,")
		return strconv.AppendInt(nil, 0, 10), nil
	case 5, 6, 246:
		//log.Info("5 double, 6 decimal ")
		return strconv.AppendFloat(nil, float64(0.00), 'f', -1, 64), nil
	case 7:
		//log.Info("7 :timestamp")
		return hack.Slice("0000-00-00 00:00:00"), nil
	case 10:
		//log.Info("10 :date")
		return hack.Slice("00:00:00"), nil
	case 252, 253, 12:
		//log.Info("253 string; 12 date,252  tinytext ")
		return hack.Slice(""), nil
	default:
		log.Error(fmt.Sprintf("FormatDefaultValue : invalid type %T %d ", typeId, typeId))
		return nil, fmt.Errorf("invalid type %T", typeId)
	}
}

func GetDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func JsonPrint(d interface{}) {
	b, err := json.Marshal(d)
	log.Info("JsonPrint: \n\terror: ", err, " \n\tdata : ", string(b))
}

func Trim(str string) string {
	str = strings.TrimRight(strings.Trim(strings.Trim(str, "\n"), ""), ";")
	return strings.ToLower(str)
}
