package api_entities

type Measure struct {
	ID           int    `json:"ID"`
	Code         int    `json:"Code"`
	FullNameRus  string `json:"FullNameRus"`
	ShortNameRus string `json:"ShortNameRus"`
	ShortNameInt string `json:"ShortNameInt"`
	CodeNameInt  string `json:"codeNameInt"`
}
