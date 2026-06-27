package gotify

import (
	"github.com/gotify/plugin-api"
)

func ToMessage(markdownMessage string, titleOverride string) plugin.Message {
	return plugin.Message{
		// If this is empty-string (zero value), the application name will show up instead.
		Title:   titleOverride,
		Message: markdownMessage,
		Extras:  map[string]interface{}{"client::display": map[string]string{"contentType": "text/markdown"}},
	}
}
