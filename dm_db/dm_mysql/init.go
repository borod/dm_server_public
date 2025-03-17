package dm_mysql

import (
	h "dm_server/dm_helper"
)

const __name = "DM MySQL"

func init() {
	h.LogRoutineStart(__name)

	InitMySQLDB()

	h.LogRoutineEnd(__name)
}
