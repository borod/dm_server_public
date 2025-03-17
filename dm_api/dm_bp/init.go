package dm_bp

import (
	h "dm_server/dm_helper"
)

const __name = "DM Business Process"

func init() {
	h.LogRoutineStart(__name)

	InitRoutes()

	h.LogRoutineEnd(__name)
}
