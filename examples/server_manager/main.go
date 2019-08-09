package main

import (
	"connector"
	"encoding/json"
	"fmt"
)

func main() {
	serverManager, err := connector.NewServerManager()

	// server 1
	server1, err := connector.NewServer(connector.WithServerConfiguration(&connector.ServerConfig{Address: ":9001"}))
	if err != nil {
		panic(err)
	}
	server1.AddMethod("sayHello", HandlerSayHello)

	// server 2
	server2, err := connector.NewServer(connector.WithServerConfiguration(&connector.ServerConfig{Address: ":9002"}))
	if err != nil {
		panic(err)
	}
	server2.AddMethod("sayGoodbye", HandlerSayGoodbye)

	// server manager
	serverManager.Register("service_one", server1)
	serverManager.Register("service_two", server2)

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

	ctx.Response.WithBody([]byte("{ \"welcome\": \"" + data.Name + "\" }")).WithStatus(connector.StatusAccepted)

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

	ctx.Response.WithBody([]byte("{ \"goodbye\": \"" + data.Name + "\" }")).WithStatus(connector.StatusAccepted)

	return nil
}
