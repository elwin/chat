package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	buffer := NewStringBuffer(51)
	users := make([]string, 0)

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})

	e.GET("/ws", func(c echo.Context) error {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		username := generateName()
		users = append(users, username)
		if err := handleClient(ws, buffer, username); err != nil {
			c.Logger().Error(err)
		}
		for i, user := range users {
			if user == username {
				users = append(users[:i], users[i+1:]...)
			}
		}

		return nil
	})

	e.GET("/users", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, users, "  ")
	})

	log.Fatal(e.Start(":8888"))
}

func handleClient(ws *websocket.Conn, buffer *StringBuffer, username string) error {
	for _, msg := range buffer.Slice() {
		out, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		if err := ws.WriteMessage(websocket.TextMessage, out); err != nil {
			return err
		}
	}

	buffer.Register(username, func(msg Message) error {
		out, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		return ws.WriteMessage(websocket.TextMessage, out)
	})
	defer buffer.Unregister(username)

	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			return nil
		} else if err != nil {
			return err
		}

		buffer.Write(Message{
			Body:      string(msg),
			Username:  username,
			Timestamp: time.Now(),
		})
	}
}

type Message struct {
	Body      string    `json:"body"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

type StringBuffer struct {
	start, stop int
	msg         []Message
	subscribers map[string]func(Message) error
}

func NewStringBuffer(n int) *StringBuffer {
	return &StringBuffer{
		start:       0,
		stop:        0,
		msg:         make([]Message, n),
		subscribers: map[string]func(Message) error{},
	}
}

func (b *StringBuffer) Register(username string, f func(Message) error) {
	b.subscribers[username] = f
}

func (b *StringBuffer) Unregister(username string) {
	delete(b.subscribers, username)
}

func (b *StringBuffer) Write(msg Message) {
	b.msg[b.stop] = msg
	b.stop = (b.stop + 1) % len(b.msg)
	if b.start == b.stop {
		b.start = (b.start + 1) % len(b.msg)
	}

	for _, handler := range b.subscribers {
		if err := handler(msg); err != nil {
			fmt.Println(err)
		}
	}
}

func (b *StringBuffer) Slice() []Message {
	out := make([]Message, 0)
	for i := b.start; i != b.stop; i = (i + 1) % len(b.msg) {
		out = append(out, b.msg[i])
	}

	return out
}
