/*
main.go - 애플리케이션 진입점

이 파일은 Go 프로그램의 시작점입니다.
main() 함수가 있어야 실행 가능한 프로그램이 됩니다.

역할:
- 의존성 주입 및 초기화
- HTTP 서버 설정
- 라우터 설정
- 서버 시작
*/

package main

import (
	"Lets-Go/ex/internal/handler"
	"Lets-Go/ex/internal/repository"
	"Lets-Go/ex/internal/service"
	"log"
	"net/http"
)

func main() {
	// 의존성 주입 (Dependency Injection)
	// 각 계층을 순서대로 생성하여 연결
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// HTTP 라우터 설정
	// 각 URL 경로에 해당하는 핸들러 함수 연결
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.GetUsers)     // 모든 사용자 조회
	mux.HandleFunc("/users/", userHandler.GetUserByID) // 특정 사용자 조회

	// HTTP 서버 시작
	log.Println("서버가 포트 8080에서 시작됩니다...")
	log.Println("API 엔드포인트:")
	log.Println("  GET /users     - 모든 사용자 조회")
	log.Println("  GET /users/{id} - 특정 사용자 조회")

	// 서버 시작 및 에러 처리
	log.Fatal(http.ListenAndServe(":8080", mux))
}
