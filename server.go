package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DbWorker struct {
	//mysql data source name
	Dsn      string
	Db       *sql.DB
	UserInfo userTB
}
type userTB struct {
	Id1 sql.NullString
	Id2 sql.NullString
	Id3 sql.NullString
	Id4 sql.NullString
	Id5 sql.NullString
	Id6 sql.NullString
	Id7 sql.NullString
	Id8 sql.NullString
	Id9 sql.NullString
}

var chdb = make(chan DbWorker)
var db = make(chan DbWorker)
var dbw DbWorker

func main() {
	go dbConnectInit()

	listener, err := net.Listen("tcp", "10.101.171.173:38250")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(conn net.Conn) {
	var ch = make(chan []byte)
	go getDb()
	dbw := <-db
	buffer := make([]byte, 30)
	go clientWriter(conn, ch)
	conn.Read(buffer)
	for i, b := range buffer {
		//fmt.Printf("i:%v,b:%v", i, b)
		if b == 0 {
			buffer = buffer[0:i]
			break
		}
	}
	para := strings.Split(string(buffer), " ")
	dbw.queryData(para, ch)
	close(ch)
}
func clientWriter(conn net.Conn, ch <-chan []byte) {
	for msg := range ch {
		conn.Write(msg)
	}
	conn.Close()
}
func (dbw *DbWorker) QueryDataPre() {
	dbw.UserInfo = userTB{}
}
func dbConnectInit() {
	var err error
	dbw = DbWorker{
		Dsn: "root:qweasdzxc1@tcp(127.0.0.1:3306)/hgyd",
	}
	dbw.Db, err = sql.Open("mysql", dbw.Dsn)
	if err != nil {
		panic(err)

	}
}
func getDb() {
	for {
		db <- dbw
	}
}
func (dbw *DbWorker) queryData(para []string, ch chan<- []byte) {
	sql := ""
	time_s := "201805"
	time_e := "201903"
	if len(para) == 2 {
		sql = `SELECT * From ` + para[1] + ` where time >= ? AND time < ?`
	} else if len(para) == 3 {
		if para[2][0:4] != "last" {
			sql = `SELECT time,` + para[2] + ` From ` + para[1] + ` where time >= ? AND time < ?`
		} else {
			time_s = timeParse(para[2][4:])
			sql = `SELECT * From ` + para[1] + ` where time >= ? AND time < ?`
		}
	} else if len(para) == 4 {
		if para[3][0:4] == "last" {
			time_s = timeParse(para[3][4:])
		}
		sql = `SELECT time,` + para[2] + ` From ` + para[1] + ` where time >= ? AND time < ?`
	}
	fmt.Printf("  " + time_s + ",")
	stmt, _ := dbw.Db.Prepare(sql)
	defer stmt.Close()

	dbw.QueryDataPre()
	fmt.Printf(sql)
	rows, err := stmt.Query(time_s, time_e)
	cols, _ := rows.Columns()
	defer rows.Close()
	if err != nil {
		fmt.Printf("query data error: %v\n", err)
		return
	}
	str := ""
	onerow := make([]interface{}, len(cols))
	values := make([][]byte, len(cols))
	for i := range values {
		onerow[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(onerow...)
		fmt.Printf("  %d", len(onerow))
		//err := rows.Scan(&dbw.UserInfo.Id1, &dbw.UserInfo.Id2, &dbw.UserInfo.Id3, &dbw.UserInfo.Id4, &dbw.UserInfo.Id5, &dbw.UserInfo.Id6, &dbw.UserInfo.Id7, &dbw.UserInfo.Id8, &dbw.UserInfo.Id9)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		str += "{"
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			if k != 0 {
				str += ","
			}
			str += "\"" + cols[k] + "\":\"" + string(v) + "\""
		}
		str += "}"
		//fmt.Println("get data, id: ", dbw.UserInfo.Id1.String, " name: ", dbw.UserInfo.Id2.String, " age: ", dbw.UserInfo.Id3.String)
	}

	fmt.Println(str)
	err = rows.Err()
	if err != nil {
		fmt.Printf(err.Error())
	}
	ch <- []byte(str)
}
func timeParse(para string) string {
	tnum, _ := strconv.Atoi(para)
	year := tnum / 12
	month := tnum % 12
	time_s := ""
	if month < 3 {
		time_s = strconv.Itoa(2019-year) + "0" + strconv.Itoa(3-month)
	} else {
		year += 1
		if 15-month < 10 {
			time_s = strconv.Itoa(2019-year) + "0" + strconv.Itoa(15-month)
		} else {
			time_s = strconv.Itoa(2019-year) + strconv.Itoa(15-month)
		}
	}
	return time_s
}
