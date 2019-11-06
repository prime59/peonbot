#!/bin/sh

go test -covermode=count -coverprofile=count.out ./peonbot/
go tool cover -func=count.out