#!	/usr/local/bin/bash




###########################################################################
#	公共库

#	print colored text
#	$1 = message
#	$2 = color

#	格式化输出
export black='\E[0m\c'
export boldblack='\E[1;0m\c'
export red='\E[31m\c'
export boldred='\E[1;31m\c'
export green='\E[32m\c'
export boldgreen='\E[1;32m\c'
export yellow='\E[33m\c'
export boldyellow='\E[1;33m\c'
export blue='\E[34m\c'
export boldblue='\E[1;34m\c'
export magenta='\E[35m\c'
export boldmagenta='\E[1;35m\c'
export cyan='\E[36m\c'
export boldcyan='\E[1;36m\c'
export white='\E[37m\c'
export boldwhite='\E[1;37m\c'
export EXPORT_LANGUAGE="echo -n"

c_notify=$boldcyan
c_error=$boldred


cecho()
{
    if [ $LANGUAGE = "utf-8" ] 
    then
        message=$1
    else
        echo $1 > /tmp/deploy_tools_tmp
        message=`iconv -f "utf-8" -t $LANGUAGE /tmp/deploy_tools_tmp`
        rm -f /tmp/deploy_tools_tmp
    fi
	color=${2:-$black}

	echo -e "$color"
	echo -e "$message"
	tput sgr0			# Reset to normal.
	echo -e "$black"
	return
}

decho()
{
    if [ $LANGUAGE = "utf-8" ] 
    then
        message=$1
    else
        echo $1 > /tmp/deploy_tools_tmp
        message=`iconv -f "utf-8" -t $LANGUAGE /tmp/deploy_tools_tmp`
        rm -f /tmp/deploy_tools_tmp
    fi
	if [ $UTILS_DEBUG -eq 1 ] 
	then
		color=${2:-$black}

		echo -e "$color"
		echo -e "$message"
		tput sgr0			# Reset to normal.
		echo -e "$black"
	fi
}

cread()
{
	color=${4:-$black}

	echo -e "$color"
	read $1 "$2" $3 
	tput sgr0			# Reset to normal.
	echo -e "$black"
	return
}

#	确认用户的输入
deploy_confirm()
{
    if [ $LANGUAGE = "utf-8" ] 
    then
        message=$1
    else
        echo $1 > /tmp/deploy_tools_tmp
        message=`iconv -f "utf-8" -t $LANGUAGE /tmp/deploy_tools_tmp`
        rm -f /tmp/deploy_tools_tmp
    fi
	while [ 1 = 1 ]
	do
		cread -p "$message [y/n]: " CONTINUE $c_notify
		if [ "y" = "$CONTINUE" ]; then
		  return 1;
		fi

		if [ "n" = "$CONTINUE" ]; then
		  return 0;
		fi
	done

	return 0;
}

error_confirm()
{
    if [ $LANGUAGE = "utf-8" ]
    then
        message=$1
    else
        echo $1 > /tmp/deploy_tools_tmp
        message=`iconv -f "utf-8" -t $LANGUAGE /tmp/deploy_tools_tmp`
        rm -f /tmp/deploy_tools_tmp
    fi
	while [ 1 = 1 ]
	do
		cread -p "$message [y/n]: " CONTINUE $c_error
		if [ "y" = "$CONTINUE" ]; then
		  return 1;
		fi

		if [ "n" = "$CONTINUE" ]; then
		  return 0;
		fi
	done

	return 0;
}

#  获取当前的时间
now()
{
	date +%Y%m%d%H%M%S;
}

########################################
# 检查参数数量是否正确
# check_args_num function_name expect_num achieved_num
########################################
function check_args_num()
{
	if [ $# -ne 3 ] 
	then
		echo "function check_args_num  expect 3 args, but achieved $? args.";
		exit 1;
	fi
	local func_name=$1;
	local expect_num=$2;
	local achieve_num=$3;
	if [ $expect_num -ne $achieve_num ]
	then
		echo "function $func_name expect $expect_num args, but achieved $achieve_num args.";
		exit 1;
	fi
}


####################################################################################################



########################################
# 从svn export 代码
# export_svn $svn_url $local_dir $revision
########################################
function export_svn()           
{                      
	check_args_num $FUNCNAME 3 $# 
	cecho "=== export chunk from svn ===" $c_notify
	local svn_url=$1
	local local_dir=$2
	local revision=$3
	cmd="svn export -r$revision $svn_url $local_dir"
	$cmd > /dev/null
	local errno=$?
	if [ $errno -ne 0 ]
	then 
		echo "run $cmd failed. errno is $errno."
		exit 1
	fi
}   

function export_svn_files()           
{
  check_args_num $FUNCNAME 4 $#
  cecho "=== export chunk from svn ===" $c_notify
  local svn_url=$1
  local local_dir=$2
  local revision=$3
  local files=$4

  mkdir -p $local_dir
  cd $local_dir
  for file in $files
  do
    dir=`dirname $file`
    mkdir -p $dir
    svn export -r$revision $svn_url/$file $file
  done
}


########################################
# 从svn 获取当前head的revision
# get_svn_head_revision svn_url
########################################
function get_svn_head_revision()           
{                      
	check_args_num $FUNCNAME 1 $# 
	local svn_url=$1
	local cmd="svn --xml info $svn_url | grep 'revision' | head \-1 | awk -F '\"' '{ print \$2 }'"
	head_revision=`echo $cmd | bash`
	if [ -z $head_revision ] 
	then
		echo "get_svn_head_revision cannot get HEAD revision. run cmd is $cmd"
		exit 1
	else
		echo $head_revision
	fi
}   

function get_os()
{
	uname -s
}

function get_remote_os()
{
	check_args_num $FUNCNAME 1 $#
  host=$1
	$SSH $host "uname -s" 2>/dev/null
}

function get_remote_shell()
{
	check_args_num $FUNCNAME 1 $#
  host=$1
	$SSH $host "echo \$SHELL" 2>/dev/null
}
########################################
# 获取文件列表
# get_file_list root_dir blacklist node1 node2 node3
# node 可以是文件，也可以是文件夹
########################################
function get_file_list()
{
	local files=''
	local root_dir=$1
	local blacklist=$2
	shift;shift;
	local_root_dir_name_len=`echo "$root_dir/" | wc -m | bc`
	while [ $# -ne 0 ]
	do
		if [ $(get_os) == "Linux" ] 
		then
			file=`echo "/usr/bin/find $root_dir/$1 -regextype posix-extended -type f -not -regex '$blacklist' | cut -c '$local_root_dir_name_len-1000' | xargs echo" | bash`
		else
			# FreeBSD
			file=`echo "/usr/bin/find -E $root_dir/$1 -type f -not -regex '$blacklist' | cut -c '$local_root_dir_name_len-1000' | xargs echo" | bash`
		fi

		shift
		files="$files$file "
	done
	echo $files
}

function init()
{
  mkdir -p $LOCAL_TMP_DIR;
  chmod 777 $LOCAL_TMP_DIR > /dev/null 2>&1
  chmod 777 $LOCAL_TMP_DIR/.. > /dev/null 2>&1
  mkdir -p $LOCAL_DEPLOY_HISTORY_DIR
  chmod 777 $LOCAL_DEPLOY_HISTORY_DIR > /dev/null 2>&1
  touch $DEPLOY_HISTORY_FILE
  chmod 777 $DEPLOY_HISTORY_FILE > /dev/null 2>&1
}

function clean()
{
  rm -rf $LOCAL_TMP_DIR
}

########################################
# 获取线上代码
# get_online_src online_host online_root_dir local_dir file_list
########################################
function get_online_src()
{
	check_args_num $FUNCNAME 4 $#
	local online_host=$1
	local online_root_dir=$2
	local local_dir=$3
	local file_list=$4
	online_src="$PROJECT_NAME-$SSH_USER-$online_host.tgz"
	$SSH $host "$EXPORT_LANGUAGE;tar cvfz $ONLINE_TMP_DIR/$online_src -C $online_root_dir $file_list" > /dev/null 2>&1
	$SCP $host:$ONLINE_TMP_DIR/$online_src $LOCAL_TMP_DIR/$online_src > /dev/null 2>&1

	if [ ! -s "$LOCAL_TMP_DIR/$online_src" ]; then
		cecho "错误：获取线上代码出错. host is $online_host." $c_error
    cecho  "      线上代码路径为 $host:$ONLINE_TMP_DIR/$online_src"
    cecho  "      本地代码路径为 $LOCAL_TMP_DIR/$online_src"
		exit 1
	fi
	local_online="$local_dir"
	rm -rf $local_online && mkdir -p $local_online && chmod 777 $local_online
	tar xzf $LOCAL_TMP_DIR/$online_src -C $local_online
	decho "online_src tgz is $LOCAL_TMP_DIR/$online_src, src_dir is $local_online"
}

# 获取线上目前的所有代码， 用于diff 比较
# get_online_src_tar host online_root_dir local_dir
function get_online_src_all()
{
	check_args_num $FUNCNAME 3 $#
	host=$1
	online_root_dir=$2
	local_dir=$3

	if [ $(get_remote_os $host) == 'Linux' ]
	then
		TARCREATE="tar zcvfh"
	else
		TARCREATE="tar zcvfL"
	fi

	online_src="$ONLINE_TMP_DIR/$PROJECT_NAME-$CURRENT_TIME.tgz"
	$SSH $host "$EXPORT_LANGUAGE;$TARCREATE $online_src $TAR_EXCLUDE $online_root_dir" > /dev/null 2>&1
	local_src="$LOCAL_TMP_DIR/$PROJECT_NAME-$CURRENT_TIME.tgz"
	$SCP $host:$online_src $local_src 
	rm -rf $local_dir && mkdir -p $local_dir
	tar xz -f $local_src -C $local_dir
}


# 检查有哪些文件不同， 提示用户确认是否继续
# check_files_diff files svn_source_dir bench_hostname bench_host_src_dir this_hostname this_host_src_dir
function check_files_diff()
{
	check_args_num $FUNCNAME 6 $#
	files=$1	
	svn_src_dir=$2
	bench_hostname=$3
	bench_host_src_dir=$4
	this_hostname=$5
	this_host_src_dir=$6

	for file in $files
	do
		#   确定文件类型，只针对 text 类型
		type=`file $svn_src_dir/$file | grep "text"`
		if [ -z "$type" ]; then
			continue
		fi

		cecho "\t$file"
		diffs=`diff -Bb $svn_src_dir/$file $this_host_src_dir/$file`

		#   如果没有不同就不要确认
		if [ -z "$diffs" ]; then
			continue
		fi

		# 如果不是基准主机， 且与基准主机的内容一致 自动提交
		if [ "$this_hostname" != "$bench_hostname" ]
		then
			tmp=`diff -Bb $bench_host_src_dir/$file $this_host_src_dir/$file`
			if [ -z "$tmp" ]; then
				continue
			fi
		fi

		#   进行 vimdiff
		sleep 1
		vimdiff $svn_src_dir/$file $this_host_src_dir/$file

		deploy_confirm "    修改确认 $file ?"
		if [ 1 != $? ]; then
			exit 1;
		fi
	done
}

# upload_src host src_tgz
# 默认上传到REMOTE_DEPLOY_DIR
function upload_src()
{
	check_args_num $FUNCNAME 2 $#
	upload_src_to $1 $2 $REMOTE_DEPLOY_DIR
	return $?
}

# 上传tar的源代码到指定目录下
# upload_src host src_tgz dst_dir
function upload_src_to()
{
	check_args_num $FUNCNAME 3 $#
	host=$1
	src_tgz=$2
	realdst=$3
	#   上传源文件
	$SSH $host "$EXPORT_LANGUAGE;mkdir -p $ONLINE_BACKUP_DIR"
	uploaded_src_tgz="$ONLINE_BACKUP_DIR/$CURRENT_TIME-$CURRENT_REVISION-$PROJECT_NAME-up.tgz"
	$SCP $src_tgz $host:$uploaded_src_tgz 
	$SSH $host "test -s $uploaded_src_tgz"
	if [ 0 -ne $? ]; then
		cecho "\t错误：文件上传失败" $c_error
		exit 1
	fi
	$SSH $host "$EXPORT_LANGUAGE;tar xvfzm $uploaded_src_tgz -C $realdst" 2>&1 | sed -e 's/^/   /'

	if [ 0 != $? ]
	then
		cecho "\t错误：部署文件失败" $c_error
		deploy_confirm "    继续部署？"
		if [ 1 != $? ]; then
			exit 1;
		fi
	fi
}




# backup_online_src $host $files
function backup_online_src()
{
	check_args_num $FUNCNAME 3 $#
	host=$1
	backup_src_tgz=$2
	files=$3
	backup_dir=`dirname $backup_src_tgz`
	$SSH $host "$EXPORT_LANGUAGE;mkdir -p $backup_dir;tar cvfz ${backup_src_tgz} -C $REMOTE_DEPLOY_DIR $files"
	decho "\t $host 备份文件路径： ${backup_src_tgz}"
	$SSH $host "test -s $backup_src_tgz"
	if [ 0 -ne $? ]; then
		cecho "\t错误：远程主机原始文件备份失败" $c_error
		exit 1
	fi
}

function ssh_run()
{
	check_args_num $FUNCNAME 2 $#
	host=$1
	cmd=$2
	$SSH $host "$EXPORT_LANGUAGE;$cmd"
	result=$?
	if [ $result -ne 0 ]
	then
		cecho "FAILED $SSH $host $cmd" $c_error
	fi
	return $result
}

function sudo_ssh_run()
{
	check_args_num $FUNCNAME 2 $#
	host=$1
	cmd=$2
	if [ -z $sudo_password ]
	then
		cread -p "input your sudo_password for sudo_ssh_run command:    " sudo_password $c_notify
	fi
	ssh -t $host "echo $sudo_password | sudo -S $cmd"
	result=$?
	if [ $result -ne 0 ]
	then
		cecho "FAILED:  ssh -t $host sudo -p $sudo_password $cmd" $c_error
		return $result
	fi
}

function check_succ()
{
	check_args_num $FUNCNAME 2 $#
	if [ $1 -ne 0 ] 
	then
		error_confirm "$2" 
		if [ 1 != $? ]; then
			exit 1;
		fi
	fi
}
