package main

import (
 "fmt"
 "./mysql"
)

const tableName = "table1"

func main() {
  fmt.Println("** mysql start ** ")
  
  db := mysql.MySQL
  db.ConnectServer()
  db.CreateDB("test1")
  columnDefs []ColumnDefinition
  columnDefs = append(columnDefs, NewColumnDefinition("id", "int", true)
  columnDefs = append(columnDefs, NewColumnDefinition("name", "char(50)", false)
  db.CreateTable(tableName, columnDefs)
  data0 := {{"0"},{"user0"}}
  db.Insert(tableName, data0)
  data2 := {{"2"},{"user2"}}
  db.Insert(tableName, data2)
  db.Update(tableName, "name", "up1", "id", "2")
  db.Delete(tableName, "id", "0")
  defer db.Close()
  fmt.Println("** mysql end ** ")
}
