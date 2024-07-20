package main

import (
	"github.com/gotify/plugin-api"
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "Slack Incoming Webhook Support",
		Description: "Allows Gotify to _receive_ Slack Incoming Webhook calls. For tools that integrate with Slack by accepting a Slack Incoming Webhook URL, simply put the Gotify URL. The Plugin will accept the Webhook and create a Gotify message from it.",
		ModulePath:  "github.com/LukasKnuth/gotify-slack-webhook",
	}
}

// Plugin is plugin instance
type Plugin struct{}

// Enable implements plugin.Plugin
func (c *Plugin) Enable() error {
	return nil
}

// Disable implements plugin.Plugin
func (c *Plugin) Disable() error {
	return nil
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Plugin{}
}

func main() {
	panic("this should be built as go plugin")
}
