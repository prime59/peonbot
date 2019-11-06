package peonbot

import "strconv"

func (bot *_bot) getRid() int {
	bot.rid++
	return bot.rid
}

/* Overwrite rid with whatever is received from the server */
func (bot *_bot) setRid(rid int) {
	bot.rid = rid
}

const _REQUEST_AUTH = "Botapiauth.AuthenticateRequest"
const _REQUEST_CONN = "Botapichat.ConnectRequest"
const _REQUEST_DISC = "Botapichat.DisconnectRequest"
const _REQUEST_MSG = "Botapichat.SendMessageRequest"
const _REQUEST_WHISPER = "Botapichat.SendWhisperRequest"
const _REQUEST_BAN = "Botapichat.BanUserRequest"
const _REQUEST_UNBAN = "Botapichat.UnbanUserRequest"
const _REQUEST_KICK = "Botapichat.KickUserRequest"
const _REQUEST_DESIGN = "Botapichat.SendSetModeratorRequest"

/* XXX: Emote is not currently implemented */
// const _REQUEST_EMOTE = "Botapichat.SendEmoteRequest"

type _request struct {
	Command   string      `json:"command"`
	RequestId int         `json:"request_id"`
	Payload   interface{} `json:"payload"`
}

type _payloadAuth struct {
	ApiKey string `json:"api_key"`
}

type _payloadMessage struct {
	Message string `json:"message"`
	UserId  string `json:"user_id"`
}

/*
	For sending kick, ban, and unban requests to the server. Unban requires
	`toon_name`.
*/
type _payloadAction struct {
	UserId   int    `json:"user_id"`
	ToonName string `json:"toon_name"`
}

func (bot *_bot) createRequest(rtype string) _request {
	return _request{
		Command:   rtype,
		RequestId: bot.getRid(),
	}
}

func (bot *_bot) createRequestAuth(apikey string) _request {
	request := bot.createRequest(_REQUEST_AUTH)
	request.Payload = _payloadAuth{
		ApiKey: apikey,
	}
	return request
}

func (bot *_bot) createRequestConn() _request {
	return bot.createRequest(_REQUEST_CONN)
}

func (bot *_bot) createRequestMessage(message string) _request {
	request := bot.createRequest(_REQUEST_MSG)
	request.Payload = _payloadMessage{
		Message: message,
	}
	return request
}

func (bot *_bot) createRequestWhisper(uid int, message string) _request {
	request := bot.createRequest(_REQUEST_WHISPER)
	request.Payload = _payloadMessage{
		Message: message,
		UserId:  strconv.Itoa(uid),
	}
	return request
}

func (bot *_bot) createRequestKick(uid int) _request {
	request := bot.createRequest(_REQUEST_KICK)
	request.Payload = _payloadAction{
		UserId: uid,
	}
	return request
}

func (bot *_bot) createRequestBan(uid int) _request {
	request := bot.createRequest(_REQUEST_BAN)
	request.Payload = _payloadAction{
		UserId: uid,
	}
	return request
}

func (bot *_bot) createRequestUnban(name string) _request {
	request := bot.createRequest(_REQUEST_UNBAN)
	request.Payload = _payloadAction{
		ToonName: name,
	}
	return request
}

func (bot *_bot) createRequestDesignate(uid int) _request {
	request := bot.createRequest(_REQUEST_DESIGN)
	request.Payload = _payloadAction{
		UserId: uid,
	}
	return request
}
