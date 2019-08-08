package connector

func (c *Client) NewRequest(method string, address string) (*Request, error) {
	return &Request{
		Base: Base{
			Client:  c,
			Method:  method,
			Address: address,
			Headers: make(Headers),
		},
	}, nil
}
