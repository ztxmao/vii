#!	/usr/local/bin/bash

################################################################################
#   配置项

#   include lib
this_file=`pwd`"/"$0

DEPLOY_TOOLS_DIR=`dirname $this_file`
. $DEPLOY_TOOLS_DIR/conf.sh
. $DEPLOY_TOOLS_DIR/utils.sh

################################################################################
# 使用帮助
if [ $# -ne 2 ] || [ "-h" = "$1" ] || [ "--help" = "$1" ]
then
	cecho "用法: $0 remote_file_name local_dir"
	exit 0;
fi


hosts=$online_cluster
# 确认服务器列表
cecho "\n=== 下载服务器列表 === \n" $c_notify
no=0;
for host in $hosts
do
	no=`echo "$no + 1" | bc`
	cecho "$no\t$host";
done
echo ""
deploy_confirm "确认服务器列表？"
if [ 1 != $? ]; then
	exit 1;
fi

remote_name=$1
remote_base_name=`basename $1`
local_dir=$2
for host in $hosts
do
	cmd="$SCP -r $host:$remote_name $local_dir/$host-$remote_base_name"
	decho "cmd is $cmd"
	$cmd
done
