package event

import "github.com/10hourlabs/whisper"

var Events = []whisper.EventHandler{
	&HelloWorldEvent{},
}
