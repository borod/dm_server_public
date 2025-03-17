package dm_entity

type DMQueryStruct struct {
	EntityName string      `json:"EntityName"`
	MdbQuery   interface{} `json:"MdbQuery"`
	PayLoad    interface{} `json:"PayLoad"`
}
