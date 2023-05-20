#!/bin/bash

DIR="./release"

for FILE in ${DIR}/*
do
    if [ -e "${FILE}" ]; then
        NAME=$(basename "${FILE}")
        mkdir -p ${NAME}
        cp install/* ${NAME}
        cp release/${NAME} ${NAME}/opsw
        tar zcf ${NAME}.tar.gz ${NAME}
        rm -rf ${NAME}
    fi
done
