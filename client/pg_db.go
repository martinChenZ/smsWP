package client

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func GetConnect() (*xorm.Engine, error) {
	var engine *xorm.Engine
	//连接数据库
	//engine, err := xorm.NewEngine("mysql", "root:Transfar@2022@tcp(localhost:53306)/gpt?charset=utf8")
	engine, err := xorm.NewEngine("sqlite3", "./db/gpt")
	if err != nil {
		fmt.Println(err)
		return nil, err

	}

	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	//defer engine.Close() //延迟关闭数据库
	fmt.Println("数据库链接成功")
	return engine, err
}
