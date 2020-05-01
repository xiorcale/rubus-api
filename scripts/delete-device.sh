#!/bin/bash

delet_device()
{
    HOSTNAME=$1

    # remove the tftp folder used to boot the device
    umount /tftp/$HOSTNAME
    rm -r /tftp/$HOSTNAME

    # unmount the overlay file system and remove it
    umount /pxe/nfs/$HOSTNAME
    rm -r /pxe/nfs/$HOSTNAME*
    sed -i /$HOSTNAME/d /etc/exports
}

delet_device $1
