# Byt
Byt is file sharing service made easy

Public platform: https://wayt.ovh

## Run your own with Docker

```
docker run --name byt --restart=always -p 80:8080 maxwayt/byt:latest
```

You can also specify a data directory

```
docker run --name byt --restart=always -p 80:8080 -v /local/data:/data maxwayt/byt:latest
```

## Development


Install npm dependencies

```bash
npm install
```

Minify web js/css

```
npm install -g grunt-cli
grunt
```

Generate static file

```bash
make generate
```

Run Golang web server

```bash
   make run
```

Run tests

```bash
make test
```
