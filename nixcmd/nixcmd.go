package nixcmd

import (
  "os"
  "bufio"
  "strings"
  "fmt"

  "golang.org/x/crypto/ssh"
  "sprinter/sshlib"
)

func Nix(command string, port int, user string, hostsfile string, key string, password string, host string) {

  var err error

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

        sshConfig := &ssh.ClientConfig{
          User: user,
          Auth: []ssh.AuthMethod{
            sshlib.PublicKeyFile(key),
          },
        }

        client := &sshlib.SSHClient{
          Config: sshConfig,
          Host:   host,
          Port:   port,
        }

        cmd := &sshlib.SSHCommand{
          Path:   item,
          Env:    []string{},
          Stdin:  os.Stdin,
          Stdout: os.Stdout,
          Stderr: os.Stderr,
        }

        fmt.Println()
        fmt.Println(host, cmd.Path)
        if err := client.RunCommand(cmd); err != nil {
          fmt.Println(os.Stderr, "command run error: %s\n", err)
          os.Exit(1)
        }
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

          sshConfig := &ssh.ClientConfig{
            User: user,
            Auth: []ssh.AuthMethod{
              sshlib.PublicKeyFile(key),
            },
          }

          client := &sshlib.SSHClient{
            Config: sshConfig,
            Host:   scanner.Text(),
            Port:   port,
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
  }
}
