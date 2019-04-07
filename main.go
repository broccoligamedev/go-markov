package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type markovEntry struct {
	value string
	count int
}

func generateNextElement(source string, markovMap map[string][]*markovEntry) string {
	slice, ok := markovMap[source]
	if ok {
		total := 0
		for _, e := range slice {
			total += e.count
		}
		choice := rand.Intn(total)
		choiceSum := 0
		for _, e := range slice {
			choiceSum += e.count
			if choice < choiceSum {
				return e.value
			}
		}
	}
	return "!"
}

func seedMarkovMap(source string, target string, markovMap map[string][]*markovEntry) {
	slice, ok := markovMap[source]
	if ok {
		for _, entry := range slice {
			if entry.value == target {
				entry.count = entry.count + 1
				return
			}
		}
		markovMap[source] = append(slice, &markovEntry{value: target, count: 1})
	} else {
		newSlice := []*markovEntry{}
		markovMap[source] = append(newSlice, &markovEntry{value: target, count: 1})
	}
}

func printMarkovMap(m map[string][]*markovEntry) {
	for k, v := range m {
		for _, e := range v {
			fmt.Println(fmt.Sprintf("%s: %s (%d)", k, e.value, e.count))
		}
	}
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano())
	forwardMap := make(map[string][]*markovEntry)
	// note(ryan): input should be of the form "<filename> <max order> <number of strings to generate>"
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	maxOrder, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	numStringsToGenerate, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// seed the Markov map
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		// note(ryan): we're using # as a null symbol and ! as a terminator in this implementation
		for i := 0; i < maxOrder; i++ {
			word = "#" + word
		}
		word = word + "!"
		for i := maxOrder - 1; i < len(word)-1; i++ {
			seedMarkovMap(word[i-(maxOrder-1):i+1], string(word[i+1]), forwardMap)
		}
	}
	// printMarkovMap(markovMap)
	// generate the strings
	for x := 0; x < numStringsToGenerate; x++ {
		newMarkovString := ""
		for i := 0; i < maxOrder; i++ {
			newMarkovString = "#" + newMarkovString
		}
		for {
			nextElement := generateNextElement(newMarkovString[len(newMarkovString)-maxOrder:], forwardMap)
			if nextElement == "!" {
				break
			}
			newMarkovString += nextElement
		}
		fmt.Println(newMarkovString[maxOrder:])
	}
}
