#!/bin/bash

echo -e "================== 开始卸载 Opsw 服务器运维管理面板 =================="

echo -e "1) 停止 Opsw 服务进程..."
systemctl stop opsw.service

echo -e "2) 删除 Opsw 服务和数据目录..."
rm -rf /usr/local/bin/opsw /etc/config/opsw.yaml /etc/systemd/system/opsw.service

echo -e "3) 重新加载服务配置文件..."
systemctl daemon-reload
systemctl reset-failed

echo -e "================================== 卸载完成 =================================="