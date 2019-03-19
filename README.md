## Goals of the project

This project implements a config system where the values are reloaded at regular interval from Database. The solution is tightly coupled to MySQL database in this version. Config values are reloaded every 1 minute using a separate goroutine.

### Dependencies

[comment]: # (For Mysql driver)
go get -u github.com/go-sql-driver/mysql
go get -u github.com/robfig/cron

### Database table schema

```mysql
create table dynamic_config (
  id bigint(20) auto_increment primary key ,
  name varchar(255) not null,
  value varchar(255) not null,
  type varchar(32) not null
)
```

The important columns as name, value & type. You may add other columns as required.

### Usage

```go
import "gonfiure"
import "fmt"
params = MysqlConfigOption{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "username",
		Password: "password",
		Schema:   "testdb",
		Table:    "dynamic_config",
	}
config := NewReloadingConfig(params)

var val string = config.ResolveD("propertyName", "defaultVal")
fmt.Println(val) // Should print the value configured in the database, else defaultVal

```