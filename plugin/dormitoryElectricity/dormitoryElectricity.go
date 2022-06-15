package dormitoryElectricity

import (
	"FvtiTools/config"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	buildUrl = "http://h5cloud.17wanxiao.com:8080/CloudPayment/user/getRoom.do?payProId=398&schoolcode=3071&optype=2&areaid=99&buildid=0&unitid=0&levelid=0&businesstype=2"
	floorUrl = "http://h5cloud.17wanxiao.com:8080/CloudPayment/user/getRoom.do?payProId=398&schoolcode=3071&optype=3&areaid=99&buildid=%v&unitid=0&levelid=0&businesstype=2"
	roomUrl  = "http://h5cloud.17wanxiao.com:8080/CloudPayment/user/getRoom.do?payProId=398&schoolcode=3071&optype=4&areaid=99&buildid=%v&unitid=0&levelid=%v&businesstype=2"
	stateUrl = "http://h5cloud.17wanxiao.com:8080/CloudPayment/user/getRoomState.do?payProId=398&schoolcode=3071&businesstype=2&roomverify=%v"
)

func init() {
	c := cron.New()
	id, err := c.AddFunc("0 */24 * * *", func() {
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
			if v.Room == "0" || v.Room == "" {
				continue
			}
			quantity, err := getRoomQuantity(v.Room)
			if err != nil {
				log.Println(err)
				bctx.SendPrivateMessage(v.QqNumber, message.Text("出错啦", err))
				continue
			}
			n := v.RoomYda - quantity
			if v.QqGroupNumber == 0 {
				if n < 0 {
					bctx.SendPrivateMessage(v.QqNumber, message.Text("昨天:", v.RoomYda, "\n今天:", quantity, "\n多了:", -n))
				} else {
					bctx.SendPrivateMessage(v.QqNumber, message.Text("昨天:",
						v.RoomYda, "\n今天:", quantity, "\n用了:", n))
				}
			} else {
				if n < 0 {
					bctx.SendGroupMessage(v.QqGroupNumber, message.Text("昨天:", v.RoomYda, "\n今天:", quantity, "\n多了:", -n))
					bctx.SendPrivateMessage(v.QqNumber, message.Text("昨天:", v.RoomYda, "\n今天:", quantity, "\n多了:", -n))
				} else {
					bctx.SendGroupMessage(v.QqGroupNumber, message.Text("昨天:", v.RoomYda, "\n今天:", quantity, "\n用了:", n))
					bctx.SendPrivateMessage(v.QqNumber, message.Text("昨天:", v.RoomYda, "\n今天:", quantity, "\n用了:", n))
				}
			}
			v.RoomYda = quantity
			config.SaveUser(v)
		}
	})

	if err != nil {
		log.Println(err)
		return
	}
	c.Start()
	zero.OnCommand("电费").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.QqNumber != ctx.Event.UserID {
			return
		}
		if u.Room == "0" || u.Room == "" {
			return
		}
		quantity, err2 := getRoomQuantity(u.Room)
		if err2 != nil {
			log.Println(err2)
			ctx.Send(message.Text("出错啦,", err2))
			return
		}
		ctx.Send(message.Text(quantity))
	})
	zero.OnCommand("绑定宿舍").Handle(func(ctx *zero.Ctx) {
		u := config.GetUser(ctx.Event.UserID)
		if u.QqNumber != ctx.Event.UserID {
			return
		}
		arguments := shell.Parse(ctx.State["args"].(string))
		if len(arguments) != 1 {
			ctx.Send(message.Text("格式错误,正确的格式是\n/绑定宿舍 1号楼#1层#101\n关闭的话/绑定宿舍 0"))
			return
		}
		u.Room = arguments[0]
		config.SaveUser(u)
		quantity, err2 := getRoomQuantity(u.Room)
		if err2 != nil {
			log.Println(err2)
			ctx.Send(message.Text("出错啦", err2))
			return
		}
		ctx.Send(message.Text("更改成功\n", u.Room, "\n当前电量:", quantity))
	})
	zero.OnCommand("roomClose", zero.SuperUserPermission).Handle(func(ctx *zero.Ctx) { //隐藏命令
		c.Remove(id)
	})
}

func getRoomQuantity(Room string) (float32, error) {
	split := strings.Split(Room, "#")
	if len(split) != 3 {
		return 0, fmt.Errorf("格式错误")
	}
	resp, err := http.Get(buildUrl)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var c buildCode
	err = json.Unmarshal(readAll, &c)
	log.Println(string(readAll))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var b string
	for _, v := range c.Roomlist {
		if v.Name == split[0] {
			b = v.ID
		}
	}
	resp, err = http.Get(fmt.Sprintf(floorUrl, b))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	readAll, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = json.Unmarshal(readAll, &c)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var b2 string
	for _, v := range c.Roomlist {
		if v.Name == split[1] {
			b2 = v.ID
		}
	}
	resp, err = http.Get(fmt.Sprintf(roomUrl, b, b2))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	readAll, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = json.Unmarshal(readAll, &c)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var b3 string
	for _, v := range c.Roomlist {
		if v.Name == split[2] {
			b3 = v.ID
		}
	}
	resp, err = http.Get(fmt.Sprintf(stateUrl, b3))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	readAll, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	var cx roomStats
	err = json.Unmarshal(readAll, &cx)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	v1, err := strconv.ParseFloat(cx.Quantity, 32)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return float32(v1), nil
}
