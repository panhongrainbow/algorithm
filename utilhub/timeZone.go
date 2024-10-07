package utilhub

import (
	"io/ioutil"
	"path/filepath"
)

// ListTimezones returns all available timezones in the zoneinfo directory.
func ListTimezones(timezoneDir string) ([]string, error) {
	var timezones []string

	files, err := ioutil.ReadDir(timezoneDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			zones, err := filepath.Glob(timezoneDir + file.Name() + "/*")
			if err != nil {
				return nil, err
			}
			for _, zone := range zones {
				timezones = append(timezones, file.Name()+"/"+filepath.Base(zone))
			}
		} else {
			timezones = append(timezones, file.Name())
		}
	}
	return timezones, nil
}
