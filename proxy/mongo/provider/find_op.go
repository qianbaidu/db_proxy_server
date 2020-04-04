package provider

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/qianbaidu/db-proxy/util"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"reflect"
	"strings"
)

var RunCmdTimeoutError = errors.New("operation exceeded time limit")

type FindOp struct {
	Find                string      `bson:"find,omitempty"`
	Filter              interface{} `bson:"filter,omitempty"`
	Sort                interface{} `bson:"sort,omitempty"`
	Projection          interface{} `bson:"projection,omitempty"`
	Hint                interface{} `bson:"hint,omitempty"`
	Skip                int         `bson:"skip,omitempty"`
	Limit               int         `bson:"limit,omitempty"`
	Batchsize           int         `bson:"batchSize,omitempty"`
	Singlebatch         bool        `bson:"singleBatch,omitempty"`
	Comment             string      `bson:"comment,omitempty"`
	Maxtimems           int         `bson:"maxTimeMS,omitempty"`
	Readconcern         interface{} `bson:"readConcern,omitempty"`
	Max                 interface{} `bson:"max,omitempty"`
	Min                 interface{} `bson:"min,omitempty"`
	Returnkey           bool        `bson:"returnKey,omitempty"`
	Showrecordid        bool        `bson:"showRecordId,omitempty"`
	Tailable            bool        `bson:"tailable,omitempty"`
	Oplogreplay         bool        `bson:"oplogReplay,omitempty"`
	Nocursortimeout     bool        `bson:"noCursorTimeout,omitempty"`
	Awaitdata           bool        `bson:"awaitData,omitempty"`
	Allowpartialresults bool        `bson:"allowPartialResults,omitempty"`
	//Collation           interface{} `bson:"collation,omitempty"`
}

type GetMoreOp struct {
	GetMore    int64  `bson:"getMore"`
	Collection string `bson:"collection"`
	BatchSize  int    `bson:"batchSize"`
	MaxTimeMS  int    `bson:"maxTimeMS"`
}

func NewFindOp() *FindOp {
	var m bson.M
	return &FindOp{
		Find:                "$cmd",
		Filter:              nil,
		Sort:                m,
		Projection:          m,
		Hint:                m,
		Skip:                0,
		Limit:               0,
		Batchsize:           0,
		Singlebatch:         false,
		Comment:             "",
		Maxtimems:           util.DEFAULT_EXEC_SQL_TIME_OUT,
		Readconcern:         m,
		Max:                 m,
		Min:                 m,
		Returnkey:           false,
		Showrecordid:        false,
		Tailable:            false,
		Oplogreplay:         false,
		Nocursortimeout:     false,
		Awaitdata:           false,
		Allowpartialresults: false,
		//Collation:           m,
	}
}

// doc https://docs.mongodb.com/manual/reference/command/find/
func BindFindOp(query bson.D) *FindOp {
	f := NewFindOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			fmt.Println("to bson d ok, k : ", k, " v : ", m.Map())
			v = m.Map()
		}
		switch strings.ToLower(k) {
		case "find":
			f.Find = v.(string)
		case "filter":
			f.Filter = v
		case "sort":
			f.Sort = v
		case "projection":
			f.Projection = v
		case "hint":
			f.Hint = v
		case "skip":
			f.Skip = v.(int)
		case "limit":
			f.Limit = v.(int)
		case "batchsize":
			f.Batchsize = v.(int)
		case "singlebatch":
			f.Singlebatch = v.(bool)
		case "comment":
			f.Comment = v.(string)
		case "maxtimems":
			// doc https://docs.mongodb.com/v3.4/reference/method/cursor.maxTimeMS/index.html
			if timeout := v.(int); timeout < 1000 {
				f.Maxtimems = util.DEFAULT_EXEC_SQL_TIME_OUT
			} else {
				f.Maxtimems = timeout
			}
		case "readconcern":
			f.Readconcern = v
		case "max":
			f.Max = v
		case "min":
			f.Min = v
		case "returnkey":
			f.Returnkey = v.(bool)
		case "showrecordid":
			f.Showrecordid = v.(bool)
		case "tailable":
			f.Tailable = v.(bool)
		case "oplogreplay":
			f.Oplogreplay = v.(bool)
		case "nocursortimeout":
			f.Nocursortimeout = v.(bool)
		case "awaitdata":
			f.Awaitdata = v.(bool)
		case "allowpartialresults":
			f.Allowpartialresults = v.(bool)
		//case "collation":
		//	f.Collation = v
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	//log.Info("findOp:")
	//util.JsonPrint(f)
	return f
}

func NewGetMoreOp() *GetMoreOp {
	return &GetMoreOp{
		GetMore:    int64(1),
		Collection: "",
		BatchSize:  10,
		MaxTimeMS:  0,
	}
}

// doc https://docs.mongodb.com/v3.4/reference/command/getMore/
func BindGetMoreOp(query bson.D) *GetMoreOp {
	f := NewGetMoreOp()
	queryMap := query.Map()
	for k, v := range queryMap {
		switch strings.ToLower(k) {
		case "getmore":
			f.GetMore = v.(int64)
		case "collection":
			f.Collection = v.(string)
		case "batchsize":
			size := v.(int)
			log.Info("batchsize size : ", size)
			if size < 1 {
				size = 10
			}
			f.BatchSize = size
		case "maxtimems":
			f.MaxTimeMS = v.(int)
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

type LastError struct {
	ConnectionId int    `bson:"connectionId"`
	Err          string `bson:"err"`
	N            int    `bson:"n"`
	OK           int    `bson:"ok"`
	SyncMillis   int    `bson:"syncMillis"`
	WrittenTo    string `bson:"writtenTo"`
}

func NewLastError() LastError {
	return LastError{
		ConnectionId: 1,
		OK:           1,
		Err:          "",
		N:            0,
		SyncMillis:   0,
		WrittenTo:    "",
	}
}

type DefaultServerReturn struct {
	Ntoreturn  int    `bson:"ntoreturn"`
	Ntoskip    int    `bson:"ntoskip"`
	Ok         int    `bson:"ok"`
	Err        string `bson:"errmsg"`
	SyncMillis int    `bson:"syncMillis"`
	WrittenTo  string `bson:"writtenTo"`
}

func NewDefaultServerReturn() DefaultServerReturn {
	return DefaultServerReturn{
		Ntoreturn: 1,
		Ntoskip:   0,
		Ok:        1,
	}
}

type BuildInfo struct {
	Version        string
	VersionArray   []int  `bson:"versionArray"` // On MongoDB 2.0+; assembled from Version otherwise
	GitVersion     string `bson:"gitVersion"`
	OpenSSLVersion string `bson:"OpenSSLVersion"`
	SysInfo        string `bson:"sysInfo"` // Deprecated and empty on MongoDB 3.2+.
	Bits           int
	Debug          bool
	MaxObjectSize  int `bson:"maxBsonObjectSize"`
	Ok             int `bson:"ok"`
}

func NewBuildInfo() BuildInfo {
	return BuildInfo{
		Version:        "3.4.23",
		VersionArray:   []int{3, 4, 23, 0},
		GitVersion:     "324017ede1dbb1c9554dd2dceb15f8da3c59d0e8",
		OpenSSLVersion: "",
		SysInfo:        "",
		Bits:           64,
		Debug:          false,
		MaxObjectSize:  16777216,
		Ok:             1,
	}
}

type IsMaster struct {
	IsMaster                     bool   `bson:"ismaster"`
	MaxBsonObjectSize            int    `bson:"maxBsonObjectSize"`
	MaxMessageSizeBytes          int    `bson:"maxMessageSizeBytes"`
	MaxWriteBatchSize            int    `bson:"maxWriteBatchSize"`
	LocalTime                    string `bson:"localTime"`
	LogicalSessionTimeoutMinutes int    `bson:"logicalSessionTimeoutMinutes"`
	MinWireVersion               int    `bson:"minWireVersion"`
	MaxWireVersion               int    `bson:"maxWireVersion"`
	ReadOnly                     bool   `bson:"readOnly"`
	Ok                           bool   `bson:"ok"`
}

func NewIsMaster() IsMaster {
	return IsMaster{
		IsMaster:            true,
		MaxBsonObjectSize:   16777216,
		MaxMessageSizeBytes: 48000000,
		MaxWriteBatchSize:   100000,
		//"localTime", ISODate("2017-12-14T17:40:28.640Z"),
		LogicalSessionTimeoutMinutes: 30,
		MinWireVersion:               0,
		MaxWireVersion:               5,
		ReadOnly:                     false,
		Ok:                           true,
	}
}

type ListCursor struct {
	Name    string      `bson:"name"`
	Ctype   string      `bson:"type"`
	Info    interface{} `bson:"info"`
	IdIndex string      `bson:"idIndex"`
}
type ListCollection struct {
	ListCollections       float64     `bson:"listCollections"`
	Filter                interface{} `bson:"filter,omitempty"`
	NameOnly              bool        `bson:"nameOnly,omitempty"`
	AuthorizedCollections bool        `bson:"authorizedCollections,omitempty"`
	//Cursor                ListCursor  `bson:"cursor,omitempty"`
}

func NewListCollection() *ListCollection {
	return &ListCollection{
		ListCollections:       1,
		Filter:                nil,
		NameOnly:              false,
		AuthorizedCollections: false,
	}
}

// doc https://docs.mongodb.com/v3.4/reference/command/listCollections/
func BindListCollectionOp(query bson.D) *ListCollection {
	f := NewListCollection()
	queryMap := query.Map()
	log.Info("BindListCollectionOp queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			v = m.Map()
		}
		switch strings.ToLower(k) {
		case "listcollections":
			if r, ok := v.(float64); ok && r > 0 {
				f.ListCollections = r
			}
		case "filter":
			f.Filter = v
		case "cursor":
		//	f.Cursor = v.(ListCursor)
		case "nameOnly":
			f.NameOnly = v.(bool)
		case "authorizedCollections":
			f.AuthorizedCollections = v.(bool)
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	log.Debug("BindListCollectionOp end . ", f)
	return f
}

type Explain struct {
	Explain   bson.D `bson:"explain"`
	Verbosity string `bson:"verbosity"`
}

type DistinctOp struct {
	Distinct    string      `bson:"distinct"`
	Key         string      `bson:"key"`
	Query       interface{} `bson:"query"`
	ReadConcern interface{} `bson:"readConcern"`
}

func NewDistinctOp() DistinctOp {
	return DistinctOp{
		Distinct:    "$cmd",
		Key:         "",
		Query:       nil,
		ReadConcern: bson.D{},
	}
}

// doc https://docs.mongodb.com/manual/reference/command/distinct/
func BindDistinctOp(query bson.D) DistinctOp {
	f := NewDistinctOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			v = m.Map()
		}
		switch strings.ToLower(k) {
		case "distinct":
			f.Distinct = v.(string)
		case "key":
			f.Key = v.(string)
		case "query":
			f.Query = v
		case "readconcern":
			f.ReadConcern = v
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

type CountOp struct {
	Count       string      `bson:"count"`
	Query       interface{} `bson:"query"`
	Limit       int         `bson:"limit"`
	Skip        int         `bson:"skip"`
	Hint        interface{} `bson:"hint"`
	ReadConcern interface{} `bson:"readConcern"`
}

func NewCountOp() CountOp {
	return CountOp{
		Count:       "",
		Query:       bson.D{},
		Limit:       10,
		Skip:        0,
		Hint:        bson.D{},
		ReadConcern: bson.D{},
	}
}

// doc https://docs.mongodb.com/v3.4/reference/command/count/
func BindCountOp(query bson.D) CountOp {
	f := NewCountOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			v = m.Map()
		}
		switch strings.ToLower(k) {
		case "count":
			f.Count = v.(string)
		case "query":
			f.Query = v
		case "limit":
			f.Limit = v.(int)
		case "skip":
			f.Skip = v.(int)
		case "hint":
			f.Hint = v
		case "readconcern":
			f.ReadConcern = v
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

type ListIndexOp struct {
	ListIndexes string      `bson:"listIndexes"`
	Cursor      interface{} `bson:"cursor,omitempty"`
}

func NewListIndexOp() ListIndexOp {
	return ListIndexOp{
		ListIndexes: "$cmd",
	}
}

// doc https://docs.mongodb.com/v3.4/reference/command/listIndexes/
func BindListIndexOp(query bson.D) ListIndexOp {
	f := NewListIndexOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		switch k {
		case "listIndexes":
			f.ListIndexes = v.(string)
		case "cursor":
			f.Cursor = v
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

type AggregateOp struct {
	Aggregate                string        `bson:"aggregate"`
	Pipeline                 []interface{} `bson:"pipeline"`
	Explain                  bool          `bson:"explain"`
	AllowDiskUse             bool          `bson:"allowDiskUse"`
	Cursor                   interface{}   `bson:"cursor"`
	MaxTimeMS                int           `bson:"maxTimeMS"`
	BypassDocumentValidation bool          `bson:"bypassDocumentValidation"`
	ReadConcern              interface{}   `bson:"readConcern"`
	Collation                interface{}   `bson:"collation"`
}

func NewAggregateOp() AggregateOp {
	d := bson.D{}
	p := make([]interface{}, 0)
	a := AggregateOp{
		Aggregate:                "$cmd",
		Pipeline:                 p,
		Explain:                  false,
		AllowDiskUse:             false,
		Cursor:                   d,
		MaxTimeMS:                30,
		BypassDocumentValidation: false,
		ReadConcern:              d,
		Collation:                d,
	}
	return a
}

// doc https://docs.mongodb.com/v3.4/reference/method/db.collection.aggregate/index.html
func BindAggregateOp(query bson.D) AggregateOp {
	f := NewAggregateOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			v = m.Map()
		}
		switch k {
		case "aggregate":
			f.Aggregate = v.(string)
		case "pipeline":
			f.Pipeline = v.([]interface{})
		case "explain":
			f.AllowDiskUse = v.(bool)
		case "allowDiskUse":
			f.AllowDiskUse = v.(bool)
		case "cursor":
			f.Cursor = v
		case "maxTimeMS":
			f.MaxTimeMS = v.(int)
		case "bypassDocumentValidation":
			f.BypassDocumentValidation = v.(bool)
		case "readConcern":
			f.ReadConcern = v
		case "collation":
			f.Cursor = v
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

type MapReduceOp struct {
	MapReduce string      `bson:"mapReduce"`
	Map       interface{} `bson:"map"`
	Reduce    interface{} `bson:"reduce"`
	Out       interface{} `bson:"out"`
	Query     interface{} `bson:"query"`
	Sort      interface{} `bson:"sort"`
	Limit     int         `bson:"limit"`
	Finalize  interface{} `bson:"finalize"`
	Scope     interface{} `bson:"scope"`
	JsMode    bool        `bson:"jsMode"`
	Collation interface{} `bson:"collation"`
}

func NewMapReduceOp() MapReduceOp {
	d := bson.D{}
	return MapReduceOp{
		MapReduce: "$cmd",
		Map:       nil,
		Reduce:    nil,
		Out:       nil,
		Query:     nil,
		Sort:      nil,
		Limit:     10,
		Finalize:  nil,
		Scope:     nil,
		JsMode:    false,
		Collation: d,
	}
}

// doc https://docs.mongodb.com/v3.4/reference/method/db.collection.mapReduce/
func BindMapReduceOp(query bson.D) MapReduceOp {
	f := NewMapReduceOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			v = m.Map()
		}
		switch k {
		case "mapreduce":
			f.MapReduce = v.(string)
		case "map":
			f.Map = v
		case "reduce":
			f.Reduce = v
		case "out":
			f.Out = v
		case "query":
			f.Query = v
		case "sort":
			f.Sort = v
		case "limit":
			f.Limit = v.(int)
		case "finalize":
			f.Finalize = v
		case "scope":
			f.Sort = v
		case "jsMode":
			f.JsMode = v.(bool)
		case "collation":
			f.Collation = v
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

func NewDbStatsResult() DbStatsResult {
	return DbStatsResult{
		AvgObjSize:  93.20588235294117,
		Collections: 1,
		DataSize:    3169,
		Db:          "local",
		IndexSize:   229376,
		Indexes:     11,
		NumExtents:  0,
		Objects:     34,
		Ok:          1,
		StorageSize: 176128,
		Views:       0,
	}
}

type DbStatsOp struct {
	DbStats float64 `bson:"dbStats"`
	Scale   int     `bson:"scale"`
}

func NewDbStatsOp() DbStatsOp {
	return DbStatsOp{
		DbStats: 1,
		Scale:   1,
	}
}

// doc https://docs.mongodb.com/v3.4/reference/command/dbStats/
func BindDbStatsOp(query bson.D) DbStatsOp {
	f := NewDbStatsOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		if m, ok := v.(bson.D); ok {
			v = m.Map()
		}
		switch k {
		case "dbStats":
			if r, ok := v.(float64); ok && r > 0 {
				f.DbStats = r
			}
		case "scale":
			f.Scale = v.(int)
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}

type CollStatsOp struct {
	CollStats string `bson:"collStats"`
	Scale     int    `bson:"scale,omitempty"`
	Verbose   bool   `bson:"verbose,omitempty"`
}

func NewCollStatsOp() CollStatsOp {
	return CollStatsOp{
		CollStats: "$cmd",
	}
}

// doc https://docs.mongodb.com/v3.4/reference/command/collStats/
func BindCollStatsOp(query bson.D) CollStatsOp {
	f := NewCollStatsOp()
	queryMap := query.Map()
	log.Info("queryMap")
	util.JsonPrint(queryMap)
	for k, v := range queryMap {
		switch k {
		case "collStats":
			f.CollStats = v.(string)
		case "scale":
			f.Scale = v.(int)
		case "verbose":
			f.Verbose = v.(bool)
		default:
			log.Error(fmt.Sprintf("not supported. type : %v , key : %s , value : %v ", reflect.TypeOf(v), k, v))
		}
	}
	return f
}
