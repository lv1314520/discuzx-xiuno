package lib

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

type Database struct {
	DBHost,
	DBUser,
	DBPassword,
	DBName,
	DBPort,
	DBChar,
	DSN string
}

func (this *Database) Setting() {
	buf := bufio.NewReader(os.Stdin)
	db := this

	var process = 0 //流程 ID
	for {
		switch process {
		case 0:
			fmt.Println("请配置数据库主机:(默认为 127.0.0.1)")
			s := Input(buf)
			if s == "" {
				s = "127.0.0.1"
			}
			db.DBHost = s
			process++
			break

		case 1:
			fmt.Println("请配置数据库用户:(默认为 root)")
			s := Input(buf)
			if s == "" {
				s = "root"
			}
			db.DBUser = s
			process++
			break

		case 2:
			fmt.Println("请配置数据库密码:")
			db.DBPassword = Input(buf)
			process++
			break

		case 3:
			fmt.Println("请配置数据库名:")
			db.DBName = Input(buf)
			process++
			break

		case 4:
			fmt.Println("请配置数据库端口:(默认为 3306)")
			s := Input(buf)
			if s == "" {
				s = "3306"
			}
			db.DBPort = s
			process++
			break

		case 5:
			fmt.Println("请配置数据库字符集:(默认为 utf8)")
			s := Input(buf)
			if s == "" {
				s = "utf8"
			}
			db.DBChar = s
			process++
			break

		default:
			break
		}

		//跳出 for
		if process > 5 {
			break
		}
	}
	fmt.Println("您配置的数据库信息: ", db)
}

/**
  连接数据库
*/
func (this *Database) MakeDSN() {
	if this.DBHost == "" {
		this.DBHost = "127.0.0.1"
	}

	if this.DBUser == "" {
		this.DBUser = "root"
	}

	if this.DBPort == "" {
		this.DBPort = "3306"
	}

	if this.DBChar == "" {
		this.DBChar = "utf8"
	}

	host := ""
	if this.DBHost != "" {
		host = "tcp(" + this.DBHost + ":" + this.DBPort + ")"
	}

	this.DSN = fmt.Sprintf("%s:%s@%s/%s?%s",
		this.DBUser,
		this.DBPassword,
		host,
		this.DBName,
		this.DBChar,
	)
}

func (this *Database) Connect() (db *sql.DB, err error) {
	this.MakeDSN()
	if this.DSN != "" {
		return sql.Open("mysql", this.DSN)
	}

	return nil, errors.New("数据库连接失败")
}

/**
  数据库字段批量加前缀
*/
func (this *Database) FieldAddPrev(prev, fieldStr string) string {
	fieldArr := strings.Split(fieldStr, ",")

	prev = prev + "."
	var newFieldArr []string
	for _, v := range fieldArr {
		newFieldArr = append(newFieldArr, prev+v)
	}
	newFieldStr := strings.Join(newFieldArr, ",")

	return newFieldStr
}

func (this *Database) FieldMakeQmark(str string, symbol string) string {
	strArr := strings.Split(str, ",")
	strLen := len(strArr)

	if symbol == "" {
		symbol = "?"
	}

	arr := make([]string, strLen)
	for i := 0; i < strLen; i++ {
		arr[i] = symbol
	}
	return strings.Join(arr, ",")
}

func (this *Database) FieldMakeValue(str string) string {
	strArr := strings.Split(str, ",")
	strLen := len(strArr)

	arr := make([]string, strLen)
	for i := 0; i < strLen; i++ {
		arr[i] = "'%s'"
	}
	return strings.Join(arr, ",")
}
