// RTCPProxyServer project RTCPProxyServer.go
package main

import (
	"KolonseWeb"
	"KolonseWeb/HttpLib"
	"KolonseWeb/Type"
)

func main() {
	LoadCfg()
	TCPServerManager.TCPServerManagerStart()
	KolonseWeb.DefaultLogs().Info("加载代理服务配置:\n%v", TCPServerManager.Dump())
	KolonseWeb.DefaultApp.Get("/GetPort", func(req *HttpLib.Request, res *HttpLib.Response, next Type.Next) {
		domain := req.URL.Query().Get("domain")
		KolonseWeb.DefaultLogs().Info("处理客户端请求 Req Domain:%v,Client Addr:%v", domain, req.RemoteAddr)
		serverInfo := TCPServerManager.GetServerInfo(domain)
		res.Json(serverInfo) // 返回服务状态
	})
	KolonseWeb.DefaultApp.Listen("0.0.0.0", *Port)
}
