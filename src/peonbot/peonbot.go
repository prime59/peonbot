package peonbot

import (
	"peonbot/verbose"
	"strings"

	"github.com/gorilla/websocket"
)

/*
	An websocket connection interface. It is defined to enable easier unit
	testing of any functions that make a call to
	`websocket.Conn.WriteJSON(interface{})`.
*/
type WebsocketClient interface {
	WriteJSON(interface{}) error
}

type _bot struct {
	Vprintf func(string, ...interface{})

	token string

	Conn      *websocket.Conn
	userTable map[int]string
	rid       int /* request id used to communicate with bot API */

	chbnt chan []byte /* responses from websocket */
	cherr chan error
	chsin chan string /* string input from stdin */

	blist     map[string]interface{}
	greetings string
	pusers    map[string]interface{}
}

const _PEONBOT_USERID = -59
const _PEONBOT_USERNAME = "*SELF"

func New(token string, blist []string, greetings string, pusers []string) *_bot {
	var bot _bot

	bot.Vprintf = verbose.Vprintf

	bot.token = token

	bot.userTable = make(map[int]string)
	bot.addSelfToUserTable()

	bot.rid = 0

	bot.chbnt = make(chan []byte)
	bot.cherr = make(chan error)
	bot.chsin = make(chan string)

	bot.blist = make(map[string]interface{})
	bot.addToBanlist(blist...)

	/* XXX: Set greetings not implemented. */
	/* Set greetings message */
	// bot.setGreetings(greetings)

	bot.pusers = make(map[string]interface{})
	bot.addPrivToSelf()
	bot.addPrivelegedUsers(pusers...)

	bot.Vprintf("%+v\n", bot)

	return &bot
}

func (bot *_bot) addSelfToUserTable() {
	bot.userTable[_PEONBOT_USERID] = _PEONBOT_USERNAME
}

func (bot *_bot) addPrivToSelf() {
	bot.pusers[strings.ToUpper(_PEONBOT_USERNAME)] = nil
}

func (bot *_bot) lookupUid(username string) int {
	for uid, user := range bot.userTable {
		if strings.Compare(strings.ToUpper(username),
			strings.ToUpper(user)) == 0 {
			return uid
		}
	}

	return -1
}

func (bot *_bot) addPrivelegedUsers(pusers ...string) {
	for _, puser := range pusers {
		bot.pusers[strings.ToUpper(puser)] = nil
	}
}

func (bot *_bot) rmPrivelegedUser(pusers ...string) {
	for _, puser := range pusers {
		delete(bot.pusers, strings.ToUpper(puser))
	}
}

func (bot *_bot) addToBanlist(busers ...string) {
	for _, buser := range busers {
		bot.blist[strings.ToUpper(buser)] = nil
	}
}

func (bot *_bot) rmFromBanlist(busers ...string) {
	for _, buser := range busers {
		delete(bot.blist, strings.ToUpper(buser))
	}
}

/* XXX: Greetings is not implemented. */
/*
func (bot *_bot) setGreetings(message string) {
	bot.greetings = message
}
*/

func (bot *_bot) Token() string {
	return bot.token
}

func (bot *_bot) Chbnt() chan []byte {
	return bot.chbnt
}

func (bot *_bot) Cherr() chan error {
	return bot.cherr
}

func (bot *_bot) Chsin() chan string {
	return bot.chsin
}

func (bot *_bot) ListenWebsocket() {
	for {
		_, data, err := bot.Conn.ReadMessage()
		if err != nil {
			bot.cherr <- err
			continue
		}

		bot.chbnt <- data
	}

}
