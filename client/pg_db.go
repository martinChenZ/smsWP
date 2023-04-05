package client

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

//var DB *xorm.Engine

func GetConnect() (*xorm.Engine, error) {
	var engine *xorm.Engine
	//连接数据库
	//engine, err := xorm.NewEngine("mysql", "root:Transfar@2022@tcp(localhost:53306)/gpt?charset=utf8")
	engine, err := xorm.NewEngine("mysql", "root:root@tcp(localhost:53306)/gpt?charset=utf8")
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
