package provider

import (
	"github.com/globalsign/mgo/bson"
)

type GetLogResult struct {
	Names             []string `bson:"names"`
	Ok                int      `bson:"ok"`
	TotalLinesWritten int      `bson:"totalLinesWritten"`
	Log               []string `bson:"log"`
}

type Ping struct {
	Ok int `bson:"ok"`
}

type DatabaseSpecification struct {
	Name       string
	SizeOnDisk int64
	Empty      bool
}

type ListDatabasesResult struct {
	Databases []DatabaseSpecification
	TotalSize int64
}

type RawDocElemCursorData struct {
	//FirstBatch []bson.M `bson:"firstBatch"`
	FirstBatch []interface{} `bson:"firstBatch"`
	//NextBatch  []bson.M      `bson:"nextBatch"`
	NextBatch []interface{} `bson:"nextBatch"`
	NS        string
	Id        int64
}

type FindReply struct {
	Ok     bool
	Code   int
	Errmsg string
	Cursor RawDocElemCursorData
}

type DistinctCmd struct {
	Collection string "distinct"
	Key        string
	Query      interface{} ",omitempty"
}

type CountCmd struct {
	Count string
	Query interface{}
	Limit int32 ",omitempty"
	Skip  int32 ",omitempty"
}

type Database struct {
	Empty      bool   `bson:"empty"`
	Name       string `bson:"name"`
	SizeOnDisk int    `bson:"sizeOnDisk"`
}

type DbListResult struct {
	Databases []Database `bson:"databases"`
	Ok        int        `bson:"ok"`
	TotalSize int        `bson:"totalSize"`
}

type CollectionResult struct {
	Type    string      `bson:"type"`
	Otions  string      `bson:"options"`
	Name    string      `bson:"name"`
	Info    interface{} `bson:"info"`
	IdIndex interface{} `bson:"idIndex"`
}

type CollectionCursor struct {
	FirstBatch []CollectionResult `bson:"firstBatch"`
	NextBatch  []bson.M           `bson:"nextBatch"`
	NS         string             `bson:"ns"`
	Id         int64              `bson:"id"`
}

type ListCollections struct {
	Cursor CollectionCursor `bson:"cursor"`
	Ok     int              `bson:"ok"`
	Code   int
	Errmsg string
}

type DbStatsResult struct {
	AvgObjSize  float64 `json:"avgObjSize"`
	Collections int     `json:"collections"`
	DataSize    int     `json:"dataSize"`
	Db          string  `json:"db"`
	IndexSize   int     `json:"indexSize"`
	Indexes     int     `json:"indexes"`
	NumExtents  int     `json:"numExtents"`
	Objects     int     `json:"objects"`
	Ok          int     `json:"ok"`
	StorageSize int     `json:"storageSize"`
	Views       int     `json:"views"`
}

type BaseReploy struct {
	Ok       int    `json:"ok"`
	Errmsg   string `json:"errmsg"`
	Code     int    `json:"code"`
	CodeName string `json:"codeName"`
}

func Unauthorized() BaseReploy {
	return BaseReploy{
		Ok:       0,
		Errmsg:   "not authorized",
		Code:     13,
		CodeName: "Unauthorized",
	}
}
