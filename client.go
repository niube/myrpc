package main

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

func main() {
	req := rpcRequestData{Name: "hello", Args: []interface{}{"Hiker", "zhumang"}}
	rpcCall(req)
}

func rpcCall(data rpcRequestData) {
	conn, err := net.Dial("tcp4", "127.0.0.1:3001")
	if err != nil {
		panic(err)
	}

	req, err2 := json.Marshal(data)
	if err2 != nil {
		panic(err2)
	}

	// 预留4个字节放置消息长度
	buf := make([]byte, 4+len(req))
	// 往buffer写入长度
	binary.BigEndian.PutUint32(buf[:4], uint32(len(req)))
	// 从第四个开始写入请求数据
	copy(buf[4:], req)
	_, err3 := conn.Write(buf)
	if err3 != nil {
		panic(err3)
	}
}

type rpcRequestData struct {
	Name string        // 函数名
	Args []interface{} // 参数
}
