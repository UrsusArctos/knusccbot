package kbconfig

import (
	"fmt"

	"github.com/UrsusArctos/dkit/pkg/aegisql"
)

const (
	// Backend
	ConfigBackend = "mysql"
	// DB structure
	// Tables
	tableBotConfig = "botconfig"
	// Columns
	colCKey   = "ckey"
	colCValue = "cvalue"
	// Config variables
	KeyToken = "token"
	// SQL queries
	queryBotConfig = "SELECT %s FROM %s.%s WHERE %s='%s';"
)

type (
	TKNUSCCConfig struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
		Protocol string `json:"Protocol"`
		Hostname string `json:"Hostname"`
		Port     uint16 `json:"Port"`
		Database string `json:"Database"`
	}
)

func (KC TKNUSCCConfig) GetDBConfigValue(dbi aegisql.TAegiSQLDB, ckey string) (string, error) {
	dataRows, err := dbi.QueryData(fmt.Sprintf(queryBotConfig, colCValue, KC.Database, tableBotConfig, colCKey, ckey))
	if err == nil {
		defer dataRows.Close()
		rowData := dataRows.UnloadNextRow()
		if rowData != nil {
			return rowData[colCValue], nil
		}
	}
	return "", err
}
