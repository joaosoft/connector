package connector

func (c *Client) NewRequest(method string, address string) *Request {
	return &Request{
		Base: Base{
			Client:  c,
			Method:  method,
			Address: address,
			Headers: make(Headers),
		},
	}
}
