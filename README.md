# sprinter

Remote executes SSH / WinRM commands


**Usage:**

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


**Examples:**

Long way, with defining Hostsfile (SSH):

    sprinter -k ~/.ssh/privatekey -f ~/Hostsfile -nc 'df -h','uname -a','ls /tmp' -u root -p 2222

	sprinter -k ~/.ssh/privatekey -h domain1.com -nc 'df -h','uname -a','ls /tmp' -u root -p 2222

Short way, without defining Hostsfile (SSH):

    sprinter -k ~/.ssh/privatekey -nc 'df -h','uname -a','ls /tmp'
    sprinter -k ~/.ssh/privatekey -h domain1.com -nc 'df -h','uname -a','ls /tmp'

Long way, with HTTPS (Windows):

	sprinter -u vagrant:vagrant -p 5985 -ca ./cert -h true -i true -h domain2.com -wc 'ipconfig /all','set'

Short way, without HTTPS (Windows):

	sprinter -u vagrant:vagrant -p 5985 -h domain2.com -wc 'ipconfig /all','set'


**Hostsfile:**
{current working directory}/Hostsfile - if '-f' not specified in CLI arguments

    domain1.com
    104.236.20.120
    domain2.com

**Example Output:**

    domain1.com df -h
    Filesystem      Size  Used Avail Use% Mounted on
    /dev/vda1        30G  5.7G   23G  21% /
    devtmpfs        491M     0  491M   0% /dev
    tmpfs           498M     0  498M   0% /dev/shm
    tmpfs           498M  240K  498M   1% /run
    tmpfs           498M     0  498M   0% /sys/fs/cgroup
    tmpfs           498M  4.0K  498M   1% /tmp

    domain1.com uname -a
    Linux bundles 3.13.9-200.fc20.x86_64 #1 SMP Fri Apr 4 12:13:05 UTC 2014 x86_64 x86_64 x86_64 GNU/Linux

    10.20.1.7 ipconfig /all
    
    Windows IP Configuration
    
       Host Name . . . . . . . . . . . . : server2008r2a
       Primary Dns Suffix  . . . . . . . : pdx.puppetlabs.demo
       Node Type . . . . . . . . . . . . : Hybrid
       IP Routing Enabled. . . . . . . . : No
       WINS Proxy Enabled. . . . . . . . : No
       DNS Suffix Search List. . . . . . : pdx.puppetlabs.demo
                                           home
    0
    
    10.20.1.7 set
    ALLUSERSPROFILE=C:\ProgramData
    APPDATA=C:\Users\vagrant\AppData\Roaming
    ChocolateyInstall=C:\ProgramData\chocolatey
    CommonProgramFiles=C:\Program Files\Common Files
    CommonProgramFiles(x86)=C:\Program Files (x86)\Common Files
    CommonProgramW6432=C:\Program Files\Common Files
    COMPUTERNAME=SERVER2008R2A
    ComSpec=C:\Windows\system32\cmd.exe
    FP_NO_HOST_CHECK=NO
    LOCALAPPDATA=C:\Users\vagrant\AppData\Local
    NUMBER_OF_PROCESSORS=1
    OS=Windows_NT
    Path=C:\Windows\system32;C:\Windows;C:\Windows\System32\Wbem;C:\Windows\System32\WindowsPowerShell\v1.0\;C:\Program Files\Puppet Labs\Puppet\bin;C:\ProgramData\chocolatey\bin;
    PATHEXT=.COM;.EXE;.BAT;.CMD;.VBS;.VBE;.JS;.JSE;.WSF;.WSH;.MSC
    PROCESSOR_ARCHITECTURE=AMD64
    PROCESSOR_IDENTIFIER=Intel64 Family 6 Model 69 Stepping 1, GenuineIntel
    PROCESSOR_LEVEL=6
    PROCESSOR_REVISION=4501
    ProgramData=C:\ProgramData
    ProgramFiles=C:\Program Files
    ProgramFiles(x86)=C:\Program Files (x86)
    ProgramW6432=C:\Program Files
    PROMPT=$P$G
    PSModulePath=C:\Windows\system32\WindowsPowerShell\v1.0\Modules\
    PUBLIC=C:\Users\Public
    SystemDrive=C:
    SystemRoot=C:\Windows
    TEMP=C:\Users\vagrant\AppData\Local\Temp
    TMP=C:\Users\vagrant\AppData\Local\Temp
    USERDOMAIN=SERVER2008R2A
    USERNAME=vagrant
    USERPROFILE=C:\Users\vagrant
    windir=C:\Windows
    0

**Build:**

    #Set GOPATH
    go get golang.org/x/crypto/ssh
    go get github.com/masterzen/winrm/winrm
    cd $GOPATH/src
    git clone https://github.com/marshyski/sprinter.git
    cd sprinter/sprinter && go build
