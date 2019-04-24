#!make

CONFIG=example_conf.toml

run:
	go run run.go --config=$(CONFIG)

build:
	go build -o rslacker run.go

test:
	go test ./...
