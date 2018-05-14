package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var flagFilename string
var flagTimelimit int

func init() {
	flag.StringVar(&flagFilename, "csvFile", "problems.csv", "Quiz problems in csv format.")
	flag.IntVar(&flagTimelimit, "timelimit", 30, "Timelimit for the quiz.")
}

func timer(end time.Time, timeUp chan bool) {
	for {
		now := time.Now()
		if now.After(end) {
			// println("Sorry time is up.\n")
			timeUp <- true
			close(timeUp)
			return
		}
	}
}

func askQuestions(lines [][]string, que chan string, ans, done chan bool) {

	for i, line := range lines {
		// fmt.Printf("%d : %s | %s\n", i, line[0], line[1])
		que <- fmt.Sprintf("Q %d: What is %s ?\n\n", i, line[0])

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		// println(scanner.Text())

		input, err := strconv.Atoi(scanner.Text())
		if err != nil {
			println("Please enter an integer.")
		}
		answer, err := strconv.Atoi(line[1])
		if err != nil {
			println("Invalid answer for question.")
		}

		if answer == input {
			// println("You are right!\n")
			ans <- true
		} else {
			// println("You are wrong.\n")
			ans <- false
		}
	}
	done <- true
}

func printResult(start time.Time, countRight, countWrong int) {
	fmt.Printf("\nYou got %d right and %d wrong.\n", countRight, countWrong)
	fmt.Printf("\nYou have completed the quiz in %s seconds.\n", time.Now().Sub(start).String())
}

func clearScreen() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	flag.Parse()

	// Open csv file
	f, err := os.Open(flagFilename)
	if err != nil {
		fmt.Printf("Error opening file %s\n", flagFilename)
		fmt.Println(err)
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		fmt.Printf("Error reading csv file. %s\n", flagFilename)
		fmt.Println(err)
	}

	// Initialize score
	countWrong, countRight := 0, 0

	// Set end time.
	println("Press enter to start!")
	startTimer := bufio.NewScanner(os.Stdin)
	startTimer.Scan()
	start := time.Now()
	end := start.Add(time.Duration(flagTimelimit) * time.Second)

	que := make(chan string)
	ans := make(chan bool, 1)
	done := make(chan bool, 1)

	go timer(end, done)
	go askQuestions(lines, que, ans, done)

	var text string
	for {
		select {
		case q, _ := <-que:
			text = q
			clearScreen()
			fmt.Printf("Total time: %d | Time remaining: %s\n%s\r\r", flagTimelimit, end.Sub(time.Now()), text)
		case a, _ := <-ans:
			if a {
				countRight++
			} else {
				countWrong++
			}
		case d, _ := <-done:
			if d {
				clearScreen()
				println("Time Up!\n")
				printResult(start, countRight, countWrong)
				return
			}
		}
	}
}
