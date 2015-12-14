package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "path/filepath"
  "sprinter/nixcmd"
  "sprinter/wincmd"
)

var nixcmdFlag = flag.String("nixcmd", "", "  NIXcommand")
var wincmdFlag = flag.String("wincmd", "", "  WINDOWScommand")
var fileFlag = flag.String("file", "", " File")
var keyFlag = flag.String("key", "", " Key")
var userFlag = flag.String("user", "root", " User")
var portFlag = flag.Int("port", 22, " Port")
var hostFlag = flag.String("host", "", " Host")
var httpsFlag = flag.Bool("https", false, " HTTPS")
var insecureFlag = flag.Bool("insecure", false, " Insecure")
var cacertFlag = flag.String("cacert", "", "CACert")

func init() {
  Hostsfile, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/Hostsfile")
  if err != nil {
    fmt.Println(err)
  }
  flag.StringVar(nixcmdFlag, "nc", "", "  NIXcommand")
  flag.StringVar(wincmdFlag, "wc", "", "  WINDOWScommand")
  flag.StringVar(fileFlag, "f", Hostsfile, "  File")
  flag.StringVar(keyFlag, "k", "", "  Key")
  flag.StringVar(userFlag, "u", "root", "  User")
  flag.IntVar(portFlag, "p", 22, "  Port")
  flag.StringVar(hostFlag, "h", "", " Host")
  flag.BoolVar(httpsFlag, "ssl", false, "HTTPS")
  flag.BoolVar(insecureFlag, "i", false, "Insecure")
  flag.StringVar(cacertFlag, "ca", "", "CACert")
}

var usage = `Usage: sprinter [options] <args>

Sprinter remote executes SSH / WinRM commands

-nc, -nixcmd                Run NIX command or commands: 'df-h','uname -a'
-wc, -wincmd                Run Windows command or commands: 'ipconfig /all','set'
-f, -file, optional         Hosts file location, default is ./Hostsfile
-h, -host, optional         Run commands on one host
-k, -key                    Private key file location: ~/.ssh/key.pem
-u, -user, optional         Username and/or password to run as: Administrator:secret
                            default is root
-p, -port, optional         Port to SSH/WinRM as: 5985, default is 22
-h, -https, optional        Use HTTPS for WinRM, default is false
-i, -insecure, optional     Use SSL validation, default is false
-ca, -cacert, optional      Use CA Certificate, default is None

Documentation:  https://github.com/marshyski/sprinter/blob/master/README.md

`

func main() {

  flag.Usage = func() {
    fmt.Println(usage)
  }

  flag.Parse()

  if *nixcmdFlag == "" && *wincmdFlag == "" {
    fmt.Println(usage)
    os.Exit(1)
  }

  if *fileFlag == "" && *hostFlag == "" {
    fmt.Println(usage)
    os.Exit(1)
  }

  username := strings.Split(*userFlag, ":")[0]
  password := ``
  if strings.Contains(*userFlag, ":") {
    password += strings.Split(*userFlag, ":")[0]
  }

  if *nixcmdFlag != "" {
    nixcmd.Nix(*nixcmdFlag, *portFlag, username, *fileFlag, *keyFlag, password, *hostFlag)
  }

  if *wincmdFlag != "" {
    wincmd.Win(*wincmdFlag, *portFlag, username, *fileFlag, *httpsFlag, *insecureFlag, password, *cacertFlag, *hostFlag)
  }
}
