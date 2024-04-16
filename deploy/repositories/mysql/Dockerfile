FROM mysql:8.3.0

ENV MYSQL_ROOT_PASSWORD "root"
ENV MYSQL_USER "real_state_admin"
ENV MYSQL_PASSWORD "real_state_pass"
ENV MYSQL_DATABASE "real_states"

RUN rm -rf /etc/localtime
RUN ln -s /usr/share/zoneinfo/UTC /etc/localtime

COPY ./deploy/repositories/mysql/ /docker-entrypoint-initdb.d/.