package peonbot

import (
	"strconv"
	"strings"
	"testing"
)

func TestCreateRequest(t *testing.T) {
	testbot := getTestbot()
	testtype := "TestType"

	expected := _request{
		Command:   testtype,
		RequestId: 1,
	}

	actual := testbot.createRequest(testtype)

	if strings.Compare(
		expected.Command, actual.Command) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			expected.Command, actual.Command)
	}

	if expected.RequestId != actual.RequestId {
		t.Errorf("Expected: %d, Actual: %d\n",
			expected.RequestId, actual.RequestId)
	}
}

func TestCreateRequestAuth(t *testing.T) {
	testbot := getTestbot()
	testApikey := "TestKey"

	expected := _payloadAuth{
		ApiKey: testApikey,
	}

	actual := testbot.createRequestAuth(testApikey)

	if strings.Compare(
		expected.ApiKey, actual.Payload.(_payloadAuth).ApiKey) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			expected.ApiKey, actual.Payload.(_payloadAuth).ApiKey)
	}
}

func TestCreateRequestConn(t *testing.T) {
	testbot := getTestbot()

	actual := testbot.createRequestConn()

	if actual.Payload != nil {
		t.Errorf("Expected: %v, Actual: %v\n",
			nil, actual.Payload)
	}
}

func TestCreateRequestMessage(t *testing.T) {
	testbot := getTestbot()
	testMessage := "TestMessage"

	expected := _payloadMessage{
		Message: testMessage,
	}

	actual := testbot.createRequestMessage(testMessage)

	if strings.Compare(
		expected.Message, actual.Payload.(_payloadMessage).Message) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			expected.Message, actual.Payload.(_payloadMessage).Message)
	}
}

func TestCreateRequestWhisper(t *testing.T) {
	testbot := getTestbot()
	_TEST_USERID_59 := 59
	testMessage := "TestWhisper"

	expected := _payloadMessage{
		Message: testMessage,
		UserId:  strconv.Itoa(_TEST_USERID_59),
	}

	actual := testbot.createRequestWhisper(_TEST_USERID_59, testMessage)

	if strings.Compare(
		expected.UserId, actual.Payload.(_payloadMessage).UserId) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			expected.UserId, actual.Payload.(_payloadMessage).UserId)
	}

	if strings.Compare(
		expected.Message, actual.Payload.(_payloadMessage).Message) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			expected.Message, actual.Payload.(_payloadMessage).Message)
	}
}

func TestCreateRequestKick(t *testing.T) {
	testbot := getTestbot()

	expected := _payloadAction{
		UserId: _TEST_USERID_59,
	}

	actual := testbot.createRequestKick(_TEST_USERID_59)

	if expected.UserId != actual.Payload.(_payloadAction).UserId {
		t.Errorf("Expected: %d, Actual: %d\n",
			expected.UserId, actual.Payload.(_payloadAction).UserId)
	}
}

func TestCreateRequestBan(t *testing.T) {
	testbot := getTestbot()

	expected := _payloadAction{
		UserId: _TEST_USERID_59,
	}

	actual := testbot.createRequestBan(_TEST_USERID_59)

	if expected.UserId != actual.Payload.(_payloadAction).UserId {
		t.Errorf("Expected: %d, Actual: %d\n",
			expected.UserId, actual.Payload.(_payloadAction).UserId)
	}
}

func TestCreateRequestUnban(t *testing.T) {
	testbot := getTestbot()

	expected := _payloadAction{
		ToonName: _TEST_USERNAME_TESTUSER60,
	}

	actual := testbot.createRequestUnban(_TEST_USERNAME_TESTUSER60)

	if strings.Compare(
		expected.ToonName, actual.Payload.(_payloadAction).ToonName) != 0 {

		t.Errorf("Expected: %s, Actual: %s\n",
			expected.ToonName, actual.Payload.(_payloadAction).ToonName)
	}
}

func TestCreateRequestDesignate(t *testing.T) {
	testbot := getTestbot()

	expected := _payloadAction{
		UserId: _TEST_USERID_59,
	}

	actual := testbot.createRequestDesignate(_TEST_USERID_59)

	if expected.UserId != actual.Payload.(_payloadAction).UserId {
		t.Errorf("Expected: %d, Actual: %d\n",
			expected.UserId, actual.Payload.(_payloadAction).UserId)
	}
}
