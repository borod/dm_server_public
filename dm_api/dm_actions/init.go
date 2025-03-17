package dm_actions

import (
	h "dm_server/dm_helper"
)

const __name = "DM Actions"

func init() {
	h.LogRoutineStart(__name)

	InitRoutes()

	h.LogRoutineEnd(__name)
}
