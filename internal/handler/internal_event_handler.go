package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/10hourlabs/whisper"
	"github.com/labstack/echo/v4"
)

type InternalEventHandler struct {
}

func NewInternalEventHandler() *InternalEventHandler {
	return &InternalEventHandler{}
}

func (*InternalEventHandler) ResourceName() string {
	return "internal_event"
}

func (*InternalEventHandler) Search(c echo.Context) error {
	return nil
}

func (h *InternalEventHandler) ReadAll(c echo.Context) error {
	return nil
}

func (h *InternalEventHandler) ReadByID(c echo.Context) error {
	return nil
}

func (h *InternalEventHandler) CreateOne(c echo.Context) error {
	client, err := whisper.GetClient(os.Getenv("PUBSUB_EVENT_CONNECTION_NAME"))
	if err != nil {
		return err
	}
	body := c.Request().Body
	defer body.Close()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return c.JSON(400, err)
	}
	// TODO: validate payload before publishing
	err = client.Publish(context.Background(), "tentn-example-topic-dev", b)
	if err != nil {
		fmt.Println("hello------>>", err.Error())
		return c.JSON(400, err)
	}
	fmt.Println("hello------>>", "No error")
	return nil
}

func (h *InternalEventHandler) UpdateByID(c echo.Context) error {
	return nil
}

func (h *InternalEventHandler) DeleteOne(c echo.Context) error {
	return nil
}
