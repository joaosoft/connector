package impl

import (
	"encoding/json"
	"fmt"
	"github.com/joaosoft/connector"
)

type DummyTest struct{}

func (t *DummyTest) Test(ctx *connector.Context) error {
	fmt.Println("HELLO I'M THE HELLO HANDER FOR Teste")

	data := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	json.Unmarshal(ctx.Request.Body, &data)
	fmt.Printf("DATA: %+v", data)

	ctx.Response.WithBody([]byte("{ \"implemented\": \"" + data.Name + "\" }"))

	return nil
}
