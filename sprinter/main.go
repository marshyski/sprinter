package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "bufio"
  "path/filepath"

  "golang.org/x/crypto/ssh"
  "sprinter/sshlib"
)

var commandFlag = flag.String("command", "", "  Command")
var hostsFlag = flag.String("hosts", "", " Hosts")
var keyFlag = flag.String("key", "", " Key")
var userFlag = flag.String("user", "root", " User")
var portFlag = flag.Int("port", 22, " Port")

func init() {
  Hostsfile, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/Hostsfile")
  if err != nil {
    fmt.Println(err)
  }
  flag.StringVar(commandFlag, "c", "", "  Command")
  flag.StringVar(hostsFlag, "h", Hostsfile, "  Hosts")
  flag.StringVar(keyFlag, "k", "", "  Key")
  flag.StringVar(userFlag, "u", "root", "  User")
  flag.IntVar(portFlag, "p", 22, "  Port")
}

var usage = `Usage: sprinter [options] <args>

Sprinter executes commands on systems by reading a file

-c, -command             Run command or commands: 'df-h','uname -a'
-h, -hosts, optional     Hosts file location, default is ./Hostsfile
-k, -key                 PEM key file location: ~/.ssh/key.pem
-u, -user, optional      Username to SSH as, default is root
-p, -port, optional      Port to SSH as, default is 22

Documentation:  https://github.com/marshyski/sprinter/blob/master/README.md

`

func main() {

  flag.Usage = func() {
    fmt.Println(usage)
  }

  flag.Parse()

  if *commandFlag == "" {
    fmt.Println(usage)
    os.Exit(1)
  }

  if *keyFlag == "" {
    fmt.Println(usage)
    os.Exit(1)
  }

  argsStr := strings.Split(*commandFlag, ",")
  argsList := make([]string, 0)

  for _, item := range argsStr {
    argsList = append(argsList, item)
  }

  if file, err := os.Open(*hostsFlag); err == nil {

    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {

      for _, item := range argsList {

        if err != nil {
          fmt.Println(item)
          fmt.Println(err)
          } else {


            sshConfig := &ssh.ClientConfig{
              User: *userFlag,
              Auth: []ssh.AuthMethod{
                sshlib.PublicKeyFile(*keyFlag),
              },
            }

            client := &sshlib.SSHClient{
              Config: sshConfig,
              Host:   scanner.Text(),
              Port:   *portFlag,
            }

            cmd := &sshlib.SSHCommand{
              Path:   item,
              Env:    []string{},
              Stdin:  os.Stdin,
              Stdout: os.Stdout,
              Stderr: os.Stderr,
            }

            fmt.Println()
            fmt.Println(scanner.Text(), cmd.Path)
            if err := client.RunCommand(cmd); err != nil {
              fmt.Println(os.Stderr, "command run error: %s\n", err)
              os.Exit(1)
            }
          }
        }
      }

      if err = scanner.Err(); err != nil {
        fmt.Println(err)
      }

      } else {
        fmt.Println(err)
      }
}
