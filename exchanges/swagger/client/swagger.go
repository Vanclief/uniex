package client

import (
  "encoding/json"
  "github.com/vanclief/ez"
  "io/ioutil"
  "net/http"
  "net/url"
)

type Client struct {
  AccountID string
  Token string
  http *http.Client
}

func New(accountID, token string) *Client {
  return &Client{
    AccountID: accountID,
    Token:     token,
    http:      &http.Client{},
  }
}

func (c *Client) httpRequest(method, URL string, data url.Values, responseType interface{}) error {
  op := "Swagger.Client.httpRequest"
  if data == nil {
    data = url.Values{}
  }

  request, err := http.NewRequest(method, URL + "?" + data.Encode(), nil)
  if err != nil {
    return ez.Wrap(op, err)
  }

  request.Header.Add("auth-token", c.Token)

  response, err := c.http.Do(request)
  if err != nil {
    return ez.Wrap(op, err)
  }

  if response.StatusCode != 200 {
    errorType := ez.HTTPStatusToError(response.StatusCode)
    return ez.New(op, errorType, "Error during Swagger UI API request", nil)
  }

  responseBody, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return ez.Wrap(op, err)
  }

  swaggerResponse := &SwaggerResponse{}
  if responseType != nil {
    swaggerResponse.Payload = responseType
  }

  err = json.Unmarshal(responseBody, swaggerResponse.Payload)
  if err != nil {
    return ez.Wrap(op, err)
  }

  defer response.Body.Close()

  return nil
}
