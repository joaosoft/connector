# connector
[![Build Status](https://travis-ci.org/joaosoft/connector.svg?branch=master)](https://travis-ci.org/joaosoft/connector) | [![codecov](https://codecov.io/gh/joaosoft/connector/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/connector) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/connector)](https://goreportcard.com/report/github.com/joaosoft/connector) | [![GoDoc](https://godoc.org/github.com/joaosoft/connector?status.svg)](https://godoc.org/github.com/joaosoft/connector)

A simple and fast tcp server and client that make their communication by method, headers and body.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for 
* Methods
* Headers
* Middlewares


>### Go
```
go get github.com/joaosoft/connector
```

## Usage 
This examples are available in the project at [connector/examples](https://github.com/joaosoft/connector/tree/master/examples)

### Server
```
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

	return ctx.Response.WithBody([]byte("{ \"welcome\": \""+data.Name+"\" }"),
	)
}
```

### Client
```
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
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
