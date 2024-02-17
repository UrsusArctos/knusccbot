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
	tableMenu      = "menu"
	// Columns
	colCKey   = "ckey"
	colCValue = "cvalue"
	colCBID   = "cbid"
	colTitle  = "title"
	colParent = "parent"
	// Config variables
	KeyToken = "token"
	// SQL queries
	queryBotConfig    = "SELECT %s FROM %s.%s WHERE %s='%s';"
	queryLoadRootMenu = "SELECT %s,%s FROM %s.%s WHERE %s IS NULL;"
	queryLoadSubMenu  = "SELECT %s,%s FROM %s.%s WHERE %s=%d;"
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

func (KC TKNUSCCConfig) LoadMenuItems(dbi aegisql.TAegiSQLDB, parent *uint64) (menu []aegisql.TAegiSQLDataRow, err error) {
	var dataRows aegisql.TAegiSQLRows
	if parent == nil {
		dataRows, err = dbi.QueryData(fmt.Sprintf(queryLoadRootMenu, colCBID, colTitle, KC.Database, tableMenu, colParent))
	} else {
		a := fmt.Sprintf(queryLoadSubMenu, colCBID, colTitle, KC.Database, tableMenu, colParent, *parent)
		dataRows, err = dbi.QueryData(a)
	}
	//fmt.Printf("datarows = %+v\n", dataRows)

	if err == nil {
		defer dataRows.Close()
		var rowData aegisql.TAegiSQLDataRow = dataRows.UnloadNextRow()
		for ; rowData != nil; rowData = dataRows.UnloadNextRow() {
			menu = append(menu, rowData)
		}
		//fmt.Printf("%+v\n", menu)
		return menu, nil

		/*
			dataRows, err := dbi.QueryData(fmt.Sprintf(queryListTopics, colTKey, colTDescr, colMemo, colFileHash, colFileHash, FC.Database, tableHelpTopics))
			if err == nil {
				defer dataRows.Close()
				var rowData aegisql.TAegiSQLDataRow = dataRows.UnloadNextRow()
				for ; rowData != nil; rowData = dataRows.UnloadNextRow() {
					rowTop := THelpTopic{Key: rowData[colTKey], Description: rowData[colTDescr], Memo: rowData[colMemo], FileHash: rowData[colFileHash]}
					htop = append(htop, rowTop)
				}
				return htop, nil
			}
			return nil, err
		*/
	}
	return nil, nil
}
