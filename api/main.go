package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	buffer := NewStringBuffer()
	users := make([]int, 0)

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

		id := rand.Int()
		users = append(users, id)
		if err := handleClient(ws, buffer, id); err != nil {
			c.Logger().Error(err)
		}
		for i, user := range users {
			if user == id {
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

func handleClient(ws *websocket.Conn, buffer *StringBuffer, id int) error {
	for _, msg := range buffer.Slice() {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			return err
		}
	}

	buffer.Register(id, func(msg string) error {
		return ws.WriteMessage(websocket.TextMessage, []byte(msg))
	})
	defer buffer.Unregister(id)

	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			return nil
		} else if err != nil {
			return err
		}
		buffer.Write(string(msg))
	}
}

type StringBuffer struct {
	msg         []string
	subscribers map[int]func(string) error
}

func NewStringBuffer() *StringBuffer {
	return &StringBuffer{
		msg:         make([]string, 0),
		subscribers: map[int]func(string) error{},
	}
}

func (b *StringBuffer) Register(id int, f func(string) error) {
	b.subscribers[id] = f
}

func (b *StringBuffer) Unregister(id int) {
	delete(b.subscribers, id)
}

func (b *StringBuffer) Write(msg string) {
	b.msg = append(b.msg, msg)
	for _, handler := range b.subscribers {
		if err := handler(msg); err != nil {
			fmt.Println(err)
		}
	}
}

func (b *StringBuffer) Slice() []string {
	return b.msg
}