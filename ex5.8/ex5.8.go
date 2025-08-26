/*
표준 입력
Scanln() 함수는 표준 입력 스트림에서 데이터를 읽어서 변수에 저장합니다
& : 변수의 주소를 전달합니다
nil : 값이 없다는 뜻
*/

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)

	var a int
	var b int

	n, err := fmt.Scanln(&a, &b)
	if err != nil {
		fmt.Println(err)
		stdin.ReadString('\n') // \n이 나올때까지 읽어라 (키보드 버퍼를 비움)
	} else {
		fmt.Println(n, a, b)
	}
}
