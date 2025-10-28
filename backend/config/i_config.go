package config

type IConfiguration interface {
	IsDevelopment() bool
	IsProduction() bool
	Server() Server
	Database() Database
	Authentication() Authentication
	Scopes() string
}
