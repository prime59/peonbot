package peonbot

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
	A test echo server to echo back requests to the client.

	If a new test requires an echo server asset, the echo server must bind
	to a unique path and port. Various unit tests make use of this mock
	echo server, so be sure to not use a conflicting path and port. The
	paths and ports are listed as constants before the test.
*/

const _ECHO_SERVER_ADDR = "ws://localhost"

func echoServer(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	/* Listen for messages and echo them back */
	for {
		mtype, raw, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		if err := conn.WriteMessage(mtype, raw); err != nil {
			panic(err)
		}
	}
}

/*
	To make use of the echo server in a test, start this function on its
	own goroutine.
*/
func startEchoServer(port string, endpoint string) {
	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	/* Start mock websocket server */
	http.HandleFunc(fmt.Sprintf("/%s", endpoint), echoServer)
	log.Fatal(server.ListenAndServe())
}

/*
	Use this function to connect the test instance of the bot to the mock
	echo server, instead of an actual websocket server.
*/
const _ECHO_SERVER_MAXRETRY = 10

func connectToEchoServer(testbot *_bot, port string, endpoint string) error {
	var conn *websocket.Conn
	var err error
	dialer := getDialer()

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, port, endpoint)
	for i := 0; i < _ECHO_SERVER_MAXRETRY; i++ {
		conn, _, err = dialer.Dial(
			addr, nil)

		if err == nil {
			testbot.Conn = conn
			return nil
		}
	}

	return err
}

type mockDialerBeforeEcs struct{}

func (d *mockDialerBeforeEcs) Dial(addr string, header http.Header) (*websocket.Conn, *http.Response, error) {
	return nil, nil, nil
}

const _TEST_WSEP_TESTCHECKSERVERCERTBEFOREECS = "TestCheckServerCertBeforeEcs"
const _ECHO_SERVER_PORT_5962 = "5962"

func TestCheckServerCertBeforeErrConnSuceeded(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5962,
		_TEST_WSEP_TESTCHECKSERVERCERTBEFOREECS)

	dialer := mockDialerBeforeEcs{}

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5962,
		_TEST_WSEP_TESTCHECKSERVERCERTBEFOREECS)

	if err := checkServerCertBefore(&dialer, addr); err == nil {
		t.Errorf("Expected an error, but got nil.")
	}
}

type mockDialerBeforeEcf struct{}

func (d *mockDialerBeforeEcf) Dial(addr string, header http.Header) (*websocket.Conn, *http.Response, error) {
	return nil, nil, fmt.Errorf("Error: connection failed.")
}

const _TEST_WSEP_TESTCHECKSERVERCERTBEFOREECF = "TestCheckServerCertBeforeEcf"
const _ECHO_SERVER_PORT_5963 = "5963"

func TestCheckServerCertBeforeErrConnFailed(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5963,
		_TEST_WSEP_TESTCHECKSERVERCERTBEFOREECF)

	dialer := mockDialerBeforeEcf{}

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5963,
		_TEST_WSEP_TESTCHECKSERVERCERTBEFOREECF)

	if err := checkServerCertBefore(&dialer, addr); err == nil {
		t.Errorf("Expected an error, but got nil.")
	}
}

type mockDialerBeforeUA struct{}

func (d *mockDialerBeforeUA) Dial(addr string, header http.Header) (*websocket.Conn, *http.Response, error) {
	return nil, nil, fmt.Errorf(_X509_UNKNOWN_AUTH)
}

const _TEST_WSEP_TESTCHECKSERVERCERTBEFOREUA = "TestCheckServerCertBeforeUa"
const _ECHO_SERVER_PORT_5964 = "5964"

func TestCheckServerCertBeforeUa(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5964,
		_TEST_WSEP_TESTCHECKSERVERCERTBEFOREUA)

	dialer := mockDialerBeforeUA{}

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5964,
		_TEST_WSEP_TESTCHECKSERVERCERTBEFOREUA)

	if err := checkServerCertBefore(&dialer, addr); err != nil {
		t.Errorf("Checking server cert before should not have returned as error, but it did: %v\n", err)
	}
}

type mockDialerAfterEcs struct{}

func (d *mockDialerAfterEcs) Dial(network string, addr string, tlsConfig *tls.Config) (*tls.Conn, error) {
	return nil, nil
}

const _TEST_WSEP_TESTCHECKSERVERCERTAFTERECS = "TestServerCertAfterEcs"
const _ECHO_SERVER_PORT_5965 = "5965"

func TestCheckServerCertAfterErrConnSucceeded(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5965,
		_TEST_WSEP_TESTCHECKSERVERCERTAFTERECS)

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5965,
		_TEST_WSEP_TESTCHECKSERVERCERTAFTERECS)

	dialer := mockDialerAfterEcs{}

	if err := checkServerCertAfter(&dialer, addr, uuid.New().String()); err == nil {
		t.Errorf("Expected an error, but got nil.")
	}
}

type mockDialerAfterEcf struct {
	fakeServerName string
}

func (d *mockDialerAfterEcf) Dial(network string, addr string,
	tlsConfig *tls.Config) (*tls.Conn, error) {

	return nil, fmt.Errorf("Error: connection failed")
}

const _TEST_WSEP_TESTCHECKSERVERCERTAFTERECF = "TestServerCertAfterEcf"
const _ECHO_SERVER_PORT_5966 = "5966"

func TestCheckServerCertAfterErrConnFailed(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5966,
		_TEST_WSEP_TESTCHECKSERVERCERTAFTERECF)

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5966,
		_TEST_WSEP_TESTCHECKSERVERCERTAFTERECF)

	dialer := mockDialerAfterEcf{fakeServerName: uuid.New().String()}

	if err := checkServerCertAfter(&dialer, addr, dialer.fakeServerName); err == nil {
		t.Errorf("Expected an error, but got nil.")
	}
}

type mockDialerAfterSN struct {
	fakeServerName string
}

func (d *mockDialerAfterSN) Dial(network string, addr string,
	tlsConfig *tls.Config) (*tls.Conn, error) {

	return nil, fmt.Errorf(fmt.Sprintf("%s %s", _X509_SERVER_NAME, d.fakeServerName))
}

const _TEST_WSEP_TESTCHECKSERVERCERTAFTERSN = "TestServerCertAfterSN"
const _ECHO_SERVER_PORT_5967 = "5967"

func TestCheckServerCertAfterServerName(t *testing.T) {
	go startEchoServer(
		_ECHO_SERVER_PORT_5967,
		_TEST_WSEP_TESTCHECKSERVERCERTAFTERSN)

	addr := fmt.Sprintf("%s:%s/%s", _ECHO_SERVER_ADDR, _ECHO_SERVER_PORT_5967,
		_TEST_WSEP_TESTCHECKSERVERCERTAFTERSN)

	dialer := mockDialerAfterSN{fakeServerName: uuid.New().String()}

	if err := checkServerCertAfter(&dialer, addr, dialer.fakeServerName); err != nil {
		t.Errorf("Checking server cert after should not have returned as error, but it did: %v\n", err)
	}
}

func TestAuthenticate(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()
	token := uuid.New()

	var actual _request
	expectedPayload := _payloadAuth{ApiKey: token.String()}
	expected := getExpectedRequest(_REQUEST_AUTH, 1, expectedPayload)

	testbot.authenticate(client, token.String())

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if err := assertDeepEqualsPayloadAuth(expectedPayload, actual.Payload.(_payloadAuth)); err != nil {
		t.Error(err)
	}
}

func TestConnectBot(t *testing.T) {
	client := getEchoClient()
	testbot := getTestbot()

	var actual _request
	expected := getExpectedRequest(_REQUEST_CONN, 1, nil)

	testbot.connectBot(client)

	actual = client.request

	if err := assertDeepEqualsRequest(expected, actual); err != nil {
		t.Error(err)
	}

	if actual.Payload != nil {
		t.Errorf("Expected payload of request to be nil, but it was not: %+v\n", actual.Payload)
	}
}
