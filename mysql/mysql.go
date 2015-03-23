package mysql

import (
  "fmt"
  "strings"
  "os"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "log"
)

/* 起動前に変更 */
const User = "root"
const Pass = ""
const Port = "3306"

type MySQL struct {
  db *sql.DB
  dbName string
}

type ColumnDefinition struct {
  name string
  dataType string
  primaryKey bool
}

type Column struct {
  name string
  value string
}

func NewColumnDefinition(n string, dt string, pk bool) *ColumnDefinition {
  columnDef := &ColumnDefinition{name: n, dataType: dt, primaryKey: pk}
  return columnDef
}

func NewColumn(n string, v string) *Column {
  column := &Column{name: n, value: v}
  return column
}

func (this *MySQL) ConnectServer() {
  var err error
  oSql := User + ":" + Pass + "@tcp(" + os.Getenv("TEST_DB_PORT_3306_TCP_ADDR") + ":" + Port + ")/"
  fmt.Println("ConnectServer:" + oSql)
  this.db, err = sql.Open("mysql", oSql)
  if err != nil {
    log.Fatal(err)
  }
}

func (this *MySQL) UseDB(dbName string) {
  uSql := "USE " + dbName
  fmt.Println("UseDB:" + uSql)
  var err error
  _, err = this.db.Exec(uSql)
  if err != nil {
    fmt.Println(err)
  } else {
    this.dbName = dbName
  }
}

func (this *MySQL) CreateDBSetChar(dbName string, charType string) {
  createDBSql := "CREATE DATABASE " + dbName
  if len(charType) != 0 {
    createDBSql += " CHARACTER SET " + charType
  }
  fmt.Println("CreateDBSetChar:" + createDBSql)
  var err error
  _, err = this.db.Exec(createDBSql)
  if err != nil {
    fmt.Println(err)
  }
  this.UseDB(dbName)
}

func (this *MySQL) CreateDB(dbName string) {
  this.CreateDBSetChar(dbName, "")
}
 
func (this *MySQL) CreateTable(tableName string, columns []ColumnDefinition) {
  cSql := "("
  var pkeys []string
  for _, column := range columns {
    cSql += column.name + " " + column.dataType
    if column.primaryKey {
      cSql += " NOT NULL"
      pkeys = append(pkeys, column.name)
    }
    cSql += ","
  }
  cSql = strings.Trim(cSql, ",")
  cSql += ")"
  cSql = "CREATE TABLE " + tableName + cSql
  fmt.Println("CreateTable:" + cSql)
  var cterr error
  _, cterr = this.db.Exec(cSql)
  if cterr != nil {
    fmt.Println(cterr)
  } else if len(pkeys) != 0 {
    pSql := "ALTER TABLE " + tableName + " ADD PRIMARY KEY("
    for _, key := range pkeys {
      pSql += key + ","
    }
    pSql = strings.Trim(pSql, ",")
    pSql += ")"
    fmt.Println("CreateTable: " + pSql)
    _, perr := this.db.Exec(pSql)
    if perr != nil {
      fmt.Println(perr)
    }
  }
}

func (this *MySQL) Close() {
  var err error
  err = this.db.Close()
  if err != nil {
    log.Fatal(err)
  }
}

func (this *MySQL) Insert (tableName string, columns []Column) {
  iSql := "INSERT INTO " + tableName + " "

  names := "("
  for _, column := range columns {
    names += column.name + ","
  }
  names = strings.Trim(names, ",")
  names += ")"

  values := " VALUES ("
  for _, column := range columns {
    values += column.value + ","
  }
  values = strings.Trim(values, ",")
  values += ")"

  iSql += names + values

  // fmt.Println("Insert:" + iSql)
  var err error
  _, err = this.db.Exec(iSql)
  if err != nil {
      fmt.Print("insert error: ")
      fmt.Println(err)
  }
}

func (this *MySQL) Update(tableName string, setColumn string, setValue string, whereColumn string, whereValue string) {
  uSql := "UPDATE " + tableName + " SET " + setColumn + " = " + setValue + " WHERE " + whereColumn + " = " + whereValue
  fmt.Println("Update:" + uSql)
  var err error
  _, err = this.db.Exec(uSql)
  if err != nil {
    fmt.Print("update error: ")
    fmt.Println(err)
  }
}

func (this *MySQL) Delete(tableName string, whereColumn string, whereValue string) {
  dSql := "DELETE FROM " + tableName + " WHERE " + whereColumn + " = " + whereValue
  fmt.Println("Delete:" + dSql)
  var err error
  _, err = this.db.Exec(dSql)
  if err != nil {
    fmt.Print("delete error: ")
    fmt.Println(err)
  }
}

func (this *MySQL) Query(qSql string) *sql.Rows {
  rows, err := this.db.Query(qSql)
  if err != nil {
    fmt.Print("query error: ")
    fmt.Println(err)
  }
  return rows;
}

func (this *MySQL) QueryRow(qSql string) *sql.Row {
  row := this.db.QueryRow(qSql)
  return row;
}

func (this *MySQL) FetchAll(qSql string) [][]interface{} {
  rows, err := this.db.Query(qSql)
  if err != nil {
    fmt.Print("fetch all error: ")
    fmt.Println(err)
  }
  defer rows.Close()

  columns, _ := rows.Columns()
  count := len(columns)
  valuePtrs := make([]interface{}, count)

  ret := make([][]interface{}, 0)
  for rows.Next() {

    values := make([]interface{}, count)
    for i, _ := range columns {
      valuePtrs[i] = &values[i]
    }
    rows.Scan(valuePtrs...)

    for i, _ := range columns {
      var v interface{}
      val := values[i]
      b, ok := val.([]byte)
      if (ok) {
        v = string(b)
      } else {
        v = val
      }
      values[i] = v
    }
    ret = append(ret, values)
  }

  return ret;
}
