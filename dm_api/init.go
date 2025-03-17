package dm_api

import (
	h "dm_server/dm_helper"
)

const __name = "DM API"

func init() {
	h.LogRoutineStart(__name)

	InitRoutes()

	h.LogRoutineEnd(__name)
}
