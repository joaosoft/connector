# connector
[![Build Status](https://travis-ci.org/joaosoft/connector.svg?branch=master)](https://travis-ci.org/joaosoft/connector) | [![codecov](https://codecov.io/gh/joaosoft/connector/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/connector) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/connector)](https://goreportcard.com/report/github.com/joaosoft/connector) | [![GoDoc](https://godoc.org/github.com/joaosoft/connector?status.svg)](https://godoc.org/github.com/joaosoft/connector)

A simple and fast tcp server and client that make their communication by method, headers and body with support for server manager and client manager.

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

### Configuration
```javascript
{
  "server": {
    "address": ":9001",
    "log": {
      "level": "error"
    }
  },
  "client": {
    "log": {
      "level": "error"
    }
  },
  "server_manager": {
    "services": {
      "service_one": {
        "address": ":9001"
      },
      "service_two": {
        "address": ":9002"
      }
    },
    "log": {
      "level": "error"
    }
  },
  "client_manager": {
    "services": {
      "service_one": {
        "address": ":9001"
      },
      "service_two": {
        "address": ":9002"
      }
    },
    "log": {
      "level": "error"
    }
  }
}
```

### Usage of Server and Client

#### Server
```go
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
```

#### Client
```go
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
}

```

### Usage of Server Manager and Client Manager

#### Server Manager
```go
func main() {
	serverManager, err := connector.NewServerManager()

	// server 1
	server1, err := connector.NewServer()
	if err != nil {
		panic(err)
	}
	server1.AddMethod("sayHello", HandlerSayHello)

	// server 2
	server2, err := connector.NewServer()
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
```

#### Client Manager
```go
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
	response, err := clientManager.Invoke("service_one", "sayHello", nil, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("service_one response: %+v\n", string(response.Body))

	// invoke service_two
	response, err = clientManager.Invoke("service_two", "sayGoodbye", nil, bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("service_two response: %+v\n", string(response.Body))
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
