package library

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlExecutor struct {
	db  *sql.DB
	dsn string

	dbname         string
	dbuser         string
	dbpwd          string
	dbport         int
	dbhost         string
	charset        string
	protocol       string
	dboption       string
	connectTimeout time.Duration
}
type DbRow struct {
	id    int
	title string
}

//Create an empty configuration file
func NewMysqlExecutor(dbname, dbuser, dbpwd, dbhost string, dbport int) *MysqlExecutor {
	mysql := new(MysqlExecutor)
	mysql.Init(dbname, dbuser, dbpwd, dbhost, dbport)
	return mysql
}

func (this *MysqlExecutor) Init(dbname, dbuser, dbpwd, dbhost string, dbport int) {
	/*{{{*/
	this.charset = "utf8"
	this.protocol = "tcp"
	this.connectTimeout = time.Millisecond * 500
	this.dboption = ""
	this.SetEnv(dbname, dbuser, dbpwd, dbhost, dbport, this.charset, this.protocol, this.dboption, this.connectTimeout)
} /*}}}*/

func (this *MysqlExecutor) OpenConnet() {
	/*{{{*/
	if this.db == nil {
		this.db = this.CreateDbLink()
	}
	err := this.db.Ping()
	if err != nil {
		this.db = this.CreateDbLink()
	}
} /*}}}*/

func (this *MysqlExecutor) CloseConnect() {
	/*{{{*/
	if this.db != nil {
		err := this.db.Close()
		if err != nil {
			this.checkErr(err)
		}
	}
} /*}}}*/

func (this *MysqlExecutor) ExecSql(sql string, args ...interface{}) int {
	/*{{{*/
	this.OpenConnet()
	op := strings.ToLower(string(sql[0:6]))
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		this.checkErr(err)
	}
	defer stmt.Close()
	execRst, err := stmt.Exec(args...)
	if err != nil {
		this.checkErr(err)
	}
	if op == "insert" {
		rst, err := execRst.LastInsertId()
		if err != nil {
			this.checkErr(err)
		}
		return int(rst)
	} else if op == "delete" || op == "update" {
		rst, err := execRst.RowsAffected()
		if err != nil {
			this.checkErr(err)
		}
		return int(rst)
	}

	return 0
} /*}}}*/
func (this *MysqlExecutor) QueryRow(sql string, args ...interface{}) map[string]interface{} {
	/*{{{*/
	rows := this.Query(sql, args...)
	if len(rows) > 0 {
		return rows[0]
	}

	return nil
} /*}}}*/
func (this *MysqlExecutor) Query(sql string, args ...interface{}) []map[string]interface{} {
	/*{{{*/
	this.OpenConnet()
	stmt, err := this.db.Prepare(sql)
	this.checkErr(err)

	rows, err := stmt.Query(args...)
	this.checkErr(err)

	defer stmt.Close()
	defer rows.Close()

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	scanVals := make([]interface{}, len(columns))
	for i := range scanVals {
		scanArgs[i] = &scanVals[i]
	}

	var result []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for i, col := range scanVals {
			field := columns[i]
			switch col.(type) {
			case []byte:
				record[field] = string(col.([]byte)) //string(col)//string(col.([]byte))
			default:
				record[field] = col
			}
		}
		result = append(result, record)
	}

	return result
} /*}}}*/

func (this *MysqlExecutor) GenDsn() string {
	/*{{{*/
	if len(this.dbname) <= 0 || len(this.dbuser) <= 0 || len(this.dbpwd) <= 0 || len(this.dbhost) <= 0 {
		this.checkErr(errors.New("invalid config: dbname, dbuser, dbpwd, dbhost can not empty"))
	}
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&timeout=%v", this.dbuser, this.dbpwd, this.protocol, this.dbhost, this.dbport,
		this.dbname, this.charset, this.connectTimeout)
	if len(this.dboption) > 0 {
		dsn = fmt.Sprintf("%s&%s", dsn, this.dboption)
	}
	this.dsn = dsn
	return dsn
} /*}}}*/

func (this *MysqlExecutor) SetEnv(dbname, dbuser, dbpwd, dbhost string, dbport int, charset, protocol, dboption string, connectTimeout time.Duration) {
	/*{{{*/
	if len(dbname) > 0 {
		this.dbname = dbname
	}
	if len(dbuser) > 0 {
		this.dbuser = dbuser
	}
	if len(dbpwd) > 0 {
		this.dbpwd = dbpwd
	}
	if dbport > 0 {
		this.dbport = dbport
	}
	if len(dbhost) > 0 {
		this.dbhost = dbhost
	}
	if len(charset) > 0 {
		this.charset = charset
	}
	if len(protocol) > 0 {
		this.protocol = protocol
	}
	if len(dboption) > 0 {
		this.dboption = dboption
	}
	if connectTimeout >= 0 {
		this.connectTimeout = connectTimeout
	}

	this.GenDsn()
} /*}}}*/

func (this *MysqlExecutor) CreateDbLink() *sql.DB {
	/*{{{*/
	db, err := sql.Open("mysql", this.dsn)
	if err != nil {
		this.checkErr(err)
	}
	return db
} /*}}}*/

func (this *MysqlExecutor) checkErr(err error) {
	/*{{{*/
	if err != nil {
		panic("Error is :" + err.Error())
	}
} /*}}}*/
