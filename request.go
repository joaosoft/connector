package connector

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"
)

func (r *Request) WithBody(body []byte) *Request {
	r.Body = body

	return r
}

func (r *Request) read() error {
	reader := bufio.NewReader(r.conn)

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

func (r *Request) readHeader(reader *bufio.Reader) error {

	// read one line (ended with \n or \r\n)
	r.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	line, _, err := reader.ReadLine()
	if err != nil {
		return fmt.Errorf("invalid header request: %s", err)
	}
	r.Method = string(line)

	return nil
}

func (r *Request) readHeaders(reader *bufio.Reader) error {
	for {
		r.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 5))
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}

		if split := bytes.SplitN(line, []byte(`: `), 2); len(split) > 0 {
			r.Headers[strings.Title(string(split[0]))] = []string{string(bytes.TrimSpace(split[1]))}
		}
	}

	return nil
}

func (r *Request) readBody(reader *bufio.Reader) error {
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


func (r *Request) build() ([]byte, error) {
	var buf bytes.Buffer

	headers, err := r.buildHeaders()
	if err != nil {
		return nil, err
	}
	buf.Write(headers)

	body, err := r.buildBody()
	if err != nil {
		return nil, err
	}
	buf.Write(body)

	return buf.Bytes(), nil
}

func (r *Request) buildHeaders() ([]byte, error) {
	var buf bytes.Buffer

	// header
	buf.WriteString(fmt.Sprintf("%s\r\n", r.Method))

	// headers
	r.Headers[HeaderHost] = []string{r.Address}
	if _, ok := r.Headers[HeaderUserAgent]; !ok {
		r.Headers[HeaderUserAgent] = []string{"client"}
	}
	r.Headers[HeaderDate] = []string{time.Now().Format(TimeFormat)}

	for key, value := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value[0]))
	}

	buf.WriteString("\r\n")

	return buf.Bytes(), nil
}

func (r *Request) buildBody() ([]byte, error) {
	return r.Body, nil
}

