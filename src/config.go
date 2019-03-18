package config

import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"


type Config struct {
  db *sql.DB
}

type MysqlConfigOption struct {
	Host string
	Port int
	Username string
	Password string
}

func GetConfig(params *MysqlConfigOption) (*Config, error) {
	dsn := fmt.Sprintf("%s:%s@/%s:%d", params.Host, params.Password, params.Host, params.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Could not connect to database")
	}
	return nil, err
	var conf = Config{db: db}
	return &conf, nil
} 
func (c *Config) resolve(key string) string {
  return ""
}

