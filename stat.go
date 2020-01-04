// Package stat prometheus for http service,redis,mysql and rpc.
// the stat package is mainly used for prometheus to do RBI monitoring services
// for the Go program. Some metrics that need to be monitored are pulled through
// prometheus metrics, which can be displayed in real
// time through the grafana graphical tool.
package stat

import (
	"github.com/daheige/stat/prom"
	"time"
)

// Stat interface.
type Stat interface {
	Timing(name string, time int64, extra ...string)
	Incr(name string, extra ...string) // name,ext...,code
	State(name string, val int64, extra ...string)
}

// default stat struct,you can add other stat.
var (
	// http client and server
	HTTPClient Stat = prom.HTTPClient
	HTTPServer Stat = prom.HTTPServer

	// redis and db
	Cache Stat = prom.LibClient
	DB    Stat = prom.LibClient

	// db query
	DBQuery Stat = prom.DBQuery

	// cache hit
	CacheHit Stat = prom.CacheHit

	// cache miss
	CacheMiss Stat = prom.CacheMiss

	// rpc
	RPCClient Stat = prom.RPCClient
	RPCServer Stat = prom.RPCServer
)

// PromBeginTime prometheus start time.
func PromBeginTime() time.Time{
	return time.Now()
}

// PromTimeSub return end time - start time interval.
func PromTimeSub(promTime time.Time) int64 {
	t := time.Now().Sub(promTime).Microseconds() // time interval us
	return t
}

// DBQueryEndTime db query end time - start time.
func DBQueryEndTime(promTime time.Time,method string,name string){
	t := PromTimeSub(promTime)

	DBQuery.Timing(method,t,name)
	DBQuery.State(method,t,name)
}
