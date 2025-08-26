/*
타입 변환
타입 변환은 타입 변환 연산자를 사용하여 수행
*/

package main

import "fmt"

func main() {
	a := 3
	var b float64 = 3.5

	var c int = int(b) // 3.5 -> 3
	d := float64(a) * b

	var e int64 = 7
	f := a * int(e) // int와 int64는 같은 타입이 아니므로 타입 변환 필요 (실제 타입이 같아도 타입 이름이 다르면 다름)

	fmt.Println(a, b, c, d, e, f)
}
