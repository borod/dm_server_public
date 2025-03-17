package dm_db

import (
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"
)

const __name = "DM DB"

func init() {
	h.LogRoutineStart(__name)

	h.LogRoutineEnd(__name)
}
