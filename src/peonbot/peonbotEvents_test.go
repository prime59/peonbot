package peonbot

import (
	"encoding/json"
	"strings"
	"testing"
)

const _TEST_REQUESTID = 55

func getEncodedEventMessage(etype string) ([]byte, _event) {
	event := _event{
		Command:   _EVENT_MSG,
		RequestId: _TEST_REQUESTID,
		Payload: _payload{
			Type:    etype,
			UserId:  _TEST_USERID_59,
			Message: "TestMessage",
		},
	}

	raw, _ := json.Marshal(&event)

	return raw, event
}

func getEncodedEventUserUpdate(username string, userid int) ([]byte, _event) {
	event := _event{
		Command:   _EVENT_USERUPDATE,
		RequestId: _TEST_REQUESTID,
		Payload: _payload{
			ToonName: username,
			UserId:   userid,
		},
	}

	raw, _ := json.Marshal(&event)

	return raw, event
}

func getEncodedEventUserExit() ([]byte, _event) {
	event := _event{
		Command:   _EVENT_USEREXIT,
		RequestId: _TEST_REQUESTID,
		Payload: _payload{
			UserId: _TEST_USERID_59,
		},
	}

	raw, _ := json.Marshal(&event)

	return raw, event
}

const _TEST_UNKNOWN_EVENT = "Unknown_Event"

func getEncodedEventUnknownEvent() ([]byte, _event) {
	event := _event{
		Command:   _TEST_UNKNOWN_EVENT,
		RequestId: 0,
	}

	raw, _ := json.Marshal(&event)

	return raw, event
}

func TestHandleEventMessage(t *testing.T) {
	testbot := getTestbot()
	raw, event := getEncodedEventMessage(_MSG_CHAN)

	if err := testbot.HandleEvent(raw); err != nil {
		t.Logf("Got error handling message event: %+v\n", event)
		t.Errorf("Error: %v\n", err)
	}
}

func TestHandleEventWhisper(t *testing.T) {
	testbot := getTestbot()
	raw, event := getEncodedEventMessage(_MSG_WHISPER)

	if err := testbot.HandleEvent(raw); err != nil {
		t.Logf("Got error handling whisper event: %+v\n", event)
		t.Errorf("Error: %v\n", err)
	}
}

func TestHandleEventUserUpdate(t *testing.T) {
	testbot := getTestbot()
	raw, event := getEncodedEventUserUpdate("TestUser60", 60)

	if err := testbot.HandleEvent(raw); err != nil {
		t.Logf("Got error handling user update event: %+v\n", event)
		t.Errorf("Error: %v\n", err)
	}

	if strings.Compare(
		strings.ToUpper(testbot.userTable[event.Payload.UserId]),
		strings.ToUpper(event.Payload.ToonName)) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			strings.ToUpper(testbot.userTable[event.Payload.UserId]),
			strings.ToUpper(event.Payload.ToonName))
	}
}

const _TEST_USERID_59_BANNED = 159
const _TEST_USERNAME_BANNED_BANNEDUSER159 = "BannedUser159"
const _TEST_WSEP_TESTHANDLEEVENTUSERUPDATEBANLISTUSER = "TestHandleEventUserUpdateBanlistUser"
const _ECHO_SERVER_PORT_5959 = "5959"

func TestHandleEventUserUpdateBanlistUser(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5959,
		_TEST_WSEP_TESTHANDLEEVENTUSERUPDATEBANLISTUSER)

	var actual _event
	expected := _request{
		Command:   _REQUEST_BAN,
		RequestId: _TEST_REQUESTID + 1,
		Payload: _payloadAction{
			UserId: _TEST_USERID_59_BANNED,
		},
	}

	testbot := getTestbot()

	if err := connectToEchoServer(
		testbot,
		_ECHO_SERVER_PORT_5959,
		_TEST_WSEP_TESTHANDLEEVENTUSERUPDATEBANLISTUSER); err != nil {

		t.Errorf("Error connecting to echo server: %v\n", err)
	}

	/* Get event that a user on the ban list has joined the channel */
	rawreq, event := getEncodedEventUserUpdate(
		_TEST_USERNAME_BANNED_BANNEDUSER159, _TEST_USERID_59_BANNED)

	if err := testbot.HandleEvent(rawreq); err != nil {
		t.Logf("Got error handling user update event: %+v\n", event)
		t.Errorf("Error: %v\n", err)
	}

	_, rawresp, err := testbot.Conn.ReadMessage()
	if err != nil {
		t.Errorf("Error reading echo: %v\n", err)
	}

	if err := json.Unmarshal(rawresp, &actual); err != nil {
		t.Errorf("Error unmarshalling echo: %v\n", err)
	}

	if strings.Compare(
		expected.Command, actual.Command) != 0 {

		t.Errorf("Expected: %s, Actual: %s",
			expected.Command, actual.Command)
	}

	if expected.RequestId != actual.RequestId {
		t.Errorf("Expected: %d, Actual: %d",
			expected.RequestId, actual.RequestId)
	}

	if expected.Payload.(_payloadAction).UserId != actual.Payload.UserId {
		t.Errorf("Expected: %d, Actual: %d",
			expected.Payload.(_payloadAction).UserId, actual.Payload.UserId)
	}
}

func TestHandleEventUserExit(t *testing.T) {
	testbot := getTestbot()
	raw, event := getEncodedEventUserExit()

	if err := testbot.HandleEvent(raw); err != nil {
		t.Logf("Got error handling user exit event: %+v\n", event)
		t.Errorf("Error: %v\n", err)
	}

	if _, ok := testbot.userTable[event.Payload.UserId]; ok {
		t.Errorf("Expected: %v, Actual: %v\n",
			nil, testbot.userTable[event.Payload.UserId])
	}
}

func TestHandleEventUnknownEvent(t *testing.T) {
	testbot := getTestbot()
	raw, _ := getEncodedEventUnknownEvent()

	if err := testbot.HandleEvent(raw); err == nil {
		t.Errorf("Expected error from sending unknown event. Got nil.")
	}
}
