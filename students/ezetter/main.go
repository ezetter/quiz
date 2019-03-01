package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type resultCounts struct {
	total   int
	correct int
}

func questionLoop(records [][]string, results *resultCounts, done chan bool) {
	reader := bufio.NewReader(os.Stdin)

	for _, rec := range records {
		fmt.Printf("%s ", rec[0])
		inp, _ := reader.ReadString('\n')
		sanatizedInp := strings.ToLower(strings.TrimSpace(inp))
		sanatizedAns := strings.ToLower(strings.TrimSpace(rec[1]))
		if sanatizedInp == sanatizedAns {
			results.correct++
			fmt.Println("Right!")
		} else {
			fmt.Printf("Wrong! The answer is %s.\n", rec[1])
		}
	}
	done <- true
}

func main() {
	var probpath = flag.String("probpath", "problems.csv", "Path to the problems file.")
	var timeout = flag.Int("timeout", 30, "Seconds until a quiz times out. 0 for no timeout")
	flag.Parse()
	f, err := os.Open(*probpath)
	defer f.Close()
	check(err)
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	check(err)
	results := resultCounts{len(records), 0}

	fmt.Print("Press enter to begin. ")
	bufio.NewReader(os.Stdin).ReadString('\n')

	done := make(chan bool)
	go questionLoop(records, &results, done)
	if *timeout > 0 {
		select {
		case <-done:
		case <-time.After(time.Duration(*timeout) * time.Second):
			fmt.Printf("\nSorry, time's up!\n")
		}
	} else {
		<-done
	}
	fmt.Printf("You got %v out of %v correct.\n", results.correct, results.total)
}
