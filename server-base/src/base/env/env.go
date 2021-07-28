package env

import (
	"encoding/json"
	"glog"
	"io/ioutil"
)

var configData map[string]map[string]string

func Load(path string) bool {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Error("[Config] Load error.", path, ", ", err )
		return false
	}
	err = json.Unmarshal(file, &configData)
	if err != nil {
		glog.Error("[Config] Parse failed. ", path, ", ", err)
		return false
	}
	return true
}


func Get(table, key string) string {
	t, ok := configData[table]
	if !ok {
		return ""
	}
	val, ok := t[key]
	if !ok {
		return ""
	}
	return val
}