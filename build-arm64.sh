#! /bin/bash

env GOOS=linux GOARCH=arm64 go build -o radarr-sync-go .