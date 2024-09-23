package main

import (
	"fmt"
	"net"
	"time"
)

// The first message is the SYN, the second message is the ACK
func main() {
	go netServer()
	time.Sleep(1 * time.Second) // Wait for the server to start
	go netClient()

	for { // Application runs indefinitely
		time.Sleep(10 * time.Second)
	}
}

func netClient() {
	pipeClient, err := net.Dial("tcp", "localhost:8080")
	printError(err)
	defer pipeClient.Close()
	fmt.Println("Connected to server")

	syn := 100 // syn is 100 for now
	state := 0 // 0 = step 1 handshake , 1 = step 2 handshake, 3 = connected
	for {
		if state == 0 {
			// Simulate sending SYN-like message
			arr := []byte{byte(syn), 0} // ack is null
			_, err = pipeClient.Write(arr)
			printError(err)
			fmt.Println("Client (step 1) send success, syn:", syn)
			state++
		} else if state == 1 {
			buf := make([]byte, 2)
			_, err = pipeClient.Read(buf)
			printError(err)
			seq := int(buf[0])
			ack := int(buf[1])

			fmt.Println("Client (step 2) received, syn:", seq, "ack:", ack)
			if ack == syn+1 { // correct handshake received, respond
				arr := []byte{byte(syn + 1), byte(seq + 1)} // ack is null
				_, err = pipeClient.Write(arr)
				printError(err)
				fmt.Println("Client (step 2) send success, syn:", syn+1, " ack:", seq+1)
				state++
			} else {
				fmt.Println("Client step 2 failed")
			}
		} else {
			fmt.Println("Client connected")
			time.Sleep(30 * time.Second)
		}
	}

	// Read server's response (SYN-ACK equivalent)
	/*buf := make([]byte, 1024)
	n, err := pipeClient.Read(buf)
	printError(err)
	fmt.Printf("Received from server: %s\n", string(buf[:n]))*/
}

func netServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	printError(err)
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		pipeServer, err := listener.Accept()
		printError(err)
		fmt.Println("Client connected:", pipeServer.RemoteAddr())

		// Simulate response to client and close the connection
		// Simulate handshake by receiving and sending data
		defer pipeServer.Close()

		syn := 0
		ackx := 0
		state := 0 // 0 = step 1 handshake , 1 = connected
		for {
			if state == 0 {
				buf := make([]byte, 2)
				_, err = pipeServer.Read(buf)
				printError(err)
				seq := int(buf[0])
				ack := int(buf[1])
				fmt.Println("Server (step 1) received, syn:", seq, "ack:", ack)

				syn = 3 //(ack + 100) * 3 // random for now
				ackx = seq + 1
				arr := []byte{byte(syn), byte(ackx)} // ack is null
				_, err = pipeServer.Write(arr)
				printError(err)
				fmt.Println("Server (step 1) send success, syn:", syn, "ack:", ackx)
				state++
			} else if state == 1 {
				time.Sleep(1 * time.Second) // Wait for the server to start
				buf := make([]byte, 2)
				_, err = pipeServer.Read(buf)
				printError(err)
				seq := int(buf[0])
				ack := int(buf[1])
				fmt.Println("Server (step 2) received, syn:", seq, "ack:", ack)
				fmt.Println(syn, ackx)
				if seq == ackx && ack == syn+1 {
					fmt.Println("Server step 2 success")
					state++
				}
			} else {
				fmt.Println("Server connected")
				time.Sleep(30 * time.Second)
			}
		}
		/*'

		buf := make([]byte, 2)
		_, err = conn.Read(buf)
		printError(err)
		fmt.Printf("Received from client: ", buf[0], buf[1])

		// Simulate server response
		response := "ACK from server"
		_, err = conn.Write([]byte(response))
		printError(err)
		fmt.Println("Sent ACK to client")*/
	}
}
func printError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
