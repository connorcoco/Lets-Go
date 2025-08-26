/*
변수 선언 (var)
Golang은 강타입 언어이므로 타입이 맞지 않으면 에러
*/

package main

import "fmt"

func main() {
	var a int = 10
	var msg string = "Hello variable"

	a = 20
	msg = "Good morning"
	fmt.Println(msg, a)
}