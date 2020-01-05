package stats

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Stats data structure
type Stats struct {
	mu                  sync.RWMutex
	Hostname            string
	Uptime              time.Time
	Pid                 int
	TotalResponseCounts map[string]map[string]int
	TotalResponseTime   map[string]map[string]time.Time
}

type customWriter struct {
	http.ResponseWriter
	status   int
	method   string
	basePath string
	length   int
}

func NewCustomWriter(w http.ResponseWriter) *customWriter {
	return &customWriter{ResponseWriter: w}
}

func (w *customWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *customWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// New constructs a new Stats structure
func New() *Stats {
	name, _ := os.Hostname()

	stats := &Stats{
		Uptime:              time.Now(),
		Pid:                 os.Getpid(),
		TotalResponseCounts: map[string]map[string]int{},
		TotalResponseTime:   map[string]map[string]time.Time{},
		Hostname:            name,
	}

	return stats
}

// Handler is a MiddlewareFunc makes Stats implement the Middleware interface.
func (mw *Stats) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beginning, sw := mw.Begin(w)
		sw.basePath = strings.Split(r.URL.Path, "/")[1]
		sw.method = r.Method + ":" + sw.basePath
		h.ServeHTTP(sw, r)

		mw.End(beginning, *sw)
	})
}

// Begin starts a recorder
func (mw *Stats) Begin(w http.ResponseWriter) (time.Time, *customWriter) {
	start := time.Now()
	sw := NewCustomWriter(w)
	return start, sw
}

// End closes the recorder with a specific status
func (mw *Stats) End(start time.Time, sw customWriter) {

	responseTime := time.Since(start)

	mw.mu.Lock()

	defer mw.mu.Unlock()

	if sw.status != 0 {
		statusCode := fmt.Sprintf("%d", sw.status)
		if _, ok := mw.TotalResponseCounts[sw.method]; ok {
			mw.TotalResponseCounts[sw.method][statusCode]++
		} else {
			mw.TotalResponseCounts[sw.method] = map[string]int{statusCode: 1}
		}
		if _, ok := mw.TotalResponseTime[sw.method]; ok {
			mw.TotalResponseTime[sw.method][statusCode] = mw.TotalResponseTime[sw.method][statusCode].Add(responseTime)
		} else {
			mw.TotalResponseTime[sw.method] = map[string]time.Time{statusCode: time.Time{}.Add(responseTime)}
		}
	}
}

// Data serializable structure
type Data struct {
	Pid                        int                           `json:"pid"`
	UpTime                     string                        `json:"uptime"`
	UpTimeSec                  float64                       `json:"uptime_sec"`
	Time                       string                        `json:"time"`
	TimeUnix                   int64                         `json:"unixtime"`
	TotalMethodStatusCodeCount map[string]map[string]int     `json:"total_method_status_code_count"`
	TotalCount                 int                           `json:"total_count"`
	TotalMethodResponseTime    map[string]map[string]float64 `json:"total_method_response_time_sec"`
	TotalResponseTime          float64                       `json:"total_response_time_sec"`
	AverageMethodResponseTime  map[string]map[string]float64 `json:"average_method_response_time_sec"`
	AverageResponseTime        float64                       `json:"average_response_time_sec"`
}

// Data returns the data serializable structure
func (mw *Stats) Data() *Data {
	mw.mu.RLock()

	totalResponseCounts := make(map[string]map[string]int)
	totalResponseTimes := make(map[string]map[string]float64)
	averageResponseTimes := make(map[string]map[string]float64)

	now := time.Now()

	uptime := now.Sub(mw.Uptime)

	totalCount := 0
	for method, codeMap := range mw.TotalResponseCounts {
		for code, count := range codeMap {
			if _, ok := totalResponseCounts[method]; ok {
				totalResponseCounts[method][code] = count
			} else {
				m := make(map[string]int)
				m[code] = count
				totalResponseCounts[method] = m
			}
			totalCount += count
		}
	}
	totalResponseTime := float64(0)
	averageResponseTime := float64(0)
	for method, codeMap := range mw.TotalResponseTime {
		for code, atime := range codeMap {
			//fmt.Println("TotalResponseTime", method, code, atime)
			if _, ok := totalResponseTimes[method]; ok {
				totalResponseTimes[method][code] = atime.Sub(time.Time{}).Seconds()
			} else {
				me := make(map[string]float64)
				me[code] = atime.Sub(time.Time{}).Seconds()
				totalResponseTimes[method] = me
			}
			totalResponseTime += atime.Sub(time.Time{}).Seconds()
			if totalResponseCounts[method][code] > 0 {
				avgNs := atime.Sub(time.Time{}).Seconds() / float64(totalResponseCounts[method][code])
				if _, ok := averageResponseTimes[method]; ok {
					averageResponseTimes[method][code] = atime.Sub(time.Time{}).Seconds()
				} else {
					t := make(map[string]float64)
					t[code] = avgNs
					averageResponseTimes[method] = t
				}
			}
		}
	}
	averageResponseTime = totalResponseTime / float64(totalCount)

	mw.mu.RUnlock()
	r := &Data{
		Pid:                        mw.Pid,
		UpTime:                     uptime.String(),
		UpTimeSec:                  uptime.Seconds(),
		Time:                       now.String(),
		TimeUnix:                   now.Unix(),
		TotalMethodStatusCodeCount: totalResponseCounts,
		TotalCount:                 totalCount,
		TotalMethodResponseTime:    totalResponseTimes,
		TotalResponseTime:          totalResponseTime,
		AverageMethodResponseTime:  averageResponseTimes,
		AverageResponseTime:        averageResponseTime,
	}
	return r
}
