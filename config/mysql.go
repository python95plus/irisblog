package config

type mysqlConfig struct {
	Database string `json:"database" form:"database"`
	User     string `json:"user" form:"user"`
	Password string `json:"password" form:"password"`
	Host     string `json:"host" form:"host"`
	Port     int    `json:"port" form:"port"`
	Url      string `json:"-"`
}
