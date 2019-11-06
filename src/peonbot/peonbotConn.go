package peonbot

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const _BNET_BOT_ADDR = "wss://connect-bot.classic.blizzard.com/v1/rpc/chat"
const _DIAL_TCP = "tcp"
const _OS_WINDOWS = "WINDOWS"
const _OS_LINUX = "LINUX"
const _X509_EXPECTED_NAME = "classic.blizzard.com"
const _X509_SERVER_NAME = "x509: certificate is valid for *.classic.blizzard.com, classic.blizzard.com, not"
const _X509_UNKNOWN_AUTH = "x509: certificate signed by unknown authority"

func errConnSucceeded() error {
	return fmt.Errorf(
		"Connection attempt to '%s' succeeded, but it should have failed it is signed by an unknown authority.",
		_BNET_BOT_ADDR)
}

func errConnFailed(expectedErr string, actualErr string) error {
	return fmt.Errorf(
		"Connection attempt to '%s' failed, but not for the expected reason. Expected: '%s', Got: '%s'.",
		_BNET_BOT_ADDR, expectedErr, actualErr)
}

func (bot *_bot) Start() error {
	switch strings.ToUpper(runtime.GOOS) {
	case _OS_WINDOWS:
		if err := bot.dialWindows(_BNET_BOT_ADDR, _X509_EXPECTED_NAME); err != nil {
			return err
		}
	default:
		if err := bot.dial(); err != nil {
			return err
		}
	}
	log.Printf("Connected: %s\n", bot.Conn.UnderlyingConn().RemoteAddr())

	if err := bot.authenticate(bot.Conn, bot.Token()); err != nil {
		return err
	}
	/*
		It appears that you must acknowledge the authentication response
		before sending a connection request
	*/
	if _, _, err := bot.Conn.ReadMessage(); err != nil {
		for {
			if _, _, err := bot.Conn.ReadMessage(); err == nil {
				break
			}

			time.Sleep(time.Duration(3) * time.Second)
		}
	}

	if err := bot.connectBot(bot.Conn); err != nil {
		return err
	}

	return nil
}

func getDialer() *websocket.Dialer {
	return &websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: time.Second * time.Duration(45),
	}
}

/*
	A "placeholder" object to satisfy the tlsDialer interface. The only
	reason to define this is to facilite unit testing.
*/
type sanityChecker struct{}

func (checker *sanityChecker) Dial(network string, addr string, tlsConfig *tls.Config) (*tls.Conn, error) {
	return tls.Dial(network, addr, tlsConfig)
}

func (bot *_bot) dial() error {
	dialer := getDialer()

	/*
		Set expected server name per documentation:
		https://s3-us-west-1.amazonaws.com/static-assets.classic.blizzard.com/public/Chat+Bot+API+Alpha+v3.pdf
	*/
	dialer.TLSClientConfig = &tls.Config{
		ServerName: _X509_EXPECTED_NAME,
	}

	/*
		First connection attempt should fail. The server name should match
		but the cert is still signed by an unknown authority.
	*/
	if err := checkServerCertBefore(dialer, _BNET_BOT_ADDR); err != nil {
		return err
	}

	/*
		XXX: If anyone has a better way to do this, please let me know.

		Unfortunately, even if the server name on the cert is verified, you
		still can't establish a connection due to the cert not being a
		trusted. So instead, disable verification, then extract the underlying
		connection and reconnect to verify the cert's server name.
	*/
	dialer.TLSClientConfig.InsecureSkipVerify = true
	conn, _, err := dialer.Dial(_BNET_BOT_ADDR, nil)
	if err != nil {
		return err
	}

	if err := checkServerCertAfter(
		&sanityChecker{},
		conn.UnderlyingConn().RemoteAddr().String(),
		_X509_SANITY_CHECK_SERVER_NAME); err != nil {

		return err
	}

	bot.Conn = conn
	return nil
}

/*
	A websocket dialer interface. It is defined to enable easier unit
	testing of any functions that makes a call to
	`websocket.Dial(...)`.
*/
type wsDialer interface {
	Dial(string, http.Header) (*websocket.Conn, *http.Response, error)
}

func checkServerCertBefore(dialer wsDialer, addr string) error {
	// _, _, err := dialer.Dial(_BNET_BOT_ADDR, nil)
	_, _, err := dialer.Dial(addr, nil)
	if err == nil {
		return errConnSucceeded()
	}

	if err != nil {
		if strings.Compare(
			strings.ToUpper(err.Error()),
			strings.ToUpper(_X509_UNKNOWN_AUTH)) != 0 {

			return errConnFailed(_X509_UNKNOWN_AUTH, err.Error())
		}
	}

	return nil
}

/*
	A tls dialer interface. It is defined to enable easier unit
	testing of any functions that makes a call to
	`tls.Dial(...)`.
*/
type tlsDialer interface {
	Dial(string, string, *tls.Config) (*tls.Conn, error)
}

/*
	Generate a random uuid for testing cert, and to greatly reduce
	likelyhood that a MITM attack cannot guess the fake server name that
	this function certifies against.
*/
var _X509_SANITY_CHECK_SERVER_NAME = uuid.New().String()

func checkServerCertAfter(dialer tlsDialer, addr string, fakeServerName string) error {

	_, err := dialer.Dial(_DIAL_TCP, addr, &tls.Config{ServerName: fakeServerName})

	if err == nil {
		return errConnSucceeded()
	}

	/* Verify the server name on the cert is *.classic.blizzard.com */
	if err != nil {
		expectedErrMsg := fmt.Sprintf("%s %s", _X509_SERVER_NAME, fakeServerName)

		if strings.Compare(
			strings.ToUpper(expectedErrMsg),
			strings.ToUpper(err.Error())) != 0 {

			return errConnFailed(expectedErrMsg, err.Error())
		}
	}

	return nil
}

func (bot *_bot) authenticate(client WebsocketClient, token string) error {
	request := bot.createRequestAuth(token)

	return client.WriteJSON(request)
}

func (bot *_bot) connectBot(client WebsocketClient) error {
	request := bot.createRequestConn()

	return client.WriteJSON(request)
}
