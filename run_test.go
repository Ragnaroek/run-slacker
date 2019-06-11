package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigOrPanic(t *testing.T) {

	exampleConf := Config{}
	bytes, err := ioutil.ReadFile("./example_conf.toml")
	assert.Nil(t, err)
	exampleToml := string(bytes)
	_, err = toml.Decode(exampleToml, &exampleConf)
	assert.Nil(t, err)

	tcs := []struct {
		desc           string
		configFilePath string
		configEnvValue string
		expectedPanic  bool
		expectedConfig Config
	}{
		{
			desc:           "no config file exists and not env var set => panic",
			configFilePath: "./nonExistingConfigFile.toml",
			configEnvValue: "",
			expectedPanic:  true,
		},
		{
			desc:           "takes config from file if set",
			configFilePath: "./example_conf.toml",
			configEnvValue: "",
			expectedConfig: exampleConf,
		},
		{
			desc:           "takes config from env if set",
			configFilePath: "./nonExistingConfigFile.toml",
			configEnvValue: exampleToml,
			expectedConfig: exampleConf,
		},
		{
			desc:           "config file toml broken => panic",
			configFilePath: "",
			configEnvValue: "",
			expectedPanic:  true,
		},
		{
			desc:           "env toml broken => panic",
			configFilePath: "./broken_conf.toml",
			configEnvValue: "invalid toml",
			expectedPanic:  true,
		},
		{
			desc:           "env variable overwrites config file",
			configFilePath: "./example_conf.toml",
			configEnvValue: `prog = "overwritten_env_config"`,
			expectedConfig: Config{
				Prog: "overwritten_env_config",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			os.Setenv(CONFIG_ENV_VAR, tc.configEnvValue)
			if tc.expectedPanic {
				assert.Panics(t, func() { _ = getConfigOrPanic(tc.configFilePath) }, "expected panic")
			} else {
				config := getConfigOrPanic(tc.configFilePath)
				assert.Equal(t, tc.expectedConfig, config)
			}
		})
	}
}

func TestGetLevelOrPanic(t *testing.T) {

	tcs := []struct {
		desc          string
		config        Config
		expectedPanic bool
		expectedLevel Level
	}{
		{
			desc: "default level",
			config: Config{
				Level: nil,
			},
			expectedLevel: Always,
		},
		{
			desc: "explicit always",
			config: Config{
				Level: ptrToStr("Always"),
			},
			expectedLevel: Always,
		},
		{
			desc: "error_or_output",
			config: Config{
				Level: ptrToStr("Error_Or_Output"),
			},
			expectedLevel: Error_Or_Output,
		},
		{
			desc: "error",
			config: Config{
				Level: ptrToStr("Error"),
			},
			expectedLevel: Error,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.expectedPanic {
				assert.Panics(t, func() { _ = getLevelOrPanic(tc.config) }, "expected panic")
			} else {
				lvl := getLevelOrPanic(tc.config)
				assert.Equal(t, tc.expectedLevel, lvl)
			}
		})
	}
}

func TestGetProgName(t *testing.T) {
	tcs := []struct {
		desc         string
		conf         *Config
		expectedName string
	}{
		{
			desc: "no name in config",
			conf: &Config{
				Prog: "ls",
			},
			expectedName: "`ls`",
		},
		{
			desc: "name in config",
			conf: &Config{
				Name: ptrToStr("list files"),
				Prog: "ls",
			},
			expectedName: "list files",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			result := getProgName(tc.conf)
			assert.Equal(t, tc.expectedName, result)
		})
	}
}

func ptrToStr(str string) *string {
	return &str
}
