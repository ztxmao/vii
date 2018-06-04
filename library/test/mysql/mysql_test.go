package mysql_test

import (
	"fmt"
	"library"
	"testing"
	"time"
)

var (
	p = fmt.Println
)

func TestMysql(t *testing.T) {

	dbname := "videocore_online"
	dbuser := "mysqldev"
	dbpwd := "daohang157dbpasswd"
	dbhost := "10.16.15.157"
	dbport := 3306

	op := make(map[string]interface{})
	op["dbname"] = dbname
	op["dbuser"] = dbuser
	op["dbpwd"] = dbpwd
	op["dbhost"] = dbhost
	op["dbport"] = dbport

	var ops []map[string]interface{}
	ops = append(ops, op)
	fmt.Println(ops)

	fmt.Println(fmt.Sprintf("%v", time.Second))

	var sql string
	mysql := library.NewMysqlExecutor(dbname, dbuser, dbpwd, dbhost, dbport)
	sql = "select * from test limit 0, 5"
	rows := mysql.Query(sql)
	fmt.Println("query result:", rows)

	fmt.Println("---------------")
	row := mysql.QueryRow(sql)
	fmt.Println(row)

	var flag int
	sql = "insert into test set entid=?, title=?,status=?"
	flag = mysql.ExecSql(sql, 3, "阅兵0904", 1)
	fmt.Println("insert result:", flag)

	sql = "update test set title=? where id=?"
	flag = mysql.ExecSql(sql, "测试update", 7)
	fmt.Println("update result:", flag)

	sql = "delete from test where id=?"
	flag = mysql.ExecSql(sql, 6)
	fmt.Println("delete result:", flag)
}
