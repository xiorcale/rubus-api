#!/bin/bash

set -e

setup_new_device()
{
    HOSTNAME=$1

    # initialize the tftp folder for the device
    mkdir /tftp/$HOSTNAME

    # copy the upper directory for the overlay
    cp -ar /pxe/nfs/NFS-TEMPLATE-UPPER /pxe/nfs/$HOSTNAME

    # personilize the upper directory of the overlay (i.e. set the hosname)
    echo "$HOSTNAME" > /pxe/nfs/$HOSTNAME/etc/hostname
    sed -i.back /raspberrypi/d /pxe/nfs/$HOSTNAME/etc/hosts
    echo -e "127.0.1.1\t$HOSTNAME" >> /pxe/nfs/$HOSTNAME/etc/hosts

    # configure the NFS share by passing a command to the kernel
    echo "console=serial0,115200 console=tty1 root=/dev/nfs \
    nfsroot=172.29.0.100:/pxe/nfs/$HOSTNAME,vers=3 rw ip=dhcp rootwait elevator=deadline" > /pxe/nfs/$HOSTNAME/boot/disabled/cmdline.txt

    # configure the NFS share on the server
    echo "/pxe/nfs/$HOSTNAME *(rw,sync,no_subtree_check,no_root_squash)" >> /etc/exports

    # mount the overlay file system
    mkdir /pxe/nfs/$HOSTNAME-work
    mount -t overlay overlay -o lowerdir=/pxe/nfs/NFS-TEMPLATE-LOWER,upperdir=/pxe/nfs/$HOSTNAME,workdir=/pxe/nfs/$HOSTNAME-work /pxe/nfs/$HOSTNAME

    # remove the temporary directory used for creating the overlay
    # rm -rf /pxe/nfs/$HOSTNAME-work
}

setup_new_device $1
