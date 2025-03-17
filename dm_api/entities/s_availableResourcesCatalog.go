package api_entities

type Work struct {
	ID        int        `json:"ID"`
	Name      string     `json:"Name"`
	Resources []Resource `json:"Resources"`
}

type Resource struct {
	ID         int        `json:"ID"`
	Type       string     `json:"Type"`
	Name       string     `json:"Name"`
	Measure    string     `json:"Measure"`
	Qty        float64    `json:"Qty"`
	Expediture float64    `json:"Expediture"`
	Note       string     `json:"Note"`
	Analog     bool       `json:"Analog"`
	Resources  []Resource `json:"Resources"`
}

type AvailableResourcesCatalog struct {
	ID    int    `json:"ID"`
	Name  string `json:"Name"`
	Works []Work `json:"Works"`
}

type ResourceList struct {
	Items []Resource `json:"Items"`
}
