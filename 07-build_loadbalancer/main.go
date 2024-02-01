package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime"
)

type Server interface {
	Address() string // 地址函数
	IsAlive() bool   // 返回布尔值的实时方法
	Serve(w http.ResponseWriter, r *http.Request)
}

// 简易服务器
type simpleServer struct {
	addr  string                 // 地址
	proxy *httputil.ReverseProxy // 代理
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	port            string   // 端口
	roundRobinCount int      // 循环计数
	Servers         []Server // 服务器集合
}

func newSimpleServer(addr string) *simpleServer {
	// 解析地址 并将其转换为一个 *url.URL 类型的对象
	// serverUrl: 目标服务器地址
	serverUrl, err := url.Parse(addr)
	handlerErr(err)
	return &simpleServer{
		addr: addr,
		// 创建一个反向代理对象 (反向代理的作用在 README 中有说明)
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0, // 计数从 0 开始
		Servers:         servers,
	}
}

func handlerErr(err error) {
	if err != nil {
		// 打印堆栈信息 显示文件名和行号
		_, file, line, _ := runtime.Caller(1)
		log.Printf("error at %s:%d: %v", file, line, err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string { return s.addr }
func (s *simpleServer) IsAlive() bool   { return true }
func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

// 下一个可用服务器
func (l *LoadBalancer) getNexAvailableServer() Server {
	server := l.Servers[l.roundRobinCount%len(l.Servers)]
	for !server.IsAlive() {
		l.roundRobinCount++
		server = l.Servers[l.roundRobinCount%len(l.Servers)]
	}
	l.roundRobinCount++
	return server
}

// 反向代理
func (l *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := l.getNexAvailableServer()
	fmt.Printf("请求转发到地址: %q\n", targetServer.Address())
	targetServer.Serve(w, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.baidu.com"),
		newSimpleServer("https://www.bing.com"),
		newSimpleServer("https://www.sogou.com"),
	}

	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}

	// 当访问根路由时开始负载均衡
	http.HandleFunc("/", handleRedirect)
	fmt.Printf("serveing requests at 'localhost:%s'\n", lb.port)

	addr := http.ListenAndServe("127.0.0.1:"+lb.port, nil)
	fmt.Printf("访问的网络地址为: %v", addr)
}
