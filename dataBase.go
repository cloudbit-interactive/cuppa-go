package cuppago

/*
Libraries
	go get -u github.com/go-sql-driver/mysql
Example:
	db := cuppago.NewDataBase("localhost","golang", "root", "", "")
	// Insert
	data := map[string]interface{}{"name":"Tufik","age":50, "date":"NOW()"}
	db.Insert("users", data, "")
	// Get
	data := map[string]interface{}{"name":"Tufik","age":50, "date":"NOW()"}
	db.Insert("users", data, "")
*/

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// DataBase : DataBase Structure
type DataBase struct {
	Host     string
	Port     string
	Db       string
	Username string
	Password string
	Conn     *sql.DB
}

// Example:
// db := cuppago.NewDataBase("localhost","golang", "root", "", "")
func NewDataBase(host string, db string, username string, password string, port string) DataBase {
	if port == "" {
		port = "3306"
	}
	dataBase := DataBase{host, port, db, username, password, nil}
	url := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	conn, err := sql.Open("mysql", url)
	if err != nil {
		panic(err.Error())
	}
	dataBase.Conn = conn
	if !dataBase.TestConnection() {
		Error("ERROR ESTABLISHING DATABASE CONNECTION")
	}
	return dataBase
}

// Add : Check if record exist if exist update it or create one
func (db *DataBase) Add(table string, data map[string]interface{}, condition string, columnsToReturn string) map[string]interface{} {
	var result map[string]interface{}
	record := db.GetRow(table, condition, "", "id")
	if record == nil {
		result = db.Insert(table, data, columnsToReturn)
	} else {
		result = db.Update(table, data, condition, columnsToReturn)
	}
	return result
}

// Example:
// data := map[string]interface{}{"name":"Tufik","age":50, "date":"NOW()"}
// db.Insert("users", data, "")
func (db *DataBase) Insert(table string, data map[string]interface{}, columnsToReturn string) map[string]interface{} {
	keys := make([]string, 0)
	values := make([]string, 0)
	for key, value := range data {
		value = strings.ReplaceAll(String(value), `'`, `\'`)
		if value != "NOW()" {
			value = "'" + String(value) + "'"
		}
		keys = append(keys, key)
		values = append(values, String(value))
	}
	sql := "INSERT INTO " + table + " ( `"
	sql += strings.Join(keys, "` , `")
	sql += "` ) VALUES ( "
	sql += strings.Join(values, " , ")
	sql += " ) "
	sql = strings.TrimSpace(sql)
	insert, err := db.Conn.Exec(sql)
	if err != nil {
		error := make(map[string]interface{})
		error["error"] = -1
		error["errorMessage"] = err.Error()
		Error(error)
		return error
	}
	id, _ := insert.LastInsertId()
	rowToReturn := db.GetRow(table, "id = "+String(id), "", columnsToReturn)
	return rowToReturn
}

// Update : Update data in a specific table
func (db *DataBase) Update(table string, data map[string]interface{}, condition string, columnsToReturn string) map[string]interface{} {
	sql := "UPDATE " + table + " SET "
	index := 0
	totalData := len(data)
	for key, value := range data {
		value = strings.ReplaceAll(String(value), `'`, `\'`)
		if value != "NOW()" {
			value = "'" + String(value) + "'"
		}
		if index == totalData-1 {
			sql += "`" + String(key) + "`=" + String(value)
		} else {
			sql += "`" + String(key) + "`=" + String(value) + " , "
		}
		index++
	}
	sql += " WHERE " + condition
	sql = strings.TrimSpace(sql)
	_, err := db.Conn.Exec(sql)
	if err != nil {
		error := make(map[string]interface{})
		error["error"] = -1
		error["errorMessage"] = err.Error()
		Error(error)
		return error
	}
	rowToReturn := db.GetRow(table, condition, "", columnsToReturn)
	return rowToReturn
}

// Example:
// data := db.GetList("users", "", "", "", "")
func (db *DataBase) GetList(table string, condition string, limit string, orderBy string, columns string) []map[string]interface{} {
	sql := "SELECT * FROM " + table
	if columns != "" {
		sql = "SELECT " + columns + " FROM " + table
	}
	if condition != "" {
		sql += " WHERE " + condition
	}
	if orderBy != "" {
		sql += " ORDER BY " + orderBy
	}
	if limit != "" {
		sql += " LIMIT " + limit
	}
	rows, _ := db.Conn.Query(sql)
	result := db.GetMap(rows)
	if len(result) == 0 {
		return nil
	}
	return result
}

// GetRow : Get a map of data from a specific table
func (db *DataBase) GetRow(table string, condition string, orderBy string, columns string) map[string]interface{} {
	sql := "SELECT * FROM " + table
	if columns != "" {
		sql = "SELECT " + columns + " FROM " + table
	}
	if condition != "" {
		sql += " WHERE " + condition
	}
	if orderBy != "" {
		sql += " ORDER BY " + orderBy
	}
	sql += " LIMIT 1 "

	rows, err := db.Conn.Query(sql)
	if err != nil {
		return nil
	}
	result := db.GetMap(rows)
	if len(result) == 0 {
		return nil
	}
	return result[0]
}

// Delete : Delete a record from a specific table
func (db *DataBase) Delete(table string, condition string) map[string]interface{} {
	sql := "DELETE FROM " + table + " WHERE " + condition
	deleted, err := db.Conn.Exec(sql)
	if err != nil {
		error := make(map[string]interface{})
		error["error"] = -1
		error["errorMessage"] = err.Error()
		Error(error)
		return error
	}
	rowsAffected, _ := deleted.RowsAffected()
	result := make(map[string]interface{})
	result["rowsAffected"] = rowsAffected
	return result
}

// SQL : Execute any sql and return it result
func (db *DataBase) SQL(query string) []map[string]interface{} {
	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil
	}
	result := db.GetMap(rows)
	if len(result) == 0 {
		return nil
	}
	return result
}

// GetMap : Extract a map of a sql.Rows
func (db *DataBase) GetMap(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	finalResult := make([]map[string]interface{}, 0)
	resultID := 0
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		tmpStruct := map[string]interface{}{}
		for i, col := range columns {
			val := values[i]
			if val == nil {
				val = ""
			}
			tmpStruct[col] = fmt.Sprintf("%s", val)
		}
		finalResult = append(finalResult, tmpStruct)
		resultID++
	}
	return finalResult
}

func (db *DataBase) GetTables() []map[string]interface{} {
	return db.SQL("SHOW TABLES")
}

func (db *DataBase) TestConnection() bool {
	_, err := db.Conn.Query("SHOW TABLES")
	if err != nil {
		Error(err.Error())
		return false
	}
	return true
}
