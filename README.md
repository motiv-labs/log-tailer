# Application Purpose
This application is to be used to tail logs for the logger service. 

It can be used to get logs from the logger service, but it is intended to allow for tailing logs. 
It will make the REST call to the service and tail those logs by repeatedly calling the given endpoint. 

Logs retrieved will be delayed approximately 10 seconds.  
# How to use
## Docker container
You can either get the Docker image from here: https://hub.docker.com/repository/docker/motivlabs/log-tailer

Or build it using a command like this: 
```
docker build -f docker-tools/Dockerfile -t log-tailer:local .
```

Run the docker image with a command like this:

```
docker run --rm=true --network host --name log-tailer log-tailer:test <log tailer arguments>
``` 

## Arguments

### Required

* --url <url path>
    * Pass the full url path for the logger service and the endpoint you wish to use
    * i.e. http://localhost:8080/private/logger/pullLogs 

### Optional
* --jwt \<jwt>
    * Pass the JWT to be used with the REST request. 
    * If a request is made and a 401 status code is returned, the tailer will wait and ask for a new jwt.  
* --follow
    * Set this to tail the logs from the provided endpoint
* --service \<service name>
    * Pass the name of the service you want to follow when using the `pullLogsByService` endpoint
* --transactionid \<transaction ID> 
    * Pass the transaction ID of the transaction you want to follow when using the `pullLogsByTransaction` endpoint
* --starttime \<timestamp or timeuuid>
    * Pass the timestamp for the starting time to pull logs from that time forward 
    * Default is 2 minutes ago
    * Timestamp must follow this layout: `yyyy-MM-dd'T'HH:mm:ss.SSS'Z'` i.e. 2006-01-02T15:04:05.000Z
    * If a timeUUID is passed in the logs will return all logs after it. If timestamp, it will return logs matching that timestamp foward.
* --endtime \<timestamp or timeuuid>
    * Pass the timestamp for the ending time that logs should be pulled until.
    * --follow will continue to pull logs after this timestamp. 
    * Timestamp must follow this layout: `yyyy-MM-dd'T'HH:mm:ss.SSS'Z'` i.e. 2006-01-02T15:04:05.000Z
* --showtimeuuid
    * Set this to show timeUUID for each log. 