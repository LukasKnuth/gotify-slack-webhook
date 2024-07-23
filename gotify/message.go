package gotify

import "github.com/gotify/plugin-api"

func ToMessage() (plugin.Message, error) {
	return plugin.Message{
		Title:   "From Slack",
		Message: "Recieved this from Slack...\n\nSome **Markdown** going _down_ right here!",
		Extras:  map[string]interface{}{"client::display": map[string]string{"contentType": "text/markdown"}},
	}, nil
}
