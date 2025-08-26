/*
타입 변환 주의사항
큰 타입에서 작은 타입으로 변환할 때 값이 잘릴 수 있음
*/

package main

import "fmt"

func main() {
	var a int16 = 3456
	var b int8 = int8(a)

	fmt.Println(a, b) // 3456 -128
}
