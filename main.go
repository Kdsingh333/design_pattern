package main

import (
	"fmt"
	"sync"
)

var (
	noOfRoutines = 2
	wg           sync.WaitGroup
	myChannel    = make(chan string)
	lock         sync.Mutex
	condition    = sync.NewCond(&lock)
	flag         = true
)

func printA() {
	defer wg.Done()

	lock.Lock()

	for !flag {
		condition.Wait()
	}

	myChannel <- "Inside A1"
	myChannel <- "Inside A2"

	flag = false

	condition.Signal()

	lock.Unlock()
}

func printB() {
	defer wg.Done()

	lock.Lock()

	for flag {
		condition.Wait()
	}

	myChannel <- "Inside B"

	flag = true

	condition.Signal()

	lock.Unlock()
}


func main() {
	fmt.Println("Start of Main")

	wg.Add(noOfRoutines)

	go printA()
	go printB()


	wg.Wait()

   close(myChannel)

	for data := range myChannel {
		fmt.Println(data)
	}

	fmt.Println("End of Main")
}
