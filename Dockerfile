FROM golang:1.11.10

WORKDIR /go
COPY src/ peonbot/
COPY bot/config/ peonbot/config/
COPY bot/tokens/ peonbot/tokens/

WORKDIR /go/peonbot
COPY scripts/build.sh .
COPY scripts/coverage.sh .
COPY scripts/coverage-report.sh .

ENTRYPOINT ["./build.sh"]