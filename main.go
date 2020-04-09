package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	fabcfg "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"io/ioutil"
	"log"
	"net/http"
)

var router *gin.Engine
var client *channel.Client

var value = `{
"_id": 66,
"trace_info": {
"store_id": "9025",
"submitter": "北京朝兴海缘商贸有限公司",
"entry_time": "2018-10-27 16:00:00",
"supplier_address": "北京市海淀区中关村大街250号",
"overseas_area": "英国",
"customs_quarantine_ticket": [{
"name": "测试证书",
"url": "https://sy.yonghui.cn/statics/group1/M00/00/E7/CgBaNVuPldCAXdUOAADAxW0Lg1M397.jpg"
}],
"submit_type": "ORDER",
"catch_seas_area": "FAO27",
"churchyard_agent_company": "北京泰格瑞迪科技有限公司",
"supplier_certifies": null,
"orginal_delivery_time": "2018-08-27 16:00:00",
"store_name": "永辉福建福州市--碧水芳洲店",
"supplier_name": "北京朝兴海缘商贸有限公司",
"supplier_id": 6,
"breed_type": "野生养殖",
"overseas_supplier_name": "SCRABSTER SEAFOODS LTD"
},
"insert_time": "2018-10-22 07:37:20",
"look_record_count": 0,
"last_look_time": "2018-10-22 07:37:20",
"batch_code": "00040142",
"modify_time": "2018-10-22 07:37:20",
"submit_type": "DELIVERY"
}`

func init() {
	router = gin.Default()
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true

	router.Use(cors.New(cfg))
	InitRouter()
}
func InitRouter() {
	router.POST("/invoke", invoketes)
	router.POST("/query", querytes)
}
func invoketes(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body) //读取请求body

	resp, err := client.Execute(channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "putValue",
		Args:        [][]byte{[]byte(body), []byte(value)},
		//TransientMap:transientDataMap,
	}, channel.WithTargets(YHPeer0))
	fmt.Println("1111")
	if err != nil {
		//fmt.Println("error to execute key number",  "error", err)
		c.JSON(http.StatusOK, gin.H{
			"result": "invoke failed",
			"key":    string(body),
			"info":   err.Error(),
		})

	} else {
		//fmt.Println("success execute key number", resp.TransactionID)
		c.JSON(http.StatusOK, gin.H{
			"result":  "invoke success",
			"key":     string(body),
			"payload": string(resp.Payload),
		})
	}

	return
}
func querytes(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body) //读取请求body

	resp, err := client.Execute(channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "query",
		Args:        [][]byte{[]byte(body), []byte("a")},
		//TransientMap:transientDataMap,
	}, channel.WithTargets(JYWYPeer0))
	if err != nil {
		//fmt.Println("error to query key number",  "error", err)
		c.JSON(http.StatusOK, gin.H{
			"result": "query failed",
			"key":    string(body),
		})

	} else {
		//fmt.Println("success query key number", string(resp.Payload))
		c.JSON(http.StatusOK, gin.H{
			"result":  "query success",
			"key":     string(body),
			"payload": string(resp.Payload),
		})
	}

	return
}

var (
	configPath = "./config/config_e2e.yaml"
	logger     *log.Logger
	ch         = make(chan int)
	JYWYPeer0  fab.Peer
	YHPeer0    fab.Peer
	//ZFJGPeer0  fab.Peer config_e2e.yaml

	//logger   *logging.Logger
)

func main() {
	err := InitConfig(configPath)
	if err != nil {
		fmt.Errorf("can't read configuration")
	}
	configProvider := fabcfg.FromFile(GetConfigFile())
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Println("error to create fabsdk", err)
		return
	}

	clientChannelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("admin"), fabsdk.WithOrg("Org1"))
	client, err = channel.New(clientChannelContext)
	if err != nil {
		fmt.Println("create channel failed", err)
		return
	}

	loadOrgPeers(sdk.Context(fabsdk.WithUser("admin"), fabsdk.WithOrg("Org1")))

	//file := "./" + time.Now().Format("2006-01-02-15:04:05") + ".log"
	//logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	//if nil != err {
	//	panic(err)
	//}

	//logger = log.New(logFile, "test_", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	router.Run(":8889")
}

func loadOrgPeers(ctxProvider contextAPI.ClientProvider) {

	ctx, err := ctxProvider()
	if err != nil {
		fmt.Println("make ctxPRovider failed", err)
		return
	}

	JYWYPeers, joined := ctx.EndpointConfig().PeersConfig("Org1")
	if joined != true {
		fmt.Println("make endpoint target failed", err)

	}

	//YHPeers, joined := ctx.EndpointConfig().PeersConfig("YH")
	//if joined != nil {
	//	fmt.Println("ner", err)
	//
	//}

	//ZFJGPeers, joined := ctx.EndpointConfig().PeersConfig("ZFJG")
	//if joined != nil {
	//	fmt.Println("make endpoint target failed", err)
	//	//return
	//}

	//fmt.Println("YHPeers",YHPeers,"JYWYPeers",JYWYPeers,"ZFJGPeers",ZFJGPeers)

	JYWYPeer0, err = ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: JYWYPeers[0]})
	if err != nil {
		fmt.Println("make endpoint target peer failed", err)
		return
	}

	//YHPeer0, err = ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: YHPeers[0]})
	//if err != nil {
	//	fmt.Println("make endpoint target peer failed", err)
	//	return
	//}

	//ZFJGPeer0, err = ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: ZFJGPeers[0]})
	//if err != nil {
	//	fmt.Println("make endpoint target peer failed", err)
	//	return
	//}
	//
	//fmt.Println("JYWYPeer0",JYWYPeer0)
	//fmt.Println("YHPeer0",YHPeer0)
	//fmt.Println("ZFJGPeer0",ZFJGPeer0)
}
