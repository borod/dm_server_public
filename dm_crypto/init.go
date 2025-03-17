package dm_crypto

import (
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"
)

const __name = "DM Crypto"

func init() {
	h.LogRoutineStart(__name)

	InitCrypto()

	h.LogRoutineEnd(__name)
}
