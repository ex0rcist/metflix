package security

import "github.com/stretchr/testify/mock"

type MockSigner struct {
	mock.Mock
}

func (m *MockSigner) CalculateSignature(data []byte) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockSigner) VerifySignature(data []byte, hash string) (bool, error) {
	args := m.Called(data, hash)
	return args.Bool(0), args.Error(1)
}
