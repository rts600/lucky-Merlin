//Copyright 2022 gdy, 272288813@qq.com

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gdy666/lucky/config"
	"github.com/gdy666/lucky/ddns"
	"github.com/gdy666/lucky/reverseproxy"
	"github.com/gdy666/lucky/socketproxy"
)

var (
	listenPort       = flag.Int("port", 16601, "http Admin Web listen port ")
	configureFileURL = flag.String("c", "", "configure file url")

	// 华硕梅林脚本兼容参数
	cdFlag             = flag.String("cd", "", "config dir")
	infoFlag           = flag.Bool("info", false, "output info")
	baseConfInfoFlag   = flag.Bool("baseConfInfo", false, "output base config info")
	rResetUser         = flag.Bool("rResetUser", false, "reset user")
	rSetHttpAdminPort  = flag.Int("rSetHttpAdminPort", 0, "set http admin port")
	rSetHttpsAdminPort = flag.Int("rSetHttpsAdminPort", 0, "set https admin port")
)

var (
	runMode = "prod"
	version = "dev"
	commit  = "none"
	date    = "2022-07-27T17:54:45Z"
)

var runTime time.Time

func init() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}

func main() {
	flag.Parse()

	if *infoFlag {
		fmt.Printf(`{"Version":"%s"}`+"\n", version)
		os.Exit(0)
	}

	if *cdFlag != "" {
		if (*cdFlag)[len(*cdFlag)-1] != '/' && (*cdFlag)[len(*cdFlag)-1] != '\\' {
			*cdFlag += "/"
		}
		*configureFileURL = *cdFlag + "lucky.conf"
	}

	if *baseConfInfoFlag {
		err := config.Read(*configureFileURL)
		if err != nil {
			config.LoadDefault(*listenPort)
		}
		gcf := config.GetConfig()
		fmt.Printf(`{"AdminWebListenPort":%d, "SafeURL":""}`+"\n", gcf.BaseConfigure.AdminWebListenPort)
		os.Exit(0)
	}

	if *rSetHttpAdminPort > 0 || *rSetHttpsAdminPort > 0 || *rResetUser {
		config.Read(*configureFileURL)
		gcf := config.GetConfig()
		
		if *rSetHttpAdminPort > 0 {
			gcf.BaseConfigure.AdminWebListenPort = *rSetHttpAdminPort
		}
		if *rSetHttpsAdminPort > 0 {
			gcf.BaseConfigure.AdminWebListenHttpsPort = *rSetHttpsAdminPort
		}
		if *rResetUser {
			gcf.BaseConfigure.AdminAccount = "666"
			gcf.BaseConfigure.AdminPassword = "666"
		}
		
		config.Save()
		// 重置完成后退出，由梅林脚本负责重新启动进程
		os.Exit(0)
	}



	config.InitAppInfo(version, date)

	err := config.Read(*configureFileURL)
	if err != nil {
		log.Printf("%s", err.Error())
		log.Printf("载入默认配置以及命令行设定的参数")
		config.LoadDefault(*listenPort)
		if len(*configureFileURL) > 0 {
			err = config.Save()
			if err != nil {
				log.Printf("保存配置到%s出错:%s", *configureFileURL, err.Error())
			}
		}
	}

	gcf := config.GetConfig()

	// 强制用命令行传入的 -port 覆盖配置文件中的端口
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "port" {
			gcf.BaseConfigure.AdminWebListenPort = *listenPort
		}
	})
	config.BlackListInit()
	config.WhiteListInit()
	config.SSLCertficateListInit()

	//fmt.Printf("*gcf:%v\n", *gcf)

	socketproxy.SetSafeCheck(config.SafeCheck)
	//socketproxy.SetGlobalMaxConnections(gcf.BaseConfigure.GlobalMaxConnections)
	//socketproxy.SetGlobalMaxProxyCount(gcf.BaseConfigure.ProxyCountLimit)
	config.SetRunMode(runMode)
	config.SetVersion(version)
	log.Printf("RunMode:%s\n", runMode)
	log.Printf("version:%s\tcommit %s, built at %s\n", version, commit, date)

	RunAdminWeb(&gcf.BaseConfigure)

	runTime = time.Now()

	//LoadRuleFromConfigFile(gcf)

	config.PortForwardsRuleListInit()

	//config.DDNSTaskListTaskDetailsInit()
	config.DDNSTaskListConfigureCheck()
	ddnsConf := config.GetDDNSConfigure()
	if ddnsConf.Enable {
		go ddns.Run(time.Duration(ddnsConf.FirstCheckDelay)*time.Second, time.Duration(ddnsConf.Intervals)*time.Second)
	}

	reverseproxy.InitReverseProxyServer()

	//ddns.RunTimer(time.Second, time.Second*30)

	//initProxyList()

	//*****************
	// time.Sleep(time.Microsecond * 50)
	// cruuentPath, _ := fileutils.GetCurrentDirectory()

	// panicFile := fmt.Sprintf("%s/relayport_panic.log", cruuentPath)
	// fileutils.PanicRedirect(panicFile)
	//*****************

	//main goroutine wait
	sigs := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		exit <- true
	}()
	<-exit
}

// func LoadRuleFromConfigFile(pc *config.ProgramConfigure) {
// 	if pc == nil {
// 		return
// 	}
// 	for i := range pc.RelayRuleList {
// 		relayRule, err := rule.CreateRuleByConfigureAndOptions(
// 			pc.RelayRuleList[i].Name,
// 			pc.RelayRuleList[i].Configurestr,
// 			pc.RelayRuleList[i].Options)
// 		if err != nil {
// 			continue
// 		}
// 		relayRule.From = "configureFile" //规则来源
// 		relayRule.IsEnable = pc.RelayRuleList[i].Enable

// 		_, e := rule.AddRuleToGlobalRuleList(false, *relayRule)
// 		if e != nil {
// 			log.Printf("%s\n", e)
// 		}
// 	}
// }
