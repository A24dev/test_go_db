package main

import (
 "time"
 "fmt"
 "strconv"
 "./mysql"
)

const databaseName = "tdb"
const tableName = "ttable"

func main() {
  fmt.Println("** mysql start ** ")
  
  var m mysql.MySQL
  m.ConnectServer()
  m.CreateDB(databaseName)
  
  var columnDefs []mysql.ColumnDefinition
  var columnd1 mysql.ColumnDefinition
  var columnd2 mysql.ColumnDefinition
  columnd1 = *mysql.NewColumnDefinition("id", "int", true)
  columnd2 = *mysql.NewColumnDefinition("name", "char(50)", false)
  columnDefs = append(columnDefs, columnd1, columnd2)
  m.CreateTable(tableName, columnDefs)
  
  startTime := time.Nanoseconds()
  for i := 0; i < 10000; ++i {
    var column []mysql.Column
    id := strconv.Itoa(i)
    name := "'user" + id + "'"
    column = append(column, *mysql.NewColumn("id", id), *mysql.NewColumn("name", name))
    m.Insert(tableName, column)
  }
  endTime := time.Nanoseconds()
  
  fmt.Print("time = ")
  fmt.Print(endTime - startTime)
  fmt.Println("[nsec]")
  
  m.Update(tableName, "name", "'up1'", "id", "2")
  
  m.Delete(tableName, "id", "0")
  
  defer m.Close()
  fmt.Println("** mysql end ** ")
}
