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
const Port = 3306

type MySQL struct {
  db *sql.DB
  dbName string
}

type ColumnDefinition struct {
  name string
  dataType string
  primaryKey bool
}

func NewColumnDefinition(n string, dt string, pk bool) *ColumnDefinition {
  columnDef := &Column{name: n, dataType: dt, primaryKey: pk}
  return columnDef
}

func (this *MySQL) ConnectServer() {
  var err error
  this.db, err = sql.Open("mysql", User + ":" + Pass + "@(" + os.Getenv("$TEST_DB_PORT_3306_TCP_ADDR") + ":" + Port + ")/")
  if err != nil {
    log.Fatal(err)
  }
}

func (this *MySQL) UseDB(dbName string) {
  this.db, err = sql.Exec("USE " + dbName)
  if err != nil {
    log.Fatal(err)
  } else {
    this.dBName = dbName
  }
}

func (this *MySQL) CreateDBSetChar(dbName string, charType string) {
  createDBSql := "CREATE DATABASE " + dbName
  if charSet != nil {
    createDBSql += " CHARACTER SET " + charType
  }
  var err error
  this.db, err = sql.Exec(createDBSql)
  if err != nil {
    log.Fatal(err)
  }
  UseDB(dbName)
}

func (this *MySQL) CreateDB(dbName string) {
  CreateDBSetChar(dbName, nil)
}

func (this *MySQL) CreateTable(tableName string, columns []ColumnDefinition) {
  cSql = "("
  pkeys []string
  for _, column := range columns {
    cSql += column.name + " " + column.dataType
    if column.primaryKey {
      cSql += " NOT NULL"
      pkeys = append(pkeys, column.name)
    }
    cSql += ","
  }
  cSql = Trim(cSql, ",")
  cSql += ") ENGINE=InnoDB"
  
  this.db, cterr = sql.Exec("CREATE TABLE " + tableName + cSql)
  if cterr != nil {
    log.Fatal(cterr)
  }
  
  if len(pkeys) != 0 {
    pSql := "ALTER TABLE " + tableName + " ADD PRIMARY KEY("
    for _, key := range pkeys {
      pSql += key ","
    }
    pSql = Trim(pSql, ",")
    pSql += ")"
    this.db, pkerr = sql.Exec(pSql)
  }
}

func (this *MySQL) Close() {
  var err error
  err = this.db.Close()
  if err != nil {
    log.Fatal(err)
  }
}

func (this *MySQL) Insert (tableName string, datas []string) {
  iSql := "INSERT INTO " + tableName + " "
  params := "values("
  for _, data := range datas {
    params += data + ","
  }
  params = Trim(params, ",")
  params += ")"
  iSql += params
  this.db, err = sql.Exec(iSql)
  if err != nil {
      log.Fatal("insert error: ", err)
  }
}

func (this *MySQL) Update(tableName string, setColumn string, setValue string, whereColumn string, whereValue string) {
  uSql := "UPDATE " + tableName + " SET " + setColumn + "=" + setValue + " WHERE " + whereColumn + "=" + whereValue
  this.db, err = sql.Exec(uSql)
  if err != nil {
    log.Fatal("update error: ", err)
  }
}

func (this *MySQL) Delete(tableName string, whereColumn string, whereValue string) {
  dSql := "DELETE FROM " + tableName + " " + whereColumn + "=" + whereValue
  this.db, err = sql.Exec(dSql)
  if err != nil {
    log.Fatal("delete error: ", err)
  }
}

func (this *MySQL) Query(qSql string) *sql.Rows {
  rows, err := this.db.Query(qSql)
  if err != nil {
    log.Fatal(err)
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
