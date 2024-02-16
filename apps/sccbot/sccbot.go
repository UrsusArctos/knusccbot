package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"projects/knusccbot/internal/kbconfig"
	"time"

	"github.com/UrsusArctos/dkit/pkg/aegisql"
	"github.com/UrsusArctos/dkit/pkg/daemonizer"
	"github.com/UrsusArctos/dkit/pkg/kotobot"
	"github.com/UrsusArctos/dkit/pkg/logmeow"
)

const (
	// General info
	projectName = "knuscc"
)

type (
	TKNUSCCBot struct {
		LinuxDaemon daemonizer.TLinuxDaemon
		Logger      logmeow.TLogMeow
		Config      kbconfig.TKNUSCCConfig
		SQLDB       aegisql.TAegiSQLDB
		Bot         kotobot.TKotoBot
		//Logic       dekalogic.TDekalogic
	}
)

func (KB *TKNUSCCBot) BotInit() (err error) {
	// Check if config file exists
	_, errInit := os.Stat(KB.LinuxDaemon.ConfFile)
	if errInit == nil {
		KB.Logger.LogEventInfo(fmt.Sprintf("Config file found: %s", KB.LinuxDaemon.ConfFile))
		// Load config file
		jsonConfig, _ := os.ReadFile(KB.LinuxDaemon.ConfFile)
		errInit = json.Unmarshal(jsonConfig, &KB.Config)
		if errInit == nil {
			KB.Logger.LogEventInfo("Config file loaded successfully")
			// Initialize DB connection
			dbptr, errDB := sql.Open(kbconfig.ConfigBackend,
				aegisql.MakeDSN(kbconfig.ConfigBackend,
					KB.Config.Username,
					KB.Config.Password,
					KB.Config.Protocol,
					KB.Config.Hostname,
					fmt.Sprint(KB.Config.Port),
					KB.Config.Database))
			if (dbptr != nil) && (errDB == nil) {
				KB.SQLDB = aegisql.TAegiSQLDB{DB: dbptr}
				KB.Logger.LogEventInfo("Database connection opened")
				// Load actual bot config: TG Bot API Token
				svalue, errGetVal := KB.Config.GetDBConfigValue(KB.SQLDB, kbconfig.KeyToken)
				if errGetVal == nil {
					KB.Logger.LogEventInfo("Token loaded")
					// Initialize TG Bot instance
					KB.Bot, errInit = kotobot.NewInstance(svalue)
					if errInit == nil {
						KB.Bot.ParseMode = kotobot.PMPlainText
						KB.Bot.Updates_StartWatch()
						KB.Logger.LogEventInfo(fmt.Sprintf("Bot started as @%s", KB.Bot.BotInfo.UserName))
						// TG Bot post-init
						// KB.Bot.MessageHandler = KB.Logic.MessageDispatcher
						// KB.Bot.CallbackHandler = KB.Logic.CallbackDispatcher
						// KB.Logic.RandGen = rand.New(rand.NewSource(time.Now().UnixNano()))
						// KB.Logic.Logger = &KB.Logger
						// KB.Logic.SQLDB = &KB.SQLDB
						// KB.Logic.Config = &KB.Config
						// KB.Logic.Bot = &KB.Bot
						// Load Dekalogic operational data
						// svalue, errGetVal = FB.Config.GetDBConfigValue(FB.SQLDB, fit2config.KeyHelpDesk)
						// slvalue, errGetSVal := FB.Config.GetDBConfigValue(FB.SQLDB, fit2config.KeyHelpDeskLegacy)
						// if (errGetVal == nil) && (errGetSVal == nil) {
						// 	FB.Logic.HelpDeskChatID, errGetVal = strconv.ParseInt(svalue, 10, 64)
						// 	FB.Logic.HelpDeskLegacyChatID, errGetSVal = strconv.ParseInt(slvalue, 10, 64)
						// 	if (errGetVal == nil) && (errGetSVal == nil) {
						// 	}
						// }
						if errGetVal != nil {
							KB.Logger.LogEventError(fmt.Sprintf("Error reading config: %s", errGetVal))
						}
					}
				} else {
					return errGetVal
				}
			} else {
				return errDB
			}
		}
	}
	return errInit
}

func (KB *TKNUSCCBot) BotClose() (err error) {
	KB.SQLDB.Close()
	return nil
}

func (KB *TKNUSCCBot) BotMain() (err error) {
	// Check for new TG Bot API events and process them
	if KB.Bot.Updates_ProcessAll() {
		// Some events were processed, now restart waiting for new events
		KB.Bot.Updates_StartWatch()
	} else {
		// Do some periodic things in the logic
		//KB.Logic.Periodic()
		// And sleep
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func main() {
	const strExiting = "Exiting"
	// Init daemon
	fitbot := TKNUSCCBot{LinuxDaemon: daemonizer.NewLinuxDaemon(projectName)}
	defer fitbot.LinuxDaemon.Close()
	fitbot.LinuxDaemon.FuncInit = fitbot.BotInit
	fitbot.LinuxDaemon.FuncClose = fitbot.BotClose
	fitbot.LinuxDaemon.FuncMain = fitbot.BotMain
	// Init logger
	var enfac uint8 = logmeow.FacFile
	if fitbot.LinuxDaemon.Foreground {
		enfac |= logmeow.FacConsole
	}
	fitbot.Logger = logmeow.NewLogMeow(projectName, enfac, fitbot.LinuxDaemon.LogPath)
	defer fitbot.Logger.Close()
	// Run daemon
	derror := fitbot.LinuxDaemon.Run()
	if derror != nil {
		fitbot.Logger.LogEventError(fmt.Sprintf("%s: %v", strExiting, derror))
	} else {
		fitbot.Logger.LogEventInfo(strExiting)
	}
}
