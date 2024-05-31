package logger

import (
	"booksapi/api/router"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var lgr *slog.Logger

func init() {
	f, err := os.OpenFile("./log.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	wr := io.MultiWriter(os.Stdout, f)

	lgr = slog.New(slog.NewJSONHandler(wr, nil))
}

func Info(msg string) {
	lgr.Info(msg)
}

func Error(msg string) {
	lgr.Error(msg)
}

func Warn(msg string) {
	lgr.Warn(msg)
}

func Debug(msg string) {
	lgr.Debug(msg)
}

type requestResponseLog struct {
	Req        requestLog  `json:"Req"`
	Resp       responseLog `json:"Resp"`
	StatusCode int         `json:"StatusCode"`
}

type requestLog struct {
	Route  string              `json:"Route"`
	Method string              `json:"Method"`
	Body   string              `json:"Body"`
	Params map[string][]string `json:"Params"`
}

type responseLog struct {
	Header map[string][]string `json:"Header"`
	Body   string              `json:"Body"`
}

func getRequestLog(r *http.Request) requestLog {
	var result requestLog

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	bodyStr := string(body[:])

	result.Body = bodyStr
	result.Method = r.Method

	result.Route = r.URL.String()
	result.Params = r.URL.Query()

	return result
}

func getResponseLog(rww router.ResponseWriterWrapper) responseLog {
	var result responseLog

	var buf bytes.Buffer
	buf.WriteString(rww.Body.String())

	result.Header = (*rww.W).Header()
	result.Body = buf.String()

	return result
}

// returns json
func GetRequestResponseLog(rww router.ResponseWriterWrapper, r *http.Request) string {
	rrl := requestResponseLog{
		Req:        getRequestLog(r),
		Resp:       getResponseLog(rww),
		StatusCode: *(rww.StatusCode),
	}

	bytes, err := json.Marshal(rrl)
	if err != nil {
		panic(fmt.Sprintf("can't decode model for logging, panicking, err: %s", err.Error()))
	}

	return string(bytes[:])
}
