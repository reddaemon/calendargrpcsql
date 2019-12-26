package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/labstack/gommon/log"
	"github.com/reddaemon/calendargrpcsql/protofiles"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

var createResponse *api.CreateResponse
var updateResponse *api.UpdateResponse

//var readResponse api.UpdateResponse
//var deleteResponse api.UpdateResponse

func main() {
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := api.NewEventServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	inCreate := &api.CreateRequest{
		Event: &api.Event{
			Title:       "Test",
			Description: "testing",
			Date: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}

	inUpdate := &api.UpdateRequest{
		Id: 21,
		Event: &api.Event{
			Title:       "Test1",
			Description: "Behavior testing1",
			Date: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}
	inRead := api.ReadRequest{
		Id: 21,
	}
	inDelete := api.DeleteRequest{
		Id: 20,
	}
	createResponse, err = c.Create(ctx, inCreate)
	if err != nil {
		fmt.Println(createResponse, err)
	}

	fmt.Printf("%v\n", createResponse)
	readResponse, err := c.Read(ctx, &inRead)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", readResponse)
	updateResponse, err = c.Update(ctx, inUpdate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", updateResponse)
	readResponse, err = c.Read(ctx, &inRead)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", readResponse)
	deleteResponse, err := c.Delete(ctx, &inDelete)
	if err != nil {
		fmt.Println("DELETEERROR: ", err)
	}
	fmt.Printf("%v\n", deleteResponse)
}
