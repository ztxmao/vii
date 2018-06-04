#!/bin/bash


####################################################################################################
#   配置项

#   include lib
this_file=`pwd`"/"$0

DEPLOY_TOOLS_DIR=`dirname $this_file`
. $DEPLOY_TOOLS_DIR/conf.sh
. $DEPLOY_TOOLS_DIR/utils.sh


####################################################################################################
# 使用帮助
if [ $# -lt 1 ] || [ "-h" = "$1" ] || [ "--help" = "$1" ]
then
	cecho "用法: $0 FILE";
	cecho "FILE* :  revert 所使用的备份文件"
	cecho "FILE* :  /home/sync/deploy_history/xxxx" 
	cecho "用法: $0 ls";
  cecho "列出部署的历史，每次备份的文件位置"

	exit 0;
fi

###################################################################################################
if [ $1 == "ls" ]
then
	no=0
	echo "deploy history file is $DEPLOY_HISTORY_FILE"
  while read line
  do
		no=`echo "$no + 1" | bc`
		cecho "$no\t$line";
  done < $DEPLOY_HISTORY_FILE
	exit 0;
fi


revert_src_tgz=$1
#   开始回滚
hosts=$online_cluster
echo $hosts
for host in ${hosts}
do
	cecho "\n=== ${host} ===\n" $c_notify
    echo $SSH
	$SSH $host "test -s $revert_src_tgz"
	if [ 0 -ne $? ]; then
		cecho "\t错误：远程主机原始备份文件不存在" $c_error
		deploy_confirm "    是否继续 ?"
		if [ 1 != $? ]; then
			exit 1;
		else
			continue
		fi
	fi
	cecho "\n=== ${host} 回滚文件列表===\n" $c_notify
	$SSH $host tar xvfz $revert_src_tgz -C $REMOTE_DEPLOY_DIR
	if [ 0 -ne $? ]; then
		cecho "\t错误：$host  回滚失败 " $c_error
	fi
	if [ "$AUTORUN_RELEASE_CMD" != '' ]
    then
        ssh_run $host "$AUTORUN_RELEASE_CMD"
        check_succ $? "运行脚本失败： $AUTORUN_RELEASE_CMD, 是否继续?"
    fi
done


deploy_confirm " revert 完成， 是否删除deploy_history中的这个记录?"
if [ 1 != $? ]; then
	exit 0;
else
	cat $DEPLOY_HISTORY_FILE > $DEPLOY_HISTORY_FILE_BAK
	cat $DEPLOY_HISTORY_FILE_BAK | grep -v "$revert_src_tgz" > $DEPLOY_HISTORY_FILE
fi
