package httputil

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var DefaultHttp = HttpClient{Timeout: 10 * time.Second}

type HttpClient struct {
	Timeout time.Duration
}
type HttpResult struct {
	Err      error
	Response *http.Response
	Body     []byte
}
type HttpOption func(req *http.Request)

func HttpWithHeader(key, value string) HttpOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

func (c HttpClient) PutJson(url string, data interface{}, options ...HttpOption) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	result, err := c.OriginJson(url, "PUT", body, options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) PostJson(url string, data interface{}, options ...HttpOption) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.PostJsonBytes(url, body, options...)
}
func (c HttpClient) PostJsonBytes(url string, data []byte, options ...HttpOption) ([]byte, error) {
	result, err := c.OriginJson(url, "POST", data, options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) PostMapFrom(url string, data map[string]string, options ...HttpOption) ([]byte, error) {
	body := c.encode(data)
	result, err := c.OriginForm(url, "POST", []byte(body), options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) PostFrom(url string, data url.Values, options ...HttpOption) ([]byte, error) {
	body := data.Encode()

	result, err := c.OriginForm(url, "POST", []byte(body), options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) Get(url string, options ...HttpOption) ([]byte, error) {
	result, err := c.OriginGet(url, options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (c HttpClient) OriginJson(url string, method string, data []byte, options ...HttpOption) (*HttpResult, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("content-type", "application/json")

	for _, op := range options {
		op(req)
	}

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		return &HttpResult{Response: res, Err: errors.New("status:" + res.Status)}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &HttpResult{Response: res, Err: err}, errors.WithStack(err)
	}

	return &HttpResult{
		Response: res,
		Body:     body,
	}, nil
}
func (c HttpClient) OriginForm(url string, method string, data []byte, options ...HttpOption) (*HttpResult, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	for _, op := range options {
		op(req)
	}

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		return &HttpResult{Response: res, Err: errors.New("status:" + res.Status)}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &HttpResult{Response: res, Err: err}, errors.WithStack(err)
	}
	return &HttpResult{Response: res, Body: body}, nil
}
func (c HttpClient) OriginGet(url string, options ...HttpOption) (*HttpResult, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, op := range options {
		op(req)
	}

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		return &HttpResult{Response: res, Err: errors.New("status:" + res.Status)}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &HttpResult{Response: res, Err: err}, errors.New("status:" + res.Status)
	}
	return &HttpResult{Response: res, Body: body}, nil
}
func (c HttpClient) encode(data map[string]string, options ...HttpOption) string {
	var buf strings.Builder
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(data[k]))
	}
	return buf.String()
}

func PutJson(url string, data interface{}, options ...HttpOption) ([]byte, error) {
	return DefaultHttp.PutJson(url, data, options...)
}
func PostJson(url string, data interface{}, options ...HttpOption) ([]byte, error) {
	return DefaultHttp.PostJson(url, data, options...)
}
func PostJsonBytes(url string, data []byte, options ...HttpOption) ([]byte, error) {
	return DefaultHttp.PostJsonBytes(url, data, options...)
}
func PostMapFrom(url string, data map[string]string, options ...HttpOption) ([]byte, error) {
	return DefaultHttp.PostMapFrom(url, data, options...)
}
func PostFrom(url string, data url.Values, options ...HttpOption) ([]byte, error) {
	return DefaultHttp.PostFrom(url, data, options...)
}
func Get(url string, options ...HttpOption) ([]byte, error) {
	return DefaultHttp.Get(url, options...)
}
func OriginJson(url string, method string, data []byte, options ...HttpOption) (*HttpResult, error) {
	return DefaultHttp.OriginJson(url, method, data, options...)
}
func OriginForm(url string, method string, data []byte, options ...HttpOption) (*HttpResult, error) {
	return DefaultHttp.OriginForm(url, method, data, options...)
}
func OriginGet(url string, options ...HttpOption) (*HttpResult, error) {
	return DefaultHttp.OriginGet(url, options...)
}
