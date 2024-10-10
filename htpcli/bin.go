package htpcli

import (
	"bytes"
	"fmt"
	"net/http"
)

type (
	HTTPclient struct {
		Socket  *http.Client
		Method  string
		URL     string
		Body    []byte
		Headers map[string]string
	}
)

// client types init
func (cli *HTTPclient) NewGet(url string, headers map[string]string) *HTTPclient {
	cli.Method = http.MethodGet
	cli.URL = url
	cli.Body = nil
	cli.Headers = headers
	return cli
}

func (cli *HTTPclient) NewPost(url string, headers map[string]string, body []byte) *HTTPclient {
	cli.URL = url
	cli.Method = http.MethodPost
	cli.Body = body
	cli.Headers = headers
	return cli
}

// sent client request
func (cli *HTTPclient) Do() (resp *http.Response, err error) {
	reader := bytes.NewReader(cli.Body)
	req, err := http.NewRequest(cli.Method, cli.URL, reader)
	if err != nil {
		err = fmt.Errorf("request building error: %w", err)
		return
	}
	if len(cli.Headers) != 0 {
		for n, v := range cli.Headers {
			if v != "" && n != "" {
				req.Header.Add(n, v)
			}
		}
	}

	resp, err = cli.Socket.Do(req)
	return
}
