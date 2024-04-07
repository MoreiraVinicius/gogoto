.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/createToDestinationUrl functions/createShortenedUrl/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/redirectToDestinationUrl.exe functions/redirectToDestinationUrl/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/deleteToDestinationUrl.exe functions/deleteShortenedUrl/main.go
	

build-win: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/createShortenedUrl.exe functions/createShortenedUrl/main.go
	env GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/redirectToDestinationUrl.exe functions/redirectToDestinationUrl/main.go
	env GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/deleteShortenedUrl.exe functions/deleteShortenedUrl/main.go


clean:
	rm -rf ./bin ./vendor *.exe

# deploy: build
# 	chmod u+x setEnv.sh
# 	./setEnv.sh
# 	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
