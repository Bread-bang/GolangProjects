package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func main() {
	args := os.Args

	targetWord := args[1]
	textFiles := make([]string, len(args) - 2)	// 2: 실행파일(args[0]), targetWord(args[1])

	for i := 2; i < len(args); i++ {
		textFiles = append(textFiles, args[i])
	}

	fmt.Println(targetWord)
	fmt.Println(textFiles)


	for _, fileName := range textFiles {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("[%s] doesn't open\n", fileName)
			continue
		}

		defer func() {
			if err := file.Close(); err != nil {
				panic(fmt.Sprintf("[%s]file couldn't be closed", fileName))
			}
		}()

		lineNum := 0
		scanner := bufio.NewScanner(file)

		fmt.Println(fileName)
		fmt.Println("-----------------------------------")
		
		for scanner.Scan() {
			line := scanner.Text()
			lineNum++
			if strings.Contains(line, targetWord) {
				fmt.Printf("%d\t%s\n", lineNum, line)
			}
		}
		fmt.Println("-----------------------------------")
	}
}
