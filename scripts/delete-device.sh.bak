#!/bin/bash

HOSTNAME=$1

# remove the tftp folder used to boot the device
umount /tftp/$HOSTNAME
rm -r /tftp/$HOSTNAME

# remove the nfs share
rm -r /pxe/nfs/$HOSTNAME
sed -i /$HOSTNAME/d /etc/exports
