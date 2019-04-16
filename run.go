package main

import (
	"os"
  "os/exec"
  "fmt"
  "net/http"
  "strings"

	"github.com/BurntSushi/toml"
)

type Slack struct {
  Hook string `toml:"hook"`
}

type Config struct {
  Prog string
  Args []string
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

   out, err := exec.Command(config.Prog, config.Args...).CombinedOutput()
   if err != nil {
     //TODO slack panic
     panic(fmt.Sprintf("Error running program: %#v", err))
   }

   slack(&config, string(out))
}

func slack(conf *Config, msg string) {
  payload := fmt.Sprintf("{\"text\": \"%s\"}", msg)
	resp, err := http.Post(conf.Slack.Hook, "application/json", strings.NewReader(payload))
  if err != nil {
    panic(fmt.Sprintf("Unexpected error from slack %#v", resp))
  }
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("Got unexpected response from slack %#v", resp))
	}
}
