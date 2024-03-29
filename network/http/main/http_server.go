package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
)

func echo(w http.ResponseWriter, r *http.Request) {
	log.Printf("Receive request %s from %s", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "hello "+r.URL.Path)
}

// WithServerHeader decorator example
func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("-->WithServerHeader()")
		w.Header().Set("Server", "EchoSever 0.0.1")
		for k, v := range r.Header {
			for _, s := range v {
				w.Header().Add(k, s)
			}
		}

		version := os.Getenv("JAVA_HOME")
		if version != "" {
			w.Header().Set("VERSION", version)
		} else {
			w.Header().Set("VERSION", "UnKnown")
		}

		w.WriteHeader(200)

		h(w, r)
	}
}

// WithAuthCookie example
func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithAuthCookie()")
		cookie := &http.Cookie{Name: "Auth", Value: "Pass", Path: "/"}
		http.SetCookie(w, cookie)
		h(w, r)
	}
}

// HTTPHandlerDecorator for http.HandlerFunc
type HTTPHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func handler(h http.HandlerFunc, decorators ...HTTPHandlerDecorator) http.HandlerFunc {
	for i := range decorators {
		d := decorators[len(decorators)-i-1] // in reverse
		h = d(h)
	}
	return h
}

// post 请求测试
func postHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read request error")
		return
	}
	log.Println("response Body:", string(body))
	_, err = io.WriteString(w, `{"RETCODE":"Success"}`)
	if err != nil {
		log.Println("write response error")
	}
}

func main() {
	// http.ListenAndServe(addr string, handler Handler) 这种方法只有一个url
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler(echo, WithServerHeader, WithAuthCookie))
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/post", postHandler)
	server := &http.Server{
		Addr:    ":8088",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
