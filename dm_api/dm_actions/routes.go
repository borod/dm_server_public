package dm_actions

import (
	h "dm_server/dm_helper"

	"net/http"
)

var Routes = map[string]http.HandlerFunc{}

func InitRoutes() {
	h.Log("Инициализация карты маршрутов...")

	Routes["/api/actions/user_chats"] = HandleUserChats
	h.Log("Добавлен маршрут /api/actions/user_chats -> HandleUserChats")

	Routes["/api/actions/add_verifier"] = HandleAddVerifier
	h.Log("Добавлен маршрут /api/actions/add_verifier -> InitVerifications")

	Routes["/api/actions/invoice_to_verify"] = HandleInvoicetToVerify
	h.Log("Добавлен маршрут /api/actions/invoice_to_verify -> HandleInvoicetToVerify")

	Routes["/api/actions/request_to_verify"] = HandleRequestToVerify
	h.Log("Добавлен маршрут /api/actions/request_to_verify -> HandleRequestToWork")

	Routes["/api/actions/init_verifications"] = HandleInitVerifications
	h.Log("Добавлен маршрут /api/actions/init_verifications -> InitVerifications")

	Routes["/api/actions/page_data"] = GetPageData
	h.Log("Добавлен маршрут /api/actions/page_data -> GetPageData")

	Routes["/api/actions/chat_participants"] = GetChatParticipants
	h.Log("Добавлен маршрут /api/actions/chat_participants -> GetChatParticipants")

	Routes["/api/actions/chat_messages"] = GetChatMessages
	h.Log("Добавлен маршрут /api/actions/chat_messages -> GetChatMessages")

	Routes["/api/actions/verify"] = HandleVerify
	h.Log("Добавлен маршрут /api/actions/verify -> Verify")

	Routes["/api/actions/exclude_chat_participants"] = HandleExcludeChatParticipants
	h.Log("Добавлен маршрут /api/actions/exclude_chat_participants -> HandleExcludeChatParticipants")

	Routes["/api/actions/include_chat_participants"] = HandleIncludeChatParticipants
	h.Log("Добавлен маршрут /api/actions/include_chat_participants -> HandleIncludeChatParticipants")

	Routes["/api/actions/create_chat"] = HandleCreateChat
	h.Log("Добавлен маршрут /api/actions/create_chat -> CreateChat")

	Routes["/api/actions/residue_request"] = HandleGetResidueRequest
	h.Log("Добавлен маршрут /api/actions/residue_request -> GetResidue")

	Routes["/api/actions/residue_invoice"] = HandleGetResidueInvoice
	h.Log("Добавлен маршрут /api/actions/residue_invoice -> GetResidueInvoice")

	Routes["/api/actions/user_divisions"] = GetUserDivisionsResponse
	h.Log("Добавлен маршрут /api/actions/user_divisions -> GetUserDivisions")

	Routes["/api/actions/enums"] = GetEnums
	h.Log("Добавлен маршрут /api/actions/enums -> GetResidue")
}
