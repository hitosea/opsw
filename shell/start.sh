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

mkdir -p /tmp/opsw
if [ $? != 0 ];then
    echo "创建临时目录失败，请检查权限"
    exit 1
fi

#####################################
############# opsw 部分 #############
#####################################

function olog() {
    echo -e "[Opsw Log]: $1 "
}

VERSION=$(curl -s https://api.github.com/repos/hitosea/opsw/releases/latest | grep 'tag_name' | cut -d\" -f4)
if [[ "x${VERSION}" == "x" ]];then
    olog "获取最新版本失败，请稍候重试"
    exit 1
fi

olog "开始下载 ${VERSION} 版本在线安装包"

package_name=$(echo "opsw-${VERSION}-${PLATFORM}-${FRAMEWORK}" | tr '[A-Z]' '[a-z]')
package_file_name="${package_name}.tar.gz"
package_download_url="https://github.com/hitosea/opsw/releases/download/${VERSION}/${package_file_name}"

olog "安装包下载地址： ${package_download_url}"

cd /tmp/opsw
curl -LOk -o ${package_file_name} ${package_download_url}
if [ ! -f ${package_file_name} ];then
	olog "下载安装包失败，请稍候重试。"
	exit 1
fi

tar zxvf ${package_file_name}
if [ $? != 0 ];then
	olog "下载安装包失败，请稍候重试。"
	rm -f ${package_file_name}
	exit 1
fi

cd ${package_name}
sed -i "/mode:/c mode: server" opsw.yaml
sed -i "/url:/c url: {{.URL}}" opsw.yaml
sed -i "/token:/c token: {{.TOKEN}}" opsw.yaml

/bin/bash tool.sh {{.ACTION}}

#####################################
############# panel 部分 #############
#####################################

echo "--------------------"

function plog() {
    echo -e "[Opsw Log] [Panel]: $1 "
}

VERSION=$(curl -s https://api.github.com/repos/hitosea/1Panel/releases/latest | grep 'tag_name' | cut -d\" -f4)
if [[ "x${VERSION}" == "x" ]];then
    plog "获取最新版本失败，请稍候重试"
    exit 1
fi

plog "开始下载 ${VERSION} 版本在线安装包"

package_name=$(echo "1panel-${VERSION}-${PLATFORM}-${FRAMEWORK}" | tr '[A-Z]' '[a-z]')
package_file_name="${package_name}.tar.gz"
package_download_url="https://github.com/hitosea/1Panel/releases/download/${VERSION}/${package_file_name}"

plog "安装包下载地址： ${package_download_url}"

cd /tmp/opsw
curl -LOk -o ${package_file_name} ${package_download_url}
if [ ! -f ${package_file_name} ];then
	plog "下载安装包失败，请稍候重试。"
	exit 1
fi

tar zxvf ${package_file_name}
if [ $? != 0 ];then
	olog "下载安装包失败，请稍候重试。"
	rm -f ${package_file_name}
	exit 1
fi

cd ${package_name}
mv ./1panel ./opspanel
cp ./opspanel /usr/local/bin && chmod +x /usr/local/bin/opspanel
if [[ ! -f /usr/bin/opspanel ]]; then
    ln -s /usr/local/bin/opspanel /usr/bin/opspanel >/dev/null 2>&1
fi

mkdir -p /opt/manage-panel
cat > /usr/bin/manage-panel <<-EOF
BASE_DIR=/opt/manage-panel
ORIGINAL_PORT={{.PANEL_PORT}}
ORIGINAL_VERSION=v1.2.4
ORIGINAL_ENTRANCE=entrance
ORIGINAL_USERNAME={{.PANEL_USERNAME}}
ORIGINAL_PASSWORD={{.PANEL_PASSWORD}}
EOF
cat > /etc/systemd/system/opspanel.service <<-EOF
[Unit]
Description=Panel, a modern open source linux cluster management
After=syslog.target network.target

[Service]
ExecStart=/usr/bin/opspanel
ExecReload=/bin/kill -s HUP \$MAINPID
Restart=always
RestartSec=5
LimitNOFILE=1048576
LimitNPROC=1048576
LimitCORE=1048576
Delegate=yes
KillMode=process

[Install]
WantedBy=multi-user.target
EOF

systemctl enable opspanel; systemctl daemon-reload 2>&1

plog "启动服务"
systemctl start opspanel

for b in {1..30}
do
    sleep 3
    service_status=`systemctl status opspanel 2>&1 | grep Active`
    if [[ $service_status == *running* ]];then
        plog "服务启动成功!"
        break;
    else
        plog "服务启动出错!"
        exit 1
    fi
done

plog "安装完成"

cd ~
rm -rf /tmp/opsw