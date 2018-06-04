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
	cecho "用法: $0 FILE1 [ FILE2 ... ]";
	cecho "FILE* : 需要上传的文件/目录；注意：每一个文件必须是相对于 $PROJECT_HOME 根目录的相对路径"
	exit 0;
fi


init

####################################################################################################
# 从当前目录获取代码
CURRENT_REVISION=$VERSION
CURRENT_TIME=$(now)
LOCAL_SOURCE_DIR=$DEPLOY_TOOLS_DIR/..  # deploy需要在home目录的一级子目录
# 获取当前要上线的代码文件列表
files=$(get_file_list $LOCAL_SOURCE_DIR $BLACKLIST $@)
if [ 0 -ne `expr "$files" : ' *'` ]; then
	cecho "\n没有找到要上传的文件，请调整输入参数" $c_error
	exit 1;
fi

# 确认文件列表
cecho "\n=== 上传文件列表 === \n" $c_notify
no=0;
for file in $files
do
	no=`echo "$no + 1" | bc`
	cecho "$no\t$file";
done
echo ""
deploy_confirm "确认文件列表？"
if [ 1 != $? ]; then
	exit 1;
fi


# 待上线的代码打包

#   源文件打包
cecho "\n=== 上线文件打包 === \n" $c_notify
src_tgz="$LOCAL_TMP_DIR/patch.${PROJECT_NAME}-${CURRENT_REVISION}-${CURRENT_TIME}.tgz"
decho $LOCAL_SOURCE_DIR
decho $files
tar cvfz $src_tgz -C $LOCAL_SOURCE_DIR $files > /dev/null 2>&1
decho "打包文件:   $src_tgz"
if [ ! -s "$src_tgz" ]; then
	cecho "错误：文件打包失败" $c_error
	exit 1
fi

#   开始上线代码
if [ "$DEPLOY_BETA" != "beta" ]
then
	hosts=$online_cluster
else
	hosts=$beta_cluster
fi


#记录当前的更新日志
mkdir -p $LOCAL_DEPLOY_HISTORY_DIR
backup_src_tgz="$ONLINE_BACKUP_DIR/$CURRENT_TIME-$PROJECT_NAME-bak.tgz"
echo $backup_src_tgz $USER >> $DEPLOY_HISTORY_FILE


for host in ${hosts}
do
  if [ $(get_remote_os $host) == "Linux" ]
  then
    LINK="ln -T -s"
  else
    LINK="ln -s"
  fi
  if [[ $(get_remote_shell $host) == *csh ]]
  then
    EXPORT_LANGUAGE="setenv LANGUAGE $LANGUAGE"
  else
    EXPORT_LANGUAGE="export LANGUAGE=$LANGUAGE"
  fi
	cecho "\n=== ${host} ===\n" $c_notify
	
	# 备份线上代码
	backup_online_src $host $backup_src_tgz "$files"
	# 上传需要更新的代码
	upload_src $host $src_tgz
    
    #   记录基准主机
    if [ "" = "$bench_host" ]; then
        bench_host="$host"
    fi
    
	##########################################################################################################
	# 在这里添加更新代码后需要执行的程序。 单元测试，建立link等
	for script in $AUTORUN_RELEASE
	do
	    cecho "    执行此命令恢复原始版本： $SSH $host \"tar xvfz $backup_src_tgz -C $REMOTE_DEPLOY_DIR;sh  $REMOTE_DEPLOY_DIR/serverctl reload\"";
		ssh_run $host "$REMOTE_DEPLOY_DIR/$script"
		check_succ $? "运行脚本失败： $REMOTE_DEPLOY_DIR/$script, 是否继续?"
	done

    if [ "$AUTORUN_RELEASE_CMD" != '' ]
    then
	    cecho "    执行此命令恢复原始版本： $SSH $host \"tar xvfz $backup_src_tgz -C $REMOTE_DEPLOY_DIR;sh  $REMOTE_DEPLOY_DIR/serverctl reload\"";
        ssh_run $host "$AUTORUN_RELEASE_CMD"
		check_succ $? "运行脚本失败： $AUTORUN_RELEASE_CMD, 是否继续?"
    fi
	##########################################################################################################

	verify="    --- 上线完毕，执行此命令恢复原始版本： $SSH $host \"tar xvfz $backup_src_tgz -C $REMOTE_DEPLOY_DIR;sh  $REMOTE_DEPLOY_DIR/serverctl reload\"";

	if [ "$host" == "$bench_host" ]
	then
		echo ""
		deploy_confirm "$verify，请验证效果"
		if [ 1 != $? ]; then
			exit 1;
		fi
	else
		cecho "\n$verify \n" $c_notify
	fi
done
cecho "\n===============================================================\n" $c_notify
cecho "\n上线完毕，执行此命令回滚所有机器：sh deploy-revert.sh $backup_src_tgz\n" $c_notify
clean
