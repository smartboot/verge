#!/bin/bash

dir=$(dirname "$0")
cd ${dir}

# 定义需要打包的程序目录
app_name="verge"
VERSION="dev"
if [ -n $1 ]; then
  #截取/末尾的字符串
  VERSION=${1##*/}
fi
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
go mod tidy
go mod vendor
buildTime=`date +%Y%m%d%H%M%S`
version_flag="-X 'pkg.Version=${VERSION}'"
build_time_flag="-X 'pkg.BuildTime=${buildTime}'"


deploy_dir="deploy"
rm -rf ${deploy_dir}
mkdir ${deploy_dir}
build(){
    GOOS=$1
    GOARCH=$2
    rm -rf ${app_name}
    output_flag="${app_name}/${app_name}"

    GOOS=$GOOS GOARCH=${GOARCH} go build -ldflags "-s -w ${version_flag} ${build_time_flag}" -o ${output_flag} cmd/main.go
    # 添加错误处理，确保在构建失败时脚本能够退出
    if [ $? -ne 0 ]; then
      echo "构建失败: ${output_flag}"
      exit 1
    fi

    echo "成功构建: ${output_flag}"

    echo "开始打包: ${app_name}-${GOOS}-${GOARCH}-${VERSION}"
    cp -R res ${app_name}
    cp -R platform/${GOOS}/* ${app_name}
    chmod +x ${app_name}/*.sh
    chmod +x ${app_name}/*.bat

    deploy_file=''
    if [[ "${GOOS}" == "windows" ]]; then
      deploy_file="${app_name}-${GOOS}-${GOARCH}-${VERSION}.zip"
      zip -r ${deploy_file} ${app_name}
    else
      deploy_file="${app_name}-${GOOS}-${GOARCH}-${VERSION}.tar.gz"
      tar -czf ${deploy_file} ${app_name}
    fi
    mv ${deploy_file} ${deploy_dir}/
}
#make build VERSION=${VERSION} BuildTime=$(date +%Y%m%d%H%M%S)
build linux arm64
build linux amd64
build linux arm

build windows amd64
build windows arm64
build darwin amd64
build darwin arm64

rm -rf ${app_name}