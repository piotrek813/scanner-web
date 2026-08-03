package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed index.html
var index string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = map[*websocket.Conn]bool{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(index))
	})

	http.HandleFunc("/ws", wsHandler)

	http.HandleFunc("/capture", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("name")
		if filename == "" {
			filename = "photo.jpg"
		}

		for c := range clients {
			c.WriteJSON(map[string]string{
				"capture": filename,
			})
		}

		w.Write([]byte("sent"))
	})

	log.Println("Open http://localhost:8080")
	log.Println("Trigger: http://localhost:8080/capture?name=test.jpg")

	log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clients[conn] = true
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
