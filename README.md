# Kibana log viewer for Golang application (example)

## How to see in work?

### 1. Run docker container
```sh
docker-compose up
```

### 2. Run golang application which writes a few log records in JSON format
```sh
go run .
```
output:
```
{"level":"debug","module":"world","msg":"Hello!!!","time":"2019-10-01T17:09:19+03:00"}
{"level":"info","module":"updater","msg":"updating info...","time":"2019-10-01T17:09:19+03:00"}
{"level":"error","module":"user","msg":"failed to fetch user info","time":"2019-10-01T17:09:19+03:00"}
```

### 3. see logs on [http://localhost:5601](http://localhost:5601)

![](https://lh4.googleusercontent.com/xFuYni1J91lkSFCdiNzi8FNiMCxPNkh2T4_ipiXlK-eqoUAHG1etBbq09QTbVC6QTydMt2jv-R-gLte9O9FSwzNPJ3vjZX5mUot1DUASAx5aFwi6a3chTXYv1FHReSKvgqp5YJtE)
