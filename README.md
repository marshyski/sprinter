# sprinter

Executes SSH commands on systems by reading a file line by line


**Usage:**

    -c, -command             Run command or commands: 'df-h','uname -a'
    -h, -hosts, optional     Hosts file location, default is ./Hostsfile
    -k, -key                 PEM key file location: /Users/.ssh/key.pem
    -u, -user, optional      Username to SSH as, default is root
    -p, -port, optional      Port to SSH as, default is 22


**Examples:**

Long way:

    sprinter -k /Users/timski/.ssh/secretkey.pem -h ~/Hostsfile -c 'df -h','uname -a','ls /tmp' -u root -p 2222

Short way:

    sprinter -k /Users/timski/.ssh/secretkey.pem -c 'df -h','uname -a','ls /tmp'


**Hostsfile:**
{current working directory}/Hostsfile - if '-h' not specified in CLI arguments

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

    104.236.20.120 df -h
    Filesystem      Size  Used Avail Use% Mounted on
    udev            486M     0  486M   0% /dev
    tmpfs           100M  580K   99M   1% /run
    /dev/vda1        30G  5.4G   23G  20% /
    tmpfs           497M  120K  497M   1% /dev/shm
    tmpfs           5.0M     0  5.0M   0% /run/lock
    tmpfs           497M     0  497M   0% /sys/fs/cgroup
    tmpfs           100M     0  100M   0% /run/user/0

    104.236.20.120 uname -a
    Linux trindzy 3.19.0-22-generic #22-Ubuntu SMP Tue Jun 16 17:15:15 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux

    domain2.com df -h
    Filesystem      Size  Used Avail Use% Mounted on
    /dev/vda1        40G  4.6G   33G  13% /
    devtmpfs        996M     0  996M   0% /dev
    tmpfs          1002M     0 1002M   0% /dev/shm
    tmpfs          1002M   17M  985M   2% /run
    tmpfs          1002M     0 1002M   0% /sys/fs/cgroup

    domain2.com uname -a
    Linux jenkins 3.10.0-123.8.1.el7.x86_64 #1 SMP Mon Sep 22 19:06:58 UTC 2014 x86_64 x86_64 x86_64 GNU/Linux

**Build:**

    #Set GOPATH
    go get golang.org/x/crypto/ssh
    cd $GOPATH/src
    git clone https://github.com/marshyski/sprinter.git
    cd sprinter/sprinter && go build
