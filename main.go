package main

import (
	"FvtiTools/config"
	_ "FvtiTools/plugin/command"
	_ "FvtiTools/plugin/health"
	//_ "FvtiTools/plugin/timeTable"
	//_ "embed"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"strconv"
)

func main() {
	setting := config.GetSetting()
	LiveState := setting.LiveState
	switch LiveState {
	case "1":
		log.Println("在家")
	case "2":
		log.Println("在校")
	case "3":
		log.Println("实习")
	default:
		log.Println("liveState输入错误，1 在家 2在校 3实习")
		return
	}
	zero.Run(zero.Config{
		NickName:      []string{setting.NickName},
		CommandPrefix: "/",
		SuperUsers:    []string{strconv.FormatInt(setting.AdminQq, 10)},
		Driver: []zero.Driver{
			driver.NewWebSocketClient(setting.QqApi, ""),
		},
	})
	select {}
}

//func createYaml() {
//	log.Println("开始创建yaml")
//	file, err := os.Create(configName)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	info := config{
//		S: setting{
//			QqApi:       "api",
//			HealthCron:  "0 6 * * *",
//			HealthCron2: "0 12 * * *",
//			HealthCron3: "0 18 * * *",
//			HealthCron4: "0 20 * * *",
//			LiveState:   "1"},
//		U: nil,
//	}
//	info.U = append(info.U, U{
//		User:     "学号",
//		Pass:     "密码",
//		User2:    "身份证号码",
//		Pass2:    "身份证号码",
//		QqNumber: "qq号码",
//		PushPlus: "PushPlus的token",
//		Remark:   "备注",
//	})
//	marshal, err := yaml.Marshal(info)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	_, err = file.Write(marshal)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	file.Close()
//	log.Println("创建完毕,请重新启动")
//}
