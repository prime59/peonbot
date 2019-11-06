package peonbot

import (
	"crypto/tls"
)

/*
	On Windows, the dialer doesn't seem to complain about the cert being
	signed by an unknown authority, so all the error checking (except for
	the server name verification) has been removed.
*/

func (bot *_bot) dialWindows(addr string, serverName string) error {
	dialer := getDialer()

	dialer.TLSClientConfig = &tls.Config{
		ServerName: serverName,
	}

	conn, _, err := dialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	bot.Conn = conn
	return nil
}
