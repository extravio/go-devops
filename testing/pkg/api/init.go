package api

import (
	"io"
	"net/http"
)

type ClientIface interface {
	Get(url string) (resp *http.Response, err error)
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Options struct {
	Password string
	LoginURL string
}

type APIIface interface {
	DoGetRequest(requestURL string) (Response, error)
}

// Uppercase => can be initialized / used outside of the package
// Lowercase => can only be instantiated ussing New
type API struct {
	Options Options
	Client  ClientIface
}

func New(options Options) APIIface {
	return API{
		Options: options,
		Client: &http.Client{
			Transport: MyJWTTransport{
				transport:  http.DefaultTransport,
				password:   options.Password,
				loginURL:   options.LoginURL,
				HTTPClient: &http.Client{},
			},
		},
	}
}
