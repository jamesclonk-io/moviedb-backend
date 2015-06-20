#!/bin/bash

export JCIO_DATABASE_TYPE=postgres
export JCIO_DATABASE_URI=postgres://secret:supersecret@localhost:5432/secret?sslmode=disable

docker run --name moviedb-postgres -p 5432:5432 -e POSTGRES_USER=secret -e POSTGRES_PASSWORD=supersecret -d postgres
