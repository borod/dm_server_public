package dm_configuration

import (
	h "dm_server/dm_helper"
)

const __name = "DM Configuration"

const confFilePath = "configuration.json"

func init() {
	h.LogRoutineStart(__name)

	ReloadConfig()

	h.LogRoutineEnd(__name)
}
