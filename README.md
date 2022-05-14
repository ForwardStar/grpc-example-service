## Install Prerequisites
Install Go and etcd.

## Initialize etcd Server and Service
```sh
etcd
```

Open another terminal and run the scripts:
```sh
go run server/server.go
```

## Run Client
```sh
go run client/client.go
```