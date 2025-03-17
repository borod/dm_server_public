package dm_entity

import (
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"
)

const __name = "DM Entity"

func init() {
	h.LogRoutineStart(__name)

	InitRoutes()

	h.LogRoutineEnd(__name)
}
