package internal

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

type HTTPHandlers struct {
	MainHandler
}

func NewHTTPHandlers(wsServer *WebSocketServer) *HTTPHandlers {
	return &HTTPHandlers{
		MainHandler: MainHandler{
			wsServer: wsServer,
		},
	}
}

type MainHandler struct {
	wsServer *WebSocketServer
}

func (mh *MainHandler) HandleIndex(c *fiber.Ctx) error {
	uid := c.Cookies("uid")

	if uid == "" {
		cookie := new(fiber.Cookie)
		cookie.Name = "uid"
		cookie.Value = uuid.New().String()
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(365 * 24 * time.Hour)

		c.Cookie(cookie)
	}

	context := fiber.Map{}
	return c.Render("index", context)
}

func (mh *MainHandler) HandleMessage(c *fiber.Ctx) error {
	uid := c.Cookies("uid")

	if uid == "" {
		return c.Redirect("/")
	}

	newMessage := Message{
		UID: uid,
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["images"]

	for _, file := range files {
		imagePath := fmt.Sprintf("./static/images/%s", uuid.New().String()+"__"+file.Filename)

		err := c.SaveFile(file, imagePath)
		if err != nil {
			return err
		}

		newMessage.Images = append(newMessage.Images, imagePath)
	}

	newMessage.Text = form.Value["text"][0]

	mh.wsServer.broadcast <- &newMessage

	return nil
}
