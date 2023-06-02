package main

import (
	"encoding/json"
	"fmt"
	"github.com/joaosoft/connector"
)

func main() {
	// create a new client
	c, err := connector.NewClient()
	if err != nil {
		panic(err)
	}

	request(c)
}

func request(c *connector.Client) {
	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "joao",
		Age:  30,
	}

	bytes, _ := json.Marshal(data)

	response, err := c.Invoke("sayHello", "localhost:9001", nil, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(response.Body))

	response, err = c.Invoke("Test", "localhost:9001", nil, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(response.Body))
}
