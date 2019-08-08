package main

import (
	"connector"
	"encoding/json"
	"fmt"
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
	request, err := c.NewRequest("sayHello", "localhost:9001")
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

	bytes, _ := json.Marshal(data)

	response, err := request.WithBody(bytes).Send()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", string(response.Body))
}
