package wincmd

import (
  "os"
  "bufio"
  "strings"
  "fmt"
  "io/ioutil"

  "github.com/masterzen/winrm"
)

func Win(command string, port int, user string, hostsfile string, https bool, insecure bool, password string, cacert string, host string) {

  var certBytes []byte
  var err error
  if cacert != "" {
    certBytes, err = ioutil.ReadFile(cacert)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  } else {
    certBytes = nil
  }

  argsStr := strings.Split(command, ",")
  argsList := make([]string, 0)

  for _, item := range argsStr {
    argsList = append(argsList, item)
  }

  if host != "" {

    for _, item := range argsList {

      if err != nil {
        fmt.Println(item)
        fmt.Println(err)
      } else {

      client, err := winrm.NewClient(&winrm.Endpoint{Host: host, Port: port, HTTPS: https, Insecure: insecure, CACert: certBytes}, user, password)
      if err != nil {
        fmt.Println(err)
      }

      fmt.Println()
      fmt.Println(host, item)

      run, err := client.RunWithInput(item, os.Stdout, os.Stderr, os.Stdin)
      if err != nil {
        fmt.Println(err)
      }

      fmt.Println(run)

      }
    }
  } else {
    file, err := os.Open(hostsfile)
    defer file.Close()
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {

      for _, item := range argsList {

        if err != nil {
          fmt.Println(item)
          fmt.Println(err)
        } else {

            client, err := winrm.NewClient(&winrm.Endpoint{Host: scanner.Text(), Port: port, HTTPS: https, Insecure: insecure, CACert: certBytes}, user, password)
            if err != nil {
              fmt.Println(err)
            }

            fmt.Println()
            fmt.Println(scanner.Text(), item)

            run, err := client.RunWithInput(item, os.Stdout, os.Stderr, os.Stdin)
            if err != nil {
              fmt.Println(err)
            }

            fmt.Println(run)

          }
        }
      }

      if err = scanner.Err(); err != nil {
        fmt.Println(err)
      }
  }
}
