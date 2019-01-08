package main

import (
	//"github.com/tokenme/tmm/tools/orderbook-server"
	//"github.com/tokenme/tmm/tools/transferwatcher"
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/fvbock/endless"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/mkideal/log"
	"github.com/shopspring/decimal"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/router"
	"github.com/tokenme/tmm/tools/articleclassifier"
	"github.com/tokenme/tmm/tools/etherscanspider"
	"github.com/tokenme/tmm/tools/gc"
	"github.com/tokenme/tmm/tools/invitebonus"
	"github.com/tokenme/tmm/tools/qutoutiaospider"
	"github.com/tokenme/tmm/tools/redpacket"
	"github.com/tokenme/tmm/tools/tmmwithdraw"
	"github.com/tokenme/tmm/tools/tokenprofile"
	"github.com/tokenme/tmm/tools/toutiaospider"
	"github.com/tokenme/tmm/tools/txaccelerate"
	"github.com/tokenme/tmm/tools/videospider"
	"github.com/tokenme/tmm/tools/wechatspider"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"
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
		addArticlesFlag            bool
		addVideoFlag               string
		updateVideoFlag            bool
		accelerateTxFlag           string
		accelerateGasFlag          int64
		dataFlag                   string
		nonceFlag                  uint64
		ucoinHoldersFlag           bool
		activeBonusFlag            bool
		fixInviteBonusFlag         bool
		redpacketFlag              bool
		redpacketTokensFlag        int64
		redpacketRecipientsFlag    uint
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
	flag.StringVar(&dataFlag, "data", "", "tx data")
	flag.Uint64Var(&nonceFlag, "nonce", 0, "tx nonce")
	flag.BoolVar(&articleClassifierTrainFlag, "train-article-classifier", false, "enable article classifer training")
	flag.BoolVar(&articleClassifyFlag, "classify-articles", false, "enable articles classify")
	flag.BoolVar(&updateVideoFlag, "update-videos", false, "enable update videos")
	flag.StringVar(&addVideoFlag, "add-video", "", "add video")
	flag.BoolVar(&addArticlesFlag, "add-articles", false, "enable add articles")
	flag.BoolVar(&ucoinHoldersFlag, "update-holders", false, "enable update ucoin holders")
	flag.BoolVar(&activeBonusFlag, "active-bonus", false, "enable check active bonus")
	flag.BoolVar(&redpacketFlag, "redpacket", false, "enable redpacket service")
	flag.BoolVar(&fixInviteBonusFlag, "fix-invite-bonus", false, "enable fix invite bonus")
	flag.Int64Var(&redpacketTokensFlag, "rp-tokens", 0, "set redpacket tokens")
	flag.UintVar(&redpacketRecipientsFlag, "rp-users", 0, "set redpacket recipients")
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
		err := txaccelerate.Accelerate(service, config, accelerateTxFlag, accelerateGasFlag, dataFlag, nonceFlag)
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
		videoSpider := videospider.NewClient(service, config)
		video, err := videoSpider.Get(addVideoFlag)
		if err != nil {
			log.Error(err.Error())
		}
		err = videoSpider.Save(video)
		if err != nil {
			log.Error(err.Error())
		}
		spew.Dump(video)
		return
	}

	if updateVideoFlag {
		videoSpider := videospider.NewClient(service, config)
		go videoSpider.StartUpdateVideosService()
		exitChan := make(chan struct{}, 1)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
			<-ch
			exitChan <- struct{}{}
			close(ch)
		}()
		<-exitChan
		videoSpider.Stop()
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

	if redpacketTokensFlag > 0 && redpacketRecipientsFlag > 0 {
		rp, err := common.NewRedpacket(service, 0, decimal.New(redpacketTokensFlag, 0), redpacketRecipientsFlag)
		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("New Redpacket:%d, Token:%s, Recipients:%d", rp.Id, rp.Tmm.String(), rp.Recipients)
		}
		return
	}

	if addArticlesFlag {
		addWxArticlesCh := make(chan struct{}, 1)
		addToutiaoArticlesCh := make(chan struct{}, 1)
		addQutoutiaoArticlesCh := make(chan struct{}, 1)
		exitChan := make(chan struct{}, 1)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
			<-ch
			exitChan <- struct{}{}
			close(ch)
		}()
		go addQutoutiaoArticles(addQutoutiaoArticlesCh, service, config)
		go addToutiaoArticles(addToutiaoArticlesCh, service, config)
		go addWxArticles(addWxArticlesCh, service, config)
		for {
			select {
			case <-addWxArticlesCh:
				go addWxArticles(addWxArticlesCh, service, config)
			case <-addToutiaoArticlesCh:
				go addToutiaoArticles(addToutiaoArticlesCh, service, config)
			case <-addQutoutiaoArticlesCh:
				go addQutoutiaoArticles(addQutoutiaoArticlesCh, service, config)
			case <-exitChan:
				close(addWxArticlesCh)
				close(addToutiaoArticlesCh)
				return
			}
		}
	}

	if ucoinHoldersFlag {
		updateHoldersCh := make(chan struct{}, 1)
		exitChan := make(chan struct{}, 1)
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTERM)
			<-ch
			exitChan <- struct{}{}
			close(ch)
		}()
		updateHoldersCh <- struct{}{}
		for {
			select {
			case <-updateHoldersCh:
				go func(ch chan struct{}) {
					defer func() {
						log.Info("Sleep for 1 hour")
						time.Sleep(1 * time.Hour)
						ch <- struct{}{}
					}()
					etherscanspider.GetHolders(service)
				}(updateHoldersCh)
			case <-exitChan:
				close(updateHoldersCh)
				return
			}
		}
		return
	}

	gcHandler := gc.New(service, config)
	if config.EnableGC {
		go gcHandler.Start()
	}
	/*
		orderbookServer, err := obs.NewServer(service, config, handler.GlobalLock)
		if config.EnableOrderBook {
			if err != nil {
				log.Error(err.Error())
				return
			}
			go orderbookServer.Start()
		}
	*/
	/*
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
	*/
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
		handler.InitHandler(service, config)
		handler.Start()
		if config.Debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		activeBonusService := invitebonus.NewService(service, config, handler.GlobalLock)
		if fixInviteBonusFlag {
			activeBonusService.FixBonus()
			return
		}
		if activeBonusFlag {
			go activeBonusService.Start()
		}
		redpacketService := redpacket.NewService(service, config, handler.GlobalLock)
		if redpacketFlag {
			go redpacketService.Start()
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
		if activeBonusFlag {
			activeBonusService.Stop()
		}
		if redpacketFlag {
			redpacketService.Stop()
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
	if config.EnableWeb {
		handler.Close()
	}
	gcHandler.Stop()
	/*
		if config.EnableOrderBook {
			orderbookServer.Stop()
		}
	*/
	if config.EnableTokenWithdraw {
		tokenWithdraw.Stop()
	}
	//transferWatcher.Stop()
}

func addWxArticles(addWxArticlesCh chan<- struct{}, service *common.Service, config common.Config) {
	defer func() {
		time.Sleep(3 * time.Hour)
		addWxArticlesCh <- struct{}{}
	}()
	crawler := wechatspider.NewCrawler(service, config)
	crawler.Run()
	err := crawler.Publish()
	if err != nil {
		log.Error(err.Error())
	}
	classifier := articleclassifier.NewClassifier(service, config)
	err = classifier.LoadModel()
	if err != nil {
		log.Error(err.Error())
		return
	}
	classifier.ClassifyDocs()
}

func addToutiaoArticles(addToutiaoArticlesCh chan<- struct{}, service *common.Service, config common.Config) {
	defer func() {
		time.Sleep(5 * time.Minute)
		addToutiaoArticlesCh <- struct{}{}
	}()
	crawler := toutiaospider.NewCrawler(service, config)
	num, err := crawler.Run()
	if err != nil {
		log.Error(err.Error())
		return
	}
	if num == 0 {
		return
	}
	err = crawler.Publish()
	if err != nil {
		log.Error(err.Error())
	}
	classifier := articleclassifier.NewClassifier(service, config)
	err = classifier.LoadModel()
	if err != nil {
		log.Error(err.Error())
		return
	}
	classifier.ClassifyDocs()
}

func addQutoutiaoArticles(addQutoutiaoArticlesCh chan<- struct{}, service *common.Service, config common.Config) {
	defer func() {
		time.Sleep(2 * time.Minute)
		addQutoutiaoArticlesCh <- struct{}{}
	}()
	crawler := qutoutiaospider.NewCrawler(service, config)
	num, err := crawler.Run()
	if err != nil {
		log.Error(err.Error())
		return
	}
	if num == 0 {
		return
	}
	err = crawler.Publish()
	if err != nil {
		log.Error(err.Error())
	}
	classifier := articleclassifier.NewClassifier(service, config)
	err = classifier.LoadModel()
	if err != nil {
		log.Error(err.Error())
		return
	}
	classifier.ClassifyDocs()
}
