package command

import (
	"FvtiTools/config"
	"FvtiTools/plugin/health"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"strings"
)

var (
	helpText = strings.Join([]string{
		"指令前面记得加/",
		"区分大小写",
		"帮助 - 帮助",
		"申报账号测试 - 测试早中晚申报的账号密码",
		"手动签到 - 早中晚签到，一遍全签",
		//"查课表 - 查询课表",
		//"关闭课表推送 - 关闭课表推送",
		//"飞翔登录 - 登录飞翔系统",
		//"检测飞翔 - 测试飞翔系统cookie是否失效",
		"电费 - 查询宿舍用电量",
		"绑定宿舍 - 绑定宿舍",
		"绑定群 - 绑定QQ群",
		"用户信息 - 查看用户信息",
		"",
	}, "\n")
)

func init() {
	//zero.GetBot(config.GetSetting().NickQQNumber).SetGroupAddRequest()

	zero.OnCommand("帮助").Handle(func(ctx *zero.Ctx) {
		ctx.Send(message.Text(helpText))
	})
	zero.OnRequest().SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)

		switch ctx.Event.RequestType {
		case "friend":
			if u.QqNumber != ctx.Event.UserID {
				return
			}
			ctx.SetFriendAddRequest(ctx.Event.Flag, true, "")
		case "group":
			if u.QqGroupNumber != ctx.Event.GroupID {
				return
			}
			if ctx.Event.SubType == "invite" {
				ctx.SetGroupAddRequest(ctx.Event.Flag, ctx.Event.SubType, true, "")
			}
		}
	})

	zero.OnCommand("绑定群").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.QqNumber != ctx.Event.UserID {
			return
		}
		arguments := shell.Parse(ctx.State["args"].(string))
		if len(arguments) != 1 {
			ctx.Send(message.Text("格式错误,正确格式\n/绑定群 114514\n关闭群推送请用/绑定群 0"))
			return
		}
		x, err := strconv.ParseInt(arguments[0], 10, 64)
		if err != nil {
			log.Println(err)
			ctx.Send(message.Text("出错啦", err))
		}
		u.QqGroupNumber = x
		config.SaveUser(u)
		ctx.Send(message.Text("绑定好了，拉机器人入群吧"))
	})

	zero.OnCommand("用户信息").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.QqNumber != ctx.Event.UserID {
			return
		}
		i := strings.Join([]string{
			"微信健康账号:" + u.User,
			"微信健康密码:" + u.Pass,
			//"飞翔bot账号:" + u.User2,
			//"飞翔bot密码:" + u.Pass2,
			"绑定宿舍号:" + u.Room,
			"绑定QQ群:" + strconv.FormatInt(u.QqGroupNumber, 10),
		}, "\n")
		ctx.Send(message.Text(i))
	})

	zero.OnCommand("timeTableClose").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.QqNumber != ctx.Event.UserID {
			return
		}
		u.Cookie2 = ""
		config.SaveUser(u)
	})

	zero.OnCommand("一键全签", zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) {
		user := config.GetUserAll()
		for _, v := range user {
			health.HealthGo(v, "1", false)
			health.HealthGo(v, "2", false)
			health.HealthGo(v, "3", false)
		}
		ctx.Send(message.Text("ok"))
	})
	zero.OnCommand("debug", zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) {
		user := config.GetUserAll()
		health.HealthGo(user[0], "1", false)
		health.HealthGo(user[0], "2", false)
		health.HealthGo(user[0], "3", false)
		ctx.Send(message.Text("ok"))
	})
}
