package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
