package gonfigure

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
)

// Config - DataStructure to hold the Database and other related information
type ReloadingConfig struct {
	db           *sql.DB
	tableName    string
	lastReloaded time.Time
	props        map[string]string
	types        map[string]string
}

// MysqlConfigOption - DataStructure that represents the Mysql configuration parameters
type MysqlConfigOption struct {
	Host     string
	Port     int
	Username string
	Password string
	Schema   string
	Table    string
}

// GetDevelopmentConfigMysql returns a config that is typically used in the development phase
func GetDevelopmentConfigMysql() *MysqlConfigOption {
	return &MysqlConfigOption{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "cloud",
		Password: "scape",
		Schema:   "testdb",
		Table:    "config",
	}
}

// NewReloadingConfig returns an instance of the Config
func NewReloadingConfig(params *MysqlConfigOption) (*ReloadingConfig, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", params.Username, params.Password, params.Host, params.Port, params.Schema)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Could not connect to database: ", err)
		return nil, err
	}
	log.Println("Database connection is successful")
	var conf = ReloadingConfig{db: db, tableName: params.Table}
	conf.props = make(map[string]string)
	conf.types = make(map[string]string)
	conf.ReloadProperties() // initial loading of all properties
	c := cron.New()
	c.AddFunc("@every 1m", func() { conf.ReloadProperties() })
	c.Start()
	return &conf, nil
}

// NewReloadingConfigFromDB - returns a new Instance from existing DB
func NewReloadingConfigFromDB(db *sql.DB, tableName string) (*ReloadingConfig, error) {
	log.Println("Database connection is successful")
	var conf = ReloadingConfig{db: db, tableName: tableName}
	conf.props = make(map[string]string)
	conf.types = make(map[string]string)
	conf.ReloadProperties() // initial loading of all properties
	c := cron.New()
	c.AddFunc("@every 1m", func() { conf.ReloadProperties() })
	c.Start()
	return &conf, nil
}

// ReloadProperties - fetches all the properties from Database
func (c *ReloadingConfig) ReloadProperties() error {
	log.Println("Commencing Reload Properties")
	query := fmt.Sprintf("select name, value, type from %s", c.tableName)
	rows, err := c.db.Query(query)
	if err != nil {
		log.Println("Error while trying to fetch all config from database: ", err)
		return err
	}
	// load the properties from the rows
	// cursor variables
	var (
		curName  string
		curValue string
		curType  string
	)
	for rows.Next() {
		err := rows.Scan(&curName, &curValue, &curType)
		if err != nil {
			log.Println("Error while reading row from config table")
			return err
		}
		c.props[curName] = curValue
		// c.types[curName] = curType // type cannot change
	}
	c.lastReloaded = time.Now()
	log.Println("Completed reloading properties")
	return nil
}

// Check function ensures that we are able to reach the database
func (c *ReloadingConfig) Check() error {
	var echoVal int
	const echoTest int = 100
	err := c.db.Ping()
	if err != nil {
		log.Println("Ping failed: ", err)
		return err
	}
	log.Println("Ping is successful")
	err = c.db.QueryRow(fmt.Sprintf("select %d", echoTest)).Scan(&echoVal)
	log.Println("Value is ", echoVal)
	if err != nil {
		return err
	}
	// ensure that there is a row of data
	log.Println("Value is ", echoVal)
	if echoTest == echoVal {
		return nil
	}
	return errors.New("Database did not echo the result correctly")
}

// Resolve a given key & return the value
func (c *ReloadingConfig) Resolve(key string) (string, error) {
	if val, ok := c.props[key]; ok {
		return val, nil
	}
	return "", errors.New("Cannot find the given key")
}

// ResolveD - resolve a given key. If not found, then return the Default Value supplied
func (c *ReloadingConfig) ResolveD(key string, defaultVal string) string {
	if val, ok := c.props[key]; ok {
		return val
	}
	return defaultVal
}

// ResolveInt - resolves the key to an int value
func (c *ReloadingConfig) ResolveInt(key string, defaultVal int) (int, error) {
	val, err := c.Resolve(key)
	if err != nil {
		return strconv.Atoi(val)
	}
	return defaultVal, nil
}

// ResolveInt64 - resolves the value to an int64 value
func (c *ReloadingConfig) ResolveInt64(key string, defaultVal int64) (int64, error) {
	val, err := c.Resolve(key)
	if err != nil {
		return strconv.ParseInt(val, 10, 64)
	}
	return defaultVal, nil
}

// ResolveFloat - resolve the value to a float variable
func (c *ReloadingConfig) ResolveFloat(key string, defaultVal float32) (float32, error) {
	val, err := c.Resolve(key)
	if err != nil {
		f64, e1 := strconv.ParseFloat(val, 32)
		if e1 != nil {
			return 0.0, e1
		}
		f32 := float32(f64)
		return f32, nil
	}
	return defaultVal, nil
}
