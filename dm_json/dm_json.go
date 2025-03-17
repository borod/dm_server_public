package dm_json

import "encoding/json"

func UpdateStructWithJSON(existingStruct interface{}, jsonData []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	existingData, err := json.Marshal(existingStruct)
	if err != nil {
		return err
	}

	err = json.Unmarshal(existingData, &data)
	if err != nil {
		return err
	}

	newData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(newData, existingStruct)
	if err != nil {
		return err
	}

	return nil
}
