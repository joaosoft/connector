package connector

func (c *Client) NewRequest(method string, address string, headers ...Headers) *Request {
	var h Headers
	if len(headers) > 0 && headers[0] != nil {
		h = headers[0]
	} else {
		h = make(Headers)
	}

	return &Request{
		Base: Base{
			Client:  c,
			Method:  method,
			Address: address,
			Headers: h,
		},
	}
}
