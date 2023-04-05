package main

import (
	"flag"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	uuid "github.com/satori/go.uuid"
	"gpt-backend/business"
	"gpt-backend/client"
	"gpt-backend/remote"
	"strconv"
	"time"
)

type RegisterRes struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

type GptResp struct {
	Messages []remote.ChoiceMsg `json:"messages"`
	Reply    string             `json:"reply"`
	ErrorMsg string             `json:"errorMsg"`
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
			Code: 0,
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
	defer session.Close()
	id := uuid.NewV4()
	_, err = session.Insert(&business.GptUser{
		id.String(),
		user.OrderId,
		3,
		time.Now(),
	})
	if err != nil {
		fmt.Println("新增错误", err)
	}

}

func main() {
	//
	//http.HandleFunc("/compotion", Compotion)
	//
	//http.ListenAndServe(":8086", nil)
	var key string
	//flag.StringVar(&name, "name", "Go语言编程之旅", "帮助")
	flag.StringVar(&key, "k", "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", "帮助")

	var port string
	//flag.StringVar(&name, "name", "Go语言编程之旅", "帮助")
	flag.StringVar(&port, "p", "8086", "帮助")

	flag.Parse()
	remote.Appkey = key
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
