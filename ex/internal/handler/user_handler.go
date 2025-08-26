/*
handler/user_handler.go - HTTP 요청 처리 계층

이 파일은 HTTP 요청을 받아서 처리하고 응답을 반환하는 역할을 합니다.
Service 계층의 메서드를 호출하여 비즈니스 로직을 실행하고,
결과를 HTTP 응답으로 변환합니다.

역할:
- HTTP 요청 파싱
- 요청 데이터 검증
- Service 계층 호출
- HTTP 응답 생성
- 에러 처리 및 상태 코드 설정
*/

package handler

import (
	"Lets-Go/ex/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// UserHandler는 사용자 관련 HTTP 요청을 처리하는 구조체입니다
type UserHandler struct {
	userService *service.UserService // 사용자 비즈니스 로직 서비스
}

// NewUserHandler는 새로운 UserHandler 인스턴스를 생성합니다
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers는 모든 사용자 목록을 조회하는 HTTP 핸들러입니다
// GET /users 요청을 처리합니다
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// HTTP 메서드 검증
	if r.Method != http.MethodGet {
		http.Error(w, "허용되지 않는 메서드입니다", http.StatusMethodNotAllowed)
		return
	}

	// Service 계층에서 사용자 목록 조회
	users := h.userService.GetAllUsers()

	// JSON 응답 헤더 설정
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// JSON 응답 생성
	json.NewEncoder(w).Encode(users)
}

// GetUserByID는 특정 사용자를 조회하는 HTTP 핸들러입니다
// GET /users/{id} 요청을 처리합니다
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// HTTP 메서드 검증
	if r.Method != http.MethodGet {
		http.Error(w, "허용되지 않는 메서드입니다", http.StatusMethodNotAllowed)
		return
	}

	// URL 경로에서 사용자 ID 추출: /users/123 -> 123
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "잘못된 사용자 ID입니다", http.StatusBadRequest)
		return
	}

	// Service 계층에서 사용자 조회
	user := h.userService.GetUserByID(id)
	if user == nil {
		http.Error(w, "사용자를 찾을 수 없습니다", http.StatusNotFound)
		return
	}

	// JSON 응답 헤더 설정
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// JSON 응답 생성
	json.NewEncoder(w).Encode(user)
}
