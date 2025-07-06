package config

import (
	"log"
	"os"

	"github.com/File-share/constants"
	"github.com/File-share/models"
	"gopkg.in/yaml.v3"
)

var configs models.Config

func Init() {
	filedata, err := os.ReadFile(constants.ConfigFile)

	if err != nil {
		log.Fatal("Error in reading config file", err)
	}

	err = yaml.Unmarshal(filedata, &configs)

	if err != nil {
		log.Fatal("Error getting configs")
	}

}

func GetConfig() models.Config {
	return configs
}
