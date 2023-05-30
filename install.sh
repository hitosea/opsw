#!/bin/bash

# web面板安装、卸载命令

echo
cat << EOF
 ██████╗ ██████╗ ███████╗██╗    ██╗
██╔═══██╗██╔══██╗██╔════╝██║    ██║
██║   ██║██████╔╝███████╗██║ █╗ ██║
██║   ██║██╔═══╝ ╚════██║██║███╗██║
╚██████╔╝██║     ███████║╚███╔███╔╝
 ╚═════╝ ╚═╝     ╚══════╝ ╚══╝╚══╝

EOF

function log() {
    echo -e "[Opsw Log]: $1 "
}

function Prepare_System(){
    PLATFORM=$(uname -s)
    FRAMEWORK=$(uname -m)
    if [[ ${PLATFORM} != "Linux" ]]; then
        log "不支持的系统：${PLATFORM}，仅支持 Linux"
        exit 1
    fi
    if [[ ${FRAMEWORK} == "aarch64" ]]; then
        FRAMEWORK="arm64"
    fi
    if [[ ${FRAMEWORK} == "x86_64" ]]; then
        FRAMEWORK="amd64"
    fi
    if [[ ${FRAMEWORK} != "amd64" ]] && [[ ${FRAMEWORK} != "arm64" ]]; then
        log "不支持的架构：${FRAMEWORK}，仅支持 amd64 和 arm64"
        exit 1
    fi

    if which opsweb >/dev/null 2>&1; then
        log "已安装，请勿重复安装"
        exit 1
    fi
}

function Set_Dir(){
    mkdir -p /opt/opsweb/tmp
    if [ $? != 0 ];then
        log "创建安装目录 /opt/opsweb 失败，请检查权限"
        exit 1
    fi
}

function Set_Port(){
    DEFAULT_PORT=`expr $RANDOM % 55535 + 10000`
    while true; do
        read -p "设置端口（默认为$DEFAULT_PORT）：" PANEL_PORT
        if [[ "$PANEL_PORT" == "" ]];then
            PANEL_PORT=$DEFAULT_PORT
        fi
        if ! [[ "$PANEL_PORT" =~ ^[1-9][0-9]{0,4}$ && "$PANEL_PORT" -le 65535 ]]; then
            echo "错误：输入的端口号必须在 1 到 65535 之间"
            continue
        fi
        log "您设置的端口为：$PANEL_PORT"
        break
    done
}

function Set_Firewall(){
    if which firewall-cmd >/dev/null 2>&1; then
        if systemctl status firewalld | grep -q "Active: active" >/dev/null 2>&1;then
            log "防火墙开放 $PANEL_PORT 端口"
            firewall-cmd --zone=public --add-port=$PANEL_PORT/tcp --permanent
            firewall-cmd --reload
        else
            log "防火墙未开启，忽略端口开放"
        fi
    fi
    if which ufw >/dev/null 2>&1; then
        if systemctl status ufw | grep -q "Active: active" >/dev/null 2>&1;then
            log "防火墙开放 $PANEL_PORT 端口"
            ufw allow $PANEL_PORT/tcp
            ufw reload
        else
            log "防火墙未开启，忽略端口开放"
        fi
    fi
}

function Init_Web(){
    VERSION=$(curl -s https://api.github.com/repos/hitosea/opsw/releases/latest | grep 'tag_name' | cut -d\" -f4)
    if [[ "x${VERSION}" == "x" ]];then
        log "获取最新版本失败，请稍候重试"
        exit 1
    fi

    log "开始下载 ${VERSION} 版本在线安装包"

    package_name=$(echo "opsw-${VERSION}-${PLATFORM}-${FRAMEWORK}" | tr '[A-Z]' '[a-z]')
    package_file_name="${package_name}.tar.gz"
    package_download_url="https://github.com/hitosea/opsw/releases/download/${VERSION}/${package_file_name}"

    log "安装包下载地址： ${package_download_url}"

    cd /opt/opsweb/tmp
    curl -LOk -o ${package_file_name} ${package_download_url}
    if [ ! -f ${package_file_name} ];then
        log "下载安装包失败，请稍候重试。"
        exit 1
    fi

    tar zxf ${package_file_name}
    if [ $? != 0 ];then
        log "下载安装包失败，请稍候重试。"
        rm -f ${package_file_name}
        exit 1
    fi

    cd ${package_name}
    mv ./opsw ./opsweb
    cp ./opsweb /usr/local/bin && chmod +x /usr/local/bin/opsweb
    if [[ ! -f /usr/bin/opsweb ]]; then
        ln -s /usr/local/bin/opsweb /usr/bin/opsweb >/dev/null 2>&1
    fi

    cat > /etc/systemd/system/opsweb.service <<-EOF
[Unit]
Description=Opsw, a modern open source linux cluster management
After=syslog.target network.target

[Service]
ExecStart=/usr/bin/opsweb --port $PANEL_PORT --cache /opt/opsweb/cache
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

    systemctl enable opsweb; systemctl daemon-reload 2>&1

    log "启动服务"
    systemctl start opsweb

    for b in {1..30}
    do
        sleep 3
        service_status=`systemctl status opsweb 2>&1 | grep Active`
        if [[ $service_status == *running* ]];then
            log "服务启动成功!"
            break;
        else
            log "服务启动出错!"
            exit 1
        fi
    done

    cd ~
    rm -rf /opt/opsweb/tmp
}

function Get_Local_Ip(){
    LOCAL_IP=$(curl -sSL4 ip.sb)
    if [[ "x${LOCAL_IP}" == "x" ]];then
        LOCAL_IP=$(curl -sSL4 ipinfo.io/ip)
    fi
    if [[ "x${LOCAL_IP}" == "x" ]];then
        LOCAL_IP=$(curl -sSL4 ifconfig.me)
    fi
    if [[ "x${LOCAL_IP}" == "x" ]];then
        LOCAL_IP="LOCAL_IP"
    fi
}

function Show_Result(){
    log ""
    log "=================感谢您的耐心等待，安装已经完成=================="
    log ""
    log "请用浏览器访问面板:"
    log "面板地址: http://$LOCAL_IP:$PANEL_PORT"
    log ""
    log "代码仓库: https://github.com/hitosea/opsw"
    log ""
    log "如果使用的是云服务器，请至安全组开放 $PANEL_PORT 端口"
    log ""
    log "================================================================"
}

function Remove_Web(){
    log "停止服务进程..."
    systemctl stop opsweb.service

    log "删除服务和数据目录..."
    rm -rf /usr/local/bin/opsweb /etc/systemd/system/opsweb.service /opt/opsweb

    log "重新加载服务配置文件..."
    systemctl daemon-reload
    systemctl reset-failed

    log ""
    log "=================卸载已完成，欢迎再次使用=================="
    log ""
    log "代码仓库: https://github.com/hitosea/opsw"
    log ""
    log "================================================================"
}

function install(){
    log "======================= 开始安装 ======================="
    Prepare_System
    Set_Dir
    Set_Port
    Set_Firewall
    Init_Web
    Get_Local_Ip
    Show_Result
}

function uninstall(){
    log "======================= 开始卸载 ======================="
    Remove_Web
}

ACTION=$1
if [[ -z "$ACTION" ]]; then
    ACTION="install"
fi
$ACTION