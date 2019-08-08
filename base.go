package connector

import (
	"reflect"
	"strings"
)

func (r *Base) SetHeader(name string, header []string) {
	r.Headers[strings.Title(name)] = header
}

func (r *Base) GetHeader(name string) string {
	if header, ok := r.Headers[strings.Title(name)]; ok {
		return header[0]
	}
	return ""
}

func (r *Base) BindHeaders(obj interface{}) error {
	if len(r.Headers) == 0 {
		return nil
	}

	data := make(map[string]string)
	for name, values := range r.Headers {
		data[name] = values[0]
	}

	return readData(reflect.ValueOf(obj), data)
}
