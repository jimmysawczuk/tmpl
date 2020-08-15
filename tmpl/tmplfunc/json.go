package tmplfunc

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetJSON reads data from the provided (local) path, and attempts to unmarshal it.
func GetJSON(path string) (interface{}, error) {
	abspath, _ := filepath.Abs(path)

	by, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "read file (path: %s, abs path: %s)", path, abspath)
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
