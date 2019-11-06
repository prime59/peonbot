package peonbot

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func getAction(command string, payload _payload) _event {
	return _event{
		Command:   command,
		RequestId: 1,
		Payload:   payload,
	}
}

func getExpectedRequest(command string, rid int, payload interface{}) _request {
	return _request{
		Command:   command,
		RequestId: rid, /* rid is never set from handleAction */
		Payload:   payload,
	}
}

func assertDeepEqualsRequest(expected _request, actual _request) error {
	if strings.Compare(expected.Command, actual.Command) != 0 {
		return fmt.Errorf("Expected: %s, Actual: %s", expected.Command, actual.Command)
	}

	if expected.RequestId != actual.RequestId {
		return fmt.Errorf("Expected: %d, Actual: %d", expected.RequestId, actual.RequestId)
	}

	return nil
}

func assertDeepEqualsPayloadAuth(expected _payloadAuth, actual _payloadAuth) error {
	if strings.Compare(
		expected.ApiKey,
		actual.ApiKey) != 0 {

		return fmt.Errorf(
			"Expected: %s, Actual: %s", expected.ApiKey, actual.ApiKey)
	}

	return nil
}

func assertDeepEqualsPayloadMessage(expected _payloadMessage, actual _payloadMessage) error {
	if strings.Compare(expected.Message, actual.Message) != 0 {
		return fmt.Errorf("Expected: %s, Actual: %s", expected.Message, actual.Message)
	}

	if strings.Compare(expected.UserId, actual.UserId) != 0 {
		return fmt.Errorf("Expected: %s, Actual: %s", expected.UserId, actual.UserId)
	}

	return nil
}

func assertDeepEqualsPayloadAction(expected _payloadAction, actual _payloadAction) error {
	if expected.UserId != actual.UserId {
		return fmt.Errorf("Expected: %d, Actual: %d", expected.UserId, actual.UserId)
	}

	if strings.Compare(expected.ToonName, actual.ToonName) != 0 {
		return fmt.Errorf("Expected: %s, Actual: %s", expected.ToonName, actual.ToonName)
	}

	return nil
}

func handleAndReturn(client *echoClient, testbot *_bot, action _event) error {
	if err := handleAction(client, testbot, action); err != nil {
		return fmt.Errorf("Error handling action: %v\n", err)
	}

	return nil
}

func TestActionNonPrivUser(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_59,
		Message: ".say hi",
	})

	testbot := getTestbot()

	if err := handleAction(nil, testbot, action); err == nil {
		t.Errorf("Expected error, but got nil. Action originated from an unpriveleged user.")
	}
}

func TestActionNoAction(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: "hi",
	})

	testbot := getTestbot()

	if err := handleAction(nil, testbot, action); err == nil {
		t.Errorf("Expected error, but got nil. User was priveleged, but message did not contain an action.")
	}
}

/* Test correct code path when a priveleged user sends only an action, such as `.say`. */
func TestActionIncompleteAction(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".say",
	})

	testbot := getTestbot()

	if err := handleAction(nil, testbot, action); err == nil {
		t.Errorf("Expected error, but got nil. User was priveleged, but action was incomplete.")
	}
}

func TestGetTarget(t *testing.T) {
	expectedTarget := "user#gateway"

	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".kick " + expectedTarget,
	})

	testbot := getTestbot()

	actualTarget, acceptable := getTarget(nil, testbot, expectedTarget, action)

	if strings.Compare(expectedTarget, actualTarget) != 0 {
		t.Errorf("Expected: %s, Actual: %s", expectedTarget, actualTarget)
	}

	if acceptable != true {
		t.Errorf("username's target should have been acceptable.")
	}
}

func TestGetTargetNoGatewayMessage(t *testing.T) {
	expectedTarget := ""
	unacceptableTarget := "userNogateway"

	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".kick " + unacceptableTarget,
		Type:    _MSG_CHAN,
	})

	var actual _request
	expectedPayload := _payloadMessage{Message: _NOTIFICATION_NO_GATEWAY}
	expected := _request{
		Command:   _REQUEST_MSG,
		RequestId: 1,
		Payload:   expectedPayload,
	}

	client := getEchoClient()
	testbot := getTestbot()

	actualTarget, acceptable := getTarget(client, testbot, unacceptableTarget, action)

	if strings.Compare(expectedTarget, actualTarget) != 0 {
		t.Errorf("Expected: '%s', Actual: '%s'", expectedTarget, actualTarget)
	}

	if acceptable != false {
		t.Errorf("username's target should have been unacceptable.")
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}

func TestGetTargetNoGatewayWhisper(t *testing.T) {
	expectedTarget := ""
	unacceptableTarget := "userNogateway"

	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".kick " + unacceptableTarget,
		Type:    _MSG_WHISPER,
	})

	var actual _request
	expectedPayload := _payloadMessage{
		Message: _NOTIFICATION_NO_GATEWAY,
		UserId:  strconv.Itoa(_TEST_USERID_155),
	}
	expected := _request{
		Command:   _REQUEST_WHISPER,
		RequestId: 1,
		Payload:   expectedPayload,
	}

	client := getEchoClient()
	testbot := getTestbot()

	actualTarget, acceptable := getTarget(client, testbot, unacceptableTarget, action)

	if strings.Compare(expectedTarget, actualTarget) != 0 {
		t.Errorf("Expected: '%s', Actual: '%s'", expectedTarget, actualTarget)
	}

	if acceptable != false {
		t.Errorf("username's target should have been unacceptable.")
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}

func TestGetTargetToonNameTooLongMessage(t *testing.T) {
	expectedTarget := ""
	unacceptableTarget := "usernameMoreThan15Chars#Gateway"

	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".kick " + unacceptableTarget,
		Type:    _MSG_CHAN,
	})

	var actual _request
	expectedPayload := _payloadMessage{Message: _NOTIFICATION_NAME_TOO_LONG}
	expected := _request{
		Command:   _REQUEST_MSG,
		RequestId: 1,
		Payload:   expectedPayload,
	}

	client := getEchoClient()
	testbot := getTestbot()

	actualTarget, acceptable := getTarget(client, testbot, unacceptableTarget, action)

	if strings.Compare(expectedTarget, actualTarget) != 0 {
		t.Errorf("Expected: '%s', Actual: '%s'", expectedTarget, actualTarget)
	}

	if acceptable != false {
		t.Errorf("username's target should have been unacceptable.")
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}

func TestGetTargetToonNameTooLongWhisper(t *testing.T) {
	expectedTarget := ""
	unacceptableTarget := "usernameMoreThan15Chars#Gateway"

	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".kick " + unacceptableTarget,
		Type:    _MSG_WHISPER,
	})

	var actual _request
	expectedPayload := _payloadMessage{
		Message: _NOTIFICATION_NAME_TOO_LONG,
		UserId:  strconv.Itoa(_TEST_USERID_155),
	}
	expected := _request{
		Command:   _REQUEST_WHISPER,
		RequestId: 1,
		Payload:   expectedPayload,
	}

	client := getEchoClient()
	testbot := getTestbot()

	actualTarget, acceptable := getTarget(client, testbot, unacceptableTarget, action)

	if strings.Compare(expectedTarget, actualTarget) != 0 {
		t.Errorf("Expected: '%s', Actual: '%s'", expectedTarget, actualTarget)
	}

	if acceptable != false {
		t.Errorf("username's target should have been unacceptable.")
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}

func TestHandleActionKickUserDne(t *testing.T) {
	testbot := getTestbot()

	if err := handleActionKick(nil, testbot, _TEST_USERNAME_TESTUSER60); err == nil {
		t.Errorf("Expected an error from kicking a user that dne. Got nil.")
	}
}

func TestHandleActionKick(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".kick " + _TEST_USERNAME_PRIVUSER155,
	})

	var actual _request
	expectedPayload := _payloadAction{UserId: _TEST_USERID_155}
	expected := getExpectedRequest(_REQUEST_KICK, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAction(expectedPayload, actual.Payload.(_payloadAction)); err != nil {
		t.Error(err)
	}
}

func TestHandleActionBanUserDne(t *testing.T) {
	testbot := getTestbot()

	if err := handleActionBan(nil, testbot, _TEST_USERNAME_TESTUSER60); err == nil {
		t.Errorf("Expected an error from banning a user that dne. Got nil.")
	}
}

func TestHandleActionBan(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".ban " + _TEST_USERNAME_PRIVUSER155,
	})

	var actual _request
	expectedPayload := _payloadAction{UserId: _TEST_USERID_155}
	expected := getExpectedRequest(_REQUEST_BAN, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAction(expectedPayload, actual.Payload.(_payloadAction)); err != nil {
		t.Error(err)
	}
}

func TestHandleActionUnban(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".unban " + _TEST_USERNAME_PRIVUSER155,
	})

	var actual _request
	expectedPayload := _payloadAction{ToonName: _TEST_USERNAME_PRIVUSER155}
	expected := getExpectedRequest(_REQUEST_UNBAN, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAction(expectedPayload, actual.Payload.(_payloadAction)); err != nil {
		t.Error(err)
	}
}

func TestHandleActionSay(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".say hi",
	})

	var actual _request
	expectedPayload := _payloadMessage{
		Message: "hi",
	}
	expected := getExpectedRequest(_REQUEST_MSG, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}

func TestHandleActionWhisperUserDne(t *testing.T) {
	testbot := getTestbot()

	if err := handleActionWhisper(nil, testbot, _TEST_USERNAME_TESTUSER60); err == nil {
		t.Errorf("Expected an error from whispering a user that dne. Got nil.")
	}
}

func TestHandleActionWhisper(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".whisper " + _TEST_USERNAME_PRIVUSER155 + " blah blah blah...",
	})

	var actual _request
	expectedPayload := _payloadMessage{
		Message: "blah blah blah...",
		UserId:  strconv.Itoa(_TEST_USERID_155),
	}
	expected := getExpectedRequest(_REQUEST_WHISPER, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadMessage(expectedPayload, actual.Payload.(_payloadMessage)); err != nil {
		t.Error(err)
	}
}

func TestHandleActionDesignateUserDne(t *testing.T) {
	testbot := getTestbot()

	if err := handleActionDesignate(nil, testbot, _TEST_USERNAME_TESTUSER60); err == nil {
		t.Errorf("Expected an error from designating a user that dne. Got nil.")
	}
}

func TestHandleActionDesignate(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".designate " + _TEST_USERNAME_PRIVUSER155,
	})

	var actual _request
	expected := getExpectedRequest(_REQUEST_DESIGN, 1, nil)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}
}

func TestHandleActionAddPriv(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".addpriv " + _TEST_USERNAME_TESTUSER61_GATEWAY,
	})

	testbot := getTestbot()

	if err := handleAction(nil, testbot, action); err != nil {
		t.Errorf("Error handling action: %v\n", err)
	}

	if _, ok := testbot.pusers[strings.ToUpper(_TEST_USERNAME_TESTUSER61_GATEWAY)]; !ok {
		t.Errorf("User privelege should have been added, but was not.")
	}
}

func TestHandleActionRmPriv(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".rmpriv " + _TEST_USERNAME_PRIVUSER155,
	})

	testbot := getTestbot()

	if err := handleAction(nil, testbot, action); err != nil {
		t.Errorf("Error handling action: %v\n", err)
	}

	if _, ok := testbot.pusers[strings.ToUpper(_TEST_USERNAME_PRIVUSER155)]; ok {
		t.Errorf("User privelege should have been removed, but was not.")
	}
}

func TestHandleActionAddBan(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".addban " + _TEST_USERNAME_TESTUSER61_GATEWAY,
	})

	var actual _request
	expectedPayload := _payloadAction{UserId: _TEST_USERID_61}
	expected := getExpectedRequest(_REQUEST_BAN, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	if _, ok := testbot.blist[strings.ToUpper(_TEST_USERNAME_TESTUSER61_GATEWAY)]; !ok {
		t.Errorf("User should have been aded to the banlist, but was not.")
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

func TestHandleActionRmBan(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".rmban " + _TEST_USERNAME_TESTUSER61_GATEWAY,
	})

	var actual _request
	expectedPayload := _payloadAction{ToonName: _TEST_USERNAME_TESTUSER61_GATEWAY}
	expected := getExpectedRequest(_REQUEST_UNBAN, 1, expectedPayload)

	client := getEchoClient()
	testbot := getTestbot()

	/* Add user to ban list */
	testbot.blist[strings.ToUpper(_TEST_USERNAME_TESTUSER61_GATEWAY)] = nil

	if err := handleAndReturn(client, testbot, action); err != nil {
		t.Error(err)
	}

	if _, ok := testbot.blist[strings.ToUpper(_TEST_USERNAME_TESTUSER61_GATEWAY)]; ok {
		t.Errorf("Should should have been removed from the ban list, but was not.")
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

func TestActionUnhandledAction(t *testing.T) {
	action := getAction(_EVENT_MSG, _payload{
		UserId:  _TEST_USERID_155,
		Message: ".unrecognizedaction " + _TEST_USERNAME_TESTUSER61_GATEWAY,
	})

	testbot := getTestbot()

	if err := handleAndReturn(nil, testbot, action); err == nil {
		t.Errorf("Expected an error from sending an unrecognized action, but got nil.")
	}
}
