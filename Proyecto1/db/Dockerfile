FROM mysql:latest

COPY ./init.sql /tmp

EXPOSE 3306

CMD [ "mysqld", "--init-file=/tmp/init.sql" ]