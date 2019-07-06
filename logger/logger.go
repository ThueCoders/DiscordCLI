package logger

import (
    "flag"
    "os"
    "fmt"
    "io/ioutil"
    "log"
    "encoding/json"
    "go/build"
)

var (
    Log      *log.Logger
)


func init() {
    configFile, err := os.Open("config.json")
    if (err != nil) {
        fmt.Println(err)
    }
    byteValue, _ := ioutil.ReadAll(configFile)
    configFile.Close()
    var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)
    path := result["log-path"].(string)

    // set location of log file
    var logpath = build.Default.GOPATH + path

    flag.Parse()
    var file, err1 = os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)

    if err1 != nil {
        panic(err1)
    }
    Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
    Log.Println("LogFile : " + logpath)
}
