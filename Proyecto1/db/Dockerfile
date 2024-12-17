FROM mysql:latest

COPY ./init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306

# The official MySQL image already has a default CMD and ENTRYPOINT, 
# so you typically do not need to specify them again.