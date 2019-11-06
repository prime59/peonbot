package peonbot

import (
	"peonbot/verbose"
	"strings"
	"testing"

	"github.com/google/uuid"
)

/*
	Similar to the test echo server defined in `peonbotConn_test.go`, this
	test client "echoes" back requests that client publishes to the
	websocket server via `websocket.Conn.WriteJSON(interface{})`. This mock
	client should be used for testing any call to `WriteJSON`(interface{})`.
*/
type echoClient struct {
	request _request
}

func getEchoClient() *echoClient {
	return &echoClient{}
}

func (e *echoClient) WriteJSON(v interface{}) error {
	e.request = v.(_request)
	return nil
}

const _TEST_USERID_59 = 59
const _TEST_USERID_60 = 60
const _TEST_USERID_61 = 61
const _TEST_USERID_155 = 155
const _TEST_USERID_159_BANNED = 159
const _TEST_USERNAME_TESTUSER59 = "TestUser59"
const _TEST_USERNAME_TESTUSER60 = "TestUser60"
const _TEST_USERNAME_TESTUSER61_GATEWAY = "TestUser61#Gateway"
const _TEST_USERNAME_PRIVUSER155 = "PrivUser155#Azeroth"
const _TEST_USERNAME_BANNED_BANNEDUSER159 = "BannedUser159"

func getTestbot() *_bot {
	userTable := map[int]string{
		_PEONBOT_USERID:  _PEONBOT_USERNAME,
		_TEST_USERID_59:  _TEST_USERNAME_TESTUSER59,
		_TEST_USERID_61:  _TEST_USERNAME_TESTUSER61_GATEWAY,
		_TEST_USERID_155: _TEST_USERNAME_PRIVUSER155,
	}

	blist := map[string]interface{}{
		strings.ToUpper(_TEST_USERNAME_BANNED_BANNEDUSER159): nil,
	}

	pusers := map[string]interface{}{
		strings.ToUpper(_PEONBOT_USERNAME):          nil,
		strings.ToUpper(_TEST_USERNAME_PRIVUSER155): nil,
	}

	return &_bot{
		Vprintf:   verbose.Vprintf,
		rid:       0,
		userTable: userTable,
		blist:     blist,
		pusers:    pusers,
	}
}

/*
	XXX: Leaving this here for posterity.

	I had defined a TestMain function to start a mock echo server for some
	of peonbot's unit tests. This was necessary because go was unable to
	register the echo server's request handler to localhost multiple times,
	and I could not figure out a good way to unregister the handler test
	over test.

	Unfortunately, defining a TestMain function disables the go test tool's
	ability to generate a code coverage report.

	Therefore I have modified the function signatures in peonbotActions.go.
	Instead of being bound to the _bot struct, they now accept an interface
	type whose only known API is `WriteJSON(interface{}) error`, and a
	pointer to the _bot struct itself. This enables me to actually test the
	package.
*/
// func TestMain(m *testing.M) {
// 	server := &http.Server{Addr: _ECHO_SERVER_PORT}

// 	/* Start mock websocket server */
// 	http.HandleFunc("/", echoServer)
// 	go server.ListenAndServe()
// }

func TestNewBot(t *testing.T) {
	apiKey := uuid.New()
	privUsers := []string{
		_TEST_USERNAME_PRIVUSER155,
	}
	greetings := "NOT IMPLEMENTED"
	bannedUsers := []string{
		_TEST_USERNAME_BANNED_BANNEDUSER159,
	}

	bot := New(apiKey.String(), privUsers, greetings, bannedUsers)

	if strings.Compare(apiKey.String(), bot.Token()) != 0 {
		t.Errorf("Expected: %s, Actual: %s\n", apiKey, bot.Token())
	}

	if bot.Chbnt() != bot.chbnt {
		t.Errorf("Bot's battle.net requests accessor does not equal expected channel.")
	}

	if bot.Cherr() != bot.cherr {
		t.Errorf("Bot's error channel does not equal expected channel.")
	}

	if bot.Chsin() != bot.chsin {
		t.Errorf("Bot's stdin channel does not equal expected channel.")
	}
}

func TestLookupUidUserExists(t *testing.T) {
	testbot := getTestbot()

	expected := _TEST_USERID_59
	actual := testbot.lookupUid(_TEST_USERNAME_TESTUSER59)

	if expected != actual {
		t.Errorf("Expected: %d, Actual: %d",
			expected, actual)
	}
}

func TestLookupUidUserDne(t *testing.T) {
	testbot := getTestbot()

	expected := -1
	actual := testbot.lookupUid("")

	if expected != actual {
		t.Errorf("Expected: %d, Actual: %d",
			expected, actual)
	}
}

func TestAddPrivelegeUser(t *testing.T) {
	testbot := getTestbot()

	testbot.addPrivelegedUsers(_TEST_USERNAME_TESTUSER60)

	if _, ok := testbot.pusers[strings.ToUpper(_TEST_USERNAME_TESTUSER60)]; !ok {
		t.Errorf("Bot priveleges should have been granted, but were not.")
	}
}

func TestRmPrivelegeUser(t *testing.T) {
	testbot := getTestbot()

	testbot.addPrivelegedUsers(_TEST_USERNAME_TESTUSER60)

	if _, ok := testbot.pusers[strings.ToUpper(_TEST_USERNAME_TESTUSER60)]; !ok {
		t.Errorf("Bot priveleges should have been granted, but were not.")
	}

	testbot.rmPrivelegedUser(_TEST_USERNAME_TESTUSER60)

	if _, ok := testbot.pusers[strings.ToUpper(_TEST_USERNAME_TESTUSER60)]; ok {
		t.Errorf("Bot priveleges should have been removed, but were not.")
	}
}

func TestAddBanlist(t *testing.T) {
	testbot := getTestbot()

	testbot.addToBanlist(_TEST_USERNAME_TESTUSER60)

	if _, ok := testbot.blist[strings.ToUpper(_TEST_USERNAME_TESTUSER60)]; !ok {
		t.Errorf("User should have been added to banlist, but was not.")
	}
}

func TestRmBanlist(t *testing.T) {
	testbot := getTestbot()

	testbot.addToBanlist(_TEST_USERNAME_TESTUSER60)

	if _, ok := testbot.blist[strings.ToUpper(_TEST_USERNAME_TESTUSER60)]; !ok {
		t.Errorf("User should have been added to banlist, but was not.")
	}

	testbot.rmFromBanlist(_TEST_USERNAME_TESTUSER60)

	if _, ok := testbot.blist[strings.ToUpper(_TEST_USERNAME_TESTUSER60)]; ok {
		t.Errorf("User should have been removed from banlist, but was not.")
	}
}
