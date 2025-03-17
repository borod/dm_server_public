package dm_redmine

import (
	h "dm_server/dm_helper"
)

const C_EndPointCreateUser = "/users.json"
const C_RedmineUserPrefix = "DM_"

const __name = "DM Redmine"

func init() {
	h.LogRoutineStart(__name)

	h.LogRoutineEnd(__name)
}
