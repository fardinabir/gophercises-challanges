package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go interruptCather(interrupt, &wg)

	records := csvLoader()

	var correctCount, totalCount int
	userInput := make(chan string, 1)
	go questionServer(&wg, records, userInput, &correctCount, &totalCount)

	wg.Wait()
	fmt.Printf("End of test. Result is %v out of %v.\n", correctCount, totalCount)
}

func questionServer(wg *sync.WaitGroup, records [][]string, userInput chan string, correctCount, totalCount *int) {
	defer wg.Done()
	for i := 0; i < len(records); i++ {
		timeOut := make(chan bool, 1)
		*totalCount++
		fmt.Println(records[i][0], "=?")
		go scanUserInput(userInput)
		go timeoutNotifier(timeOut)
		select {
		case val := <-userInput:
			ans, _ := strconv.Atoi(records[i][1])
			given, _ := strconv.Atoi(val)
			if ans == given {
				*correctCount++
			}
		case <-timeOut:
			fmt.Println("Times'Up..........")
			return
		}
		//close(timeOut)
	}
	fmt.Printf("End of test. Result is %v out of %v.", *correctCount, *totalCount)
}

func interruptCather(interrupt chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done()
	<-interrupt
	fmt.Println("\nReceived interrupt signal. Stopping input process...")
	//os.Exit(0)
}

func timeoutNotifier(notifier chan bool) {
	time.Sleep(5 * time.Second)
	notifier <- true
	//fmt.Println("Times'Up..........")
	close(notifier)
}

func scanUserInput(text chan string) {
	var input string
	fmt.Scanln(&input)
	text <- input
}

func csvLoader() [][]string {
	file, err := os.Open("problems.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return records
}

// //8+3,11
// //1+2,3
// //8+6,14
// //3+1,4
// //1+4,5
// //5+1,6
// //2+3,5
// //3+3,6
// //2+4,6
// //5+2,7
//func main() {
//	var wg sync.WaitGroup
//	wg.Add(1)
//	go func1(&wg)
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		time.Sleep(5 * time.Second)
//		fmt.Println("Returning from func2")
//	}()
//	wg.Wait()
//	fmt.Println("Returning from all funcs")
//}
//
//func func1(group *sync.WaitGroup) {
//	defer group.Done()
//	time.Sleep(10 * time.Second)
//	fmt.Println("Returning from func1")
//}
//
//func func2(group *sync.WaitGroup) {
//	defer group.Done()
//	time.Sleep(5 * time.Second)
//	fmt.Println("Returning from func2")
//}
