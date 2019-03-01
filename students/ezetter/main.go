package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
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

func getInds(num int, randomized bool) []int {
	var a []int
	if randomized {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		a = r.Perm(num)
	} else {
		a = make([]int, num)
		for i := range a {
			a[i] = i
		}
	}
	fmt.Println(a)
	return a
}

func questionLoop(records [][]string, results *resultCounts, randomize bool, done chan bool) {
	reader := bufio.NewReader(os.Stdin)
	inds := getInds(len(records), randomize)
	for _, i := range inds {
		fmt.Printf("%s ", records[i][0])
		inp, _ := reader.ReadString('\n')
		sanatizedInp := strings.ToLower(strings.TrimSpace(inp))
		sanatizedAns := strings.ToLower(strings.TrimSpace(records[i][1]))
		if sanatizedInp == sanatizedAns {
			results.correct++
			fmt.Println("Right!")
		} else {
			fmt.Printf("Wrong! The answer is %s.\n", records[i][1])
		}
	}
	done <- true
}

func main() {
	var probpath = flag.String("probpath", "problems.csv", "Path to the problems file.")
	var timeout = flag.Int("timeout", 30, "Seconds until a quiz times out. 0 for no timeout")
	var randomize = flag.Bool("randomize", false, "If true, present quiz questions in a random order.")
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
	go questionLoop(records, &results, *randomize, done)
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
