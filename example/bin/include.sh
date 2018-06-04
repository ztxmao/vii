#!/bin/bash

#
#全局变量首字母大写,局部变量首字母小写,驼峰式
#

PRJHOME=$(readlink -nf $0 | xargs dirname | xargs dirname)

ExeFile="{@project}"
#ServerPort=$(cat $PRJHOME/conf/server.conf | perl -ne 'print "$1" if m/^port\s*=\s*(\d+)/')
Ports=$(cat $PRJHOME/conf/server.conf | perl -ne 'print "$1 " if m/^port\s*=\s*(\d+)/')
ServerPortArr=($Ports)
ServerPort=${ServerPortArr[0]}

#started 1, not start 0
check_service_started()
{
	runFlag=$(ps aux | grep "bin/$ExeFile" | grep -v "grep" | wc -l)
	if [[ $runFlag -ge 1 ]]
	then 
		return 1
	fi

	return 0
}

check_mon_result()
{
	url=$1
	#echo "mon url: $url"
	succFlag=$(curl -s $url | grep "REV" | grep "OK" | wc -l)

	if [[ $succFlag -eq 1 ]]
	then
		echo "err=0"
	else
		echo "err=1"
	fi
}
send_msg()
{
    #TODO 发送报警信息
    echo "alert"
}
