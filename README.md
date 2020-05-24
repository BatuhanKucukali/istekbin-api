# Istekbin - Inspect HTTP Requests

Istekbin is a free service that allows you to collect http request.

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
go build
```
4 . Run
```bash
./istekbin
```

## Swagger
Regenerate swagger.

1 . Install [swag](https://github.com/swaggo/swag)
```bash
go get -u github.com/swaggo/swag/cmd/swag
```

2 . Generate swagger docs
```bash
swag init --parseDependency
```


## Getting help ##

If you're having trouble getting this project running, feel free to [open an issue](https://github.com/BatuhanKucukali/istekbin-api/issues/new)



