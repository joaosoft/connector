package connector

type Status int

const (
	StatusOk            Status = 0
	StatusCreated       Status = 1
	StatusAccepted      Status = 2
	StatusNoContent     Status = 3
	StatusFound         Status = 4
	StatusBadRequest    Status = 5
	StatusUnauthorized  Status = 6
	StatusNotFound      Status = 7
	StatusConflict      Status = 8
	StatusInternalError Status = 9
	StatusBadGateway    Status = 10
	StatusUnavailable   Status = 11
	StatusTimeout       Status = 12
)

var statusText = map[Status]string{
	StatusOk:            "Ok",
	StatusCreated:       "Created",
	StatusAccepted:      "Accepted",
	StatusNoContent:     "No Content",
	StatusFound:         "Found",
	StatusBadRequest:    "Bad Request",
	StatusUnauthorized:  "Unauthorized",
	StatusNotFound:      "Not Found",
	StatusConflict:      "Conflict",
	StatusInternalError: "Internal Error",
	StatusBadGateway:    "Bad Gateway",
	StatusUnavailable:   "Unavailable",
	StatusTimeout:       "Timeout",
}

func StatusText(code Status) string {
	return statusText[code]
}
