package main

import (
	"fmt"
	"time"
)

// The first message is the SYN, the second message is the ACK
func main() {
	var pipe = make(chan int)

	go Client(pipe)
	go Server(pipe)

	for { // Application runs indefinitely
		time.Sleep(10 * time.Second)
	}
}

func Client(pipe chan int) {
	syn := 100 // syn is 100 for now
	state := 0 // 0 = step 1 handshake , 1 = step 2 handshake, 3 = connected
	for {
		if state == 0 {
			pipe <- syn 
			pipe <- 0 // ack is null
			state++
			fmt.Println("Client (step 1) send success, syn:", syn)
		} else if state == 1 {
			seq := <-pipe
			ack := <-pipe
			fmt.Println("Client (step 2) received, syn:", seq, "ack:", ack)
			if ack == syn+1 { // correct handshake received, respond
				pipe <- syn + 1
				pipe <- seq + 1
				state++
				fmt.Println("Client (step 2) send success, syn:", syn+1, " ack:", seq+1)
			} else {
				fmt.Println("Client step 2 failed")
			}
		} else {
			fmt.Println("Client connected")
			time.Sleep(30 * time.Second)
		}
	}
}

func Server(pipe chan int) {
	syn := 0
	ackx := 0
	state := 0 // 0 = step 1 handshake , 1 = connected
	for {
		if state == 0 {
			seq := <-pipe
			ack := <-pipe
			fmt.Println("Server (step 1) received, syn:", seq, "ack:", ack)

			syn = (ack + 100) * 3 // random for now
			ackx = seq + 1
			pipe <- syn
			pipe <- ackx
			state++
			fmt.Println("Server (step 1) send success, syn:", syn, "ack:", ackx)
		} else if state == 1 {
			seq := <-pipe
			ack := <-pipe
			fmt.Println("Server (step 2) received, syn:", seq, "ack:", ack)
			if seq == ackx && ack == syn+1 {
				state++
				fmt.Println("Server step 2 success")
			}
		} else {
			fmt.Println("Server connected")
			time.Sleep(30 * time.Second)
		}
	}
}
