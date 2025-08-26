/*
model/user.go - 데이터 모델 정의

이 파일은 애플리케이션에서 사용하는 데이터 구조를 정의합니다.
JSON 태그를 사용하여 HTTP API 응답 시 JSON 형태로 변환됩니다.

역할:
- 데이터 구조 정의
- JSON 직렬화/역직렬화 지원
- 데이터 유효성 검증 규칙 정의
*/

package model

import "time"

// User 구조체는 사용자 정보를 담는 데이터 모델입니다
type User struct {
	ID        int       `json:"id"`         // 사용자 고유 ID
	Name      string    `json:"name"`       // 사용자 이름
	Email     string    `json:"email"`      // 사용자 이메일
	CreatedAt time.Time `json:"created_at"` // 계정 생성 시간
	UpdatedAt time.Time `json:"updated_at"` // 정보 수정 시간
}

// NewUser는 새로운 User 인스턴스를 생성하는 팩토리 함수입니다
func NewUser(name, email string) *User {
	now := time.Now()
	return &User{
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateInfo는 사용자 정보를 업데이트합니다
func (u *User) UpdateInfo(name, email string) {
	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now()
}
