package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

type Slack struct {
	Hook string `toml:"hook"`
}

type Config struct {
	Dir   string
	Prog  string
	Args  []string
	Slack Slack
}

func main() {

	if len(os.Args) != 2 {
		panic(fmt.Sprintf("No config file given"))
	}

	confFile := os.Args[1]
	config := Config{}
	_, err := toml.DecodeFile(confFile, &config)
	if err != nil {
		//TODO slack panic
		panic(fmt.Sprintf("Error reading config file: %#v", err))
	}

	cmd := exec.Command(config.Prog, config.Args...)
	cmd.Dir = config.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		//TODO slack panic
		errText, err := json.Marshal(err)
		if err != nil {
			panic(fmt.Sprintf("error getting error"))
		}
		panic(fmt.Sprintf("Error running program, err text: %s", errText))
	}

	slack(&config, string(out))
}

type slackPayload struct {
	Text string `json:"text"`
}

func slack(conf *Config, msg string) {
	payload := slackPayload{Text: msg}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("Could not create slack message: %#v", err))
	}
	resp, err := http.Post(conf.Slack.Hook, "application/json", strings.NewReader(string(payloadStr)))
	if err != nil {
		panic(fmt.Sprintf("Unexpected error from slack: %#v", resp))
	}
	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		body := string(bodyBytes)
		if err != nil {
			body = fmt.Sprintf("err reading body from response: %#v", resp)
		}
		panic(fmt.Sprintf("Got unexpected response from slack: %s", body))
	}
}
