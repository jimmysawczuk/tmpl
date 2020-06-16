package tmplfunc

import (
	"encoding/json"
	"io/ioutil"
)

// GetJSON reads data from the provided (local) path, and attempts to unmarshal it.
func GetJSON(path string) (interface{}, error) {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var target interface{}
	err = json.Unmarshal(by, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

// JSONify marshals the provided interface into its JSON representation.
func JSONify(v interface{}) (string, error) {
	by, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(by), nil
}
