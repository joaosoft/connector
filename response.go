package connector

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (r *Response) withStatus(status Status) *Response {
	r.Status = status
	return r
}

func (r *Response) WithBody(body []byte) *Response {
	r.Body = body
	return r
}

func (r *Response) write() error {
	var buf bytes.Buffer

	if headers, err := r.buildHeaders(); err != nil {
		return err
	} else {
		buf.Write(headers)
	}

	if body, err := r.buildBody(); err != nil {
		return err
	} else {
		buf.Write(body)

		if r.Server.logger.IsDebugEnabled() {
			r.Server.logger.Infof("[RESPONSE BODY] [%s]", string(body))
		}
	}

	r.conn.Write(buf.Bytes())

	return nil
}

func (r *Response) buildHeaders() ([]byte, error) {
	var buf bytes.Buffer

	// header
	buf.WriteString(fmt.Sprintf("%s %d\r\n", r.Method, r.Status))

	// headers
	r.Headers[HeaderServer] = []string{r.Server.name}
	r.Headers[HeaderDate] = []string{time.Now().Format(HeaderTimeFormat)}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Response) buildBody() ([]byte, error) {
	return r.Body, nil
}

func setValue(kind reflect.Kind, obj reflect.Value, newValue string) error {

	if !obj.CanAddr() {
		return nil
	}

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.Atoi(newValue)
		obj.SetInt(int64(v))
	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(newValue, 64)
		obj.SetFloat(v)
	case reflect.String:
		obj.SetString(newValue)
	case reflect.Bool:
		v, _ := strconv.ParseBool(newValue)
		obj.SetBool(v)
	}

	return nil
}

func (r *Response) read() error {
	reader := bufio.NewReader(r.Reader)

	// header
	if err := r.readHeader(reader); err != nil {
		return err
	}

	// headers
	if err := r.readHeaders(reader); err != nil {
		return err
	}

	// body
	if err := r.readBody(reader); err != nil {
		return err
	}

	return nil
}

func (r *Response) readHeader(reader *bufio.Reader) error {

	// read one line (ended with \n or \r\n)
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("invalid header send: %s", err)
	}

	if firstLine := bytes.SplitN(line, []byte(` `), 2); len(firstLine) < 2 {
		return errors.New("invalid header send")
	} else {
		r.Method = string(firstLine[0])
		status, err := strconv.Atoi(string(firstLine[1]))
		if err != nil {
			return fmt.Errorf("invalid connector response [%s]", string(line))
		}

		r.Status = Status(status)
	}

	r.Method = string(line)

	return nil
}

func (r *Response) readHeaders(reader *bufio.Reader) error {
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		if split := bytes.SplitN(line, []byte(`: `), 2); len(split) > 0 {
			r.Headers[strings.Title(string(split[0]))] = []string{string(split[1])}
		}
	}

	return nil
}

func (r *Response) readBody(reader *bufio.Reader) error {
	var buf bytes.Buffer

	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)
	}

	r.Body = buf.Bytes()

	return nil
}
