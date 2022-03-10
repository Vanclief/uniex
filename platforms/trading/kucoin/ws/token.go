package ws

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vanclief/ez"
)

const (
	baseEndpoint = "https://api.kucoin.com/api/v1/bullet-public"
)

type instanceServer struct {
	Endpoint     string `json:"endpoint"`
	Encrypt      bool   `json:"encrypt"`
	Protocol     string `json:"protocol"`
	PingInterval int    `json:"pingInterval"`
	PingTimeout  int    `json:"pingTimeout"`
}

type tokenData struct {
	Token           string           `json:"token"`
	InstanceServers []instanceServer `json:"instanceServers"`
}

type token struct {
	Code string
	Data tokenData
}

func GetToken() (foundToken token, err error) {
	const op = "kucoin.GetToken"

	resp, err := http.Post(baseEndpoint, "application/json", nil)
	if err != nil {
		return token{}, ez.New(op, ez.EINTERNAL, "error obtaining token", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token{}, ez.New(op, ez.EINTERNAL, "error reading the response", err)
	}

	err = json.Unmarshal(body, &foundToken)
	if err != nil {
		return token{}, ez.New(op, ez.EINTERNAL, "error parsing the token", err)
	}
	return
}
