#!/bin/sh

go test -covermode=count -coverprofile=count.out ./peonbot/
go tool cover -html=count.out -o /coverage/coverage.html