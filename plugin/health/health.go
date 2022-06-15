package health

import (
	"FvtiTools/config"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var hGo = func(v config.U, dType string) {
	var context string
	var context2 string
	switch dType {
	case "1":
		context = "早报第一次失败，准备十分钟后重报"
		context2 = "早报健康日报成功"
	case "2":
		context = "午报第一次失败，准备十分钟后重报"
		context2 = "午报健康日报成功"
	case "3":
		context = "晚报第一次失败，准备十分钟后重报"
		context2 = "晚报健康日报成功"
	}
	err := HealthGo(v, dType)
	if err != nil {
		pushMsg(v, context, v.User)
		go func() {
			time.Sleep(time.Minute * 10)
			err = HealthGo(v, dType)
			if err != nil {
				pushMsg(v, v.User, "第二次都失败了，没救了，等晚上吧")
				return
			}
		}()
	} else {
		pushMsg(v, context2, v.User)
	}
}

func addCron() {
	s := config.GetSetting()
	c := cron.New()
	log.Println(s)

	_, err := c.AddFunc(s.HealthCron, func() {
		user := config.GetUserAll()
		for _, v := range user {
			hGo(v, "1")
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
	_, err = c.AddFunc(s.HealthCron2, func() {
		user := config.GetUserAll()
		for _, v := range user {
			hGo(v, "2")
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
	_, err = c.AddFunc(s.HealthCron3, func() {
		user := config.GetUserAll()
		for _, v := range user {
			hGo(v, "3")
		}
	})
	if err != nil {
		log.Println(err)
		return
	}

	if s.HealthCron4 != "0" {
		_, err = c.AddFunc(s.HealthCron4, func() {
			user := config.GetUserAll()
			for _, v := range user {
				HealthGo2(v)
			}
		})
		if err != nil {
			log.Println(err)
			return
		}
	}
	_, err = c.AddFunc(s.HealthCron5, func() {
		user := config.GetUserAll()
		for _, v := range user {
			s, err := healthCheck(v)
			if err != nil {
				pushMsg(v, v.User, "晚间补签检查错误，请自行上号检查")
				continue
			}
			go func(v config.U, s []string) {
				if len(s) == 0 {
					pushMsg(v, v.User, "今天全签了")
					return
				}
				for _, v2 := range s {
					switch v2 {
					case "1":
						pushMsg(v, v.User, "早上没签，准备补签")
					case "2":
						pushMsg(v, v.User, "中午没签，准备补签")
					case "3":
						pushMsg(v, v.User, "晚上没签，准备补签")
					}
					hGo(v, v2)
				}
			}(v, s)
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
	c.Start()
	defer c.Stop()
	log.Println("cron添加成功")
	select {}
}

func healthCheck(info config.U) ([]string, error) {
	token := FvtiLogin(info)
	data, err := healthSendPost(getHealthIdApi, []byte("{\"page\":1,\"rows\":9}"), 1, token)
	if err != nil {
		return []string{}, err
	}
	healthIdx := new(healthId)
	err = json.Unmarshal(data, &healthIdx)
	if err != nil {
		return []string{}, err
	}
	var d []string
	for _, v := range healthIdx.Rows {
		if v.DataDate == time.Now().Format("2006-01-02") {
			if v.DataStatus != "2" {
				d = append(d, v.DataType)
			}
		}
	}
	return d, nil
}
func init() {
	go addCron()
	zero.OnCommand("申报账号测试").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.User == "" {
			return
		}
		ctx.Send(message.Text(FvtiLogin(u)))
	})
	zero.OnCommand("手动签到").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.User == "" {
			return
		}
		ctx.Send(message.Text("指令已发送"))
		hGo(u, "1")
		hGo(u, "2")
		hGo(u, "3")
	})
	zero.OnCommand("一键全签", zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) {
		user := config.GetUserAll()
		for _, v := range user {
			hGo(v, "1")
			hGo(v, "2")
			hGo(v, "3")
		}
		ctx.Send(message.Text("ok"))
	})
	zero.OnCommand("debug", zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) {
		user := config.GetUserAll()
		hGo(user[0], "1")
		hGo(user[0], "2")
		hGo(user[0], "3")
		ctx.Send(message.Text("ok"))
	})
}

const (
	checkApi          = "https://health.fvti.linyisong.top/api/mStuApi/getByStuEpidemicStuRoster.do"
	loginApi          = "https://health.fvti.linyisong.top/api/mStuApi/token.do"
	healthGo2Api      = "https://health.fvti.linyisong.top/api/mStuApi/queryForPhoneListEpidemicCheckingIn.do"
	healthGo2Api2     = "https://health.fvti.linyisong.top/api/mStuApi/updateForPhoneEpidemicCheckingIn.do"
	getHealthIdApi    = "https://health.fvti.linyisong.top/api/mStuApi/queryByStuEpidemicHealthReport.do"
	healthRosterIDApi = "https://health.fvti.linyisong.top/api/mStuApi/getByStuEpidemicStuRoster.do"
	healthGoApi       = "https://health.fvti.linyisong.top/api/mStuApi/updateHealthReportEpidemicHealthReport.do"
)

func pushMsg(info config.U, title, content string) {

	z := zero.GetBot(config.GetSetting().NickQQNumber)
	if z == nil {
		log.Println("GetBot nil")
		time.Sleep(time.Second * 10)
		z = zero.GetBot(config.GetSetting().NickQQNumber)
		if z == nil {
			return
		}
	}
	z.SendPrivateMessage(info.QqNumber, message.Text(title+"\n"+content+"\n"))
}
func check(accessToken string) bool {
	url := checkApi
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	req.Header.Set("accessToken", accessToken)
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	var AutoGenerated struct {
		IsSuccess bool `json:"isSuccess"`
		Data      struct {
			Model string `json:"model"`
		} `json:"data"`
	}
	err = json.Unmarshal(all, &AutoGenerated)
	if err != nil {
		log.Println(err)
		return false
	}
	if AutoGenerated.IsSuccess {
		return true
	}
	return false
}
func FvtiLogin(info config.U) string {
	if info.Cookie != "" {
		b := check(info.Cookie)
		if b {
			return info.Cookie
		}
	}
	t := url.Values{}
	t.Add("userCode", info.User)
	t.Add("userPwd", fmt.Sprintf("%x", md5.Sum([]byte(info.Pass))))
	resp, err := http.PostForm(loginApi, t)
	if err != nil {
		log.Println(err)
		return ""
	}
	tk, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	resp.Body.Close()
	var token struct {
		IsSuccess bool `json:"isSuccess"`
		Data      struct {
			StudentID   string `json:"studentId"`
			StudentName string `json:"studentName"`
			AccessToken string `json:"accessToken"`
		} `json:"data"`
	}
	err = json.Unmarshal(tk, &token)
	if err != nil {
		log.Println(err)
		return ""
	}
	info.Cookie = token.Data.AccessToken
	config.SaveUser(info)
	return token.Data.AccessToken
}

func HealthGo2(info config.U) {
	log.Println("开始晚打卡", info.User)
	var (
		title    = "晚打卡失败"
		content  = info.User
		ChickIDx = ChickID{}
	)
	defer func() {
		pushMsg(info, title, content)
	}()
	AccessToken := FvtiLogin(info)
	if AccessToken == "" {
		log.Println("登录错误")
		return
	}
	data, err := healthSendPost(healthGo2Api, []byte("theSelect=今日&page=1&rows=9"), 1, AccessToken)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(data, &ChickIDx)
	if err != nil {
		log.Println(err)
		return
	}
	if len(ChickIDx.Rows) < 1 {
		log.Println("ChickID获取失败")
		return
	}
	t := url.Values{}
	t.Add("chickId", ChickIDx.Rows[0].ChickID)
	t.Add("backAddress", "福建省福州市闽侯县上街镇113县道福州职业技术学院")
	t.Add("remark", "")
	data, err = healthSendPost(healthGo2Api2, []byte(t.Encode()), 1, AccessToken)
	if err != nil {
		log.Println("晚打卡失败")
		return
	}
	var postResp struct {
		IsSuccess bool `json:"isSuccess"`
		Data      struct {
		} `json:"data"`
	}
	err = json.Unmarshal(data, &postResp)
	if err != nil {
		log.Println(err)
		return
	}

	if postResp.IsSuccess {
		title = "晚打卡成功"
	}

}

//getHealthId dataType 1为早 2为中 3为晚
func getHealthId(AccessToken, dataType string) (HealthID string) {
	data, err := healthSendPost(getHealthIdApi, []byte("{\"page\":1,\"rows\":9}"), 1, AccessToken)
	if err != nil {
		log.Println(err)
		return
	}
	healthIdx := new(healthId)
	err = json.Unmarshal(data, &healthIdx)
	if err != nil {
		log.Println(err)
		return
	}
	atoi, err := strconv.Atoi(dataType)
	if err != nil {
		log.Println("getHealthId Atoi err", err)
		return
	}
	if atoi > len(healthIdx.Rows) {
		log.Println("getHealthId Rows err")
		return
	}

	for _, v := range healthIdx.Rows {
		if v.DataType == dataType {
			return v.HealthID
		}
	}
	return ""
}

func HealthRosterID(AccessToken string) (RosterID string) {
	data, err := healthSendPost(healthRosterIDApi, nil, 1, AccessToken)
	if err != nil {
		log.Println(err)
		return
	}
	rosterIdx := new(rosterId)
	err = json.Unmarshal(data, &rosterIdx)
	if err != nil {
		log.Println("rosterIdx", err)
		return
	}
	rosterIdModelx := new(rosterIdModel)
	err = json.Unmarshal([]byte(rosterIdx.Data.Model), rosterIdModelx)
	if err != nil {
		log.Println("rosterIdModelx", err)
		return
	}
	if len(rosterIdModelx.Rows) == 0 {
		log.Println("rosterIdModelx Rows err")
		return
	}
	return rosterIdModelx.Rows[0].RosterID
}

func HealthGo(info config.U, dataType string) error {
	time.Sleep(time.Second * 10)
	log.Println("开始", info.User)
	var endTime string
	AccessToken := FvtiLogin(info)
	if AccessToken == "" {
		return fmt.Errorf("登录失败")
	}
	switch dataType {
	case "1":
		endTime = "00:15:00"
	case "2":
		endTime = "11:00:00"
	case "3":
		endTime = "17:00:00"
	}
	x := url.Values{}
	x.Add("healthId", getHealthId(AccessToken, dataType))
	x.Add("rosterId", HealthRosterID(AccessToken))
	x.Add("isGfxReturn", "2")
	x.Add("isJwReturn", "2")
	x.Add("isContactPatient", "2")
	x.Add("isContactRiskArea", "2")
	x.Add("isHealthCodeOk", "2")
	x.Add("details", "")
	//x.Add("isSick", "0")
	x.Add("dataDate", time.Now().Format("2006-01-02"))
	x.Add("dataType", dataType) //1 早 2 中 3晚
	x.Add("endTime", endTime)
	x.Add("symptomAndHandle", "")
	x.Add("otherThing", "")
	//x.Add("dataDate", "")
	//x.Add("dataType", "") //1 早 2 中 3晚
	//x.Add("dataStatus", "2")
	x.Add("dataStatus", "")
	//x.Add("endTime", "")
	x.Add("liveState", config.GetSetting().LiveState)
	x.Add("nowAddress", "福建省福州市闽侯县上街镇")
	x.Add("nowAddressDetail", "联榕路8号")
	x.Add("nowTiwenState", "1")
	x.Add("nowHealthState", "1")
	x.Add("counsellorApprovalStatus", "")
	x.Add("temperature", "36.4")
	data, err := healthSendPost(healthGoApi, []byte(x.Encode()), 2, AccessToken)
	if err != nil {
		log.Println(err)
		return err
	}
	var postResp struct {
		IsSuccess bool `json:"isSuccess"`
		Data      struct {
		} `json:"data"`
	}
	log.Println(string(data))
	err = json.Unmarshal(data, &postResp)
	if err != nil {
		return err
	}

	if postResp.IsSuccess {
		return nil
	}
	return fmt.Errorf("失败")
}

func healthSendPost(url string, data []byte, typex int, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/Windows WindowsWechat")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	switch typex {
	case 1:
		req.Header.Set("accessToken", token)
	case 2:
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("accessToken", token)
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Referer", "https://servicewechat.com/wx56b1d7357f3df890/25/page-frame.html")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var reader io.ReadCloser
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		reader = resp.Body
	}
	database, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	return database, nil
}
