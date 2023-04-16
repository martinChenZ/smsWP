package main

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/sashabaranov/go-openai"
	uuid "github.com/satori/go.uuid"
	"gpt-backend/business"
	"gpt-backend/client"
	"gpt-backend/remote"
	"strconv"
	"time"
)

var (
	AccCode    string
	FirstTimes int
)

func main() {
	//
	//http.HandleFunc("/compotion", Compotion)
	//
	//http.ListenAndServe(":8086", nil)
	var key string
	flag.StringVar(&key, "k", "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", "帮助")

	var port string
	flag.StringVar(&port, "p", "8086", "帮助")

	var access string
	flag.StringVar(&access, "a", "Gpt_key@!@#", "帮助")

	var freeTimes int64
	flag.Int64Var(&freeTimes, "t", 9, "帮助")

	var firstTimes int
	flag.IntVar(&firstTimes, "f", 15, "帮助")

	flag.Parse()
	remote.Appkey = key
	AccCode = access
	business.FreeT = freeTimes
	FirstTimes = firstTimes
	s := g.Server()
	s.BindHandler("/completion", Completion)
	s.BindHandler("/mykey", QtyApi)
	s.BindHandler("/AddKey", AddKey)

	atoi, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("参数异常")
		return
	}
	s.SetPort(atoi)
	s.Run()
}

type RegisterRes struct {
	Code  int         `json:"code"`
	Error string      `json:"errorMsg"`
	Data  interface{} `json:"data"`
}

type GptResp struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
	Reply    string                         `json:"reply"`
	ErrorMsg string                         `json:"errorMsg"`
}

func Completion(r *ghttp.Request) {
	var q *business.Question
	if err := r.Parse(&q); err != nil {
		r.Response.WriteJsonExit(RegisterRes{
			Code:  1,
			Error: err.Error(),
		})
	}
	ip := r.GetClientIp()
	fmt.Println("clientIp: ", ip)
	ip = r.GetRemoteIp()
	fmt.Println("remoteIp: ", ip)
	gpt := business.CallGpt(q, ip)
	r.Response.WriteJsonExit(RegisterRes{
		Data: &GptResp{
			q.Messages,
			gpt,
			"",
		},
	})
}

func QtyApi(r *ghttp.Request) {
	var user *business.GptUser
	if err := r.Parse(&user); err != nil {
		r.Response.WriteJsonExit(RegisterRes{
			Code:  1,
			Error: err.Error(),
		})
	}
	session, err := client.GetConnect()
	if err != nil {
		fmt.Println("数据库连接异常", err)
	}
	has, err := session.Where("order_no=?", user.OrderId).Get(user)
	if err != nil {
		r.Response.WriteJsonExit(RegisterRes{
			Code:  1,
			Error: err.Error(),
		})
	}
	if has {
		r.Response.WriteJsonExit(RegisterRes{
			Code: 200,
			Data: &map[string]string{"key": user.ApiKey},
		})
	} else {
		r.Response.WriteJsonExit(RegisterRes{
			Code: 1,
			Data: "无数据",
		})
	}

}

func AddKey(r *ghttp.Request) {
	var user *business.GptUserReq
	if err := r.Parse(&user); err != nil {
		r.Response.WriteJsonExit(RegisterRes{
			Code:  1,
			Error: err.Error(),
		})
	}
	if user.Access != AccCode {
		fmt.Println("req: ", user.Access)
		fmt.Println("server: ", AccCode)
		r.Response.WriteJsonExit(RegisterRes{
			Code: 0,
			Data: "无权限",
		})
		return
	}

	session, err := client.GetConnect()
	session.ShowSQL(true)
	if err != nil {
		fmt.Println("数据库连接异常", err)
	}
	defer session.Close()
	// 已存在的apikey 新增
	if user.ApiKey != "" {
		curUser := &business.GptUser{}
		session.Where("api_key = '" + user.ApiKey + "'").Get(curUser)
		curUser.Balance = user.Balance
		session.Where("api_key = '" + user.ApiKey + "'").Update(curUser)
		r.Response.WriteJsonExit(RegisterRes{
			Code: 0,
			Data: "更新成功",
		})
		return
	}
	count, _ := session.Table("gpt_user").Where("order_no = " + user.OrderId).Count()
	if count > 1 {
		fmt.Println("重复订单id", user.OrderId)
		r.Response.WriteJsonExit(RegisterRes{
			Code: 1,
			Data: "重复订单id",
		})
		return
	}

	id := uuid.NewV4()
	_, err = session.Insert(&business.GptUser{
		id.String(),
		user.OrderId,
		FirstTimes,
		time.Now(),
	})
	if err != nil {
		fmt.Println("新增错误", err)
	}
	r.Response.WriteJsonExit(RegisterRes{
		Code: 0,
		Data: "新增成功",
	})

}
