/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "grpc-example-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr_greeting = flag.String("addr_greeting", "localhost:50051", "the address to connect to")
	addr_etcd     = flag.String("addr_etcd", "localhost:50052", "the address to connect to")
	name          = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	/*
		conn, err := grpc.Dial(*addr_greeting, grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		c := pb.NewGreeterClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	*/

	connEtcd, err := grpc.Dial(*addr_etcd, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer connEtcd.Close()

	c1 := pb.NewETCDWrapperClient(connEtcd)

	ctxEtcd, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var Values string

	r1, err := c1.SetKV(ctxEtcd, &pb.SetKVRequest{Key: "greeting", Value: &pb.DummyInfo{
		N: 6,
		M: 9,
	}})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r1.GetMessage())

	r2, err := c1.GetKey(ctxEtcd, &pb.GetKeyRequest{Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %v", r2.GetValue().String())

	r3, err := c1.SetKV(ctxEtcd, &pb.SetKVRequest{Key: "greet", Value: &pb.DummyInfo{
		N: 5,
		M: 10,
	}})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r3.GetMessage())

	r4, err := c1.ListValues(ctxEtcd, &pb.ListValuesRequest{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	Values = ""
	for i := 0; i < len(r4.Values); i++ {
		Values += r4.Values[i].String() + " "
	}
	log.Printf("ETCD Server: %s", Values)

	r5, err := c1.DeleteKey(ctxEtcd, &pb.DeleteKeyRequest{Key: "greeting"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r5.GetMessage())

	r6, err := c1.ListValues(ctxEtcd, &pb.ListValuesRequest{})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	Values = ""
	for i := 0; i < len(r6.Values); i++ {
		Values += r6.Values[i].String() + " "
	}
	log.Printf("ETCD Server: %s", Values)

	r7, err := c1.DeleteKey(ctxEtcd, &pb.DeleteKeyRequest{Key: "greet"})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("ETCD Server: %s", r7.GetMessage())

}
