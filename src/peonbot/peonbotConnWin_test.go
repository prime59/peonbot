package peonbot

import (
	"fmt"
	"testing"
)

const _TEST_WSEP_TESTDIALWINDOWS = "TestDialWindows"
const _ECHO_SERVER_PORT_5960 = "5960"

func TestDialWindows(t *testing.T) {
	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5960,
		_TEST_WSEP_TESTDIALWINDOWS)
	serverName := ""

	go startEchoServer(_ECHO_SERVER_PORT_5960, _TEST_WSEP_TESTDIALWINDOWS)

	testbot := getTestbot()

	if err := testbot.dialWindows(addr, serverName); err != nil {
		t.Errorf("Error dialing echo server: %v\n", err)
	}
}
