# Gotify Slack Webhook

Allows Gotify to receive [Slack Incoming Webhook](https://api.slack.com/messaging/webhooks) requests. If a service integrates with Slack by allowing you to specify the full Slack Webhook URL, simply give it the Gotify endpoint instead. The plugin will accept the request and create a message for the specified application.

## ToDo

- [] Setup Github action workflow to use official `gotify/build` containers to build multi-arch
- [] Add these outputs to a release, allowing the plugin to be downloaded

## References

- [Slack Incoming Webhook documentation](https://api.slack.com/messaging/webhooks)
- [Slack BlockKit](https://api.slack.com/block-kit)
- [Gotify Plugin - Getting started](https://gotify.net/docs/plugin)
- [Gotify Plugin - Reference](https://gotify.net/docs/plugin-write)
- [Gotify Plugin API](https://pkg.go.dev/github.com/gotify/plugin-api)
- [Gotify Plugin Template](https://github.com/gotify/plugin-template)
- [Gin Web Framework](https://gin-gonic.com/docs/examples/param-in-path/) - For the Webhook integration
