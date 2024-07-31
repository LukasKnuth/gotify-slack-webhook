package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
	"github.com/lukasknuth/gotify-slack-webhook/gotify"
	"github.com/lukasknuth/gotify-slack-webhook/webhook"
)

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "Slack Incoming Webhook Support",
		Description: "Allows Gotify to _receive_ Slack Incoming Webhook calls. For tools that integrate with Slack by accepting a Slack Incoming Webhook URL, simply put the Gotify URL. The Plugin will accept the Webhook and create a Gotify message from it.",
		ModulePath:  "github.com/LukasKnuth/gotify-slack-webhook",
		Website:     "github.com/LukasKnuth/gotify-slack-webhook",
		License:     "MIT",
		Author:      "Lukas Knuth",
	}
}

const (
	webhook_path = "/webhook/slack/:app_token"
)

// Plugin is plugin instance
type Plugin struct {
	msgHandler plugin.MessageHandler
	basePath   string
}

// Called by the SDK later on, allows us to send messages to the user.
func (c *Plugin) SetMessageHandler(msgHandler plugin.MessageHandler) {
	// TODO I _think_ this doesn't allow us to send messages "on behalv" of existing apps.
	// Could call Gotify REST API instead...
	c.msgHandler = msgHandler
}

func (c *Plugin) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	// TODO Does it make sense to use ENV to override the base path? Allows cluster-internal trafik only
	c.basePath = basePath
	mux.POST(webhook_path, func(endpoint *gin.Context) {
		body, err := io.ReadAll(endpoint.Request.Body)
		if err != nil {
			endpoint.String(http.StatusBadRequest, "Could not read body from request")
			return
		}
		payload := &webhook.WebhookBody{}
		err = payload.Parse(body)
		if err != nil {
			endpoint.String(http.StatusBadRequest, "Could not parse JSON body")
			return
		}
		rendered, err := payload.Render()
		if err != nil {
			endpoint.String(http.StatusBadRequest, "Could not render result")
		}
		msg := gotify.ToMessage(rendered)
		err = gotify.SendMessage(&msg, endpoint.Param("app_token"))
		if err != nil {
			_ = c.msgHandler.SendMessage(plugin.Message{
				Title: "Could not deliver Slack Webhook content to Gotify",
				// TODO perhaps more info _why_ this wasn't possible?
				Message: err.Error(),
			})
		}
		endpoint.String(http.StatusOK, "OK")
	})
}

func (c *Plugin) GetDisplay(location *url.URL) string {
	fullUrl := location.JoinPath(c.basePath, webhook_path)
	return fmt.Sprintf("Webhook URL: `%s`", fullUrl)
}

func (c *Plugin) Enable() error {
	return nil
}

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
