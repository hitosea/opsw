#!/bin/bash

CURRENT_DIR=$(
    cd "$(dirname "$0")"
    pwd
)

function log() {
    message="[Opsw Log]: $1 "
    echo -e "${message}" 2>&1 | tee -a ${CURRENT_DIR}/install.log
}

echo
cat << EOF
 ██████╗ ██████╗ ███████╗██╗    ██╗
██╔═══██╗██╔══██╗██╔════╝██║    ██║
██║   ██║██████╔╝███████╗██║ █╗ ██║
██║   ██║██╔═══╝ ╚════██║██║███╗██║
╚██████╔╝██║     ███████║╚███╔███╔╝
 ╚═════╝ ╚═╝     ╚══════╝ ╚══╝╚══╝
EOF

log "======================= 开始安装 ======================="

function Prepare_System(){
    is64bit=`getconf LONG_BIT`
    if [[ $is64bit != "64" ]]; then
        log "不支持 32 位系统安装 Opsw Linux 服务器运维管理面板，请更换 64 位系统安装"
        exit 1
    fi

    if which opsw >/dev/null 2>&1; then
        log "Opsw Linux 服务器运维管理面板已安装，请勿重复安装"
        exit 1
    fi
}

function Install_Docker(){
    if which docker >/dev/null 2>&1; then
        log "检测到 Docker 已安装，跳过安装步骤"
        log "启动 Docker "
        systemctl start docker 2>&1 | tee -a ${CURRENT_DIR}/install.log
    else
        log "... 在线安装 docker"

        curl -fsSL https://get.docker.com -o get-docker.sh 2>&1 | tee -a ${CURRENT_DIR}/install.log
        if [[ ! -f get-docker.sh ]];then
            log "docker 在线安装脚本下载失败，请稍候重试"
            exit 1
        fi
        if [[ $(curl -s ipinfo.io/country) == "CN" ]]; then
            sh get-docker.sh --mirror Aliyun 2>&1 | tee -a ${CURRENT_DIR}/install.log
        else
            sh get-docker.sh 2>&1 | tee -a ${CURRENT_DIR}/install.log
        fi
        
        log "... 启动 docker"
        systemctl enable docker; systemctl daemon-reload; systemctl start docker 2>&1 | tee -a ${CURRENT_DIR}/install.log

        docker_config_folder="/etc/docker"
        if [[ ! -d "$docker_config_folder" ]];then
            mkdir -p "$docker_config_folder"
        fi

        docker version >/dev/null 2>&1
        if [[ $? -ne 0 ]]; then
            log "docker 安装失败"
            exit 1
        else
            log "docker 安装成功"
        fi
    fi
}

function Install_Compose(){
    docker-compose version >/dev/null 2>&1
    if [[ $? -ne 0 ]]; then
        log "... 在线安装 docker-compose"

        curl -s -L "https://get.daocloud.io/docker/compose/releases/download/v2.16.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        if [[ ! -f /usr/local/bin/docker-compose ]];then
            log "docker-compose 下载失败，请稍候重试"
            exit 1
        fi
        chmod +x /usr/local/bin/docker-compose
        ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

        docker-compose version >/dev/null 2>&1
        if [[ $? -ne 0 ]]; then
            log "docker-compose 安装失败"
            exit 1
        else
            log "docker-compose 安装成功"
        fi
    else
        compose_v=`docker-compose -v`
        if [[ $compose_v =~ 'docker-compose' ]];then
            log "Docker Compose 版本为 $compose_v，可能会影响应用商店的正常使用"
        else
            log "检测到 Docker Compose 已安装，跳过安装步骤"
        fi
    fi
}

function Init_Service(){
    log "配置 Opsw Service"

    cd ${CURRENT_DIR}

    mkdir -p /etc/config/
    cp ./opsw.yaml /etc/config/opsw.yaml

    cp ./opsw /usr/local/bin && chmod +x /usr/local/bin/opsw
    if [[ ! -f /usr/bin/opsw ]]; then
        ln -s /usr/local/bin/opsw /usr/bin/opsw >/dev/null 2>&1
    fi

    cp ./opsw.service /etc/systemd/system

    systemctl enable opsw; systemctl daemon-reload 2>&1 | tee -a ${CURRENT_DIR}/install.log

    log "启动 Opsw 服务"
    systemctl start opsw | tee -a ${CURRENT_DIR}/install.log

    for b in {1..30}
    do
        sleep 3
        service_status=`systemctl status opsw 2>&1 | grep Active`
        if [[ $service_status == *running* ]];then
            log "Opsw 服务启动成功!"
            break;
        else
            log "Opsw 服务启动出错!"
            exit 1
        fi
    done
}

function main(){
    Prepare_System
    Install_Docker
    Install_Compose
    Init_Service
}
main
