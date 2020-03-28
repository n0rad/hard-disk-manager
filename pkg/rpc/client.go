package rpc

import (
	"net"
)

func test() {
	conn, err := net.Dial("tcp", "localhost:8222")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//c := jsonrpc.NewClient(conn)
	//c.ca
}
