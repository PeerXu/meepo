package pprof

import (
	"net/http"
	_ "net/http/pprof"
)

func Setup(lisStr string) {
	go http.ListenAndServe(lisStr, nil)
}
