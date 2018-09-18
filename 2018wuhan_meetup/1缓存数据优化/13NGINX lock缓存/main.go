package main

import (
	"fmt"
	"io"
	"net/http"
	time "time"
)

func main() {
	http.Handle("/hello", &ServeMux{})
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("err", err.Error())
	}
}

type ServeMux struct {
}

func (p *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)

	fmt.Println("get one request")
	fmt.Println(r.RequestURI)
	io.WriteString(w, "hello world")
}
