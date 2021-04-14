package main

import (
	"bufio"
	"CPS831-Final/service"

	"strings"
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
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}
			if msg.Content == "Connected New User" {
				temp := strings.Split(msg.Id, "\n")
				username := temp[1]
				fmt.Printf("\n" + username + " Connected\n")
				fmt.Printf("\nWhat do you want to say (type \"exit\" to leave chat): ")
				continue
			} 
			if msg.Content == "exit" {
				temp := strings.Split(msg.Id, "\n")
				username := temp[1]
				fmt.Printf("\n" + username + " Disconnected\n")
				fmt.Printf("\nWhat do you want to say (type \"exit\" to leave chat): ")
				continue
			} 
			fmt.Printf("\n%v: %s\n", msg.Id, msg.Content)
			fmt.Printf("\nWhat do you want to say (type \"exit\" to leave chat): ")
		}
	}(stream)
	
	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make (chan int)

	username := ""
	fmt.Printf("What is your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		username = scanner.Text()
		break
	}
 
	// id not sure if it works (testing for username)
	id := timestamp.Format("01-02-2006 15:04:05") + "\n" + username
	
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect to service: %v", err)
	}
	
	client = proto.NewBroadcastClient(conn)
	user := &proto.User{
		Id: id,
		Name: username,
	}

	connect(user)
	wait.Add(1)

	connected_user := &proto.Message{
		Id: user.Id,
		Content: "Connected New User",
		Timestamp: timestamp.Format("01-02-2006 15:04:05"),
	}
	_, errorTemp := client.BroadcastMesssage(context.Background(),connected_user)
			if errorTemp != nil{
				fmt.Printf("Error Sending Message: %v", errorTemp)
			}

	go func() {
		defer wait.Done()
		scanner := bufio.NewReader(os.Stdin)
		fmt.Printf("\nWhat do you want to say (type \"exit\" to leave chat): ")
		for (scanner.Scan()) {
			msg := &proto.Message{
				Id: user.Id,
				Content: scanner.Text(),
				Timestamp: timestamp.Format("01-02-2006 15:04:05"),
			}
			if msg.Content == "exit" {
				_, err := client.BroadcastMesssage(context.Background(),msg)
				if err != nil{
					fmt.Printf("Error Sending Message: %v", err)
					break
				}
				os.Exit(3)
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
