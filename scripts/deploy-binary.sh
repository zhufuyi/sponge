#!/usr/bin/expect

set serviceName "serverNameExample"

# parameters
set username [lindex $argv 0]
set password [lindex $argv 1]
set hostname [lindex $argv 2]

set timeout 30

spawn scp -r ./${serviceName}-binary.tar.gz ${username}@${hostname}:/tmp/
#expect "*yes/no*"
#send  "yes\r"
expect "*password:*"
send  "${password}\r"
expect eof

spawn ssh ${username}@${hostname}
#expect "*yes/no*"
#send  "yes\r"
expect "*password:*"
send  "${password}\r"

# execute a command or script
expect "*${username}@*"
send "cd /tmp && tar zxvf ${serviceName}-binary.tar.gz\r"
expect "*${username}@*"
send "bash /tmp/${serviceName}-binary/deploy.sh\r"

# logging out of a session
expect "*${username}@*"
send "exit\r"

expect eof
