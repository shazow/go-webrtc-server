package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/websocket"
	webrtc "github.com/keroserene/go-webrtc"
)

func main() {
	var upgrader = socketUpgrader{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},

		OnChannel: func(channel *webrtc.DataChannel) {
			channel.Send([]byte("Hello from server"))
			channel.OnMessage = func(msg []byte) {
				fmt.Println("channel received:", string(msg))
			}
		},
	}

	//http.Handle("/", http.RedirectHandler("/static", http.StatusSeeOther))
	http.HandleFunc("/connect", upgrader.Handler)

	path, err := filepath.Abs("./static")
	if err != nil {
		panic(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path))))

	bind := "127.0.0.1:3000"
	fmt.Println("listening on:", bind)
	http.ListenAndServe(bind, nil)
}
