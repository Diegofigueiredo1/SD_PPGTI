package main

import (
	"fmt"
	"net"
	"net/rpc"
	remotelist "ppgti/remotelist/pkg"
)

const filename = "Data_file.json"

func main() {
	list, err := remotelist.NewRemoteList(filename)
	if err != nil {
		fmt.Println("Falha em carregar o RemoteList", err)
		return
	}

	rpcs := rpc.NewServer()
	rpcs.Register(list)
	l, e := net.Listen("tcp", "[localhost]:5000")
	defer l.Close()
	if e != nil {
		fmt.Println("listen error:", e)
	}
	for {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
}
