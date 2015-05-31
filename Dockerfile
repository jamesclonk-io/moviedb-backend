FROM ubuntu:14.04

EXPOSE 4007

ADD moviedb-backend /moviedb-backend

CMD ["/moviedb-backend"]
