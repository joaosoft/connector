package connector

type Methods map[string]*Method

type Method struct {
	Name        string
	Method      string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

func NewMethod(method string, handler HandlerFunc, middleware ...MiddlewareFunc) *Method {
	return &Method{
		Method:      method,
		Handler:     handler,
		Middlewares: middleware,
		Name:        GetFunctionName(handler),
	}
}
