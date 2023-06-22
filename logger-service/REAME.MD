# Logger Service
----------------
* Serve as logger writer service to database
* Database used mongoDb (no sql)
* mongo images (https://hub.docker.com/_/mongo)
* connceting mongoDB using mongosh (shell ) or mongoDBcompass (GUI)  - https://www.mongodb.com/try/download/compass

# RPC & GRPC in golang
----------------------
* if we use RPC each server must written in go-lang
* instead in gRPC support multi-lang
* gRPC stand for Google RPC (https://grpc.io/)
* gRPC using protofile
* Advance gRPC usage in project simple bank
* gRPC plugin for golang
````console
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
````
* Protoc(ProtoCompiler) helper for compeling protoFile (https://grpc.io/docs/protoc-installation/)