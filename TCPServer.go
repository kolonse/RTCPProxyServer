// TCPServer
package main

import (
	"KolonseWeb"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

var PROXY_SERVER_MARK []byte = []byte("i'm proxy server")

type Who struct {
	Name string    // 服务者名字
	Ip   string    // 客户端 IP
	conn *net.Conn // 连接套接字
}

func (w *Who) Close() {
	if w.conn != nil {
		(*w.conn).Close()
	}
}

func (w *Who) Dump() string {
	ret := ""
	ret += fmt.Sprintln("\tName:", w.Name)
	ret += fmt.Sprintln("\tIp:", w.Ip)
	return ret
}

func (w *Who) Set(name string, ip string, conn *net.Conn) {
	w.Close()
	w.conn = conn
	w.Name = name
	w.Ip = ip
}

type ServerInfo struct {
	Domain string
	Port   uint16
	Desc   string
	Status string // 服务状态
	ForWho Who
	Parter Who
}

func (si *ServerInfo) GetStatus() string {
	return si.Status
}

func (si *ServerInfo) Dump() string {
	ret := ""
	ret += fmt.Sprintln("Domain:", si.Domain)
	ret += fmt.Sprintln("Port:", si.Port)
	ret += fmt.Sprintln("Desc:", si.Desc)
	ret += fmt.Sprintln("Status:", si.Status)
	ret += fmt.Sprintln("ForWho:")
	ret += si.ForWho.Dump()
	return ret
}

func (si *ServerInfo) ReadLessNByte(conn net.Conn, nBytes int, buff []byte) err {
	totalRecv := 0
	for {
		n, err := conn.Read(buff)
		if err != nil {
			return err
		}
		totalRecv += n
		if n >= nBytes {
			return nil
		}
	}
}

func (si *ServerInfo) handleConnection(conn net.Conn) {
	KolonseWeb.DefaultLogs().Info("Recv Conn,RemoteAddr:%v %v %v %v",
		conn.RemoteAddr().Network(), conn.RemoteAddr().String(),
		conn.LocalAddr().Network(), conn.LocalAddr().String())
	// 收到一个连接 读取开始 i'm proxy server
	buff := make([]byte, 10000) //  缓存长度
	err := si.ReadLessNByte(conn, len(PROXY_SERVER_MARK), buff)
	if err != nil {
		conn.Close()
		return
	}
}

func (si *ServerInfo) Start() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", si.Port))
	if err != nil {
		panic(err)
	}
	KolonseWeb.DefaultLogs().Info("TCP Server Listen On %v", si.Port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go si.handleConnection(conn)
	}
}

func NULLServerInfo() ServerInfo {
	return ServerInfo{
		Status: "have no this server",
	}
}

func NewServerInfo() ServerInfo {
	return ServerInfo{
		Status: "wait",
	}
}

type TypeServerManager map[string]ServerInfo

var TCPServerManager TypeServerManager

func (tm *TypeServerManager) GetServerInfo(domain string) ServerInfo {
	for _, value := range *tm {
		if value.Domain == domain {
			return value
		}
	}
	return NULLServerInfo()
}

func (tm *TypeServerManager) Dump() string {
	ret := ""
	for _, value := range *tm {
		ret += "+++++++++++++++++++++++++++++++++\n"
		ret += value.Dump()
		ret += "=================================\n"
	}
	return ret
}

func (tm *TypeServerManager) TCPServerManagerStart() {
	for _, value := range *tm {
		go value.Start()
	}
}

func LoadCfg() {
	file, err := os.Open("./cfg.rtps")
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&TCPServerManager)
	if err != nil {
		panic(err)
		return
	}
}

func init() {
	TCPServerManager = make(TypeServerManager)
}
