#!/bin/bash

CurDir=$(readlink -nf $0 | xargs dirname)
source $CurDir/include.sh
echo "--------------"
echo "[info]"`date "+%G-%m-%d-%H:%M:%S"`
#service is started
check_service_started
if [[ 1 == $? ]]
then
	echo "[info]serveice is running"
else
	Server=`hostname | awk -F. '{ print $1 }'`
	echo "[info]service is stop, now starting"
  send_msg "DOWN:{@project}-$Server" "{@project} down-$Server"
	echo ""
	/bin/bash $PRJHOME/serverctl start
fi
echo "[info]"`date "+%G-%m-%d-%H:%M:%S"`
echo "--------------"

