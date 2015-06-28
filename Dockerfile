FROM ubuntu:14.04

MAINTAINER JamesClonk

EXPOSE 4007

RUN apt-get update
RUN apt-get install -y ca-certificates

COPY moviedb-backend /moviedb-backend
COPY migrations /migrations

ENV JCIO_ENV production
ENV PORT 4007
ENV JCIO_DATABASE_TYPE postgres

CMD ["/moviedb-backend"]
