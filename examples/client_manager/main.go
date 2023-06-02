package main

import (
	"encoding/json"
	"fmt"
	"github.com/joaosoft/connector"
)

func main() {
	// client manager
	clientManager, err := connector.NewClientManager()
	if err != nil {
		panic(err)
	}

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "joao",
		Age:  30,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// invoke service_one
	response, err := clientManager.Invoke("server_one", "sayHello", nil, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server_one response: %+v\n", string(response.Body))

	// invoke service_two
	response, err = clientManager.Invoke("server_two", "sayGoodbye", nil, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("server_two response: %+v\n", string(response.Body))
}
