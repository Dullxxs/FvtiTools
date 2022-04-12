package config

import (
	_ "embed"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sync"
)

func init() {
	file, err := os.Open(configName)
	if err != nil {
		log.Println(err)
		createYaml()
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	file.Close()
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println(err)
		return
	}
	for k, v := range config.U {
		userIndex[v.QqNumber] = k
	}
}

//go:embed config.yml
var demo []byte

var (
	lck       sync.Mutex
	userIndex = make(map[int64]int)
)

const (
	configName = "config.yml"
)

type U struct {
	User          string  `yaml:"user"`
	Pass          string  `yaml:"pass"`
	Cookie        string  `yaml:"cookie"`
	User2         string  `yaml:"user2"`
	Pass2         string  `yaml:"Pass2"`
	Cookie2       string  `yaml:"cookie2"`
	Room          string  `yaml:"room"`
	RoomYda       float32 `yaml:"roomYda"`
	QqNumber      int64   `yaml:"qq-number"`
	QqGroupNumber int64   `yaml:"qq-GroupNumber"`
	Remark        string  `yaml:"remark"`
}

type Setting struct {
	QqApi        string `yaml:"qq-api"`
	AdminQq      int64  `yaml:"admin-qq"`
	NickName     string `yaml:"nick-name"`
	NickQQNumber int64  `yaml:"nick-qq-number"`
	HealthCron   string `yaml:"health-cron"`
	HealthCron2  string `yaml:"health-cron2"`
	HealthCron3  string `yaml:"health-cron3"`
	HealthCron4  string `yaml:"health-cron4"`
	HealthCron5  string `yaml:"health-cron5"`
	LiveState    string `yaml:"live-state"`
}

var config struct {
	S Setting `yaml:"s"`
	U []U     `yaml:"u"`
}

func SaveSetting(setting Setting) {
	lck.Lock()
	defer lck.Unlock()
	config.S = setting
	saveConfig()
}

func SaveUser(u U) {
	lck.Lock()
	defer lck.Unlock()
	k := userIndex[u.QqNumber]
	config.U[k] = u
	saveConfig()
}

func GetSetting() Setting {
	lck.Lock()
	defer lck.Unlock()
	return config.S
}

func GetUser(qq int64) U {
	lck.Lock()
	defer lck.Unlock()
	return config.U[userIndex[qq]]
}

func GetUserAll() []U {
	lck.Lock()
	defer lck.Unlock()
	return config.U
}

func saveConfig() {
	file, err := os.Create(configName)
	if err != nil {
		log.Println(err)
		return
	}
	out, err := yaml.Marshal(config)
	if err != nil {
		log.Println(err)
		return
	}
	file.Write(out)
	file.Close()
}

func createYaml() {
	log.Println("开始创建yaml")
	file, err := os.Create(configName)
	if err != nil {
		log.Println(err)
		return
	}
	file.Write(demo)
	file.Close()
}
