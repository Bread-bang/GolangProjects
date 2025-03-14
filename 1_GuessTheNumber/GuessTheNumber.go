package main

import (
	"fmt"
	"bufio"
	"math/rand"
	"os"
	"time"
)

var stdin = bufio.NewReader(os.Stdin)

func InputIntValue() (int, error) {
	var n int

	_, err := fmt.Scanln(&n)
	if err != nil {
		stdin.ReadString('\n')
	}

	return n, err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	answer := rand.Intn(100)
	cnt := 1
	var guesses []int  // 이전 추측들을 저장할 슬라이스

	for true {
		fmt.Print("숫자값을 입력하세요:")
		guessNum, err := InputIntValue()

		if err != nil {
			fmt.Println("숫자만 입력하세요.")
		} else {
			guesses = append(guesses, guessNum)  // 추측 저장
			if answer == guessNum {
				fmt.Printf("숫자를 맞췄습니다. 축하합니다. 시도한 횟수: %d\n", cnt)
				fmt.Printf("추측 기록 길이: %d, 용량: %d\n", len(guesses), cap(guesses))
				break
			} else if answer > guessNum {
				fmt.Println("입력하신 숫자가 더 작습니다.")
			} else if answer < guessNum {
				fmt.Println("입력하신 숫자가 더 큽니다.")
			}
			cnt++
		} 
	}
}