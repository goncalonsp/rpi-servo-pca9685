#!/usr/bin/expect

# For the expect interpreter do
# sudo apt-get install expect

set executable [lindex $argv 0];
set remote_executable [lindex $argv 1];
set user [lindex $argv 2];
set host [lindex $argv 3];
set pass [lindex $argv 4];

spawn scp "$executable" "$user@$host:$remote_executable"
expect {
	password: {send "$pass\r"; exp_continue}
}