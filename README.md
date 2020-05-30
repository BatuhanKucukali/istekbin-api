# Istekbin - Inspect HTTP Requests

Istekbin is a free service that allows you to collect http request.

[![GitHub license](https://img.shields.io/github/license/BatuhanKucukali/istekbin-api)](https://github.com/BatuhanKucukali/istekbin-api/blob/master/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/BatuhanKucukali/istekbin-api)](https://github.com/BatuhanKucukali/istekbin-api/issues)
[![Build Status](https://travis-ci.org/BatuhanKucukali/istekbin-api.svg?branch=master)](https://travis-ci.org/BatuhanKucukali/istekbin-api)
[![Coverage Status](https://coveralls.io/repos/github/BatuhanKucukali/istekbin-api/badge.svg?branch=master)](https://coveralls.io/github/BatuhanKucukali/istekbin-api?branch=master)

## Clients
[Istekbin Web](https://github.com/BatuhanKucukali/istekbin-frontend)    

## API Documentation
[Click to see Swagger](https://api.istekbin.com/swagger/index.html)   

## Run this project

1 . Clone project on your machine
```bash
git clone git@github.com:BatuhanKucukali/istekbin-api.git
```
2 . Change directory
```bash
cd istekbin-api
```
3 . Build
```bash
go build ./cmd/istekbin-api
```
4 . Run
```bash
./istekbin-api
```

### Run Test ###
```bash
go test ./internal/api/
```

## Swagger
Regenerate swagger.

1 . Install [swag](https://github.com/swaggo/swag)
```bash
go get -u github.com/swaggo/swag/cmd/swag
```
If command not working go [here](https://github.com/swaggo/swag/issues/97#issuecomment-543134010).

2 . Generate swagger docs
```bash
swag init --parseDependency -dir cmd/istekbin-api/
```

## Docker Compose
```bash
docker-compose up --build
```

## Getting help ##

If you're having trouble getting this project running, feel free to [open an issue](https://github.com/BatuhanKucukali/istekbin-api/issues/new)



