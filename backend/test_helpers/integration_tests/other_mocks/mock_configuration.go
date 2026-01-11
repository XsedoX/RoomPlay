package othermocks

import "github.com/XsedoX/RoomPlay/config"

type MockConfiguration struct{}

func (m *MockConfiguration) IsDevelopment() bool { return false }
func (m *MockConfiguration) IsProduction() bool  { return false }
func (m *MockConfiguration) IsTesting() bool     { return true }
func (m *MockConfiguration) Server() config.Server {
	return config.Server{Host: "localhost", Port: "8080"}
}
func (m *MockConfiguration) Database() config.Database { return config.Database{} }
func (m *MockConfiguration) Authentication() config.Authentication {
	return config.Authentication{EncryptionKey: "12345678901234567890123456789012"}
}
func (m *MockConfiguration) Scopes() string { return "" }
