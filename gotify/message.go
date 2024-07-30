package gotify

import (
	"github.com/gotify/plugin-api"
)

func ToMessage(markdownMessage string) plugin.Message {
	return plugin.Message{
		Message: markdownMessage,
		Extras:  map[string]interface{}{"client::display": map[string]string{"contentType": "text/markdown"}},
	}
}
