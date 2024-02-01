package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandle(w http.ResponseWriter, r *http.Request) {
	// 解析表单
	if err := r.ParseForm(); err != nil {
		if _, err = fmt.Fprintf(w, "表单解析出错: %v\n", err); err != nil {
			log.Printf("表单响应出错: %v\n", err)
			return
		}
	}
	_, _ = fmt.Fprintf(w, "响应成功\n")
	name := r.FormValue("name")
	age := r.FormValue("age")
	_, _ = fmt.Fprintf(w, "名字是: %s\n", name)
	_, _ = fmt.Fprintf(w, "年龄是: %v\n", age)
}

func helloHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		_, err := fmt.Fprintf(w, "请使用 GET 方法请求(请求方法错误)")
		if err != nil {
			return
		}
	}
	_, err := fmt.Fprintf(w, "Hello!")
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", helloHandle)
	http.HandleFunc("/form", formHandle)

	fmt.Printf("启动端口为 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
