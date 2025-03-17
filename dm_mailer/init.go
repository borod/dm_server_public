package dm_mailer

import (
	h "dm_server/dm_helper"
)

const __name = "DM Mailer"

const ANCHOR_TOKEN = "&TOKEN&"

func init() {
	h.LogRoutineStart(__name)

	h.LogRoutineEnd(__name)
}
