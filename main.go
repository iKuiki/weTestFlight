package main

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ikuiki/wwdk"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/mdp/qrterminal"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
	"wegate/common"
	"wegate/common/test"
	"wegate/wechat"
)

func main() {
	w := commontest.Work{}
	opts := w.GetDefaultOptions("tcp://sz.kuiki.cn:3563")
	opts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		fmt.Println("ConnectionLost", err.Error())
	})
	opts.SetOnConnectHandler(func(client MQTT.Client) {
		fmt.Println("OnConnectHandler")
	})
	err := w.Connect(opts)
	if err != nil {
		panic(err)
	}
	loginStatusChannel := make(chan bool)
	w.On("LoginStatus", func(client MQTT.Client, msg MQTT.Message) {
		var loginItem wwdk.LoginChannelItem
		json.Unmarshal(msg.Payload(), &loginItem)
		fmt.Println("LoginStatus: ", loginItem.Code)
		switch loginItem.Code {
		case wwdk.LoginStatusWaitForScan:
			qrterminal.Generate(loginItem.Msg, qrterminal.L, os.Stdout)
		case wwdk.LoginStatusGotBatchContact:
			loginStatusChannel <- true
		}
	})
	w.On("AddPlugin", func(client MQTT.Client, msg MQTT.Message) {
		var pluginDesc wechat.PluginDesc
		e := json.Unmarshal(msg.Payload(), &pluginDesc)
		if e != nil {
			fmt.Println("addplugin: json.Unmarshal(msg.Payload(),&pluginDesc) error: ", e)
			return
		}
		fmt.Println("add plugin: ", pluginDesc)
	})
	w.On("RemovePlugin", func(client MQTT.Client, msg MQTT.Message) {
		var pluginDesc wechat.PluginDesc
		e := json.Unmarshal(msg.Payload(), &pluginDesc)
		if e != nil {
			fmt.Println("removeplugin: json.Unmarshal(msg.Payload(),&pluginDesc) error: ", e)
			return
		}
		fmt.Println("remove plugin: ", pluginDesc)
	})
	pass := "hello" + time.Now().Format(time.RFC822)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	resp, _ := w.Request("Login/HD_Login", []byte(`{"username":"abc","password":"`+string(hashedPass)+`"}`))
	if resp.Ret != common.RetCodeOK {
		panic(fmt.Sprintf("登录失败: %s", resp.Msg))
	}
	resp, _ = w.Request("Wechat/HD_Wechat_RegisterMQTTPlugin", []byte(fmt.Sprintf(
		`{"name":"%s","description":"%s","loginListenerTopic":"%s","contactListenerTopic":"%s","msgListenerTopic":"%s","addPluginListenerTopic":"%s","removePluginListenerTopic":"%s"}`,
		"testPlugin",   // name
		"测试模块",         // description
		"LoginStatus",  // loginListenerTopic
		"",             // contactListenerTopic
		"",             // msgListenerTopic
		"AddPlugin",    // addPluginListenerTopic
		"RemovePlugin", // removePluginListenerTopic
	)))
	if resp.Ret != common.RetCodeOK {
		panic(fmt.Sprintf("注册plugin失败: %s", resp.Msg))
	}
	token := resp.Msg
	fmt.Printf("获取到token：%s\n", token)
	// 获取已注册的plugin
	resp, _ = w.Request("Wechat/HD_Plugin_GetPluginList", []byte(`{"token":"`+token+`"}`))
	if resp.Ret != common.RetCodeOK {
		panic(fmt.Sprintf("GetPluginList失败: %s", resp.Msg))
	}
	var pluginDescList []wechat.PluginDesc
	json.Unmarshal([]byte(resp.Msg), &pluginDescList)
	fmt.Println("pluginList: ")
	for _, p := range pluginDescList {
		fmt.Println(p)
	}
	<-loginStatusChannel
	fmt.Println("remote wechat login success")
	// 测试调用wechat方法
	resp, _ = w.Request("Wechat/HD_Wechat_CallWechat", []byte(`{"fnName":"GetRunInfo","token":"`+token+`"}`))
	if resp.Ret != common.RetCodeOK {
		panic(fmt.Sprintf("GetRunInfo失败: %s", resp.Msg))
	}
	fmt.Println("RunInfo: ", resp.Msg)
	resp, _ = w.Request("Wechat/HD_Wechat_CallWechat", []byte(`{"fnName":"GetContactList","token":"`+token+`"}`))
	if resp.Ret != common.RetCodeOK {
		panic(fmt.Sprintf("GetContactList失败: %s", resp.Msg))
	}
	var contacts []datastruct.Contact
	json.Unmarshal([]byte(resp.Msg), &contacts)
	for _, contact := range contacts {
		if contact.IsChatroom() {
			// fmt.Printf("%+v\n", contact)
			fmt.Println("Chatroom: ", contact.NickName)
			fmt.Println(contact.HeadImgURL)
			fmt.Println("MemberCount: ", len(contact.MemberList))
			var headCount int
			for _, m := range contact.MemberList {
				if m.KeyWord != "" {
					headCount++
					if headCount < 5 {
						fmt.Println("member ", m.NickName, " head: ", m.KeyWord)
					}
				}
			}
			fmt.Println("total head ", headCount)
			fmt.Println("")
		}
	}
}
