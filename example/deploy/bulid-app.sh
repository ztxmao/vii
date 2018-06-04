#!/bin/sh

PRJHOME="$GOPATH/src/{@project}"
ExeFile="{@project}"


echo "App root: $PRJHOME"
rm -f $PRJHOME/bin/$ExeFile
cd $PRJHOME/src
#main pkg dir name
#PkgName=server
#go build -o $RJHOME/bin/$ExeFile $PkgName
export CGO_LDFLAGS="-L$PRJHOME/lib/"
export LD_LIBRARY_PATH="$PRJHOME/lib/"

go build -o $PRJHOME/bin/$ExeFile
if [[ "$?" == "0" ]]
then
    echo "build succ: "$PRJHOME/bin/$ExeFile
else
    echo "build fail"
fi
#export GOPATH=$RawGoPath

op=$1
if [[ "X$op" == "Xrun" ]]
then 
    $PRJHOME/bin/$ExeFile
fi

