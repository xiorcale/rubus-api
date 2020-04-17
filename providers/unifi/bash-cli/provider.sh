#!/bin/bash

# this script takes ideas from: https://github.com/firizki/wailea

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


####################
# REST API
####################


rm -f response
mkfifo response

METHOD=""
ACTION=""
DEVICE_ID=""
RESPONSE=""

router() 
{
    while test $# -gt 0; do
        case "$1" in
        device) # /device
            if test "$METHOD" == "GET"; then 
                ACTION="poe-filter | poe-csv"; 
            fi
            shift
            ;;
        ''|*[0-9]*) # /device/:id
            DEVICE_ID=$1
            if test "$METHOD" == "GET"; then 
                ACTION="poe-filter 0/${1} | poe-csv"
            fi
            shift
            ;;
        on) # /device/:id/on
            ACTION=""
            if test "$METHOD" == "POST"; then 
                ACTION="poe-set 0/${DEVICE_ID} auto; echo this,is,OK"
            fi
            break
            ;;
        off) # /device/:id/off
            ACTION=""
            if test "$METHOD" == "POST"; then
                ACTION="poe-set 0/${DEVICE_ID} shutdown; echo this,is,OK"
            fi
            break
            ;;
        *)
            break
            ;;
        esac
    done

    RESPONSE=$(eval $ACTION)
}


to-json()
{
    json=""
    N="$#"

    # json array if we have more than one object to process
    if [[ "$N" -gt "1" ]]; then
        json="[\n"
    fi

    # format each line as a JSON object
    while test $# -gt 0; do
        IFS=',' read -ra line_arr <<< "$1"
        
        # filter broken ports
        id=${line_arr[0]}
        if [[ "$id" == "Broken" ]]; then
            shift
            continue
        fi

        # filter devices without any hostname
        hostname=${line_arr[2]}
        if [[ -z "$hostname" ]]; then
            shift
            continue
        fi

        if [[ "${line_arr[1]}" == "Up" ]]; then
            isTurnOn="true"
        else
            isTurnOn="false"
        fi

        obj="{ \"id\": $id, \"isTurnOn\": $isTurnOn, \"hostname\": \"$hostname\" },\n"
        json=$json$obj
        shift
    done

    json=${json::-3} # remove ,\n of the last object

    # close the json array if er have more than one object to process
    if [[ "$N" -gt "1" ]]; then
        json=$json"\n]\n"
    fi

    echo -e $json
}

# Main loop
echo -e "\nRubus Provider started on http://localhost:1080/\n"

while true
do
    cat response | nc -l 1080 > >(
    if read line; then
        METHOD=$(echo "$line" | cut -d ' ' -f 1)
        URL=$(echo "$line" | cut -d ' ' -f 2)
        IFS='/' read -ra URL_ARRAY <<< "$URL"

        # logger
        d=$(date)
        echo "[ ${d} ] ${METHOD} ${URL}"

        router "${URL_ARRAY[@]:1}" # removes first empty element

        if [[ -z "${ACTION}" ]]; then
            cat not_found > response
            break
        fi

        if [[ -z "${RESPONSE}" ]]; then
            cat not_found > response
            break
        fi

        RESPONSE=$(to-json $RESPONSE)
        echo -e "$(cat header)\n\n$RESPONSE" > response
    fi
    ) 
done
