.PHONY: build clean install

# Compila para el sistema actual
build:
	go build -o bin/env-vault main.go

# Instala globalmente en el sistema (requiere que $GOPATH/bin esté en el PATH o copia a /usr/local/bin)
install:
	go install

# Compila para todos los sistemas de tu equipo
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/env-vault-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/env-vault-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/env-vault-mac-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/env-vault-mac-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/env-vault-windows-amd64.exe main.go

clean:
	rm -rf bin/
