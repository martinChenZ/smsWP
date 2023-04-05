package business

import (
	"fmt"
	"github.com/coocood/freecache"
	"github.com/go-xorm/xorm"
	"github.com/pkg/errors"
	"gpt-backend/client"
	"gpt-backend/remote"
	"math"
	"time"
)

type Question struct {
	Appkey   string             `json:"privateKey"`
	Messages []remote.ChoiceMsg `json:"messages"`
}

type GptUser struct {
	ApiKey     string
	OrderId    string `xorm:"VARCHAR(30) 'Order_no'"`
	Balance    int8
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

var Cache = freecache.NewCache(100 * 1024 * 1024)

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

func validNoApiBycache(g *GptUser, ip string) (msg string) {

	value, err := Cache.Get([]byte(ip))
	if err != nil {
		value, _ := IntToBytes(3)
		Cache.Set([]byte(ip), value, 24*60*60)
		//msg = remote.CallGpt3(q.Messages)
	}
	toInt, _ := BytesToInt(value)
	if toInt < 1 {
		msg = "无请求次数"
	} else {
		Cache.Update([]byte(ip), Update)
		//msg = remote.CallGpt3(q.Messages)
	}

	return
}

func validNoApiBySql(g *GptUser, ip string, session *xorm.Engine) (msg string) {
	session.ShowSQL(true)
	//exec, err := session.Exec("select count(1) from gpt_log where  request_ip = " + ip + " and update_time > now() + INTERVAL 1 Day ")
	count, err := session.Table("gpt_log").Where("request_ip = '" + ip + "' and update_time > now() + INTERVAL -1 Day ").Count()
	if err != nil {
		return ""
	}
	fmt.Println("@@", count)
	if count > 2 {
		msg = "无请求次数"
	}
	return
}

func Update(key []byte, found bool) (newValue []byte, replace bool, expireSeconds int) {
	value, at, _ := Cache.GetWithExpiration(key)
	toInt, _ := BytesToInt(value)
	newValue, _ = IntToBytes(toInt - 1)
	expireSeconds = int(at)
	replace = found
	return
}

func IntToBytes(a int) ([]byte, error) {
	if a > math.MaxInt32 {
		return nil, errors.New(fmt.Sprintf("a>math.MaxInt32, a is %d\n", a))
	}
	buf := make([]byte, 4)
	for i := 0; i < 4; i++ {
		var b uint8 = uint8(a & 0xff)
		buf[i] = b
		a = a >> 8
	}
	return buf, nil
}

func BytesToInt(buf []byte) (int, error) {
	if len(buf) != 4 {
		return -1, errors.New(fmt.Sprintf("BytesToInt len(buf) must be 4, but got %d\n", len(buf)))
	}
	result := 0
	for i := 0; i < 4; i++ {
		result += int(buf[i]) << (8 * i)
	}
	return result, nil
}
