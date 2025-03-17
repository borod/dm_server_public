package dm_list

import (
	h "dm_server/dm_helper"
)

const __name = "DM List"

func init() {
	h.LogRoutineStart(__name)

	InitRoutes()

	h.LogRoutineEnd(__name)
}
