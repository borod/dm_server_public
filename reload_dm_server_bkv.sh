#!/bin/bash
. /opt/dm/dmvars.conf

# сносит локальные изменения и берёт версию репозитория
echo -e "\e[35mDM Server will be pulled from Git Repository and restarted \e[0m\e[97mat \e[34m$DMSERVERPORT\e[0m"

sudo systemctl status dm_server
if [ -z "$CURRENT_USER_HOME" ]; then
    CURRENT_USER_HOME=/home/borod
fi
if [ -z "$DM_SERVER_SRC_PATH" ]; then
    DM_SERVER_SRC_PATH=$CURRENT_USER_HOME/src/dm_server
fi

cd $DM_SERVER_SRC_PATH
sudo systemctl stop dm_server
eval `ssh-agent -s`
ssh-add ~/.ssh/bkv_ed25519
echo -e "\e[35mResetting local changes to Repository HEAD version... \e[0m"
git reset --hard HEAD
git pull
chmod +x $DM_SERVER_SRC_PATH/reload_dm_server_bkv.sh
chmod +x $DM_SERVER_SRC_PATH/run_dm_server.sh
/bin/cp -rf ./reload_dm_server_bkv.sh ~/
/bin/cp -rf ./run_dm_server.sh ~/
sudo rm -rf /var/cache/nginx/
cd $DM_SERVER_SRC_PATH
go build
#echo cp $DM_SERVER_SRC_PATH/dm_server /opt/dm/dm_server/
sudo cp $DM_SERVER_SRC_PATH/dm_server /opt/dm/dm_server/
sudo cp $DM_SERVER_SRC_PATH/configuration.json /opt/dm/dm_server/
sudo chmod 770 /opt/dm/dm_server/dm_server
echo -e "\e[35mStarting DM Server \e[0m\e[97mat \e[34m$DMSERVERPORT\e[0m"
#go run main.go
sudo systemctl start dm_server