shell:
	docker run --rm -it \
		-v ${PWD}/src:/go/peonbot \
		-v ${PWD}/bot/config:/go/peonbot/config \
		-v ${PWD}/bot/tokens:/go/peonbot/tokens \
		-w /go/peonbot/ \
		golang:1.11 \
		/bin/bash

build:
	docker build -t peonbot .

bot: build
	docker run --rm \
		-v ${PWD}/bot:/bot \
		-e OUTFILE=peonbot_linux_amd64 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		peonbot

bot-win: build
	docker run --rm \
		-v ${PWD}/bot:/bot \
		-e OUTFILE=peonbot_windows_amd64.exe \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		peonbot

coverage: build
	docker run --rm \
		--entrypoint "/bin/sh" \
		peonbot \
		"coverage.sh"

coverage-report: build
	mkdir -p coverage
	touch coverage/coverage.html
	docker run --rm \
		-v ${PWD}/coverage:/coverage \
		--entrypoint "/bin/sh" \
		peonbot \
		"coverage-report.sh"