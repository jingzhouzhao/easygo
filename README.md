# EasyGo
This is a simple HTTP file server that uploads, downloads, and deletes files via the RESTful API.

### Getting Started
We need to do some simple configuration:
```text
nodes = ["10.100.210.1:26655","10.100.210.2:26655","10.100.210.3:26655"]
[log]
    logDir = "/easygo/logs"
[http]
    [http.ports]
            httpPort = 25566
            syncPort = 26655
[data]
    dataDir = "/backup"
```

It is recommended that each node be configured the same.

Start the server:
```shell script
start.sh
```
If you see `EasyGo server listening at: 25566`, it means the server is successfully started.

Let's enjoy it!

#### Upload:
```shell script
curl -X POST localhost:25566/files -F "file=@/users/filename.jpg"
```

#### Download:
```shell script
curl -O localhost:25566/files/e4d547d4ab854efa9d1f30c2abf96a03cabjajbgbbei
```

#### Delete:
```shell script
curl -X DELETE localhost:25566/files/e4d547d4ab854efa9d1f30c2abf96a03cabjajbgbbei
```