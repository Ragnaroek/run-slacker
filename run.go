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
	Name  *string
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
		panic(fmt.Sprintf("Error reading config file: %#v", err))
	}

	cmd := exec.Command(config.Prog, config.Args...)
	cmd.Dir = config.Dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			slack(&config, runFailedMessage(&config, &exitCode, string(out)))
		} else {
			slack(&config, runFailedMessage(&config, nil, string(out)))
		}
		return
	}

	slack(&config, runOkMessage(&config, string(out)))
}

func runFailedMessage(conf *Config, exitCode *int, output string) *SlackBlock {
	name := getProgName(conf)
	var exitCodeStr string
	if exitCode != nil {
		exitCodeStr = fmt.Sprintf("%d", *exitCode)
	} else {
		exitCodeStr = "unknown"
	}

	return section(markdownText(fmt.Sprintf(":red_circle: %s failed (exit code %s)\n```%s```", name, exitCodeStr, output)))
}

func runOkMessage(conf *Config, output string) *SlackBlock {
	name := getProgName(conf)
	text := ""
	if strings.TrimSpace(output) != "" {
		text = fmt.Sprintf("\n```%s```", output)
	}
	return section(markdownText(fmt.Sprintf(":large_blue_circle: %s ok%s", name, text)))
}

func getProgName(conf *Config) string {
	name := fmt.Sprintf("`%s`", conf.Prog)
	if conf.Name != nil {
		name = *conf.Name
	}
	return name
}

func slack(conf *Config, msg *SlackBlock) {
	payloadStr, err := json.Marshal(SlackMessage{Blocks: []*SlackBlock{msg}})
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

type SlackMessage struct {
	Blocks []*SlackBlock `json:"blocks"`
}

type SlackBlock struct {
	Type string `json:"type"`
	Text *Text  `json:"text"`
}

func section(text *Text) *SlackBlock {
	return &SlackBlock{
		Type: "section",
		Text: text,
	}
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func markdownText(msg string) *Text {
	return &Text{Type: "mrkdwn", Text: msg}
}
