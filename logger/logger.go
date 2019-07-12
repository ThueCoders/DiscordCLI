package logger

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(configFile)
	configFile.Close()
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	path := result["log-path"].(string)

	var logpath = build.Default.GOPATH + path

	flag.Parse()
	var file, err1 = os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY, 0664)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}
