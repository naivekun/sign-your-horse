# Sign your horse 签个马

> 签你🐎的到

Go编写的签到平台框架

```
            Provider                                     Reporter

+-------------------------------+             +------------+   +------------+
|                               |             |            |   |            |
|       chaoxing_course_1       +------------------------------------------------>  pushMessage
|                               |             |            |   |            |
+-------------------------------+             | WeChat APP |   |  Telegram  |
                                              |            |   |    Bot     |
+-------------------------------+             |            |   |   (TODO)   |
|                               |             |            |   |            |
|       chaoxing_course_2       +------------------------------------------------>  pushMessage
|                               |             |            |   |            |
+-------------------------------+             |            |   |            |
                                              |            |   |            |
+-------------------------------+             |            |   |            |
|                               |             |            |   |            |
|          teachermate          +------------------------------------------------>  pushMessage
|                               |             |            |   |            |
+-------------------------------+             +------------+   +------------+

```

## 使用方法

```
$ ./sign-your-horse -h
Usage of sign-your-horse.exe:
  -config string
        specify config file (default "config.json")
```

直接运行，会在目录下创建config.json，包含默认配置

```
{
	"provider": [
		{
			"name": "chaoxing_default",
			"module": "chaoxing",
			"config": {
				"cookie": "",
				"useragent": "",
				"uid": "",
				"courseid": "",
				"classid": "",
				"interval": 5
			}
		},
		{
			"name": "teachermate_default",
			"module": "teachermate",
			"config": {
				"server": "0.0.0.0:3000"
			}
		}
	],
	"reporter": [
		{
			"name": "console",
			"config": {}
		},
		{
			"name": "wechat",
			"config": {
				"corpID": "",
				"corpSecret": "",
				"toparty": 0,
				"agentid": 0
			}
		}
	]
}
```

## Provider

Provider适配各个签到平台，提供Init和Run方法和默认配置json

在配置文件中provider是一个列表，可指定某个平台的某个签到任务，如果有多节课需要签到可以配置多个provider即可

### chaoxing

超星签到模块，各个参数说明如下

```
alias: "别名，用于推送消息时区分各个任务",
cookie: "超星登录cookie",
useragent: "User-Agent",
uid: "超星的uid，从cookie里面扣",
courseid: "课程ID",
classid: "班级ID",
interval: 5
```

### teachermage

微助教签到模块

工作原理如下：

[CloudScan APP](https://github.com/naivekun/cloudscan-android)/CloudScan Web -> Sign-your-horse Server -> QRCode/Redirect/Text

到教室的同学使用[CloudScan APP](https://github.com/naivekun/cloudscan-android)或CloudScan Web(自带)扫描二维码 发送到后端

其他同学使用微信扫描后端提供的二维码或在微信里点击重定向链接即可跳转到签到页面

* `/static/qr.html` 3秒刷新一次验证码，可供直接使用微信扫描
* `/static/scan.html` CloudScan Web，到教室的同学使用，会使用浏览器调用摄像头，扫描二维码上传签到信息。Thanks to @EarthC
* `/url/add` 上传接口
* `/url/redirect` 重定向接口，点一下直接302跳到签到页面，适合微助教这种依赖微信登录的使用
* `/url/raw` 获取url明文

## Reporter

Reporter用于接收Provider推送的数据，一般来说是签到成功/失败的通知。Provider的通知会依次调用所有Reporter推送信息

### wechat

用于向微信企业号推送消息

实时反馈签到状态，并配合微助教完成一键url跳转签到

### console

直接把消息print到console上

## 开发

### 构建

由于使用了静态资源，你需要使用`packr`完成构建

```
$ packr build
```

### provider

在provider目录增加模块，实现Init和Run方法，init函数中使用`provider.RegisterProvider`注册模块即可。

### reporter

在reporter目录增加模块，实现Init和Report方法，init函数中使用`reporter.RegisterReporter`注册模块即可

自带的console模块会简单把结果打印到stdout，用它直接改是个不错的选择

## License

MIT