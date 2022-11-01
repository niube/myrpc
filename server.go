package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
)

func main() {
	s := NewServer()
	s.Register("hello", hello)
	s.Run()
}

func hello(args ...interface{}) {
	for _, v := range args {
		fmt.Println("hello ", v)
	}

}

type rpcData struct {
	Name string        // 函数名
	Args []interface{} // 参数
}

type server struct {
	conn net.Conn                 // socket连接
	maps map[string]reflect.Value // 函数字典
}

func NewServer() *server {
	return &server{
		maps: make(map[string]reflect.Value),
	}
}

// Register 注册函数
func (s *server) Register(fnName string, fun interface{}) {
	if _, ok := s.maps[fnName]; !ok {
		s.maps[fnName] = reflect.ValueOf(fun)
	}
}

// 开启socket连接，处理请求
func (s *server) Run() {
	listen, err := net.Listen("tcp4", ":3001")
	if err != nil {
		panic(err)
	}

	for {
		s.conn, err = listen.Accept()
		if err != nil {
			fmt.Println("监听异常： ", err)
			continue
		}
		s.handleConnect()
	}
}

// 处理消息
func (s *server) handleConnect() {
	for {
		// 前四个字节放置消息长度，知道才能正常解码消息
		header := make([]byte, 4)
		if _, err := s.conn.Read(header); err != nil {
			if err.Error() == "EOF" {
				fmt.Println("读取完毕：200")
				break
			}
			fmt.Println("读取数据失败： ", err)
			continue
		}

		bodyLen := binary.BigEndian.Uint32(header)
		body := make([]byte, bodyLen)
		if _, err := s.conn.Read(body); err != nil {
			fmt.Println("读取body失败： ", err)
			continue
		}

		var req rpcData
		if err := json.Unmarshal(body, &req); err != nil {
			fmt.Println("序列化body失败： ", err)
			continue
		}

		inArgs := make([]reflect.Value, len(req.Args))
		for i := range req.Args {
			inArgs[i] = reflect.ValueOf(req.Args[i])
		}

		fn := s.maps[req.Name]
		fn.Call(inArgs)
	}
}
