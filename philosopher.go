package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
Explanation:
The code works by having each Fork have a single channel for communication via integers.
The Philosopher have channels to it's surrounding left and right forks.

A Philosopher starts by trying to fetch it's left side Fork. If it fails then it starts thinking for a random duration of time.
If the Philosopher gets the left side Fork it then tries to fetch it's right side Fork.
If it fails the Philosopher releases both forks then starts thinking for a random duration of time.
If the Philosopher gets both forks it starts eating for a random amount of time then releases the forks and starts thinking for a random duration of time.
This process repeats indefinitely.

Deadlocks are prevented by having the Philosopher think for some time if it fails to grab both forks and right after eating.
During the thinking time this frees up the forks for other Philosophers to grab which prevents deadlocks.
By using a random amount of time it guarantees every Philosopher will eventually eat at least 3 times.

The Fork keeps track of whether it is occupied or not by a Philosopher
The Fork can receive commands and respond based on the specification below:

Commands sent to Fork:
0 = fetch, 1 = release
Responses sent to Philosopher:
0 = failed, 1 = success
*/
func main() {
	// Create 5 channels to communicate between Fork and Philosopher
	var f1, f2, f3, f4, f5 = make(chan int), make(chan int), make(chan int), make(chan int), make(chan int)

	// Create 5 Fork GO routines
	go Fork(f1)
	go Fork(f2)
	go Fork(f3)
	go Fork(f4)
	go Fork(f5)

	// Create 5 Philosopher GO routines. Each Philosopher gets it's left and right fork channels
	go Philosopher(1, f5, f1)
	go Philosopher(2, f1, f2)
	go Philosopher(3, f2, f3)
	go Philosopher(4, f3, f4)
	go Philosopher(5, f4, f5)

	for { // So application runs indefinitely
		time.Sleep(10 * time.Second)
	}
}

// PrintAndWait
/*
Helper function for philosopher to wait randomly, manage and print the Philosopher current status.
This method is used to ensure the Philosopher status is printed only when the status actually changes instead of in every iteration of the Philosopher loop
shouldThink - What the new status should be, True = thinking, False = eating
isThinking - Pointer Philosopher's WasThinking bool that represents what the old status was
*/
func PrintAndWait(shouldThink bool, isThinking *bool, delay int, id int) {
	if *isThinking != shouldThink { // If the new status is different from the current status print the change
		*isThinking = shouldThink // Change status of the Philosopher's WasThinking bool to current status
		if shouldThink {          // Print the status
			fmt.Println("Philosopher", id, "is thinking")
		} else {
			fmt.Println("Philosopher", id, "is eating")
		}
	}
	time.Sleep(time.Duration(rand.Intn(delay)) * time.Second) // Wait for a random amount of time
}

func Philosopher(id int, rightFork chan int, leftFork chan int) {
	wasThinking := false // Status of what Philosopher was doing before current iteration
	delay := 10          // Random delay multiplier in seconds
	for {
		leftFork <- 0 // Try to fetch left fork
		l := <-leftFork
		if l == 1 { // Fork fetched successfully
			rightFork <- 0 // Try to fetch right fork
			r := <-rightFork
			if r == 1 { // Philosopher has both forks
				PrintAndWait(false, &wasThinking, delay, id) // Start eating for some time
				leftFork <- 1                                // Release left fork
				rightFork <- 1                               // Release right fork
				PrintAndWait(true, &wasThinking, delay, id)  // Start thinking for some time to give other's chance to eat
			} else { // Failed to get both forks
				leftFork <- 1                               // Release left fork
				PrintAndWait(true, &wasThinking, delay, id) // Start thinking for some time
			}
		} else { // Failed to get left fork
			PrintAndWait(true, &wasThinking, delay, id) // Start thinking for some time
		}
	}
}

func Fork(c chan int) {
	isTaken := false // Status for whether fork is currently taken by a philosopher
	for {
		arg := <-c    // Wait for philosopher commands: 0 = fetch request, 1 = release fork
		if arg == 0 { // Request to fetch fork, check if it is available and respond
			if isTaken {
				c <- 0 // Return failed
			} else {
				isTaken = true
				c <- 1 // Return success
			}
		} else if arg == 1 { // Fork is released
			isTaken = false
		}
	}
}
