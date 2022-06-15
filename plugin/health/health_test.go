package health

import (
	"FvtiTools/config"
	"crypto/md5"
	"fmt"
	"log"
	"testing"
)

func TestHealthGo(t *testing.T) {
	utest := config.U{
		User: "",
		Pass: "",
	}
	utest.Pass = fmt.Sprintf("%x", md5.Sum([]byte(utest.Pass)))
	log.Println(FvtiLogin(utest) == "")
	//liveState = "1"
	//healthGo(utest, "3")

}

func TestPushPlus(t *testing.T) {
	//pushPlus("","212","6256")
}
