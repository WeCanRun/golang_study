package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"study/tools/internal/sql2go"
)

var (
	host      string
	username  string
	password  string
	dbType    string
	dbName    string
	tableName string
	charset   string
)

func init() {
	sql2GoCmd.Flags().StringVarP(&host, "host", "", "localhost:3306", "请输入数据库服务地址")
	sql2GoCmd.Flags().StringVarP(&username, "username", "u", "root", "请输入数据库用户名")
	sql2GoCmd.Flags().StringVarP(&password, "password", "p", "root", "请输入数据库密码")
	sql2GoCmd.Flags().StringVarP(&dbType, "type", "T", "mysql", "请输入数据库类型")
	sql2GoCmd.Flags().StringVarP(&dbName, "db", "d", "mysql", "请输入数据库名")
	sql2GoCmd.Flags().StringVarP(&tableName, "table", "t", "", "请输入数据表名")
	sql2GoCmd.Flags().StringVarP(&charset, "charset", "c", "utf8mb4", "请输入数据库字符集")
	sqlCmd.AddCommand(sql2GoCmd)
}

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "sql 转 go 结构体",
	Long:  "sql 转 go 结构体",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var sql2GoCmd = &cobra.Command{
	Use:   "2go",
	Short: "sql 转结构体",
	Long:  "sql 转结构体",
	Run: func(cmd *cobra.Command, args []string) {
		info := &sql2go.DBInfo{
			DBType:   dbType,
			Host:     host,
			UserName: username,
			Password: password,
			Charset:  charset,
		}
		log.Println("begin to connect db: ", host)

		model := sql2go.NewDBModel(info)
		if err := model.Connect(); err != nil {
			log.Fatalf("connect db fail, err: %v, host: %s, username: %s, password: %s, dbName: %s, dbName: %s\n",
				err, host, username, password, dbType, dbName)
		}

		log.Printf("connect db [%s] success", host)

		columns, err := model.GetColumns(dbName, tableName)
		if err != nil {
			log.Fatalf("get columns fail, err: %v\n", err)
		}

		log.Printf("columns: %v", columns)

		t := sql2go.NewStructTemplate()
		if err = t.Generate(dbName, tableName, columns); err != nil {
			log.Fatalf("Generate columns fail, tableName: %s, err: %v\n", tableName, err)
		}
	},
}
