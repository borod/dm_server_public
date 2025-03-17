package dm_mongo

import (
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"
)

const __name = "DM Mongo"

func init() {
	h.LogRoutineStart(__name)

	// InitMongoDB()

	h.LogRoutineEnd(__name)
}
