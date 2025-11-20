package usecase

import (
	"errors"
	"hackathon-backend/model"
	"strings"
	"testing"
)

type stubUserRepository struct {
	insertCalled bool
	lastUser     *model.User
	err          error
}

func (s *stubUserRepository) Insert(u *model.User) error {
	s.insertCalled = true
	s.lastUser = u
	return s.err
}

func TestRegisterUserUsecase_Execute_InvalidInput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name string
		age  int
	}{
		"empty name":    {"", 25},
		"name too long": {strings.Repeat("a", 51), 25},
		"age too young": {"alice", 19},
		"age too old":   {"bob", 81},
	}

	for testName, tc := range tests {
		tc := tc
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			repo := &stubUserRepository{}
			usecase := NewRegisterUserUsecase(repo)

			_, err := usecase.Execute(tc.name, tc.age)
			if !errors.Is(err, ErrInvalidUser) {
				t.Fatalf("expected ErrInvalidUser, got %v", err)
			}
			if repo.insertCalled {
				t.Fatalf("Insert should not be called on invalid input")
			}
		})
	}
}

func TestRegisterUserUsecase_Execute_ValidInput(t *testing.T) {
	t.Parallel()

	repo := &stubUserRepository{}
	usecase := NewRegisterUserUsecase(repo)

	id, err := usecase.Execute("charlie", 30)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Fatalf("expected generated id, got empty string")
	}
	if !repo.insertCalled {
		t.Fatalf("expected Insert to be called")
	}
	if repo.lastUser == nil {
		t.Fatalf("expected user to be passed to repository")
	}
	if repo.lastUser.Name != "charlie" || repo.lastUser.Age != 30 {
		t.Fatalf("unexpected user values: %+v", repo.lastUser)
	}
	if repo.lastUser.ID != id {
		t.Fatalf("expected usecase ID to match inserted user ID")
	}
}
