package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	// Synchronous call
	var reply bool
	var reply_i int

	err = client.Call("RemoteList.Append", []interface{}{"teste", 10}, &reply)
	err = client.Call("RemoteList.Append", []interface{}{"teste", 20}, &reply)
	err = client.Call("RemoteList.Append", []interface{}{"1", 100}, &reply)
	err = client.Call("RemoteList.Size", "teste", &reply_i)

	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Tamanho da lista:", reply_i)
	}

	err = client.Call("RemoteList.Size", "1", &reply_i)

	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Tamanho da lista:", reply_i)
	}

	err = client.Call("RemoteList.Remove", "teste", &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Elemento retirado:", reply_i)
	}

	err = client.Call("RemoteList.Get", []interface{}{"teste", 0}, &reply_i)

	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("O valor do index é:", reply_i)
	}
}
