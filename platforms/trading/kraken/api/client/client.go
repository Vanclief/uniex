package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vanclief/ez"
)

// Client represents a new Kraken API Client
type Client struct {
	APIKey    string
	APISecret string
	http      *http.Client
}

// New initializes a new Kraken API Client
func New(APIKey, APISecret string) *Client {

	return &Client{
		APIKey:    APIKey,
		APISecret: APISecret,
		http:      &http.Client{},
	}
}

func (c *Client) generateSignature(URL string, data url.Values) (string, error) {
	const op = "Client.generateSignature"

	sha := sha256.New()

	_, err := sha.Write([]byte(data.Get("nonce") + data.Encode()))
	if err != nil {
		return "", ez.Wrap(op, err)
	}

	hashData := sha.Sum(nil)
	s, err := base64.StdEncoding.DecodeString(c.APISecret)
	if err != nil {
		return "", ez.Wrap(op, err)
	}

	hmacObj := hmac.New(sha512.New, s)
	_, err = hmacObj.Write(append([]byte(URL[22:]), hashData...))
	if err != nil {
		return "", ez.Wrap(op, err)
	}

	hmacData := hmacObj.Sum(nil)

	return base64.StdEncoding.EncodeToString(hmacData), nil
}

func (c *Client) httpRequest(method, URL string, data url.Values, responseType interface{}) error {
	const op = "Client.httpRequest"

	// Step 0: Generate the nonce
	if data == nil {
		data = url.Values{}
	}
	data.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	// Step 1: Create the request
	request, err := http.NewRequest(method, URL, strings.NewReader(data.Encode()))
	if err != nil {
		return ez.Wrap(op, err)
	}

	// Generate the signature
	signature, err := c.generateSignature(URL, data)
	if err != nil {
		return ez.Wrap(op, err)
	}

	request.Header.Add("API-Key", c.APIKey)
	request.Header.Add("API-Sign", signature)

	// Step 2: Make the request
	response, err := c.http.Do(request)
	if err != nil {
		return ez.Wrap(op, err)
	}

	// Step 3: Parse the response
	if response.StatusCode != 200 {
		errorType := ez.HTTPStatusToError(response.StatusCode)
		return ez.New(op, errorType, "Error during Kraken API request", nil)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ez.Wrap(op, err)
	}

	krakenResponse := &KrakenResponse{}
	if responseType != nil {
		krakenResponse.Result = responseType
	}

	err = json.Unmarshal(responseBody, krakenResponse)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if len(krakenResponse.Error) > 0 {
		errMsg := fmt.Sprintf("Kraken request returned an error: %s", krakenResponse.Error)
		return ez.New(op, ez.EINTERNAL, errMsg, nil)
	}

	defer response.Body.Close()

	return nil
}
