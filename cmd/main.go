package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"htmx-with-go/internal"
	"log"
)

func main() {
	viewsEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})

	app.Static("/static/", "./static")

	server := internal.NewWebSocketServer()
	httpHandlers := internal.NewHTTPHandlers(server)

	app.Get("/", httpHandlers.HandleIndex)
	app.Post("/messages", httpHandlers.HandleMessage)

	app.Get("/ws", websocket.New(func(ctx *websocket.Conn) {
		server.HandleWebSocket(ctx)
	}))

	go server.HandleMessages()

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
