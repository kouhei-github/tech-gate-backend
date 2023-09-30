FROM nginx:1.19-alpine

ENV TZ Asia/Tokyo

COPY ./docker/nginx/default.conf /etc/nginx/conf.d/default.conf

COPY ./src /var/www/html

