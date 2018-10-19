package main

import (
	"flag"
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/fvbock/endless"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/router"
	"github.com/tokenme/tmm/tools/gc"
	"github.com/tokenme/tmm/tools/orderbook-server"
	"github.com/tokenme/tmm/tools/tokenprofile"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		config         common.Config
		configFlag     common.Config
		configPath     string
		parseTokenFlag bool
	)

	os.Setenv("CONFIGOR_ENV_PREFIX", "-")

	flag.StringVar(&configPath, "config", "config.toml", "configuration file")
	flag.IntVar(&configFlag.Port, "port", 0, "set port")
	flag.StringVar(&configFlag.LogPath, "log", "", "set log file path without filename")
	flag.BoolVar(&configFlag.Debug, "debug", false, "set debug mode")
	flag.BoolVar(&configFlag.EnableWeb, "web", false, "enable http web server")
	flag.BoolVar(&configFlag.EnableGC, "gc", false, "enable gc")
	flag.BoolVar(&configFlag.EnableTx, "tx", false, "enable tx queue handler")
	flag.BoolVar(&parseTokenFlag, "parse-token", false, "enable parse token")
	flag.BoolVar(&configFlag.EnableOrderBook, "orderbook", false, "enable orderbook handler")
	flag.Parse()

	configor.New(&configor.Config{Verbose: configFlag.Debug, ErrorOnUnmatchedKeys: true, Environment: "production"}).Load(&config, configPath)

	if configFlag.Port > 0 {
		config.Port = configFlag.Port
	}
	if configFlag.LogPath != "" {
		config.LogPath = configFlag.LogPath
	}
	if configFlag.EnableWeb {
		config.EnableWeb = configFlag.EnableWeb
	}

	if configFlag.EnableGC {
		config.EnableGC = configFlag.EnableGC
	}

	if configFlag.EnableTx {
		config.EnableTx = configFlag.EnableTx
	}

	if configFlag.EnableOrderBook {
		config.EnableOrderBook = configFlag.EnableOrderBook
	}

	if configFlag.Debug {
		config.Debug = configFlag.Debug
	}

	os.Setenv("PHONE_DATA_DIR", config.PhonedataPath)

	wd, err := os.Getwd()
	if err != nil {
		log.Error(err.Error())
		return
	}
	var logPath string
	if path.IsAbs(config.LogPath) {
		logPath = config.LogPath
	} else {
		logPath = path.Join(wd, config.LogPath)
	}
	defer log.Uninit(log.InitMultiFileAndConsole(logPath, "tmm.log", log.LvERROR))

	raven.SetDSN(config.SentryDSN)
	service := common.NewService(config)
	defer service.Close()
	service.Db.Reconnect()

	if parseTokenFlag {
		tokenprofile.Update(service, config)
		return
	}

	handler.InitHandler(service, config)
	handler.Start()

	gcHandler := gc.New(service, config)
	if config.EnableGC {
		go gcHandler.Start()
	}
	orderbookServer, err := obs.NewServer(service, config, handler.GlobalLock)
	if config.EnableOrderBook {
		if err != nil {
			log.Error(err.Error())
			return
		}
		go orderbookServer.Start()
	}
	//queueManager := sqs.NewManager(config.SQS)
	//queues := make(map[string]sqs.Queue)
	//queues = map[string]sqs.Queue{
	//	config.SQS.TxQueue:  sqs.NewTxQueue(queueManager, service, config),
	//	config.SQS.GasQueue: sqs.NewGasQueue(queueManager, service, config),
	//}
	//if config.EnableTx {
	//	for _, queue := range queues {
	//		go queue.Start()
	//	}
	//}

	if config.EnableWeb {
		if config.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		//gin.DisableBindValidation()
		templatePath := path.Join(config.Template, "./*")
		log.Info("Template path: %s", templatePath)
		r := router.NewRouter(templatePath, config)
		log.Info("%s started at:0.0.0.0:%d", config.AppName, config.Port)
		defer log.Info("%s exit from:0.0.0.0:%d", config.AppName, config.Port)
		srv := endless.NewServer(fmt.Sprintf(":%d", config.Port), r)
		srv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT] = append(
			srv.SignalHooks[endless.PRE_SIGNAL][syscall.SIGINT],
			func() {
				handler.ExitCh <- struct{}{}
			})
		err := srv.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
			return
		}
	} else {
		exitChan := make(chan struct{}, 1)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
			<-ch
			exitChan <- struct{}{}
			close(ch)
		}()
		<-exitChan
	}
	//if config.EnableTx {
	//	for _, queue := range queues {
	//		queue.Stop()
	//	}
	//}
	handler.Close()
	gcHandler.Stop()
	if config.EnableOrderBook {
		orderbookServer.Stop()
	}
}
