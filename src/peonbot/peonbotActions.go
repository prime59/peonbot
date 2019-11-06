package peonbot

import (
	"fmt"
	"log"
	"strings"
)

const _ACTION_KICK = ".KICK"
const _ACTION_BAN = ".BAN"
const _ACTION_UNBAN = ".UNBAN"
const _ACTION_SAY = ".SAY"
const _ACTION_WHISPER = ".WHISPER"
const _ACTION_DESIGNATE = ".DESIGNATE"
const _ACTION_ADDPRIV = ".ADDPRIV"
const _ACTION_RMPRIV = ".RMPRIV"
const _ACTION_ADDBAN = ".ADDBAN"
const _ACTION_RMBAN = ".RMBAN"

/*
	XXX: Not implementing this any further until I can figure out how to
	differentiate between initial and subsequent user update events.
*/
// const _ACTION_SETGREET = ".SETGREET"

func errActionIgnoreIncomplete() error {
	return fmt.Errorf("Ignoring. Incomplete action.")
}

func (bot *_bot) errActionUserDne(errmsg string) error {
	bot.Vprintf("Dumping user table: %v\n", bot.userTable)
	return fmt.Errorf(errmsg)
}

func handleAction(client WebsocketClient, bot *_bot, event _event) error {
	if _, ok := bot.pusers[strings.ToUpper(
		bot.userTable[event.Payload.UserId])]; !ok {
		return fmt.Errorf("Ignoring. User is not priveleged.")
	}

	/* All actions begin with a ".", e.g. `.say hi` */
	if strings.Compare(string(event.Payload.Message[0]), ".") != 0 {
		return fmt.Errorf("Ignoring. No action to handle.")
	}

	parts := strings.Split(event.Payload.Message, " ")
	/*
		Ensure action command has a sufficient amount of information. All
		actions will have a .action verb, and at least one other parameter.
	*/
	if len(parts) < 2 {
		return errActionIgnoreIncomplete()
	}

	action := parts[0]

	switch strings.ToUpper(action) {
	case _ACTION_KICK:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}

		if err := handleActionKick(client, bot, target); err != nil {
			return err
		}
		break
	case _ACTION_BAN:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		if err := handleActionBan(client, bot, target); err != nil {
			return err
		}
		break
	case _ACTION_UNBAN:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		if err := handleActionUnban(client, bot, target); err != nil {
			return err
		}
		break
	case _ACTION_SAY:
		message := parts[1:len(parts)]
		if err := handleActionSay(client, bot, message...); err != nil {
			return err
		}
		break
	case _ACTION_WHISPER:
		/*
			A whisper action must be of the form: `.whisper user message`.
		*/
		if len(parts) < 3 {
			return errActionIgnoreIncomplete()
		}

		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		message := parts[2:len(parts)]
		if err := handleActionWhisper(client, bot, target, message...); err != nil {
			return err
		}
		break
	case _ACTION_DESIGNATE:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		if err := handleActionDesignate(client, bot, target); err != nil {
			return err
		}
		break
	case _ACTION_ADDPRIV:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		handleActionAddpriv(bot, target)
		break
	case _ACTION_RMPRIV:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		handleActionRmpriv(bot, target)
		break
	case _ACTION_ADDBAN:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		handleActionAddBan(client, bot, target)
		break
	case _ACTION_RMBAN:
		target, acceptable := getTarget(client, bot, parts[1], event)
		if !acceptable {
			return errActionIgnoreIncomplete()
		}
		handleActionRmBan(client, bot, target)
		break
	default:
		return fmt.Errorf("Unrecognized action: %v", event)
	}

	return nil
}

/*
	Add a quick check ensure the gateway is specified before passing
	transmitting the request. Also checks if usernames are within
	battle.net character limit.
*/

const _NOTIFICATION_NO_GATEWAY = "Must specify gateway with username. E.g. name#Azeroth."
const _NOTIFICATION_NAME_TOO_LONG = "Username cannot be >15 characters."

func getTarget(client WebsocketClient, bot *_bot, target string, event _event) (string, bool) {
	/*
		Attempt to let user know of lexicographical error before proceeding
		with the action. Ignore errors.
	*/
	if !strings.Contains(target, "#") {
		sendNotification(
			client, bot,
			_NOTIFICATION_NO_GATEWAY,
			event)
		return "", false
	}

	if len(strings.Split(target, "#")[0]) > 15 {
		sendNotification(
			client, bot,
			_NOTIFICATION_NAME_TOO_LONG,
			event)
		return "", false
	}

	return target, true
}

func sendNotification(client WebsocketClient, bot *_bot, message string, event _event) {
	switch strings.ToUpper(event.Payload.Type) {
	case _MSG_CHAN:
		_ = handleActionSay(client, bot, message)
	case _MSG_WHISPER:
		_ = _handleActionWhisper(client, bot, event.Payload.UserId, message)
	}
}

func handleActionKick(client WebsocketClient, bot *_bot, username string) error {
	uid := bot.lookupUid(username)
	if uid == -1 {
		return bot.errActionUserDne(
			fmt.Sprintf("Cannot kick. Username '%s' does not exist in my user table.\n",
				username))
	}

	return _handleActionKick(client, bot, uid)
}

func _handleActionKick(client WebsocketClient, bot *_bot, uid int) error {
	request := bot.createRequestKick(uid)
	bot.Vprintf("Sending request: %+v\n", request)

	return client.WriteJSON(request)
}

func handleActionBan(client WebsocketClient, bot *_bot, username string) error {
	uid := bot.lookupUid(username)
	if uid == -1 {
		return bot.errActionUserDne(
			fmt.Sprintf("Cannot ban. Username '%s' does not exist in my user table.\n",
				username))
	}

	return _handleActionBan(client, bot, uid)
}

func _handleActionBan(client WebsocketClient, bot *_bot, uid int) error {
	request := bot.createRequestBan(uid)
	bot.Vprintf("Sending request: %+v\n", request)

	return client.WriteJSON(request)
}

func handleActionUnban(client WebsocketClient, bot *_bot, username string) error {
	request := bot.createRequestUnban(username)
	bot.Vprintf("Sending request: %+v\n", request)

	return client.WriteJSON(request)
}

func handleActionSay(client WebsocketClient, bot *_bot, message ...string) error {
	mstring := strings.Join(message, " ")
	request := bot.createRequestMessage(mstring)
	bot.Vprintf("Sending request: %+v\n", request)

	return client.WriteJSON(request)
}

func handleActionWhisper(client WebsocketClient, bot *_bot, username string, message ...string) error {
	uid := bot.lookupUid(username)
	if uid == -1 {
		return bot.errActionUserDne(
			fmt.Sprintf("Cannot whisper. Username '%s' does not exist in my user table.\n",
				username))
	}

	return _handleActionWhisper(client, bot, uid, message...)
}

func _handleActionWhisper(client WebsocketClient, bot *_bot, uid int, message ...string) error {
	mstring := strings.Join(message, " ")
	request := bot.createRequestWhisper(uid, mstring)
	bot.Vprintf("Sending request: %+v\n", request)

	return client.WriteJSON(request)
}

func handleActionDesignate(client WebsocketClient, bot *_bot, username string) error {
	uid := bot.lookupUid(username)
	if uid == -1 {
		return bot.errActionUserDne(
			fmt.Sprintf("Cannot designate. Username '%s' does not exist in my user table.\n",
				username))
	}

	return _handleActionDesignate(client, bot, uid)
}

func _handleActionDesignate(client WebsocketClient, bot *_bot, uid int) error {
	request := bot.createRequestDesignate(uid)
	bot.Vprintf("Sending request: %+v\n", request)

	return client.WriteJSON(request)
}

func handleActionAddpriv(bot *_bot, target string) {
	bot.addPrivelegedUsers(target)
	log.Printf("[Bot log message] Privelege added: %s\n", target)
}

func handleActionRmpriv(bot *_bot, target string) {
	bot.rmPrivelegedUser(target)
	log.Printf("[Bot log message] Privelege removed: %s\n", target)
}

func handleActionAddBan(client WebsocketClient, bot *_bot, target string) {
	bot.addToBanlist(target)
	log.Printf("[Bot log message] Added to banlist: %s\n", target)
	if _, ok := bot.blist[strings.ToUpper(target)]; ok {
		_ = handleActionBan(client, bot, target)
	}
}

func handleActionRmBan(client WebsocketClient, bot *_bot, target string) {
	bot.rmFromBanlist(target)
	log.Printf("[Bot log message] Removed from banlist: %s\n", target)
	_ = handleActionUnban(client, bot, target)
}
