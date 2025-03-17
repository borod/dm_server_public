package api_entities

// import (
// 	h "dm_server/dm_helper"

// 	"encoding/json"
// 	"os"
// )

// type Request1 struct {
// }

// type VerifiedStatus1 struct {
// 	DivisionName string `json:"divisionName"`
// 	Status       int    `json:"status"`
// }

// type Percentages1 struct {
// 	Ordered   float64 `json:"Ordered"`
// 	Payed     float64 `json:"Payed"`
// 	Delivered float64 `json:"Delivered"`
// }

// type Resource1 struct {
// 	RequestID    int     `json:"RequestID"`
// 	ResourceId   int     `json:"ResourceId"`
// 	ResourceVid  int     `json:"ResourceVid"` // 0 - материал, 1 - персонал, 2 - техника, 3 - услуги
// 	ResourceName string  `json:"ResourceName"`
// 	Measure      string  `json:"Measure"`
// 	Qty          float64 `json:"Qty"`
// }

// const requestsFilePath = "/files/requests.json"

// func GetDummyRequests() []Request1 {
// 	result, err := GetRequestsFromFile(h.DMFilePath(requestsFilePath))
// 	if err != nil {

// 	}

// 	return result
// }

// func GetRequestsFromFile(path string) ([]Request1, error) {
// 	fileData, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var requests []Request1
// 	err = json.Unmarshal(fileData, &requests)
// 	if err != nil {
// 		h.Err(err.Error())
// 		return nil, err
// 	}

// 	for i := range requests {
// 		if requests[i].IsTzr == 0 {
// 			requests[i].IsTzr = nil
// 		}
// 	}

// 	return requests, nil
// }

// type Request struct {
// 	ID         int            `json:"ID"`
// 	ResourceID int            `json:"ResourceID"`
// 	ObjectID   int            `json:"ObjectID"`
// 	ARCID      int            `json:"ARCID"`
// 	WorkID     int            `json:"WorkID"`
// 	Fields     []RequestField `json:"Resources"`
// }

// type RequestField struct {
// 	ID          int     `json:"ID"`
// 	ResourceID  int     `json:"ResourceID"`
// 	NewName     string  `json:"NewName"`
// 	Qty         float64 `json:"Qty"`
// 	Analogue    bool    `json:"Analogue"`
// 	Explanation string  `json:"Explanation"`
// }
