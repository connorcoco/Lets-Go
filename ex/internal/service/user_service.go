/*
service/user_service.go - 비즈니스 로직 계층

이 파일은 애플리케이션의 핵심 비즈니스 로직을 처리합니다.
Repository에서 데이터를 가져와서 비즈니스 규칙을 적용하고,
Handler에게 결과를 전달합니다.

역할:
- 비즈니스 규칙 적용
- 데이터 검증
- 트랜잭션 관리
- Repository와 Handler 사이의 중재자
*/

package service

import (
	"Lets-Go/ex/internal/model"
	"Lets-Go/ex/internal/repository"
)

// UserService는 사용자 관련 비즈니스 로직을 처리하는 구조체입니다
type UserService struct {
	userRepo *repository.UserRepository // 사용자 데이터 접근 객체
}

// NewUserService는 새로운 UserService 인스턴스를 생성합니다
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetAllUsers는 모든 사용자 목록을 반환합니다
// 비즈니스 로직이 필요한 경우 여기에 추가할 수 있습니다
func (s *UserService) GetAllUsers() []*model.User {
	return s.userRepo.GetAll()
}

// GetUserByID는 ID로 특정 사용자를 조회합니다
func (s *UserService) GetUserByID(id int) *model.User {
	// 비즈니스 규칙: ID가 0보다 작으면 nil 반환
	if id <= 0 {
		return nil
	}

	return s.userRepo.GetByID(id)
}

// CreateUser는 새로운 사용자를 생성합니다
func (s *UserService) CreateUser(name, email string) *model.User {
	// 비즈니스 규칙: 이름과 이메일이 비어있으면 nil 반환
	if name == "" || email == "" {
		return nil
	}

	user := model.NewUser(name, email)
	user.ID = s.userRepo.Create(user)
	return user
}
