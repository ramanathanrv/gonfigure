package gonfigure

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestPing(t *testing.T) {
	options := GetDevelopmentConfigMysql()
	c, err := NewReloadingConfig(options)
	if err != nil {
		log.Println("Couldn't get valid DB instance: ", err)
		t.Errorf("Can't connect to Database correctly")
	} else {
		pingErr := c.Check()
		if pingErr != nil {
			t.Errorf("There is some problem with the db connection")
		}
	}
}

func getConfigInstance() *ReloadingConfig {
	params := GetDevelopmentConfigMysql()
	var c *ReloadingConfig
	c, _ = NewReloadingConfig(params)
	return c
}

func openDatabase() *sql.DB {
	params := GetDevelopmentConfigMysql()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", params.Username, params.Password, params.Host, params.Port, params.Schema)
	db, _ := sql.Open("mysql", dsn)
	return db
}

func removeHelloEntry() {
	params := GetDevelopmentConfigMysql()
	db := openDatabase()
	deleteQuery := fmt.Sprintf("delete from %s where name = ?", params.Table)

	res, _ := db.Exec(deleteQuery, testKey)
	log.Println("Delete query executed ", res)
	log.Println(res.RowsAffected())
	db.Close()
}

func addHelloEntry() {
	params := GetDevelopmentConfigMysql()
	db := openDatabase()
	insertQuery := fmt.Sprintf("insert into %s (name, value, type) values (?, ?, ?)", params.Table)

	res, _ := db.Exec(insertQuery, testKey, initVal, "string")
	log.Println("Insert query executed ", res)
	log.Println(res.RowsAffected())
	db.Close()
}

func updateDatabase(key string, value string) {
	params := GetDevelopmentConfigMysql()
	db := openDatabase()
	updateQuery := fmt.Sprintf("update %s set value = ? where name = ?", params.Table)

	res, err := db.Exec(updateQuery, value, key)
	if err != nil {
		log.Println("Error while trying to execute update query: ", err)
		os.Exit(1)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error while trying to execute update query: ", err)
		os.Exit(1)
	}
	log.Println("Rows affected: ", rowsAffected)
	db.Close()
}

const resultNotFound string = "noworldforyou"
const initVal string = "beforebigbang"
const newVal string = "newworld"
const testKey string = "hello"
const randomKey string = "__random"

func testResolve(c *ReloadingConfig, key string) string {
	val := c.ResolveD(key, resultNotFound)
	return val
}

func assertS(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Fatalf("Expected %s but actually got %s", expected, actual)
	}
}
func TestResolve(t *testing.T) {
	var actual string
	// test when the entry is not present at all
	removeHelloEntry()
	c := getConfigInstance()
	actual = testResolve(c, testKey)
	assertS(resultNotFound, actual, t)

	// test when the entry is freshly added to the DB
	addHelloEntry()
	c.ReloadProperties()
	actual = testResolve(c, testKey)
	assertS(initVal, actual, t)

	// test when the entry is updated in the DB
	updateDatabase(testKey, newVal)
	c.ReloadProperties()
	actual = testResolve(c, testKey)
	assertS(newVal, actual, t)

	// test non-existent key
	actual = testResolve(c, randomKey)
	assertS(resultNotFound, actual, t)

	time.Sleep(120 * time.Second)
}
