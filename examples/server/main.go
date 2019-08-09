package main

import (
	"connector"
	"encoding/json"
	"fmt"
)

func main() {
	// create a new server
	w, err := connector.NewServer()
	if err != nil {
		panic(err)
	}

	w.AddMiddlewares(MyMiddlewareOne(), MyMiddlewareTwo())

	w.AddMethods(
		connector.NewMethod("sayHello", HandlerSayHello),
	)

	// start the server
	if err := w.Start(); err != nil {
		panic(err)
	}
}

func MyMiddlewareOne() connector.MiddlewareFunc {
	return func(next connector.HandlerFunc) connector.HandlerFunc {
		return func(ctx *connector.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE ONE")
			return next(ctx)
		}
	}
}

func MyMiddlewareTwo() connector.MiddlewareFunc {
	return func(next connector.HandlerFunc) connector.HandlerFunc {
		return func(ctx *connector.Context) error {
			fmt.Println("HELLO I'M THE MIDDLEWARE TWO")
			return next(ctx)
		}
	}
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
