#!	/usr/local/bin/bash

################################################################################
#
#   配置项
#
#   demo: ./tools/all-exec.sh /home/q/system/video_flow/serverctl start
#
#   include lib
this_file=`pwd`"/"$0

DEPLOY_TOOLS_DIR=`dirname $this_file`
. $DEPLOY_TOOLS_DIR/conf.sh
. $DEPLOY_TOOLS_DIR/utils.sh

################################################################################
# 使用帮助
if [ $# -lt 1 ] || [ "-h" = "$1" ] || [ "--help" = "$1" ]
then
	cecho "用法: $0 cmd"
	exit 0;
fi


hosts=$online_cluster
# 确认服务器列表
cecho "\n=== 服务器列表 === \n" $c_notify
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


for host in $hosts
do
	cecho "\n=== $host  === \n" $c_notify
	$SSH $host $@
done
