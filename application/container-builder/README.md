





# how to build a container image
Ensure the docker daemon is running on your machine. Then run the following command from the root of the project directory:
```
export ENVIRONMENT=dev && go run main.go
```

you can now access the grpc server on port 4040 


## Creating a build 
Ensure you have grpcurl installed on your machine. Then run the following command to create a build:
```
grpcurl -plaintext -d '{"repoURL": "https://github.com/s1ntaxe770r/pong-server", "type": "go", "customerID": "jjj"}' localhost:4040 container.ContainerService/CreateContainer
``` 
sample response:
```
{
  "buildID": "b1b2b3b4b5b6b7b8b9b0"
}
```

it should respond with a build ID. You can then check the status of the build by running the following command:
```
grpcurl -plaintext -d '{"buildID": "b1b2b3b4b5b6b7b8b9b0"}' localhost:4040 container.ContainerService/ImageStatus
```
```
sample response:
{
  "status": "building"
}
```




