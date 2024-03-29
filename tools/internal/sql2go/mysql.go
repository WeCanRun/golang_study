package sql2go

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func NewDBModel(info *DBInfo) *dbModel {
	return &dbModel{
		DBInfo: info,
	}
}

type DBInfo struct {
	DBType   string
	Host     string
	UserName string
	Password string
	Charset  string
}

type TableColumn struct {
	ColumnName    string
	DataType      string
	IsNullable    string
	ColumnKey     string
	ColumnType    string
	ColumnComment string
}

type dbModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

func (db *dbModel) Connect() error {
	var err error
	s := fmt.Sprintf("%s:%s@tcp(%s)/information_schema?charset=%s&parseTime=True",
		db.DBInfo.UserName,
		db.DBInfo.Password,
		db.DBInfo.Host,
		db.DBInfo.Charset)

	db.DBEngine, err = sql.Open(db.DBInfo.DBType, s)
	return err
}

func (db *dbModel) GetColumns(dbName, tableName string) (columns []*TableColumn, err error) {
	sql := " select column_name, data_type, column_key, is_nullable, column_type, column_comment " +
		"from information_schema.columns where table_schema = ? and table_name = ? "

	rows, err := db.DBEngine.Query(sql, dbName, tableName)
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return nil, errors.New("have not data")
	}

	defer rows.Close()

	for rows.Next() {
		col := TableColumn{}
		if err := rows.Scan(&col.ColumnName, &col.DataType, &col.ColumnKey, &col.IsNullable, &col.ColumnType,
			&col.ColumnComment); err != nil {
			log.Fatalf("rows.Scan err: %v", err)
		}
		columns = append(columns, &col)
	}

	return
}
