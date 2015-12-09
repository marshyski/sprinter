package main

import (
  "flag"
  "fmt"
  "os"
  "os/exec"
  "os/user"
  "strings"
  "io"
  "io/ioutil"
  "net/http"
  "net"
  "bytes"
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

var usage = `Usage: sprinter [options] <args>

   Sprinter your builds, tests and deploys

    -c, -config              Config
    -n, -name, optional      Override name in build config
    -v, -version, optional   Override version of your app in build config

Documentation:  https://github.com/marshyski/sprinter/blob/master/README.md

`

// HTTP client timeout
var timeout = time.Duration(300 * time.Millisecond)

func dialTimeout(network, addr string) (net.Conn, error) {
  return net.DialTimeout(network, addr, timeout)
}

func main() {

  start := time.Now()

  flag.Usage = func() {
    fmt.Println(usage)
  }

  flag.Parse()

  viper.SetConfigName("sprinter")
  if *configFlag != "" {
    viper.AddConfigPath(*configFlag)
  } else {
    viper.AddConfigPath("/opt/sprinter/conf")
  }

    viperReadConfig := viper.ReadInConfig()
    viper.SetDefault("elastic_host", "localhost")
    viper.SetDefault("elastic_port", "9200")
    viper.SetDefault("elastic_username", "undef")
    viper.SetDefault("elastic_password", "undef")
    viper.SetDefault("name", "undef")
    viper.SetDefault("build", "undef")
    viper.SetDefault("test", "undef")
    viper.SetDefault("version", "undef")
    viper.SetDefault("test", "undef")
    viper.SetDefault("zip", "undef")
    viper.SetDefault("upload", "undef")
    viper.SetDefault("upload_username", "undef")
    viper.SetDefault("upload_password", "undef")

    elastic_host := viper.GetString("elastic_host")
    elastic_port := viper.GetString("elastic_port")
    elastic_username := viper.GetString("elastic_username")
    elastic_password := viper.GetString("elastic_password")
    name := viper.GetString("name")
    build := viper.GetString("build")
    test := viper.GetString("test")
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
    buildScriptOut, buildERR := buildScript.Output()
    execBuildCommand := exec.Command(build)
    execBuildOut, buildERR := execBuildCommand.Output()
    execBuildString := string(execBuildOut)
    execBuildSlice := strings.Split(execBuildString,"\n")
    execBuildJson,_ := json.Marshal(execBuildSlice)


    filename := "sprinter-output.json"

    // Create JSON file
    f, err := os.Create(filename)
    if err != nil {
      fmt.Println(err.Error())
      return
    }
    n, err := io.WriteString(f, "{")
    if err != nil {
      fmt.Println(n, err)
      return
    }

    if buildERR != nil {
      build_status := `
      "build_status": "fail",
      "build_exit": "%s",`
      //buildStatus := fmt.Sprintf(build_status, "fail")
      //fmt.Printf(buildStatus)
//      build_exit := `
//      "build_exit": "%s",`
      //buildExit := fmt.Sprintf(build_exit, err)
      //fmt.Printf(buildExit)

      buildStatusLine := fmt.Sprintf(build_status, buildERR)
      writeBuildStatus, err := io.WriteString(f, buildStatusLine)
      if err != nil {
        fmt.Println(writeBuildStatus, err)
        return
      }

    } else {
      build_status := `
      "build_status": "success",`
      //buildStatus := fmt.Sprintf(build_status, "success")
      //fmt.Printf(build_status)
      buildStatusLine := fmt.Sprintf(build_status)
      writeBuildStatus, err := io.WriteString(f, buildStatusLine)
      if err != nil {
        fmt.Println(writeBuildStatus, err)
        return
      }

    }

    build_name := `
      "build_name": "%s",`

    if *nameFlag != "" {
      buildNameLine := fmt.Sprintf(build_name, *nameFlag)
      //fmt.Printf(buildNameLine)
      writeBuildName, err := io.WriteString(f, buildNameLine)
      if err != nil {
        fmt.Println(writeBuildName, err)
        return
      }
    } else {
      buildNameLine := fmt.Sprintf(build_name, name)
      //fmt.Printf(buildNameLine)
      writeBuildName, err := io.WriteString(f, buildNameLine)
      if err != nil {
        fmt.Println(writeBuildName, err)
        return
      }
    }

    if string(buildScriptOut) != "" {
      execBuild := `
      "build_output": %s,`

      execBuildLine := fmt.Sprintf(execBuild, string(execBuildJson))
      execBuildReplace := strings.Replace(execBuildLine, ",\"\"]", "]", -1)

      //fmt.Println(execBuildReplace)
      writeExecBuild, err := io.WriteString(f, execBuildReplace)
      if err != nil {
        fmt.Println(writeExecBuild, err)
        return
      }
    }


    if *verFlag != "" && version != "undef" {
      build_version, _ := fmt.Printf("      \"build_version\": \"%s\",", *verFlag)
      writeBuildVer, err := io.WriteString(f, build_version)
      if err != nil {
        fmt.Println(writeBuildVer, err)
        return
      }
    }

    if *verFlag == "" && version != "undef" {
      build_version, _ := fmt.Sprintf("      \"build_version\": \"%s\",\n", version)
      writeBuildVer, err := io.WriteString(f, build_version)
      if err != nil {
        fmt.Println(writeBuildVer, err)
        return
      }
    }

    timeN := time.Now()
    buildTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", timeN.Year(), timeN.Month(), timeN.Day(), timeN.Hour(), timeN.Minute(), timeN.Second())
    hostname, _ := os.Hostname()
    usr, _ := user.Current()

    fmt.Printf("      \"build_time\": \"%s\",\n", buildTime)
    fmt.Printf("      \"build_user\": \"%s\",", usr.Username)

    testScript := exec.Command("ls", test)
    testScriptOut, err := testScript.Output()
    execTestCommand := exec.Command(test)
    execTestOut, err := execTestCommand.Output()
    execTestString := string(execTestOut)
    execTestSlice := strings.Split(execTestString,"\n")
    execTestJson,_ := json.Marshal(execTestSlice)


    if err != nil && test != "undef" {
      test_status := `
      "test_status": "%s",`
      testStatus := fmt.Sprintf(test_status, "fail")
      fmt.Printf(testStatus)
      test_exit := `
      "test_exit": "%s",`
      testExit := fmt.Sprintf(test_exit, err)
      fmt.Printf(testExit)
    }

    if err == nil && test != "" {
      test_status := `
      "test_status": "%s",`
      testStatus := fmt.Sprintf(test_status, "success")
      fmt.Printf(testStatus)
    }

    if test != "" {
      if string(testScriptOut) != "" {
        execTest := `
      "test_output": %s,`

        execTestLine := fmt.Sprintf(execTest, string(execTestJson))
        execTestReplace := strings.Replace(execTestLine, ",\"\"]", "]", -1)

        fmt.Println(execTestReplace)
      }
    }

    fmt.Printf("      \"hostname\": \"%s\",\n", hostname)

    elapsed := time.Since(start)
    fmt.Printf("      \"execution_time\": \"%s\"\n    }\n", elapsed)

    elastic_url := "http://" + elastic_host + ":" + elastic_port + "/sprinter" + "/" + `%s`

    if *nameFlag != "" {
      elasticUrlLine := fmt.Sprintf(elastic_url, *nameFlag)
      fmt.Printf(elasticUrlLine)
    } else {
      elasticUrlLine := fmt.Sprintf(elastic_url, name)
      fmt.Printf(elasticUrlLine)
    }

    transport := http.Transport{
      Dial: dialTimeout,
    }

    client := http.Client{
      Transport: &transport,
    }

    // Check to see if ElasticSearch server is up
    elasticResponse, err := client.Get(elastic_url)
    if elasticResponse != nil {
      jsonStr,err := ioutil.ReadFile(filename)
      if err != nil {
        fmt.Println(err.Error())
        return
      }
      fmt.Println(buildTime, hostname, "INFO elasticsearch endpoint:", elastic_url)

      reqPost, err := http.NewRequest("POST", elastic_url, bytes.NewBuffer(jsonStr))
      if elastic_password != "undef" {
        reqPost.SetBasicAuth(elastic_username, elastic_password)
      }
      reqPost.Header.Set("Content-Type", "application/json")

      clientReq := &http.Client{}
      respPost, err := clientReq.Do(reqPost)
      if err != nil {
        fmt.Println(err.Error())
      }
      defer respPost.Body.Close()
      fmt.Println(buildTime, hostname, "POST json elasticsearch type status:", respPost.Status)
      postBody, _ := ioutil.ReadAll(respPost.Body)
      fmt.Println(buildTime, hostname, "POST response body:", string(postBody))
      } else {
        fmt.Println(buildTime, hostname, "FAIL unable to connect to elasticeearch server:", "http://" + elastic_host + ":" + elastic_port)
      }

}
