#!/bin/bash

function log() {
    echo -e "[Opsw Log]: $1"
}

log "开始卸载 Opsw 服务器运维管理面板"

log "1) 停止 Opsw 服务进程..."
systemctl stop opsw.service

log "2) 删除 Opsw 服务和数据目录..."
rm -rf /usr/local/bin/opsw /etc/config/opsw.yaml /etc/systemd/system/opsw.service

log "3) 重新加载服务配置文件..."
systemctl daemon-reload
systemctl reset-failed

log "卸载完成"