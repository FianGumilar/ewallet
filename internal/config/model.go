package config

type Config struct {
	Database Database
	Server   Server
	Redis    Redis
}

type Database struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type Server struct {
	Host string
	Port string
}

type Redis struct {
	Addr string
	Pass string
}
