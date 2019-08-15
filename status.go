package connector

type Status int

const (
	StatusOk    Status = 0
	StatusError Status = 1
)

var statusText = map[Status]string{
	StatusOk:    "Ok",
	StatusError: "Error",
}

func StatusText(code Status) string {
	return statusText[code]
}
