/*
go.mod : 의존성 관리 파일
-> 기본 라이브러리를 사용하면 필요없지만 외부 라이브러리를 사용하면 필요
go mod init Lets-Go/hello (모듈 이름)
go get github.com/gorilla/mux (외부 라이브러리 설치)
go build (빌드)
*/

package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}