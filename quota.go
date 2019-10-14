package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/k0kubun/pp"
)

var cmdQuota = &Command{
	Run:       runQuota,
	UsageLine: "quota ",
	Short:     "qouta",
	Long: `

	`,
}

// runQuota executes quota command and return exit code.
func runQuota(args []string) int {
	pp.Println(runq())
	return 0
}

const dbSizeQuery = `SELECT
    table_schema, sum(data_length+index_length) /1024 /1024 AS MB
FROM
    information_schema.tables
GROUP BY
    table_schema
`
const enableQuery = `UPDATE db SET Insert_priv='Y', Create_priv='Y' WHERE Db=?`
const disableQuery = `UPDATE db SET Insert_priv='N', Create_priv='N' WHERE Db=?`

func runq() error {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		"root",
		"",
		"localhost",
		3306,
		"mysql",
	))
	if err != nil {
		return err
	}

	err = DBSumPerUser(db)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	Name    string
	DBs     []string
	LimitMB float64
	DBHost  string
}

func DBSumPerUser(db *gorm.DB) error {
	ret := map[string]float64{}
	users := map[string]User{
		"test-user": User{
			Name:    "test-user",
			DBs:     []string{"test"},
			LimitMB: 2000,
			DBHost:  "localhost",
		}}

	rows, err := db.Raw(dbSizeQuery).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var db string
		var size float64
		if err := rows.Scan(&db, &size); err != nil {
			return err
		}

		if contains(users, db) {
			ret[users["test-user"].Name] += size
		}
	}

	for user, size := range ret {
		query := ""
		if users[user].LimitMB < size {
			query = disableQuery
		} else {
			query = enableQuery
		}
		err := do(db, query, users[user].DBs)
		if err != nil {
			return err
		}
	}
	return nil
}

func do(db *gorm.DB, query string, dbs []string) error {
	for _, dbn := range dbs {
		if err := db.Exec(query, dbn).Error; err != nil {
			return err
		}
	}
	return nil
}

func contains(users map[string]User, e string) bool {
	for _, u := range users {
		for _, v := range u.DBs {
			if e == v {
				return true
			}
		}
	}
	return false
}
