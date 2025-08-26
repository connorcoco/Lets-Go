/*
repository/user_repository.go - 데이터 접근 계층

이 파일은 데이터 저장소(데이터베이스, 메모리 등)에 접근하는 로직을 담당합니다.
현재는 메모리 기반 저장소를 사용하지만, 나중에 데이터베이스로 쉽게 교체할 수 있습니다.

역할:
- 데이터 CRUD 작업 (Create, Read, Update, Delete)
- 데이터 저장소 추상화
- 동시성 안전성 보장 (mutex 사용)
*/

package repository

import (
	"Lets-Go/ex/internal/model"
	"sync"
)

// UserRepository는 사용자 데이터 접근을 담당하는 구조체입니다
type UserRepository struct {
	users  map[int]*model.User // 메모리 기반 사용자 저장소
	mutex  sync.RWMutex        // 읽기/쓰기 동시성 제어
	nextID int                 // 다음 사용자 ID
}

// NewUserRepository는 새로운 UserRepository 인스턴스를 생성합니다
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:  make(map[int]*model.User),
		nextID: 1,
	}
}

// GetAll은 모든 사용자 목록을 반환합니다
func (r *UserRepository) GetAll() []*model.User {
	r.mutex.RLock()         // 읽기 락 획득
	defer r.mutex.RUnlock() // 함수 종료 시 읽기 락 해제

	// map의 값들을 slice로 변환
	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

// GetByID는 ID로 특정 사용자를 조회합니다
func (r *UserRepository) GetByID(id int) *model.User {
	r.mutex.RLock()         // 읽기 락 획득
	defer r.mutex.RUnlock() // 함수 종료 시 읽기 락 해제
	return r.users[id]
}

// Create는 새로운 사용자를 생성합니다
func (r *UserRepository) Create(user *model.User) int {
	r.mutex.Lock()         // 쓰기 락 획득
	defer r.mutex.Unlock() // 함수 종료 시 쓰기 락 해제

	user.ID = r.nextID
	r.users[user.ID] = user
	r.nextID++

	return user.ID
}
