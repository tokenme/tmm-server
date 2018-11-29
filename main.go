package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/fvbock/endless"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mkideal/log"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/router"
	"github.com/tokenme/tmm/tools/articleclassifier"
	"github.com/tokenme/tmm/tools/gc"
	"github.com/tokenme/tmm/tools/orderbook-server"
	"github.com/tokenme/tmm/tools/tmmwithdraw"
	"github.com/tokenme/tmm/tools/tokenprofile"
	"github.com/tokenme/tmm/tools/transferwatcher"
	"github.com/tokenme/tmm/tools/txaccelerate"
	"github.com/tokenme/tmm/tools/videospider"
	"github.com/tokenme/tmm/tools/wechatspider"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		config                     common.Config
		configFlag                 common.Config
		configPath                 string
		parseTokenFlag             bool
		articleCrawlerFlag         bool
		articlePublishFlag         bool
		articleClassifierTrainFlag bool
		articleClassifyFlag        bool
		addVideoFlag               string
		accelerateTxFlag           string
		accelerateGasFlag          int64
	)

	os.Setenv("CONFIGOR_ENV_PREFIX", "-")

	flag.StringVar(&configPath, "config", "config.toml", "configuration file")
	flag.IntVar(&configFlag.Port, "port", 0, "set port")
	flag.StringVar(&configFlag.LogPath, "log", "", "set log file path without filename")
	flag.BoolVar(&configFlag.Debug, "debug", false, "set debug mode")
	flag.BoolVar(&configFlag.EnableWeb, "web", false, "enable http web server")
	flag.BoolVar(&configFlag.EnableGC, "gc", false, "enable gc")
	flag.BoolVar(&configFlag.EnableTx, "tx", false, "enable tx queue handler")
	flag.BoolVar(&configFlag.EnableTokenWithdraw, "withdraw", false, "enable token withdraw queue handler")
	flag.BoolVar(&parseTokenFlag, "parse-token", false, "enable parse token")
	flag.BoolVar(&articleCrawlerFlag, "crawle-articles", false, "enable crawle_articles")
	flag.BoolVar(&articlePublishFlag, "publish-articles", false, "enable publish_articles")
	flag.BoolVar(&configFlag.EnableOrderBook, "orderbook", false, "enable orderbook handler")
	flag.StringVar(&accelerateTxFlag, "accelerate-tx", "", "accelerate tx hex")
	flag.Int64Var(&accelerateGasFlag, "gas", 0, "set gas price")
	flag.BoolVar(&articleClassifierTrainFlag, "train-article-classifier", false, "enable article classifer training")
	flag.BoolVar(&articleClassifyFlag, "classify-articles", false, "enable articles classify")
	flag.StringVar(&addVideoFlag, "add-video", "", "add video")
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

	if configFlag.EnableTokenWithdraw {
		config.EnableTokenWithdraw = configFlag.EnableTokenWithdraw
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

	if accelerateTxFlag != "" && accelerateGasFlag > 0 {
		err := txaccelerate.Accelerate(service, config, accelerateTxFlag, accelerateGasFlag)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	if parseTokenFlag {
		tokenprofile.Update(service, config)
		return
	}

	if addVideoFlag != "" {
		spider := videospider.NewClient(service, config)
		video, err := spider.Get(addVideoFlag)
		if err != nil {
			log.Error(err.Error())
		}
		err = spider.Save(video)
		if err != nil {
			log.Error(err.Error())
		}
		spew.Dump(video)
		return
	}

	if articleCrawlerFlag {
		crawler := wechatspider.NewCrawler(service, config)
		crawler.Run()
		return
	}

	if articlePublishFlag {
		crawler := wechatspider.NewCrawler(service, config)
		err := crawler.Publish()
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	if articleClassifierTrainFlag {
		trainer := articleclassifier.NewClassifier(service, config)
		err := trainer.Train()
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	if articleClassifyFlag {
		classifier := articleclassifier.NewClassifier(service, config)
		err := classifier.LoadModel()
		if err != nil {
			log.Error(err.Error())
			return
		}
		classifier.ClassifyDocs()
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
	transferWatcher, err := transferwatcher.NewWatcher(config.TMMTokenAddress, service, config)
	if err != nil {
		log.Error(err.Error())
		return
	}
	go func() {
		err = transferWatcher.Start()
		if err != nil {
			log.Error(err.Error())
			//return
		}
	}()
	tokenWithdraw := tmmwithdraw.NewService(service, config)
	if config.EnableTokenWithdraw {
		go tokenWithdraw.Start()
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
	if config.EnableTokenWithdraw {
		go tokenWithdraw.Stop()
	}
	transferWatcher.Stop()
}
