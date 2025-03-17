package dm_authorization

import (
	h "dm_server/dm_helper"
)

const __name = "DM Authorization"

func init() {
	h.LogRoutineStart(__name)

	InitRoutes()

	h.LogRoutineEnd(__name)
}
