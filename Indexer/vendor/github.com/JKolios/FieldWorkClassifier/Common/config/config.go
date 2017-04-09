package config

import (
	"encoding/json"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"io/ioutil"
	"log"
	"github.com/davecgh/go-spew/spew"
)

type Settings struct {
	ApiURL          string   `json:"apiURL"`
	ElasticURL      string   `json:"elasticURL"`
	ElasticUsername string   `json:"ElasticUsername"`
	ElasticPassword string   `json:"ElasticPassword"`
	SniffCluster    bool     `json:"sniffCluster, omitempty"`
	GinDebug        bool     `json:"ginDebug, omitempty"`
}

//GetConfFromJSONFile reads application configuration from *filename* and maps it to a Settings struct
func GetConfFromJSONFile(filename string) *Settings {

	confContent, err := ioutil.ReadFile(filename)
	utils.CheckFatalError(err)
	config := &Settings{}
	err = json.Unmarshal(confContent, config)
	utils.CheckFatalError(err)
	log.Println("Configuration loaded:")
	spew.Config.Indent = "\t"
	spew.Dump(*config)
	return config
}
