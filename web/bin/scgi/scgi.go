package scgi

// https://bitbucket.org/mjl/scgi/src

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cgi"
	"os"
	"runtime"
	"strconv"
)

type responseWriter struct {
	wroteHeader bool
	header      http.Header
	conn        net.Conn
}

func (rw *responseWriter) Header() http.Header {
	return rw.header
}

func (rw *responseWriter) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.conn.Write(buf)
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		log.Println("scgi: multiple response.WriteHeader calls")
		return
	}
	rw.wroteHeader = true
	fmt.Fprintf(rw.conn, "Status: %d %s\r\n", code, http.StatusText(code))
	rw.header.Write(rw.conn)
	rw.conn.Write([]byte("\r\n"))
}

type limitedReadCloser struct {
	*io.LimitedReader
}

func (r *limitedReadCloser) Close() error {
	return nil
}

func serve(rw net.Conn, handler http.Handler) {
	defer func() {
		if err := recover(); err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("scgi: panic serving request: %v\n%s", err, buf)
		}
	}()
	defer rw.Close()

	// read request
	b := []byte{}
	c := make([]byte, 1)
Length:
	for {
		_, err := rw.Read(c)
		if err != nil {
			return
		}
		switch c[0] {
		case ':':
			break Length
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b = append(b, c[0])
		}
	}
	n, err := strconv.Atoi(string(b))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, n+1)
	_, err = io.ReadAtLeast(rw, buf, len(buf))
	if err != nil {
		return
	}
	if buf[len(buf)-1] != ',' {
		log.Println("scgi: invalid netstring in request")
		return
	}

	t := bytes.Split(buf[:len(buf)-1], []byte{0})
	if len(t)%2 != 1 || len(t[len(t)-1]) != 0 {
		log.Println("scgi: malformed request")
		return
	}
	params := map[string]string{}
	for i := 0; i+1 < len(t); i += 2 {
		k := string(t[i])
		v := string(t[i+1])
		params[k] = v
	}

	params["SERVER_PROTOCOL"] = "HTTP/1.1"
	req, err := cgi.RequestFromMap(params)
	if err != nil {
		log.Println("scgi: parsing cgi environment:", err.Error())
		return
	}

	req.Body = &limitedReadCloser{io.LimitReader(rw, req.ContentLength).(*io.LimitedReader)}

	w := &responseWriter{false, make(http.Header), rw}
	handler.ServeHTTP(w, req)
}

func Serve(l net.Listener, handler http.Handler) error {
	if l == nil {
		var err error
		l, err = net.FileListener(os.Stdin)
		if err != nil {
			return err
		}
		defer l.Close()
	}
	if handler == nil {
		handler = http.DefaultServeMux
	}
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		go serve(rw, handler)
	}
}
