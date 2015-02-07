package main

import (
	"./app"
	"./config"
	"./scgi"
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"
	"time"
)

const (
	PROTOCOL_HTTP = "http"
	PROTOCOL_FCGI = "fcgi"
	PROTOCOL_SCGI = "scgi"
)

type XcgiHandler struct{}

func printProcInfo() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	hostname, _ := os.Hostname()
	fmt.Printf("Version: %s, Goroutine: %d, NumCPU: %d, Hostname: %s\n",
		runtime.Version(), runtime.NumGoroutine(), runtime.NumCPU(), hostname)
}

func main() {

	printProcInfo()

	conf := config.GetConfig()

	protocol := conf.Protocol

	listener, err := net.Listen("tcp", conf.Bind)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(conf.Bind, protocol)

	switch protocol {
	case PROTOCOL_FCGI:

		var fcgi_handler XcgiHandler
		fcgi.Serve(listener, fcgi_handler)

	case PROTOCOL_SCGI:

		var scgi_handler XcgiHandler
		scgi.Serve(listener, scgi_handler)

	case PROTOCOL_HTTP:

		http.HandleFunc("/", handler)
		http_server := &http.Server{
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 16,
		}
		http_server.Serve(listener)

	default:

		fmt.Println("Protocol '" + protocol + "' not support")
	}
}

func (XcgiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.Handler(res, req)
}

func handler(res http.ResponseWriter, req *http.Request) {
	app.Handler(res, req)
}
