# 福州职业技术学院多功能机器人
2022.6.14 测试可用
## 功能
* 早中晚健康上报，以及晚打卡。
* 默认每天晚上18:30检查当天签到情况
* ~~宿舍电量查询以及群推送。~~ 接口加了cookie校验，而且有效期似乎不是很长，懒得修
* ~~课表查询以及群推送。~~ 没接口了，还得自己登陆vpn，麻烦懒得修
* 支持qq推送
* 支持多用户
## 使用
* 随便下载一个支持正向websocket OneBot v11的框架即可，推荐[go-cqhttp](https://github.com/Mrs4s/go-cqhttp)
* 配置好框架，打开程序，会生成一个yml.
* 填写好:qq-api,机器qq，管理员qq,福职健康小程序账号密码.

## 参考项目
[FvtiJKDK](https://github.com/AiMuC/FvtiJKDK)

