package mock_jwt_provider

import (
	"github.com/XsedoX/RoomPlay/domain/user/user_id"
	"github.com/stretchr/testify/mock"
)

type MockJwtProvider struct {
	mock.Mock
}

func (m *MockJwtProvider) GenerateToken(userId user_id.UserId) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func (m *MockJwtProvider) ValidateTokenAndGetUserId(tokenString string) (*user_id.UserId, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user_id.UserId), args.Error(1)
}
