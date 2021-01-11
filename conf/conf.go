package conf

import (
	"fmt"
	"github.com/tkanos/gonfig"
	"log"
	"os"
	"path/filepath"
)

type configuration struct {
	MongoDBName string
	MongoDBPort int32
	MongoDBHost string
}

const ENV  = "dev"

var configs *configuration

// configsFromJson private function that is used to fetch/read
// values from a json file
func configsFromJson() {
	confFile, confFileErr := filepath.Abs(fmt.Sprintf("conf/%s_config.json", ENV))
	if confFileErr != nil {
		log.Printf("ERROR | The following error occurred while fetching/reading json conf file: %v\n", confFileErr)
		os.Exit(1)
	}

	err := gonfig.GetConf(confFile, &configs)
	if err != nil {
		log.Printf("ERROR | Reading app's json conf file failed with message: %v\n", err)
		os.Exit(1)
	}
}

// GetAppConfigs public function that exposes the configs variable
func GetAppConfigs() *configuration {
	return configs
}
