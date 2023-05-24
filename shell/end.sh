#!/bin/bash

function olog() {
    echo -e "[Opsw Log]: $1"
}

function plog() {
    echo -e "[Opsw Log] [Panel]: $1"
}

olog "开始卸载"

olog "停止服务进程..."
systemctl stop opsw.service

olog "删除服务和数据目录..."
rm -rf /usr/local/bin/opsw /etc/systemd/system/opsw.service /etc/config/opsw.yaml

plog "停止服务进程..."
systemctl stop opspanel.service

plog "删除服务和数据目录..."
rm -rf /usr/local/bin/opspanel /etc/systemd/system/opspanel.service /usr/bin/manage-panel /opt/manage-panel

olog "重新加载服务配置文件..."
systemctl daemon-reload
systemctl reset-failed

olog "卸载完成"