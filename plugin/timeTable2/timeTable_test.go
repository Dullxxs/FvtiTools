package timeTable2

import (
	"FvtiTools/config"
	"testing"
)

//func TestLoginSky(t *testing.T) {
//	loginSky("", config.U{}, captcha{})
//}

func TestG(t *testing.T) {
	//c, err := getCaptcha()
	//if err != nil {
	//	t.Log(err)
	//	return
	//}
	//t.Log(c)
	//var a string
	//fmt.Scanf("%s\n", &a)
	user := config.U{
		User:          "",
		Pass:          "",
		Cookie:        "",
		User2:         "",
		Pass2:         "",
		Cookie2:       "",
		Room:          "",
		RoomYda:       0,
		QqNumber:      0,
		QqGroupNumber: 0,
		Remark:        "",
	}
	t.Log(loginSky("zbrd", user))
}
