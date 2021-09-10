package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/vanclief/ez"
)

type Client struct {
	AccountID string
	Token     string
	http      *http.Client
}

func New(accountID, token string) *Client {
	return &Client{
		AccountID: accountID,
		Token:     token,
		http:      &http.Client{},
	}
}

func (c *Client) httpRequest(method, URL string, data url.Values, responseType interface{}) error {
	op := "MetaAPI.Client.httpRequest"
	if data == nil {
		data = url.Values{}
	}

	request, err := http.NewRequest(method, URL+"?"+data.Encode(), nil)
	if err != nil {
		return ez.Wrap(op, err)
	}

	request.Header.Add("auth-token", c.Token)

	response, err := c.http.Do(request)
	if err != nil {
		return ez.Wrap(op, err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.StatusCode != 200 {
		errorType := ez.HTTPStatusToError(response.StatusCode)

		apiError := &MetaAPIError{}
		err = json.Unmarshal(responseBody, apiError)
		if err != nil {
			return ez.Wrap(op, err)
		}

		return ez.New(op, errorType, apiError.Message, nil)
	}

	apiResponse := &MetaAPIResponse{}
	if responseType != nil {
		apiResponse.Payload = responseType
	}

	err = json.Unmarshal(responseBody, apiResponse.Payload)
	if err != nil {
		return ez.New(op, ez.EINVALID, err.Error(), err)
	}

	defer response.Body.Close()

	return nil
}
