package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: searchWord <word> <file1> <file2> ...")
		os.Exit(1)
	}

	targetWord := os.Args[1]
	textFiles := os.Args[2:]

	for _, filename := range textFiles {
		if err := searchInFile(filename, targetWord); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func searchInFile(filename string, targetWord string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("[%s] doesn't open: %v", filename, err)
	}
	defer file.Close()

	lineNum := 0
	scanner := bufio.NewScanner(file)

	fmt.Println(filename)
	fmt.Println("-----------------------------------")

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		if strings.Contains(line, targetWord) {
			fmt.Printf("%d\t%s\n", lineNum, line)
		}
	}
	fmt.Println("-----------------------------------")

	return nil
}