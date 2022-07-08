package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// build file path
func buildFilePath(dir string, key string) string {
	return fmt.Sprintf("%s/%s", dir, key)
}

// creator file if not exists
func creatDir(dir string) (path string, err error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}
	return
}

// delete file if exists
func DeleteFile(dir string, key string) (err error) {
	path := buildFilePath(dir, key)
	if _, err := os.Stat(path); !os.IsExist(err) {
		err = os.Remove(path)
	}
	return
}

// cover data to file
func CoverFile(dir string, key string, data interface{}) error {
	path := buildFilePath(dir, key)
	creatDir(dir)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, jsonData, 0644)
}

// get file data
func GetFileData(dir string, key string) (v interface{}, err error) {
	path := buildFilePath(dir, key)
	data, err := ioutil.ReadFile(path)
	err = json.Unmarshal(data, &v)
	return
}

// check file exists
func ExistFile(dir string, key string) (isExists bool, err error) {
	path := buildFilePath(dir, key)
	_, err = os.Stat(path)
	return os.IsExist(err), err
}
