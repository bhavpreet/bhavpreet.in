# GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
clean:
	rm -rf bin/

all: clean
	go install web

heroku: all
	heroku container:push web
