#!/usr/bin/env bash

rm -rf completions
mkdir completions
go run . completion bash > completions/allincart-cli.bash
go run . completion zsh > completions/allincart-cli.zsh
go run . completion fish > completions/allincart-cli.fish