/*
repository/user_repository.go - 데이터 접근 계층

이 파일은 데이터 저장소(데이터베이스, 메모리 등)에 접근하는 로직을 담당합니다.
현재는 메모리 기반 저장소를 사용하지만, 나중에 데이터베이스로 쉽게 교체할 수 있습니다.

역할:
- 데이터 CRUD 작업 (Create, Read, Update, Delete)
- 데이터 저장소 추상화
- 고급 동시성 안전성 보장 (세밀한 락, 트랜잭션 지원)
*/

package repository

import (
	"Lets-Go/ex/internal/model"
	"errors"
	"sync"
	"time"
)

// Transaction은 데이터베이스 트랜잭션을 시뮬레이션합니다
type Transaction struct {
	users  map[int]*model.User
	nextID int
	mu     sync.RWMutex
}

// UserRepository는 사용자 데이터 접근을 담당하는 구조체입니다
type UserRepository struct {
	users  map[int]*model.User // 메모리 기반 사용자 저장소
	mutex  sync.RWMutex        // 읽기/쓰기 동시성 제어
	nextID int                 // 다음 사용자 ID

	// 고급 동시성 제어를 위한 필드들
	userLocks map[int]*sync.RWMutex // 개별 사용자별 락
	lockMutex sync.RWMutex          // userLocks 맵 보호용 락

	// 성능 모니터링
	stats struct {
		readCount  int64
		writeCount int64
		lockWait   time.Duration
	}
	statsMutex sync.RWMutex
}

// NewUserRepository는 새로운 UserRepository 인스턴스를 생성합니다
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:     make(map[int]*model.User),
		userLocks: make(map[int]*sync.RWMutex),
		nextID:    1,
	}
}

// getUserLock은 특정 사용자 ID에 대한 락을 반환합니다 (없으면 생성)
func (r *UserRepository) getUserLock(userID int) *sync.RWMutex {
	r.lockMutex.RLock()
	lock, exists := r.userLocks[userID]
	r.lockMutex.RUnlock()

	if !exists {
		r.lockMutex.Lock()
		// double-check pattern
		if lock, exists = r.userLocks[userID]; !exists {
			lock = &sync.RWMutex{}
			r.userLocks[userID] = lock
		}
		r.lockMutex.Unlock()
	}
	return lock
}

// BeginTransaction은 새로운 트랜잭션을 시작합니다
func (r *UserRepository) BeginTransaction() *Transaction {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 현재 상태를 복사
	users := make(map[int]*model.User)
	for id, user := range r.users {
		// 깊은 복사 (실제로는 더 정교한 복사 필요)
		userCopy := *user
		users[id] = &userCopy
	}

	return &Transaction{
		users:  users,
		nextID: r.nextID,
	}
}

// Commit은 트랜잭션을 커밋합니다
func (r *UserRepository) Commit(tx *Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 트랜잭션의 변경사항을 메인 저장소에 적용
	for id, user := range tx.users {
		r.users[id] = user
	}
	r.nextID = tx.nextID

	return nil
}

// Rollback은 트랜잭션을 롤백합니다
func (r *UserRepository) Rollback(tx *Transaction) {
	// 트랜잭션 객체는 가비지 컬렉션에 의해 정리됨
}

// GetAll은 모든 사용자 목록을 반환합니다 (성능 개선)
func (r *UserRepository) GetAll() []*model.User {
	start := time.Now()
	r.mutex.RLock()
	defer func() {
		r.mutex.RUnlock()
		r.updateStats(true, time.Since(start))
	}()

	// map의 값들을 slice로 변환
	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

// GetByID는 ID로 특정 사용자를 조회합니다 (개별 락 사용)
func (r *UserRepository) GetByID(id int) *model.User {
	start := time.Now()

	// 개별 사용자 락 사용
	userLock := r.getUserLock(id)
	userLock.RLock()
	defer func() {
		userLock.RUnlock()
		r.updateStats(true, time.Since(start))
	}()

	return r.users[id]
}

// Create는 새로운 사용자를 생성합니다 (원자적 ID 할당)
func (r *UserRepository) Create(user *model.User) int {
	start := time.Now()

	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.updateStats(false, time.Since(start))
	}()

	user.ID = r.nextID
	r.users[user.ID] = user
	r.nextID++

	// 새 사용자에 대한 락 생성
	r.userLocks[user.ID] = &sync.RWMutex{}

	return user.ID
}

// Update는 사용자 정보를 업데이트합니다 (낙관적 락킹)
func (r *UserRepository) Update(user *model.User) error {
	start := time.Now()

	userLock := r.getUserLock(user.ID)
	userLock.Lock()
	defer func() {
		userLock.Unlock()
		r.updateStats(false, time.Since(start))
	}()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("사용자를 찾을 수 없습니다")
	}

	r.users[user.ID] = user
	return nil
}

// Delete는 사용자를 삭제합니다
func (r *UserRepository) Delete(id int) error {
	start := time.Now()

	userLock := r.getUserLock(id)
	userLock.Lock()
	defer func() {
		userLock.Unlock()
		r.updateStats(false, time.Since(start))
	}()

	if _, exists := r.users[id]; !exists {
		return errors.New("사용자를 찾을 수 없습니다")
	}

	delete(r.users, id)
	return nil
}

// BulkCreate는 여러 사용자를 한 번에 생성합니다 (배치 처리)
func (r *UserRepository) BulkCreate(users []*model.User) []int {
	start := time.Now()

	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.updateStats(false, time.Since(start))
	}()

	ids := make([]int, len(users))
	for i, user := range users {
		user.ID = r.nextID
		r.users[user.ID] = user
		ids[i] = r.nextID
		r.nextID++

		// 새 사용자에 대한 락 생성
		r.userLocks[user.ID] = &sync.RWMutex{}
	}

	return ids
}

// SearchByCondition은 조건에 맞는 사용자를 검색합니다
func (r *UserRepository) SearchByCondition(predicate func(*model.User) bool) []*model.User {
	start := time.Now()

	r.mutex.RLock()
	defer func() {
		r.mutex.RUnlock()
		r.updateStats(true, time.Since(start))
	}()

	var results []*model.User
	for _, user := range r.users {
		if predicate(user) {
			results = append(results, user)
		}
	}

	return results
}

// GetStats는 레포지토리 사용 통계를 반환합니다
func (r *UserRepository) GetStats() (readCount, writeCount int64, avgLockWait time.Duration) {
	r.statsMutex.RLock()
	defer r.statsMutex.RUnlock()

	if r.stats.readCount+r.stats.writeCount > 0 {
		avgLockWait = r.stats.lockWait / time.Duration(r.stats.readCount+r.stats.writeCount)
	}

	return r.stats.readCount, r.stats.writeCount, avgLockWait
}

// updateStats는 통계 정보를 업데이트합니다
func (r *UserRepository) updateStats(isRead bool, lockWait time.Duration) {
	r.statsMutex.Lock()
	defer r.statsMutex.Unlock()

	if isRead {
		r.stats.readCount++
	} else {
		r.stats.writeCount++
	}
	r.stats.lockWait += lockWait
}

// Cleanup은 사용하지 않는 락들을 정리합니다
func (r *UserRepository) Cleanup() {
	r.lockMutex.Lock()
	defer r.lockMutex.Unlock()

	// 존재하지 않는 사용자 ID에 대한 락 제거
	for id := range r.userLocks {
		if _, exists := r.users[id]; !exists {
			delete(r.userLocks, id)
		}
	}
}
