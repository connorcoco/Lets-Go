/*
변수의 범위
변수는 선언된 블록 내에서만 접근 가능
*/

package main

import "fmt"

var g int = 10 // 패키지 전역 변수

func main() {
	var m int = 20

	{
		var s int = 50
		fmt.Println(m, s, g)
	}

	// m = s + 20 // undefined: scompilerUndeclaredName
}
