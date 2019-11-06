package peonbot

import (
	"bufio"
	"os"
	"strings"
)

func (bot *_bot) ListenStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()

		if len(msg) == 0 {
			continue
		}

		bot.chsin <- msg
	}
}

func (bot *_bot) HandleMessage(client WebsocketClient, message string) error {
	/* Handle actions sent from stdin */
	switch isAction(message) {
	case true:
		message = formatStdinAction(message)
	case false:
		message = formatStdinMessage(message)
	}

	return handleAction(client, bot, bot.getStdinAction(message))
}

const _STDIN_ACTION_DELIMITER = "/"

func isAction(message string) bool {
	if strings.Compare(
		string(message[0]), _STDIN_ACTION_DELIMITER) == 0 {

		return true
	}

	return false
}

/* Actions received from stdin start with a "/" instead of a "." */
func formatStdinAction(message string) string {
	return "." + message[1:len(message)]
}

/* Messages received from stdin need to be prepended with ".say " */
func formatStdinMessage(message string) string {
	return ".say " + message
}

func (bot *_bot) getStdinAction(formatted string) _event {
	return _event{
		Command:   _EVENT_MSG,
		RequestId: bot.rid,
		Payload: _payload{
			UserId:  _PEONBOT_USERID,
			Message: formatted,
		},
	}
}
