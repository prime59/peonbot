package peonbot

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const _STDIN_MESSAGE_STRING = "message"
const _STDIN_MESSAGE_KICK = "/kick " + _TEST_USERNAME_TESTUSER61_GATEWAY
const _STDIN_MESSAGE_BAN = "/ban " + _TEST_USERNAME_TESTUSER61_GATEWAY
const _STDIN_MESSAGE_UNBAN = "/unban " + _TEST_USERNAME_TESTUSER61_GATEWAY
const _STDIN_MESSAGE_WHISPER = "/whisper " + _TEST_USERNAME_TESTUSER61_GATEWAY + " " + _STDIN_MESSAGE_STRING
const _STDIN_MESSAGE_DESIGNATE = "/designate " + _TEST_USERNAME_TESTUSER61_GATEWAY
const _STDIN_NOT_AN_ACTION_MESSAGE = "message with / in it ///"

var _ACTUAL_MESSAGES = []string{
	_STDIN_MESSAGE_KICK,
	_STDIN_MESSAGE_BAN,
	_STDIN_MESSAGE_UNBAN,
	_STDIN_MESSAGE_WHISPER,
	_STDIN_MESSAGE_DESIGNATE,
	_STDIN_NOT_AN_ACTION_MESSAGE,
}

func TestIsAction(t *testing.T) {
	expected := []bool{
		true,
		true,
		true,
		true,
		true,
		false,
	}

	failures := make([]string, 0)

	for i, action := range _ACTUAL_MESSAGES {
		isaction := isAction(action)
		if expected[i] != isaction {
			failures = append(
				failures,
				fmt.Sprintf("Expected: %t, Actual: %t", expected[i], isaction))
		}
	}

	if len(failures) > 0 {
		t.Errorf("Failures: %+v", failures)
	}
}

func TestFormatStdinAction(t *testing.T) {
	actual := _ACTUAL_MESSAGES[0 : len(_ACTUAL_MESSAGES)-1]
	expected := []string{
		"." + _STDIN_MESSAGE_KICK[1:len(_STDIN_MESSAGE_KICK)],
		"." + _STDIN_MESSAGE_BAN[1:len(_STDIN_MESSAGE_BAN)],
		"." + _STDIN_MESSAGE_UNBAN[1:len(_STDIN_MESSAGE_UNBAN)],
		"." + _STDIN_MESSAGE_WHISPER[1:len(_STDIN_MESSAGE_WHISPER)],
		"." + _STDIN_MESSAGE_DESIGNATE[1:len(_STDIN_MESSAGE_DESIGNATE)],
	}

	failures := make([]string, 0)

	for i, action := range actual {
		formatted := formatStdinAction(action)
		if expected[i] != formatted {
			failures = append(
				failures,
				fmt.Sprintf("Expected: %s, Actual: %s", expected[i], formatted))
		}
	}

	if len(failures) > 0 {
		t.Errorf("Failures: %+v\n", failures)
	}
}

func TestFormatStdinMessage(t *testing.T) {
	actual := formatStdinMessage(_STDIN_NOT_AN_ACTION_MESSAGE)
	expected := ".say " + _STDIN_NOT_AN_ACTION_MESSAGE

	if strings.Compare(expected, actual) != 0 {
		t.Errorf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestHandleMessageKick(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expectedPayload := _payloadAction{UserId: _TEST_USERID_61}
	expected := getExpectedRequest(_REQUEST_KICK, 1, expectedPayload)

	if err := testbot.HandleMessage(client, _STDIN_MESSAGE_KICK); err != nil {
		t.Errorf("Error handling stdin message: %v\n", err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAction(
		expectedPayload, actual.Payload.(_payloadAction)); err != nil {

		t.Error(err)
	}
}

func TestHandleMessageBan(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expectedPayload := _payloadAction{UserId: _TEST_USERID_61}
	expected := getExpectedRequest(_REQUEST_BAN, 1, expectedPayload)

	if err := testbot.HandleMessage(client, _STDIN_MESSAGE_BAN); err != nil {
		t.Errorf("Error handling stdin message: %v\n", err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAction(
		expectedPayload, actual.Payload.(_payloadAction)); err != nil {

		t.Error(err)
	}
}

func TestHandleMessageUnban(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expectedPayload := _payloadAction{ToonName: _TEST_USERNAME_TESTUSER61_GATEWAY}
	expected := getExpectedRequest(_REQUEST_UNBAN, 1, expectedPayload)

	if err := testbot.HandleMessage(client, _STDIN_MESSAGE_UNBAN); err != nil {
		t.Errorf("Error handling stdin message: %v\n", err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAction(
		expectedPayload, actual.Payload.(_payloadAction)); err != nil {

		t.Error(err)
	}
}

func TestHandleMessageWhisper(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expectedPayload := _payloadMessage{
		UserId:  strconv.Itoa(_TEST_USERID_61),
		Message: _STDIN_MESSAGE_STRING,
	}
	expected := getExpectedRequest(_REQUEST_WHISPER, 1, expectedPayload)

	if err := testbot.HandleMessage(client, _STDIN_MESSAGE_WHISPER); err != nil {
		t.Errorf("Error handling stdin message: %v\n", err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(
		expectedPayload, actual.Payload.(_payloadMessage)); err != nil {

		t.Error(err)
	}
}

func TestHandleMessageDesignate(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expected := getExpectedRequest(_REQUEST_DESIGN, 1, nil)

	if err := testbot.HandleMessage(client, _STDIN_MESSAGE_DESIGNATE); err != nil {
		t.Errorf("Error handling stdin message: %v\n", err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}
}

func TestHandleMessageSay(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expectedPayload := _payloadMessage{
		Message: _STDIN_NOT_AN_ACTION_MESSAGE,
	}
	expected := getExpectedRequest(_REQUEST_MSG, 1, expectedPayload)

	if err := testbot.HandleMessage(
		client, _STDIN_NOT_AN_ACTION_MESSAGE); err != nil {

		t.Errorf("Error handling stdin message: %v\n", err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(
		expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}
