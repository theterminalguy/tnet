package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/theterminalguy/whisper"
)

type HelloWorldEventPayload struct {
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Message   string `json:"message"`
}

type HelloWorldEvent struct{}

func (*HelloWorldEvent) GetEventName() whisper.Event {
	return "hello-world"
}

func (*HelloWorldEvent) GetSubscriptionID() string {
	return "tentn-example-topic-dev-sub"
}

func (*HelloWorldEvent) GetContext() context.Context {
	return context.Background()
}

func (*HelloWorldEvent) ValidatePayload(payload []byte) error {
	var p HelloWorldEventPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return whisper.ErrInvalidPayload
	}
	return nil
}

func (*HelloWorldEvent) Handle(ctx context.Context, body []byte) error {
	var data HelloWorldEventPayload
	json.Unmarshal(body, &data) // gauranteed to not error
	fmt.Printf("%s\n", data.Message)
	return nil
}
