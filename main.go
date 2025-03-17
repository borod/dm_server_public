package main

import (
	api "dm_server/dm_api"

	mysql "dm_server/dm_db/dm_mysql"
	excel "dm_server/dm_excel"

	h "dm_server/dm_helper"
)

const shouldCreateDatabase bool = false

// const shouldCreateDatabase bool = true

func main() {
	if false {
		excel.Do()
		return
	}

	h.Log(h.TimeCurrStrMS())

	if shouldCreateDatabase {
		mysql.CreateDB(mysql.GormDB)
	} else {
		api.StartServer()
	}
}
