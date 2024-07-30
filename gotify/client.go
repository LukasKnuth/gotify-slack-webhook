package gotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gotify/plugin-api"
)

func SendMessage(msg *plugin.Message, app_token string) error {
	port := findPort()
	req, err := messageRequest(port, msg, app_token)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("Gotify API indicated message sending failed. Is the token correct?")
	} else {
		return nil
	}
}

func findPort() string {
	port, is_set := os.LookupEnv("GOTIFY_SERVER_PORT")
	if is_set {
		return port
	} else {
		return "80"
		// TODO if the port is set via config, we're out of luck... Document in Readme?
		// TODO same for listening address. Can we even do anything then?
	}
}

func messageRequest(port string, message *plugin.Message, app_token string) (*http.Request, error) {
	url := fmt.Sprintf("http://localhost:%v/message", port)
	json, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Gotify-Key", app_token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
