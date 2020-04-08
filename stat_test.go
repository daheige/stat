package stat

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestStat(t *testing.T) {
	t.Log("test stat")

	port := 6060

	createHttpService(port)

	ch := make(chan os.Signal, 1)
	log.Println("wait exit signal...")
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal support syscall.SIGUSR2
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch
	log.Println("exit signal: ", sig.String())

	log.Println("shutting down")
}

// createHttpService create a test http server.
// PProf monitoring and metrics service with one port
func createHttpService(port int) {
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	httpMux.HandleFunc("/check", check)

	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello`))
	})

	// test stat prometheus
	httpMux.HandleFunc("/test", testStat)

	// add prometheus metrics handler
	httpMux.Handle("/metrics", promhttp.Handler())

	// http server
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("PProf exec recover: ", err)
			}
		}()

		log.Println("server PProf run on: ", port)

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), httpMux); err != nil {
			log.Println("PProf listen error: ", err)
		}

	}()
}

// check PProf check
func check(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"alive": true}`))
}

// testStat test stat
func testStat(w http.ResponseWriter, r *http.Request) {
	start := PromBeginTime()
	log.Println(start.Format("2006-01-02 15:04:05"))

	// mock cache miss
	CacheMiss.Incr("get_user")

	getUser()

	DBQueryEndTime(start, "user", "user_info")

	w.Write([]byte(`ok`))
}

type user struct {
	Id   int64
	Name string
}

// getUser get user info
func getUser() {
	// DSN data source string: username: password @ protocol (address: port) / database? Parameter = parameter value
	db, err := sql.Open("mysql", "root:root@tcp(192.168.0.11:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	u := user{}

	rows3 := db.QueryRow("select id,name from users where id = ?", 1)

	err = rows3.Scan(&u.Id, &u.Name)
	log.Println("scan error: ", err)

	log.Println("user: ", u)
}

/**
=== RUN   TestStat
2020/01/04 13:33:28 wait exit signal...
2020/01/04 13:33:28 server PProf run on:  6060
2020/01/04 13:34:51 2020-01-04 13:34:51
2020/01/04 13:34:51 scan error:  <nil>
2020/01/04 13:34:51 user:  {1 heige}
2020/01/04 13:34:54 2020-01-04 13:34:54
2020/01/04 13:34:54 scan error:  <nil>
2020/01/04 13:34:54 user:  {1 heige}
2020/01/04 13:34:54 2020-01-04 13:34:54
2020/01/04 13:34:54 scan error:  <nil>
2020/01/04 13:34:54 user:  {1 heige}
2020/01/04 13:34:55 2020-01-04 13:34:55
2020/01/04 13:34:55 scan error:  <nil>
2020/01/04 13:34:55 user:  {1 heige}

browser access
1. http://localhost:6060/test

2.http://localhost:6060/metrics
# HELP go_cache_miss go_cache_miss
# TYPE go_cache_miss counter
go_cache_miss{name="get_user"} 4
# HELP go_db_query go_db_query
# TYPE go_db_query histogram
go_db_query_bucket{method="user",name="user_info",le="0.005"} 0
go_db_query_bucket{method="user",name="user_info",le="0.01"} 0
go_db_query_bucket{method="user",name="user_info",le="0.025"} 0
go_db_query_bucket{method="user",name="user_info",le="0.05"} 0
go_db_query_bucket{method="user",name="user_info",le="0.1"} 0
go_db_query_bucket{method="user",name="user_info",le="0.25"} 0
go_db_query_bucket{method="user",name="user_info",le="0.5"} 0
go_db_query_bucket{method="user",name="user_info",le="1"} 0
go_db_query_bucket{method="user",name="user_info",le="2.5"} 0
go_db_query_bucket{method="user",name="user_info",le="5"} 0
go_db_query_bucket{method="user",name="user_info",le="10"} 0
go_db_query_bucket{method="user",name="user_info",le="+Inf"} 4
go_db_query_sum{method="user",name="user_info"} 74508
go_db_query_count{method="user",name="user_info"} 4
# HELP go_db_query_state go_db_query_state
# TYPE go_db_query_state gauge
go_db_query_state{method="user",name="user_info"} 6500
*/
