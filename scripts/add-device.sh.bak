#!/bin/bash

set -e

setup_new_device()
{
    HOSTNAME=$1

    # initialize the tftp folder for the device
    mkdir /tftp/$HOSTNAME

    # initilize the nfs share for the device
    cp -ar /pxe/nfs/NFS-TEMPLATE /pxe/nfs/$HOSTNAME

    # personilize the nfs share (i.e. set the hosname)
    echo "$HOSTNAME" > /pxe/nfs/$HOSTNAME/etc/hostname
    sed -i.back /raspberrypi/d /pxe/nfs/$HOSTNAME/etc/hosts
    echo -e "127.0.1.1\t$HOSTNAME" >> /pxe/nfs/$HOSTNAME/etc/hosts

    # configure the NFS share by passing a command to the kernel
    echo "console=serial0,115200 console=tty1 root=/dev/nfs \
    nfsroot=172.29.0.100:/pxe/nfs/$HOSTNAME,vers=3 rw ip=dhcp rootwait elevator=deadline" \
    > /pxe/nfs/$HOSTNAME/boot/disabled/cmdline.txt

    # configure the NFS share on the server
    echo "/pxe/nfs/$HOSTNAME *(rw,sync,no_subtree_check,no_root_squash)" \
    >> /etc/exports
}

setup_new_device $1
