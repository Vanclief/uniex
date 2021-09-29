package ws

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func Web() {

	u := url.URL{Scheme: "https", Host: "mt-client-api-v1.agiliumtrade.agiliumtrade.ai", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	fmt.Println(c)

}
