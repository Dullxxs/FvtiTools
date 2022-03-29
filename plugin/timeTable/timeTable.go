package timeTable

import (
	"FvtiTools/config"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	checkApi   = "http://121.5.139.76:9999/cx?s=ck&p=6"
	loginApi   = "http://121.5.139.76:9999/bd?"
	captchaApi = "http://121.5.139.76:9999/getyzm"
	classApi   = "http://121.5.139.76:9999/cx?s=ck&p="
)

func init() {
	c := cron.New()
	id, err := c.AddFunc("10 */24 * * *", func() {
		us := config.GetUserAll()
		bctx := zero.GetBot(config.GetSetting().NickQQNumber)
		if bctx == nil {
			time.Sleep(time.Second * 10)
			bctx = zero.GetBot(config.GetSetting().NickQQNumber)
			if bctx == nil {
				return
			}
		}
		for _, v := range us {
			if v.Cookie2 == "0" || v.Cookie2 == "" {
				continue
			}
			x := getTimeTable(v.Cookie2)
			if v.QqGroupNumber == 0 {
				bctx.SendPrivateMessage(v.QqNumber, message.Image("base64://"+x))
			} else {
				bctx.SendGroupMessage(v.QqGroupNumber, message.Image("base64://"+x))
				bctx.SendPrivateMessage(v.QqNumber, message.Image("base64://"+x))
			}
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
	c.Start()
	zero.OnCommand("飞翔登录").Handle(func(ctx *zero.Ctx) {
		user := config.GetUser(ctx.Event.UserID)
		if user.QqNumber != ctx.Event.UserID {
			return
		}
		arguments := shell.Parse(ctx.State["args"].(string))
		if len(arguments) != 2 {
			ctx.Send(message.Text("格式错误,正确的格式是\n/飞翔登录 账户 密码"))
			return
		}
		if user.Cookie2 != "" {
			ctx.Send(message.Text("已存在cookie，是否进行替换。\n是回复1，否回复2"))
			echox, cancelx := ctx.FutureEvent("message", ctx.CheckSession()).Repeat()
			defer cancelx()
			select {
			case <-time.After(time.Minute):
				ctx.Send(message.Text("超时，请重新登录"))
				return
			case ex := <-echox:
				log.Println("收到")
				switch ex.Message.String() {
				case "1":
					ctx.Send(message.Text("ok"))
				case "2":
					ctx.Send(message.Text("登录取消"))
					return
				default:
					ctx.Send(message.Text("输入错误"))
					return
				}
			}
		}
		echo, cancel := ctx.FutureEvent("message", ctx.CheckSession()).Repeat()
		defer cancel()
		c, err := getCaptcha()
		if err != nil {
			ctx.Send(message.Text("验证码获取错误", err))
			return
		}
		ctx.Send(message.Text("2分钟内输入验证码"))
		ctx.Send(message.Image("base64://" + c.S))
		select {
		case <-time.After(time.Minute * 2):
			ctx.Send(message.Text("超时了，请重新登录"))
			return
		case e := <-echo:
			user.User2, user.Pass2 = arguments[0], arguments[1]
			token, err := loginSky(e.Message.String(), user, c)
			if err != nil {
				ctx.Send(message.Text(err.Error()))
				return
			}
			ctx.Send(message.Text("登录成功,token为", token))
			user.Cookie2 = token
			config.SaveUser(user)
		}

	})

	zero.OnCommand("检测飞翔").Handle(func(ctx *zero.Ctx) {
		user := config.GetUser(ctx.Event.UserID)
		if user.QqNumber != ctx.Event.UserID {
			return
		}
		b := check(user)
		if b {
			ctx.Send(message.Text("cookie有效"))
			return
		}
		ctx.Send(message.Text("cookie已失效"))
	})

	zero.OnCommand("查课表").Handle(func(ctx *zero.Ctx) {
		user := config.GetUser(ctx.Event.UserID)
		if user.QqNumber != ctx.Event.UserID {
			return
		}
		if user.Cookie2 == "" {
			ctx.Send(message.Text("亲，您还没登录"))
			return
		}
		ctx.Send(message.Image("base64://" + getTimeTable(user.Cookie2)))
	})

	zero.OnCommand("关闭课表推送", zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) { //隐藏命令
		if ctx.Event.UserID != config.GetSetting().AdminQq {
			return
		}
		c.Remove(id)
	})
}

type captcha struct {
	PID string `json:"PID"`
	PNS string `json:"PNS"`
	S   string `json:"S"`
}

func getTimeTable(cookie string) string {
	getP := func() int {
		t := time.Now()
		yearDay := t.YearDay()
		yearFirstDay := t.AddDate(0, 0, -yearDay+1)
		firstDayInWeek := int(yearFirstDay.Weekday())
		//今年第一周有几天
		firstWeekDays := 1
		if firstDayInWeek != 0 {
			firstWeekDays = 7 - firstDayInWeek + 1
		}
		var week int
		if yearDay <= firstWeekDays {
			week = 1
		} else {
			week = (yearDay-firstWeekDays)/7 + 2
		}
		return week - 36
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	var res []byte
	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			network.SetExtraHTTPHeaders(network.Headers{
				"accessToken": cookie,
			}).Do(ctx)
			return nil
		}),
		chromedp.Navigate(classApi+strconv.Itoa(getP())),
		chromedp.WaitVisible(".mdui-table"),
		chromedp.Screenshot(".mdui-table", &res),
	); err != nil {
		log.Println(err)
	}
	return base64.StdEncoding.EncodeToString(res)
}
func getCaptcha() (captcha, error) {
	resp, err := http.Get(captchaApi)
	if err != nil {
		log.Println("getCaptcha", err)
		return captcha{}, err
	}
	defer resp.Body.Close()
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("getCaptcha ReadAll", err)
		return captcha{}, err
	}
	var c captcha
	err = json.Unmarshal(readAll, &c)
	if err != nil {
		log.Println("getCaptcha Unmarshal", err)
		return captcha{}, err
	}
	return c, nil
}

func loginSky(writeCaptcha string, u config.U, c captcha) (string, error) {
	v := url.Values{}
	v.Add("sfz", u.User2)
	v.Add("mm", u.Pass2)
	v.Add("qq", strconv.FormatInt(u.QqNumber, 10))
	v.Add("PID", c.PID)
	v.Add("PNS", c.PNS)
	v.Add("yzm", writeCaptcha)
	log.Println(v.Encode())
	resp, err := http.Get(loginApi + v.Encode())
	if err != nil {
		log.Println(err)
		return "", err
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var result struct {
		IsSuccess   bool   `json:"isSuccess"`
		Code        int    `json:"code"`
		Message     string `json:"message"`
		AccessToken string `json:"accessToken"`
	}
	err = json.Unmarshal(readAll, &result)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if result.IsSuccess != true {
		return "", fmt.Errorf(result.Message)
	}
	return result.AccessToken, nil
}

func check(u config.U) bool {
	client := &http.Client{}
	req, err := http.NewRequest("Get", checkApi, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Set("accessToken", u.Cookie2)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	if string(all) == "accessToken错" {
		return false
	}
	return true
}
