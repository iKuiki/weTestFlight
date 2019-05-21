module weTestFlight

go 1.12

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/ikuiki/wwdk v2.3.4+incompatible
	github.com/mdp/qrterminal v1.0.1
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190509222800-a4d6f7feada5 // indirect
	wegate v0.0.0-00010101000000-000000000000
)

replace github.com/liangdas/mqant => github.com/ikuiki/mqant v1.8.1-0.20190427142930-7dabfa32d064

replace wegate => github.com/ikuiki/wegate v0.0.0-20190520102521-641d4a557a35

// 解决国内无法下载的几个包
replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190513172903-22d7a77e9e5f
	golang.org/x/net => github.com/golang/net v0.0.0-20190514140710-3ec191127204
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190514135907-3a4b5fb9f71f
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190515035509-2196cb7019cc
	google.golang.org/appengine => github.com/golang/appengine v1.6.0
)
