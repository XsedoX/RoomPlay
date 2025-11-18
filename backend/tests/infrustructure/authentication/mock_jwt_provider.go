package authentication

import (
	"github.com/stretchr/testify/mock"
	"xsedox.com/main/domain/user"
)

type MockJwtProvider struct {
	mock.Mock
}

func (m *MockJwtProvider) GenerateToken(userId user.Id) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func (m *MockJwtProvider) ValidateTokenAndGetUserId(tokenString string) (*user.Id, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.Id), args.Error(1)

}
