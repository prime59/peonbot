package peonbot

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const _EVENT_MSG = "Botapichat.MessageEventRequest"
const _EVENT_USERUPDATE = "Botapichat.UserUpdateEventRequest"
const _EVENT_USEREXIT = "Botapichat.UserLeaveEventRequest"
const _MSG_CHAN = "CHANNEL"
const _MSG_WHISPER = "WHISPER"

type _event struct {
	Command   string `json:"command"`
	RequestId int    `json:"request_id"`
	Payload   _payload
}

type _payload struct {
	ToonName string `json:"toon_name"`
	UserId   int    `json:"user_id"`
	Type     string `json:"type"`
	Message  string `json:"message"`
}

func (bot *_bot) HandleEvent(raw []byte) error {
	var event _event

	if err := json.Unmarshal(raw, &event); err != nil {
		return err
	}
	bot.setRid(event.RequestId)

	bot.Vprintf("%+v\n", event)

	switch event.Command {
	case _EVENT_MSG:
		/* Handle action if issued from a priveleged user */
		if err := handleAction(bot.Conn, bot, event); err != nil {
			bot.Vprintf("Could not process action: %v\n", err)
		}

		bot.handleUserMessage(event)
		break
	case _EVENT_USERUPDATE:
		bot.handleUserUpdate(event)
		break
	case _EVENT_USEREXIT:
		bot.handleUserExit(event)
		break
	default:
		return fmt.Errorf("Received unknown event from server: %+v\n", event)
	}

	return nil
}

func (bot *_bot) handleUserMessage(event _event) {
	switch strings.ToUpper(event.Payload.Type) {
	case _MSG_CHAN:
		log.Printf("[%s] %s\n", bot.userTable[event.Payload.UserId],
			event.Payload.Message)
	case _MSG_WHISPER:
		log.Printf(">>> [FROM: %s] %s\n", bot.userTable[event.Payload.UserId],
			event.Payload.Message)
	}
}

/*
	I don't know if events can be sent out of order. If they can be, then
	there is no good way to guarantee an exit event will be sent before a
	update event with a coinciding `user_id`.
*/
func (bot *_bot) handleUserUpdate(event _event) {
	if strings.Compare(bot.userTable[event.Payload.UserId], "") != 0 {
		goto ban_list
	}

	bot.userTable[event.Payload.UserId] = event.Payload.ToonName

	log.Printf("> %s has joined the channel.\n",
		bot.userTable[event.Payload.UserId])

	/*
		XXX: Disabling this for now. Bot whispers everyone upon joining the
		channel because it can't differentiate between initial user update
		events, and subsequent ones that correspond to actual events.
	*/
	// if len(bot.greetings) > 0 {
	// 	_ = bot._handleActionWhisper(event.Payload.UserId, bot.greetings)
	// }

ban_list:
	if _, ok := bot.blist[strings.ToUpper(
		bot.userTable[event.Payload.UserId])]; ok {

		_ = _handleActionBan(bot.Conn, bot, event.Payload.UserId)
	}
}

func (bot *_bot) handleUserExit(event _event) {
	/*
		No need to check if user exists until it is shown that spurious or
		duplicate user exit events are sent from the server
	*/
	log.Printf("< %s has left the channel.\n",
		bot.userTable[event.Payload.UserId])

	delete(bot.userTable, event.Payload.UserId)
}
