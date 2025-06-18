FROM mysql:latest

RUN chown -R mysql:root /var/lib/mysql

ARG MYSQL_ROOT_PASSWORD
ARG MYSQL_DATABASE_NAME

ENV MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
ENV MYSQL_DATABASE=${MYSQL_DATABASE_NAME}

COPY init.sql /etc/mysql/init.sql

RUN sed -i 's/MYSQL_DATABASE/'$MYSQL_DATABASE'/g' /etc/mysql/init.sql
RUN cp /etc/mysql/init.sql /docker-entrypoint-initdb.d
