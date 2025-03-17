#!/bin/bash
. /opt/dm/dmvars.conf
cd ~/src/dm_server/
echo -e "\e[35mStarting DM Server \e[0m\e[97mat \e[34m$DMSERVERPORT\e[0m"
go run main.go