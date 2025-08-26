/*
여러가지 변수 선언법 (타입 또는 값을 명시적으로 지정해야함)
var a int = 10
var a int  (초기값: 0)
var a = 10
a := 10 (선언대입문 / var a = 10 과 동일)
*/

package main

import "fmt"

func main() {
	var a int = 3
	var b int
	var c = 4
	d := 5
	var e = "Hello"
	f := 3.14

	fmt.Println(a, b, c, d, e, f)
}