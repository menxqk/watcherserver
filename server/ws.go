package server

import (
	"fmt"
	"log"
	"net/http"
	"watcherserver/watcher"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ws(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	watcher := watcher.New(htmlDir)
	evtCh, err := watcher.Start()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer watcher.Stop()

	// wait for event
	<-evtCh
	// send reload message
	err = ws.WriteMessage(1, []byte("RELOAD"))
	if err != nil {
		fmt.Println(err)
	}
	ws.Close()
}
