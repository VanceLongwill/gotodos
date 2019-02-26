FROM postgres
ENV POSTGRES_USER gotodos
ENV POSTGRES_PASSWORD gotodos
ENV POSTGRES_DB gotodos
COPY init.sql /docker-entrypoint-initdb.d/10-init.sql
EXPOSE 80
