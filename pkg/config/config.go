package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dir      string `json:"dir"`
}

// Reads json file and maps its values in struct
func (c *config) SetConfigFromJson(jsonPath string) (err error) {
	path := filepath.Join(jsonPath)
	jsonFile, err := os.Open(path)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &c)
	return
}

// NewConfig ...
func NewConfig() *config {
	return &config{}
}
