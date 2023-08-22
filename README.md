# pii-db-scanner

## Usage

```sh
export SCANNER_USERNAME=fizz
export SCANNER_PASSWORD=buzz
export SCANNER_HOST=hostage
export SCANNER_DBTYPE=postgres
export SCANNER_DATABASE=cosmos
export SCANNER_PORT=5432
go run cmd/scannercli/main.go 
```

```sh
export SCANNER_USERNAME=fizz
export SCANNER_PASSWORD=buzz
export SCANNER_HOST=hostage
export SCANNER_DBTYPE=postgres
export SCANNER_DATABASE=cosmos
export SCANNER_PORT=5432
go run cmd/scannercli/main.go --host=hostage --username=fizz --password=buzz --database=cmos --dbtype=mysql --port=3306
```

```sh
docker build .
docker run -it <By setting all env>
```