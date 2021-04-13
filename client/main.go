package main

import (
	"bufio"
	"CPS831-Final/proto"

	"flag"
	"fmt"
	"os"

	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {

	var streamerror error

	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User: user,
		Active: true,
	})
	
	if err != nil {
		return fmt.Errorf("Connection failed: %v", err)
	}

	wait.Add(1)
	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()

		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error readinf message: %v", err)
				break
			}

			fmt.Printf("%v : %s\n", msg.Id, msg.Content)
			fmt.Printf("WHAT SO YOU WANT TO SAY: ")
		}
	}(stream)
	
	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make (chan int)


	username := ""
	fmt.Printf("WHAT IS YOUR NAME NEW PERSON??: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		username = scanner.Text()
		break
	}
 

	name := flag.String("N", " " + username, "The name of the user")
	flag.Parse()
	// id not sure if it works (testing for username)
	id := timestamp.String() + *name
	
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to service: %v", err)
	}
	
	client = proto.NewBroadcastClient(conn)
	user := &proto.User{
		Id: id,
		Name: *name,
	}

	connect(user)
	wait.Add(1)

	go func() {
		defer wait.Done()
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("WHAT SO YOU WANT TO SAY: ")
		for scanner.Scan() {
			
			msg := &proto.Message{
				Id: user.Id,
				Content: scanner.Text(),
				Timestamp: timestamp.String(),
			}
			_, err := client.BroadcastMesssage(context.Background(),msg)
			if err != nil{
				fmt.Printf("Error Sending Message: %v", err)
				break
			}
		}
	}()
	
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
}
