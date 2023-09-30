#!/bin/sh
source ./.env

if [ $DEBUG = "True" ]
then
    docker compose up -d
else
    # gunicornを起動させる時はプロジェクト名を指定します
    # 今回はconfigにします [fast_api-start admin configみたいなので打ったやつ]
    docker compose -f docker-compose-prod.yml up -d
fi
