package main

import (
 "fmt"
 "./mysql"
)

const tableName = "table1"

func main() {
  fmt.Println("** mysql start ** ")
  
  var m mysql.MySQL
  m.ConnectServer()
  m.CreateDB("test1")
  var columnDefs []mysql.ColumnDefinition
  var columnd1 mysql.ColumnDefinition
  var columnd2 mysql.ColumnDefinition
  columnd1 = *mysql.NewColumnDefinition("id", "int", true)
  columnd2 = *mysql.NewColumnDefinition("name", "char(50)", false)
  columnDefs = append(columnDefs, columnd1, columnd2)
  m.CreateTable(tableName, columnDefs)
  var column0 []mysql.Column
  var column2 []mysql.Column
  column0 = append(column0, *mysql.NewColumn("id", "0"), *mysql.NewColumn("name", "'user0'"))
  column2 = append(column2, *mysql.NewColumn("id", "2"), *mysql.NewColumn("name", "'user2'"))
  m.Insert(tableName, column0)
  m.Insert(tableName, column2)
  m.Update(tableName, "name", "'up1'", "id", "2")
  m.Delete(tableName, "id", "0")
  defer m.Close()
  fmt.Println("** mysql end ** ")
}
