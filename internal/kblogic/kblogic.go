package kblogic

import (
	"fmt"
	"projects/knusccbot/internal/kbconfig"
	"strconv"

	"github.com/UrsusArctos/dkit/pkg/aegisql"
	"github.com/UrsusArctos/dkit/pkg/kotobot"
	"github.com/UrsusArctos/dkit/pkg/logmeow"
)

type (
	// KBLogic bundle to implement message processing logic in a dedicated package
	TKBLogic struct {
		Logger *logmeow.TLogMeow
		Bot    *kotobot.TKotoBot
		Config *kbconfig.TKNUSCCConfig
		SQLDB  *aegisql.TAegiSQLDB
	}
)

// === Auxiliary helper functions ===

// Check whether incoming message is a command (such as "/start") with optional parameters
func isCommand(msg kotobot.TMessage) (result bool, command string, parameter string) {
	if len(msg.Text) > 0 {
		i, err := fmt.Sscanf(msg.Text, "/%s", &command)
		// Check if parameter was supplied
		if len(msg.Text) > (len(command) + 2) {
			parameter = msg.Text[len(command)+2 : len(msg.Text)]
		} else {
			parameter = ""
		}
		return (i > 0) && (err == nil), command, parameter
	}
	return false, "", ""
}

// Formats informative user string
func userToString(u kotobot.TUser) string {
	return fmt.Sprintf("[%s %s (@%s) {%d}]", u.FirstName, u.LastName, u.UserName, u.ID)
}

// === Callback Handling ===

func (KBL *TKBLogic) CallbackDispatcher(cbq kotobot.TCallbackQuery) {
	// Call default handler
	KBL.callbackDefault(cbq)
	// Close callback query
	KBL.Bot.AnswerCallbackQuery(cbq)
}

func (KBL TKBLogic) callbackDefault(cbq kotobot.TCallbackQuery) {
	cbData, errCat := strconv.ParseUint(cbq.Data, 10, 64)
	if errCat == nil {
		KBL.Logger.LogEventInfo(fmt.Sprintf("cbData = %d\n", cbData))
		// Compose submenu
		rmenu, _ := KBL.Config.LoadMenuItems(*KBL.SQLDB, &cbData)

		ikb := kotobot.NewInlineKB(uint8(len(rmenu)), 1)
		for i, mi := range rmenu {
			cbid := mi["cbid"]
			ikb.SetButton(uint8(i), 0, mi["title"], &cbid)
		}
		KBL.Bot.SendMessage("menu", false, kotobot.RefmsgFromUID(cbq.From.ID), ikb.Egress(), nil)

	} else {
		KBL.Logger.LogEventError(fmt.Sprintf("Error interpreting callback ID: %+v", errCat))
	}
}

// === Message Handling ===

func (KBL *TKBLogic) MessageDispatcher(msg kotobot.TMessage) {
	// Detect commands
	iscmd, cmd, param := isCommand(msg)
	// Distinguish between private and group messages
	switch msg.Chat.Type {
	case "private":
		{
			// Check if this is a command or message
			if iscmd {
				KBL.commandPrivate(msg, cmd, param)
			} else {
				KBL.messagePrivate(msg)
			}
		}
	case "group", "supergroup":
		// if (msg.Chat.ID == KBL.HelpDeskChatID) || (msg.Chat.ID == KBL.HelpDeskLegacyChatID) {
		if iscmd {
			KBL.commandHelpDesk(msg, cmd, param)
		} else {
			KBL.messageHelpDesk(msg)
		}
		// }  else {
		// 	KBL.Logger.LogEventWarning(fmt.Sprintf("Message from unknown group %d", msg.Chat.ID))
		// }
	}
}

func (KBL TKBLogic) messagePrivate(msg kotobot.TMessage) {
	//
	KBL.Bot.SendMessage("OK", true, msg, nil, nil)
}

func (KBL TKBLogic) messageHelpDesk(msg kotobot.TMessage) {
	//
}

// === Command handlng ===

func (KBL TKBLogic) commandPrivate(msg kotobot.TMessage, command string, parameter string) {
	// Load topics
	// hTopics, err := KBL.Config.GetTopics(*KBL.SQLDB)
	// if err == nil {
	// Process command
	switch command {
	case "start":
		{ // This should produce root menu
			rmenu, _ := KBL.Config.LoadMenuItems(*KBL.SQLDB, nil)
			ikb := kotobot.NewInlineKB(uint8(len(rmenu)), 1)
			for i, mi := range rmenu {
				cbid := mi["cbid"]
				ikb.SetButton(uint8(i), 0, mi["title"], &cbid)
			}
			KBL.Bot.SendMessage("root_menu", false, msg, ikb.Egress(), nil)
		}
	case "test":
		{
			// 		KBL.sendWallMessage("Test wall message")
			//		fmt.Printf("param is <%s>\n", parameter)
			// 		KBL.sendInvite(60070647)
			// 		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		}
	default:
		{
		}
	} // end of switch
	// } else {
	// 	KBL.Logger.LogEventError(fmt.Sprintf("Error retrieving topic list: %+v", err))
	// }
}

func (KBL TKBLogic) commandHelpDesk(msg kotobot.TMessage, command string, parameter string) {
}
