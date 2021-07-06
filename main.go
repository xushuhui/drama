package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var db *sql.DB
var err error

type TableStruct struct {
	columnName    string
	columnDefault interface{}
	dataType      string
	columnKey     string
	columnComment string
}

func main() {
	dsn := "root@tcp(127.0.0.1:3306)/information_schema"
	schema := "mall"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dir := "./schema/"
	tables := ListTables(schema)
	is, _ := PathExists(dir)
	if !is {
		os.Mkdir(dir, os.ModePerm)
	}

	for _, v := range tables {

		columns := ListColumns(schema, v)
		var fields []string
		for _, v := range columns {
			//id 主键字段已经内置
			if v.columnKey == "PRI" {
				continue
			}
			field := makeField(v)
			fields = append(fields, field)
		}
		t := camelString(v)
		f := strings.Join(fields, "\n\t\t")
		str := fmt.Sprintf(TemplateMain, t, t, f)

		generateFile(dir+v+".go", str)
	}

}
func makeField(v TableStruct) string {

	fieldStr := `field.` + DataTypeMap[v.dataType] + `("` + v.columnName + `").Comment("` + v.columnComment + `")`
	if v.columnKey == "MUL" && DataTypeMap[v.dataType] != "Time" {
		fieldStr = fieldStr + ".Unique()"
	}
	fieldStr = fieldStr + ","
	return fieldStr

}
func generateFile(filename string, str string) {

	err := ioutil.WriteFile(filename, []byte(str), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

}
func ListColumns(schema string, table string) (tables []TableStruct) {
	//select column_name,column_default,data_type,column_comment,column_key from information_schema.columns where  TABLE_NAME="a
	rows, err := db.Query("select column_name,column_default,data_type,column_comment,column_key from information_schema.columns where table_schema = ? and TABLE_NAME=?", schema, table)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var t TableStruct
		rows.Scan(&t.columnName, &t.columnDefault, &t.dataType, &t.columnComment, &t.columnKey)
		tables = append(tables, t)
	}
	return
}
func ListTables(schema string) (tables []string) {
	rows, err := db.Query("select table_name from information_schema.tables where table_schema=?", schema)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t string
		rows.Scan(&t)

		tables = append(tables, t)
	}
	return
}

func DoQuery(sqlInfo string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	_ = rows.Close()
	fmt.Println("list", list)
	return list, nil
}
