#!/bin/bash
#title           :poe.sh
#description     :This script interacts with a Unifi Network Switch. 
#acknowledgment  :Original script from Marcelo Pasin, this is a fork.
#author			 :Quentin Vaucher
#usage		     :bash poe.sh [-j|-c|-u|-d] [-p port]
#notes           :A .env files which defines the ssh address and credentials 
#                 needs to be on the root directory.
#================================================================================

source .env

poe-describe()
{
expect << END
log_user 0
spawn ssh $POEUSER@$POESWITCH
expect "password: "
send "$POEPASSWD\r"
expect "Welcome"
expect "# "
send "telnet localhost\r"
expect "(UBNT) >"
send -- "enable\r"
expect "(UBNT) #"
log_user 1
send -- "show interfaces description $1\r"
send -- " "
expect "(UBNT) #"
log_user 0
send -- "exit\r"
END
}

poe-set()
{
expect << END
log_user 0
spawn ssh $POEUSER@$POESWITCH
expect "password: "
send "$POEPASSWD\r"
expect "Welcome"
expect "# "
send "telnet localhost\r"
expect "(UBNT) >"
send -- "enable\r"
expect "(UBNT) #"
send -- "config\r"
expect "(UBNT) "
send -- "interface $1\r"
expect "(UBNT) "
send -- "poe opmode $2 \r"
expect "(UBNT) "
send -- "exit\r"
send -- "exit\r"
send -- "exit\r"
send -- "exit\r"
END
}

poe-filter()
{
	poe-describe "$1" | tail -n +2 | grep 0/ | cut -d / -f 2 | tr -s ' ' | cut -d ' ' -f -1,3- | tr -d '\r'
}

poe-csv()
{
	# $3$4 is in case the hostname is written in two words. Can we do better ?
	awk '{ print $1","$2","$3" "$4 }'
}

usage()
{
	echo "usage:"
	echo "  $0 [-j|-c|-u|-d] [-p port]"
	echo "     -c        List status in CSV format (default)"
	echo "     -j        List status in JSON format"
	echo "     -u        Turns on (up) port (needs -p)"
	echo "     -d        Turns off (down) port (needs -p)"
	echo "     -p port   Selects port to operate"
	exit 0
}

error()
{
	echo "$1"
	exit 1
}

if test $# -eq 0; then
	usage
fi

OPC=csv
PORT=""
while test $# -gt 0; do
	case "$1" in
	-c)
		OPC=csv
		;;
	-u)
		OPC=up
		;;
	-d)
		OPC=down
		;;
	-p)
		PORT="$2"
		shift
		;;
	*)
		usage
		;;
	esac
	shift
done

if test "a$PORT" = a; then
	if test "a$OPC" = aup -o "a$OPC" = adown; then
		error "options -s and -r must use -p <port>"
	else
		PORT=
	fi
else
	PORT="0/$PORT"
fi

case "$OPC" in
csv)
	poe-filter "$PORT" # | poe-csv
	;;
up)
	poe-set "$PORT" auto
	;;
down)
	poe-set "$PORT" shutdown
	;;
esac
