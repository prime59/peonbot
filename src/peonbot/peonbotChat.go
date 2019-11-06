package peonbot

import (
	"bufio"
	"os"
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

func (bot *_bot) SendMessage(message string) {
	handleActionSay(bot.Conn, bot, message)
}
