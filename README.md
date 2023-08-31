# Test Task
# Start
```
docker-compose up
```

# Test commands
## Set
```
curl -i -L -H "Content-Type: application/json" -d '{"key":"mykey","value":"myvalue"}' 127.0.0.1:8089/set_key
```
## Get
```
curl -i -L "http://127.0.0.1:8080/get_key?key=mykey"
```
## Delete
```
curl -i -L -X DELETE -H "Content-Type: application/json" -d '{"key":"mykey"}' 127.0.0.1:8089/del_key
```