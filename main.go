package main

import (
  "flag"
  "fmt"
  "os"
  "os/exec"
  "strings"
  "encoding/json"
  "time"
  "github.com/spf13/viper"
)

func checkerror(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

var verFlag = flag.String("version", "", "  Version of your app")
var nameFlag = flag.String("name", "", "  Name of build you're running")
var configFlag = flag.String("config", "", "  Build config")

func init() {
  flag.StringVar(verFlag, "v", "", "  Version of your app")
  flag.StringVar(nameFlag, "n", "", "  Name of build you're running")
  flag.StringVar(configFlag, "c", "", "  Build config")
}

var usage = `Usage: spinkick [options] <args>

   Spinkick your builds and tests

    -c, -config              Build config
    -n, -name, optional      Override name in build config
    -v, -version, optional   Override version of your app in build config

Documentation:  https://github.com/marshyski/spinkick/blob/master/README.md

`

func main() {

  flag.Usage = func() {
    fmt.Println(usage)
  }

  flag.Parse()

  viper.SetConfigName("spinkick")
  if *configFlag != "" {
    viper.AddConfigPath(*configFlag)
  } else {
    viper.AddConfigPath("/opt/spinkick/conf")
  }

    viperReadConfig := viper.ReadInConfig()
    viper.SetDefault("elastic_host", "localhost")
    viper.SetDefault("elastic_port", "9200")
    viper.SetDefault("elastic_username", "undef")
    viper.SetDefault("elastic_password", "undef")
    viper.SetDefault("name", "undef")
    viper.SetDefault("build", "undef")
    viper.SetDefault("version", "undef")
    viper.SetDefault("test", "undef")
    viper.SetDefault("zip", "undef")
    viper.SetDefault("upload", "undef")
    viper.SetDefault("upload_username", "undef")
    viper.SetDefault("upload_password", "undef")

  //  elastic_host := viper.GetString("elastic_host")
  //  elastic_port := viper.GetString("elastic_port")
  //  elastic_username := viper.GetString("elastic_username")
  //  elastic_password := viper.GetString("elastic_password")
    name := viper.GetString("name")
    build := viper.GetString("build")
    version := viper.GetString("version")
  //  test := viper.GetString("test")
  //  zip := viper.GetString("zip")
  //  upload := viper.GetString("upload")
  //  upload_username := viper.GetString("upload_username")
  //  upload_password := viper.GetString("upload_password")

  if viperReadConfig != nil {
    fmt.Println("INFO no config file used, using default configuration")
  }

  fmt.Println(build, name, version)

  if build == "undef" {
    fmt.Println(usage)
    os.Exit(1)
  }

  if name == "undef" {
    fmt.Println(usage)
    os.Exit(1)
  }

    buildScript := exec.Command("ls", build)
    buildScriptOut, err := buildScript.Output()
    execBuildCommand := exec.Command(build)
    execBuildOut, err := execBuildCommand.Output()
    execBuildString := string(execBuildOut)
    execBuildSlice := strings.Split(execBuildString,"\n")
    execBuildJson,_ := json.Marshal(execBuildSlice)

    fmt.Printf("    {")

    if err != nil {
      build_status := `
      "build_status": "%s",`
      buildStatus := fmt.Sprintf(build_status, "fail")
      fmt.Printf(buildStatus)
    } else {
      build_status := `
      "build_status": "%s",`
      buildStatus := fmt.Sprintf(build_status, "success")
      fmt.Printf(buildStatus)
    }

    build_name := `
      "build_name": "%s",`

    buildNameLine := fmt.Sprintf(build_name, *nameFlag)
    fmt.Printf(buildNameLine)

    if string(buildScriptOut) != "" {
      execBuild := `
      "build_output": %s,`

      execBuildLine := fmt.Sprintf(execBuild, string(execBuildJson))
      execBuildReplace := strings.Replace(execBuildLine, ",\"\"]", "]", -1)

      fmt.Println(execBuildReplace)
    }

    if *verFlag != "" && version != "undef" {
      fmt.Printf("      \"build_version\": \"%s\",\n", *verFlag)
    }

    timeN := time.Now()
    buildTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", timeN.Year(), timeN.Month(), timeN.Day(), timeN.Hour(), timeN.Minute(), timeN.Second())
    hostname, _ := os.Hostname()

    fmt.Printf("      \"build_time\": \"%s\",\n", buildTime)
    fmt.Printf("      \"hostname\": \"%s\"\n    }\n", hostname)
}
