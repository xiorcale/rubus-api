#!/bin/bash

set -e

deploy_device()
{
    HOSTNAME=$1

    mount --bind /pxe/nfs/$HOSTNAME/boot /tftp/$HOSTNAME
    if [[ -f /tftp/$HOSTNAME/disabled/start4.elf ]]; then
        mv /pxe/nfs/$HOSTNAME/boot/disabled/* /pxe/nfs/$HOSTNAME/boot/
    fi
}

deploy_device $1
