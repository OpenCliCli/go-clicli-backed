#!/usr/bin/env bash
git clone
cp ./.env /
go build -o /clicli ./main.go
sudo service goweb restart