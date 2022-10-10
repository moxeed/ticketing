package common

import (
	"encoding/json"
	"os"
)

type Config struct {
	DataBaseDsn string
	ListenPort  int
	Title       string
	Schemes     []string
	Host        string
	BasePath    string
}

var Configuration Config

func init() {
	configFile, err := os.Open("./config.json")

	if err != nil {
		panic(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Configuration)
}
