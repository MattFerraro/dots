package main

import (
    "github.com/gorilla/websocket"
    "flag"
    "log"
    "net/http"
    "fmt"
)

var (
    addr = flag.String("addr", ":8080", "http service address")
)

var h = hub{
    broadcast:   make(chan []byte),
    register:    make(chan *connection),
    unregister:  make(chan *connection),
    connections: make(map[*connection]bool),
}


type connection struct {
    // The websocket connection.
    ws *websocket.Conn

    // Buffered channel of outbound messages.
    send chan []byte
}

type hub struct {
    // Registered connections.
    connections map[*connection]bool
    // Inbound messages from the connections.
    broadcast chan []byte
    // Register requests from the connections.
    register chan *connection
    // Unregister requests from connections.
    unregister chan *connection
}

func (c *connection) reader() {
    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            break
        }
        fmt.Println("msg: %s", message)
        h.broadcast <- message
    }
    c.ws.Close()
}

func (c *connection) writer() {
    for message := range c.send {
        err := c.ws.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
    c.ws.Close()
}

func main() {
    flag.Parse()
    http.HandleFunc("/ws", wsHandler)
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(w, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        return
    }
    c := &connection{send: make(chan []byte, 256), ws: ws}
    h.register <- c
    defer func() { h.unregister <- c }()
    go c.writer()
    c.reader()
}

func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
            h.connections[c] = true
        case c := <-h.unregister:
            delete(h.connections, c)
            close(c.send)
        case m := <-h.broadcast:
            for c := range h.connections {
                select {
                case c.send <- m:
                default:
                    delete(h.connections, c)
                    close(c.send)
                    go c.ws.Close()
                }
            }
        }
    }
}
