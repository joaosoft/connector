package main

import (
	"encoding/json"
	"fmt"
	"github.com/joaosoft/connector"
)

func main() {
	serverManager, err := connector.NewServerManager()

	// server 1
	server1, err := connector.NewServer(connector.WithServerName("server_one"))
	if err != nil {
		panic(err)
	}
	server1.AddMethod("sayHello", HandlerSayHello)

	// server 2
	server2, err := connector.NewServer(connector.WithServerName("server_two"))
	if err != nil {
		panic(err)
	}
	server2.AddMethod("sayGoodbye", HandlerSayGoodbye)

	// server manager
	serverManager.Register(server1)
	serverManager.Register(server2)

	serverManager.Start()
}

func HandlerSayHello(ctx *connector.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	json.Unmarshal(ctx.Request.Body, &data)
	fmt.Printf("DATA: %+v", data)

	ctx.Response.WithBody([]byte("{ \"welcome\": \"" + data.Name + "\" }"))

	return nil
}

func HandlerSayGoodbye(ctx *connector.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR POST")

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	json.Unmarshal(ctx.Request.Body, &data)
	fmt.Printf("DATA: %+v", data)

	ctx.Response.WithBody([]byte("{ \"goodbye\": \"" + data.Name + "\" }"))

	return nil
}
