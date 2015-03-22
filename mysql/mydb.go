package mysql

import (
  "os"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "log"
)

/* 起動前に変更 */
const User = "user"
const Pass = "pass"
const Addr = os.Getenv("$TEST_DB_PORT_3306_TCP_ADDR")
const Port = 3306

type MySQL struct {
  db *sql.DB
  dbName string
  tableName string
}

type ColmunDefinition struct {
  name string
  dataType string
  primaryKey bool
}

func NewColumnDefinition(n string, dt string, pk bool) *ColumnDefinition {
  columnDef := &Column{name: n, dataType: dt, primaryKey: pk}
  return columnDef
}

type ColumnDefArray []ColumnDefinition

func (cdArray *ColumnDefArray) GetCreateTableQuery string {
}

func (cdArray *ColumnDefArray) GetPrimaryKeyQuery string {
}

func (cdArray *ColumnDefArray) ExistPraimaryKey bool {
  
}

func (this *MySQL) ConnectServer() {
  var err error
  this.db, err = sql.Open("mysql", User + ":" + Pass + "@(" + Addr ":" + Port + ")/")
  if err != nil {
    log.Fatal(err)
  }
}

func (this *MySQL) CreateDB(dbName string) {
  CreateDBSetChar(dbName, nil)
}

func (this *MySQL) CreateDBSetChar(dbName string, charType string) {
  query := "CREATE DATABASE " + dbName
  if charSet != nil {
    query += " CHARACTER SET " + charType
  }
  var err error
  this.db, err = sql.Exec(query)
  if err != nil {
    log.Fatal(err)
  }
  UseDB(dbName)
}

func (this *MySQL) UseDB(dbName string) {
  this.db, err = sql.Exec("USE " + dbName)
  if err != nil {
    log.Fatal(err)
  } else {
    this.DBName = dbName
  }
}

func (this *MySQL) CreateTable(tableName string, columns ColumnDefArray) {
  query := "(" + columns.GetCreateTableQuery + ")"
  this.db, cterr = sql.Exec("CREATE TABLE " + tableName + query)
  if cterr != nil {
    log.Fatal(cterr)
  }
}

func (this *MySQL) Close() {
  var err error
  err = this.db.Close()
  if err != nil {
    log.Fatal(err)
  }
}

func (this *MySQL) Query(sql string) *sql.Rows {
  rows, err := this.db.Query(sql)
  if err != nil {
    log.Fatal(err)
  }
  return rows;
}

func (this *MySQL) QueryRow(sql string) *sql.Row {
  row := this.db.QueryRow(sql)
  return row;
}

func (this *MySQL) FetchAll(sql string) [][]interface{} {
  rows, err := this.db.Query(sql)
  if err != nil {
    log.Fatal(err)
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
