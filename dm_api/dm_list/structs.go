package dm_list

import (
	mysql "dm_server/dm_db/dm_mysql"
)

const C_pageSize = 200 // Количество элементов на странице

type ChatItemResponse struct {
	ID          int64
	Name        string
	LastMessage mysql.Message
	AuthorName  string
	AuthorID    int64
}
