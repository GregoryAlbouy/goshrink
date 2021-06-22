FROM mysql:8

COPY setup.sql /docker-entrypoint-initdb.d/
VOLUME ["./volume"]
