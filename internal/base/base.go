package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/govau/notify-client-go/notifyapi"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const NotifyBaseURL = "https://rest-api.notify.gov.au"

type Response struct {
	response *http.Response
	body     *bytes.Buffer
	Error    error
}

func (resp Response) JSON(v interface{}, at ...string) Response {
	if resp.Error != nil {
		return resp
	}

	data := resp.body.Bytes()

	for _, field := range at {
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(data, &nested); err != nil {
			return BadResponse(err)
		}

		data = nested[field]
	}

	if err := json.Unmarshal(data, v); err != nil {
		return BadResponse(err)
	}

	return resp
}

func (resp Response) JSONData(v interface{}) Response {
	return resp.JSON(v, "data")
}

func BadResponse(err error) Response {
	return Response{Error: err}
}

type Client struct {
	http.Client

	BaseURL     *url.URL
	ServiceID   string
	APIKey      string
	RouteSecret string
}

func createJWT(clientID, secret string) (string, error) {
	key := jose.SigningKey{Algorithm: jose.HS256, Key: []byte(secret)}
	sig, err := jose.NewSigner(key, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	cl := jwt.Claims{
		Issuer:   clientID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	return jwt.Signed(sig).Claims(cl).CompactSerialize()
}

func (c Client) Do(req *http.Request) (*http.Response, error) {
	token, err := createJWT(c.ServiceID, c.APIKey)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("X-Custom-Forwarder", "")
	req.Header.Add("User-agent", "NOTIFY-API-GO-CLIENT/0.0.1")

	req.URL.Host = ""
	req.URL.Scheme = ""
	req.URL = c.BaseURL.ResolveReference(req.URL)

	return c.Client.Do(req)
}

func (c Client) makeRequest(request *http.Request, options ...requestOption) Response {
	for _, option := range options {
		err := option.updateRequest(request)
		if err != nil {
			return BadResponse(err)
		}
	}

	response, err := c.Do(request)
	if err != nil {
		return BadResponse(err)
	}

	if response.StatusCode >= 400 {
		var body []byte

		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return BadResponse(err)
		}

		apiErr := &notifyapi.Error{}
		if err := json.Unmarshal(body, apiErr); err != nil {
			return BadResponse(err)
		}
		apiErr.Header = response.Header

		return BadResponse(apiErr)
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, response.Body); err != nil {
		return BadResponse(err)
	}

	if err := response.Body.Close(); err != nil {
		return BadResponse(err)
	}

	return Response{
		response: response,
		body:     &buf,
	}
}

func (c Client) Get(path string, options ...requestOption) Response {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return BadResponse(err)
	}

	return c.makeRequest(req, options...)
}

func (c Client) Post(path string, body io.Reader, options ...requestOption) Response {
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return BadResponse(err)
	}

	return c.makeRequest(req, options...)
}
