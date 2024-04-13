.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/createToDestinationUrl src/functions/createShortenedUrl/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/redirectToDestinationUrl src/functions/redirectToDestinationUrl/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/deleteToDestinationUrl src/functions/deleteShortenedUrl/main.go
	

build-win: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/createShortenedUrl.exe src/functions/createShortenedUrl/main.go
	env GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/redirectToDestinationUrl.exe src/functions/redirectToDestinationUrl/main.go
	env GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/deleteShortenedUrl.exe src/functions/deleteShortenedUrl/main.go


clean:
	rm -rf ./bin ./vendor *.exe

# deploy: build
# 	chmod u+x setEnv.sh
# 	./setEnv.sh
# 	sls deploy --verbose

oidc:
	chmod u+x create_oidc_service_principal.sh
	./create_oidc_service_principal.sh

rbac:
	chmod u+x create_rbac_contributor_scoped.sh
	./create_rbac_contributor_scoped.sh

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
