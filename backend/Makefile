
SHELL := cmd.exe
.PHONY: build

build:
	@if not exist ..\bin mkdir ..\bin
	@del /Q ..\bin\*
	@copy .env ..\bin\.env
	go build -o ..\bin\app.exe main.go

dev:
	@if not exist tmp mkdir tmp
	@copy .env .\tmp\.env
	air

setup:
	go mod tidy
