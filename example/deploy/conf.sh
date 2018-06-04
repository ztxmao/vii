#!/bin/bash

# 如果term设置的是gbk编码， 改为gbk编码
export LANGUAGE="utf-8"
#export LANGUAGE="gbk"                                                  

# 项目名
PROJECT_NAME="{@project}"
APP_NAME="{@project}"
VERSION="1.0"

# 线上集群
servers_bjyt="your hostname"

online_cluster="$servers_bjyt"
beta_cluster=""

# 项目部署的目录， link 到  $REAL_REMOTE_DEPLOY_DIR 上
REMOTE_DEPLOY_DIR="/home/q/system/$PROJECT_NAME"                    		
# 部署的真实目录
REAL_REMOTE_DEPLOY_DIR="/home/q/system/$PROJECT_NAME-$VERSION"           	

# 部署使用的账号
SSH_USER="sync"

# 设置为1的时候， 会输出debug信息
UTILS_DEBUG=0

# 安装后自动执行初始化脚本

# 运行deploy-package.sh 后自动通过全路径直接运行，脚本需要有可执行权限
#AUTORUN_PACKAGE="a.sh b.sh c.sh"                  
# 同上 root权限
#SUDO_AUTORUN_PACKAGE=""                           
REGION=`hostname | awk -F. '{ print $3 }'`
AUTORUN_PACKAGE_CMD="cd $REMOTE_DEPLOY_DIR;/bin/bash serverctl envinit $REGION;/bin/bash serverctl start;"

# 运行deploy-release.sh 后自动通过全路径直接运行，脚本需要有可执行权限
AUTORUN_RELEASE_CMD="cd $REMOTE_DEPLOY_DIR;/bin/bash serverctl reload"

# 用于diff命令  打包时过滤logs目录
DEPLOY_BASENAME=`basename $REMOTE_DEPLOY_DIR`
TAR_EXCLUDE="--exclude $DEPLOY_BASENAME/logs" 

########## 不要修改 #########################

SSH="sudo -u $SSH_USER ssh -oStrictHostKeyChecking=no "
SCP="sudo -u $SSH_USER scp -oStrictHostKeyChecking=no "

LOCAL_TMP_DIR="/tmp/deploy_tools/$USER"                                   # 保存本地临时文件的目录
BLACKLIST='(.*\.tmp$)|(.*\.log$)|(.*\.svn.*)'                             # 上传代码时过滤这些文件
ONLINE_TMP_DIR="/tmp"													  # 线上保存临时文件的目录
ONLINE_BACKUP_DIR="/home/$SSH_USER/deploy_history/$PROJECT_NAME"          # 备份代码的目录
LOCAL_DEPLOY_HISTORY_DIR="/home/$USER/deploy_history/$PROJECT_NAME"  
DEPLOY_HISTORY_FILE="$LOCAL_DEPLOY_HISTORY_DIR/deploy_history"            # 代码更新历史(本地文件）
DEPLOY_HISTORY_FILE_BAK="$LOCAL_DEPLOY_HISTORY_DIR/deploy_history.bak" 
