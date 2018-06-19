# grpc-resolver

This is a fork from https://github.com/eddyzags/resolver including some bug
fixes and updates.

## Installation

Install the resolver using the "go get" command:

`go get github.com/GFG/grpc-resolver`

Import the library into a project:

`import "github.com/GFG/grpc-resolver"`

## Usage

Resolver uses Marathon labels feature in order to identify resident services in
the cluster. The label is composed of the service unique name and the port
index. Those two informations will allow the resolver to identify the
tasks ip addresses and the port on which the grpc client should
establish a connection.  
Let's start with a simple app definition:

```json
{
  "id": "my-app",
  "cpus": 0.1,
  "mem": 64,
  "container": {
    "type": "DOCKER",
    "docker": {
      "image": "eddyzags/healthy:latest",
      "network": "BRIDGE"
    },
    "portMappings": [
      {
        "containerPort": 80,
        "hostPort": 0
      },
      {
        "containerPort": 4242,
        "hostPort": 0
      }
    ]
  },
  "labels": {
    "RESOLVER_0_NAME": "my-app-service"
  }
}
```

We have just defined an application called `my-app` with a service resolver name
of `my-app-service` which points to the port index 0.

A service resolver name can be defined using a labels map:

`"RESOLVER_{PORTINDEX}_NAME": "{NAME}"`

Once we deployed the application in Marathon, the service can be discovered through its name in the grpc client instantiation.

```golang
package main

import (
       "log"

       "github.com/GFG/grpc-resolver"

       "google.golang.org/grpc"
       pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
     resolver, err := resolver.New("marathon.mesos:8080")
     if err != nil {
        log.Fatalf("couldn't instantiate resolver: %v", err)
     }

     b := grpc.RoundRobin(resolver)

     conn, err := grpc.Dial("my-app-service", grpc.WithBalancer(b))
     if err != nil {
        log.Fatalf("couldn't dial grpc server: %v", err)
     }
     defer conn.Close()

     c := pb.NewGreeterClient(conn)

     r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "user1"})
     if err != nil {
         log.Fatalf("couldn't send say hello request: %v", err)
     }

     log.Printf("Response: %s\n", r.Message)
}
```
