package event

import "github.com/theterminalguy/whisper"

var Events = []whisper.EventHandler{
	&HelloWorldEvent{},
}
