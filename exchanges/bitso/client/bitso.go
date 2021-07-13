package client

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"github.com/vanclief/ez"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	APIKey string
	APISecret  string
	http   *http.Client
}

func New(APIKey, APISecret string) *Client {
	return &Client{
		APIKey: APIKey,
		APISecret: APISecret,
		http:   &http.Client{},
	}
}

func (c *Client) generateSignature(method, URL string, data url.Values) (string, error) {
	op := "Bitso.Client.generateSignature"

	sha := sha256.New()

	if method == "POST" {
		_, err := sha.Write([]byte(data.Get("nonce") + method + strings.Split(URL, "https://api.bitso.com")[1] + data.Encode()))
		if err != nil {
			return "", ez.Wrap(op, err)
		}
	} else {
		_, err := sha.Write([]byte(data.Get("nonce") + method + strings.Split(URL, "https://api.bitso.com")[1]))
		if err != nil {
			return "", ez.Wrap(op, err)
		}
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
	op := "Bitso.Client.httpRequest"
	if data == nil {
		data = url.Values{}
	}

	nonce := time.Now().UnixNano()

	request, err := http.NewRequest(method, URL+"?"+data.Encode(), nil)
	if err != nil {
		return ez.Wrap(op, err)
	}

	signature, err := c.generateSignature(method, URL, data)
	if err != nil {
		return ez.Wrap(op, err)
	}
	authSig := "Bitso " + c.APIKey + ":" + strconv.FormatInt(nonce, 10) + ":" + signature
	request.Header.Add("Authorization", authSig)

	response, err := c.http.Do(request)
	if err != nil {
		return ez.Wrap(op, err)
	}

	if response.StatusCode != 200 {
		errorType := ez.HTTPStatusToError(response.StatusCode)
		return ez.New(op, errorType, "Error during Bitso API request", nil)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ez.Wrap(op, err)
	}

	bitsoResponse := &BitsoResponse{}
	if responseType != nil {
		bitsoResponse.Payload = responseType
	}

	err = json.Unmarshal(responseBody, bitsoResponse)
	if err != nil {
		return ez.Wrap(op, err)
	}

	defer response.Body.Close()

	return nil
}
