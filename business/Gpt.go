package business

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"gpt-backend/client"
	"gpt-backend/remote"
	"time"
)

var FreeT int64

type Question struct {
	Appkey   string             `json:"privateKey"`
	Messages []remote.ChoiceMsg `json:"messages" `
}

type GptUserReq struct {
	Access string
	GptUser
}

type GptUser struct {
	ApiKey     string
	OrderId    string `xorm:"VARCHAR(30) 'Order_no'"`
	Balance    int
	UpdateTime time.Time
}

type GptLog struct {
	Id         int64
	ApiKey     string
	Question   string
	Response   string
	RequestIp  string
	UpdateTime time.Time
}

func CallGpt(q *Question, ip string) string {
	connect, err := client.GetConnect()
	if err != nil {
		fmt.Println()
		return ""
	}
	defer connect.Close()
	msg := "请求失败"
	if "" == q.Appkey {
		q.Appkey = "666"
	}
	g := GptUser{ApiKey: q.Appkey}

	log := &GptLog{Question: q.Messages[len(q.Messages)-1].Content}
	_, err = connect.Cols("balance").Get(&g)
	if g.Balance < 1 {
		msg = validNoApiBySql(&g, ip, connect)
		if "" != msg {
			return msg
		}
		msg = remote.CallGpt3(q.Messages)
		log.RequestIp = ip
		log.ApiKey = g.ApiKey
		log.Response = msg
		log.UpdateTime = time.Now()
	} else {

		msg = remote.CallGpt3(q.Messages)
		fmt.Println(msg)

		_, err = connect.Exec(fmt.Sprintf("update gpt_user set balance = %v where api_key = '%v' ", g.Balance-1, g.ApiKey))
		if err != nil {
			fmt.Println("更新失败", err)
		}
		log.ApiKey = g.ApiKey
		log.Response = msg
		log.UpdateTime = time.Now()
	}
	//msg := "hell"
	_, err = connect.Insert(log)
	if err != nil {
		fmt.Println("新增日志错误", err)
	}

	return msg
}

func validNoApiBySql(g *GptUser, ip string, session *xorm.Engine) (msg string) {
	session.ShowSQL(true)
	//exec, err := session.Exec("select count(1) from gpt_log where  request_ip = " + ip + " and update_time > now() + INTERVAL 1 Day ")
	count, err := session.Table("gpt_log").Where("request_ip = '" + ip + "' and update_time > date('now','-1 day') ").Count()
	if err != nil {
		return ""
	}
	fmt.Println("@@", count)
	if count > FreeT {
		msg = "无请求次数"
	}
	return
}
