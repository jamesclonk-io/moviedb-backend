TAG?=latest

all: moviedb-backend
	docker build -t jamesclonk/moviedb-backend:${TAG} .
	rm moviedb-backend

moviedb-backend: main.go
	GOARCH=amd64 GOOS=linux go build -o moviedb-backend
