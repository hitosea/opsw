#!/bin/bash

PLATFORM=$(uname -s)
FRAMEWORK=$(uname -m)
if [[ ${PLATFORM} != "Linux" ]]; then
    echo "不支持的系统：${PLATFORM}，仅支持 Linux"
    exit 1
fi
if [[ ${FRAMEWORK} == "aarch64" ]]; then
    FRAMEWORK="arm64"
fi
if [[ ${FRAMEWORK} == "x86_64" ]]; then
    FRAMEWORK="amd64"
fi
if [[ ${FRAMEWORK} != "amd64" ]] && [[ ${FRAMEWORK} != "arm64" ]]; then
    echo "不支持的架构：${FRAMEWORK}，仅支持 amd64 和 arm64"
    exit 1
fi

VERSION=$(curl -s https://api.github.com/repos/hitosea/opsw/releases/latest | grep 'tag_name' | cut -d\" -f4)
if [[ "x${VERSION}" == "x" ]];then
    echo "获取最新版本失败，请稍候重试"
    exit 1
fi

echo "开始下载 ${VERSION} 版本在线安装包"

package_file_name="opsw_${PLATFORM}_${FRAMEWORK}.tar.gz"
package_download_url="https://github.com/hitosea/opsw/releases/download/${VERSION}/${package_file_name}"

echo "安装包下载地址： ${package_download_url}"

curl -LOk -o ${package_file_name} ${package_download_url}
if [ ! -f ${package_file_name} ];then
	echo "下载安装包失败，请稍候重试。"
	exit 1
fi

tar zxvf ${package_file_name}
if [ $? != 0 ];then
	echo "下载安装包失败，请稍候重试。"
	rm -f ${package_file_name}
	exit 1
fi

cd opsw_${PLATFORM}_${FRAMEWORK}

/bin/bash install.sh