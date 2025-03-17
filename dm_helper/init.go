package dm_helper

const __name = "DM Helper"

const C_ObjType_Request = 1
const C_ObjType_Invoice = 2
const C_ObjType_Object = 3

const C_RequestCancelled = 0
const C_RequestDraft = 1
const C_RequestCreated = 2
const C_RequestAtWork = 3
const C_RequestPurchasing = 4

const C_InvoiceCancelled = 0
const C_InvoiceDraft = 1
const C_InvoiceCreated = 2
const C_InvoicePurchase = 3

const C_Verification_Created = 0
const C_Verification_WIP = 1
const C_Verification_Denied = 2
const C_Verification_Accepted = 3

const C_ChatUserUndefined = 0
const C_ChatUserSuccess = 1
const C_ChatUserAlready = 2
const C_ChatUserDoesNotExist = 3

const C_WIP = 3
const C_WorkInProgress = 4
const C_RightsFullAccess = 64

const C_RequestItemID = "RequestItemID"
const C_InvoiceID = "InvoiceID"

const C_Redirect_redmine = "redmine"

func init() {
	LogRoutineStart(__name)

	InitColors()

	LogRoutineEnd(__name)
}
