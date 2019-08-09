package main

import (
	"connector"
	"encoding/json"
	"fmt"
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
	response, err := clientManager.Invoke("service_one", "sayHello", bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("service_one response: %+v\n", string(response.Body))

	// invoke service_two
	response, err = clientManager.Invoke("service_two", "sayGoodbye", bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("service_two response: %+v\n", string(response.Body))
}
