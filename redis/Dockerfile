FROM redis:latest
COPY ./redis.conf  /usr/local/etc/redis/redis.conf
COPY ./certs/redis.crt /usr/local/etc/redis/redis.crt
COPY ./certs/redis.key /usr/local/etc/redis/redis.key
COPY ./certs/ca.crt /usr/local/etc/redis/ca.crt
RUN chown -R redis:redis /usr/local/etc/redis
USER root
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf"]
