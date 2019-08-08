package connector

func (ctx *Context) Redirect(host string) error {
	if ctx.Request.Client == nil {
		client, err := NewClient(WithClientLogger(ctx.Request.Server.logger))
		if err != nil {
			return err
		}

		ctx.Request.Client = client
	}

	ctx.Request.Address = host

	response, err := ctx.Request.Send()
	if err != nil {
		return err
	}

	return ctx.Response.WithBody(response.Body)
}
