package articlesuggest

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"github.com/davecgh/go-spew/spew"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/go-ego/gse"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/tfidf"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxDocWords int = 200
const SUGGEST_ENGIN_KEY = "share-task-suggest-%d"

type UserLog struct {
	UserId    uint64
	ArticleId uint64
	Ts        uint64
	Words     map[string]float64
}

type Task struct {
	TaskId    uint64
	ArticleId uint64
	Words     map[string]float64
}

type Engine struct {
	service   *common.Service
	config    common.Config
	Gse       *gse.Segmenter
	weighter  *tfidf.Weighter
	tasks     []*Task
	exitCh    chan struct{}
	canExitCh chan struct{}
	sync.RWMutex
}

func NewEngine(service *common.Service, config common.Config) *Engine {
	obj := &Engine{
		service:   service,
		config:    config,
		Gse:       &gse.Segmenter{},
		weighter:  tfidf.NewWeighter(),
		exitCh:    make(chan struct{}, 1),
		canExitCh: make(chan struct{}, 1),
	}
	return obj
}

func (this *Engine) Start() {
	this.Gse.LoadDict(this.config.GseDict)
	this.weighter.Load(this.config.ArticleClassifierTFIDF)
	taskTicker := time.NewTicker(3 * time.Hour)
	logTicker := time.NewTicker(1 * time.Hour)
	this.getTasks()
	this.getLogs()
	for {
		select {
		case <-taskTicker.C:
			this.getTasks()
		case <-logTicker.C:
			this.getLogs()
		case <-this.exitCh:
			taskTicker.Stop()
			logTicker.Stop()
			this.canExitCh <- struct{}{}
			return
		}
	}
}

func (this *Engine) Stop() {
	this.exitCh <- struct{}{}
	<-this.canExitCh
}

func (this *Engine) Match(userId uint64, page uint, limit uint) []uint64 {
	var taskIds []uint64
	startId := int((page - 1) * limit)
	endId := startId + int(limit)
	redisConn := this.service.Redis.Master.Get()
	defer redisConn.Close()
	infoKey := fmt.Sprintf(SUGGEST_ENGIN_KEY, userId)
	buf, err := redis.Bytes(redisConn.Do("GET", infoKey))
	if err != nil && buf != nil {
		err := json.Unmarshal(buf, &taskIds)
		if err != nil {
			log.Println(err.Error())
		} else {
			totalIds := len(taskIds)
			if startId >= totalIds {
				return nil
			}
			if endId >= totalIds {
				endId = totalIds
			}
			log.Println("Start:", startId, ", End:", endId, ", Total:", totalIds)
			return taskIds[startId:endId]
		}
	}

	this.RLock()
	tasks := this.tasks
	this.RUnlock()
	if tasks == nil {
		log.Println("tasks are empty")
		redisConn.Do("SETEX", infoKey, 10, "[]")
		return nil
	}
	db := this.service.Db
	rows, _, err := db.Query(`SELECT kw, score FROM tmm.user_reading_kws WHERE user_id=%d`, userId)
	if err != nil {
		log.Println(err.Error())
		redisConn.Do("SETEX", infoKey, 1*60, "[]")
		return nil
	}
	words := make(map[string]float64, len(rows))
	for _, row := range rows {
		words[row.Str(0)] = row.Float(1)
	}
	taskScores := make(map[uint64]float64)
	for _, task := range tasks {
		var score float64
		for w, v := range task.Words {
			if uv, found := words[w]; found {
				score += v * uv
			}
		}
		taskScores[task.TaskId] = score
	}
	sorter := NewScoreSorter(taskScores)
	sort.Sort(sort.Reverse(sorter))

	for _, i := range sorter {
		taskIds = append(taskIds, i.Key)
	}
	js, err := json.Marshal(taskIds)
	if err == nil {
		redisConn.Do("SETEX", infoKey, 1*60, string(js))
	}
	totalIds := len(taskIds)
	if startId >= totalIds {
		return nil
	}
	if endId >= totalIds {
		endId = totalIds
	}
	log.Println("Start:", startId, ", End:", endId, ", Total:", totalIds)
	return taskIds[startId:endId]
}

func (this *Engine) getLogs() {
	var (
		db        = this.service.Db
		startTime string
		endTime   string
		idMap     = make(map[uint64]struct{})
		idArr     []string
		userLogs  = make(map[uint64][]*UserLog)
		tasks     = make(map[uint64]Task)
	)
	for {
		endTime = startTime
		rows, _, err := db.Query(`SELECT l.user_id, st.id, st.link, st.id, l.updated_at FROM tmm.reading_logs AS l INNER JOIN tmm.share_tasks AS st ON (st.id=l.task_id) WHERE st.creator=0 AND l.updated_at>='%s' ORDER BY l.updated_at ASC LIMIT 1000`, startTime)
		if err != nil {
			log.Println(err.Error())
			break
		}
		for _, row := range rows {
			endTime = row.ForceLocaltime(4).Format("2006-01-02 15:04:05")
			userId := row.Uint64(0)
			//taskId := row.Uint64(1)
			link := row.Str(2)
			ts := row.Uint64(3)
			idStr := strings.Replace(link, "https://tmm.tokenmama.io/article/show/", "", -1)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil || id == 0 {
				log.Println(err.Error())
				continue
			}
			userLog := &UserLog{
				UserId:    userId,
				ArticleId: id,
				Ts:        ts,
			}
			userLogs[userId] = append(userLogs[userId], userLog)
			if _, found := idMap[id]; found {
				continue
			}
			idMap[id] = struct{}{}
			idArr = append(idArr, fmt.Sprintf("%d", id))
		}
		if len(idArr) > 0 {
			rows, _, err := db.Query(`SELECT id, content FROM tmm.articles WHERE id IN (%s)`, strings.Join(idArr, ","))
			if err != nil {
				log.Println(err.Error())
				continue
			}
			for _, row := range rows {
				articleId := row.Uint64(0)
				content := row.Str(1)
				topWords, err := this.topWords(content)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				tasks[articleId] = Task{
					ArticleId: articleId,
					Words:     topWords,
				}
			}
		}

		if startTime == endTime {
			break
		}
		startTime = endTime
	}
	logs := make(map[uint64]*UserLog)
	for userId, uLogs := range userLogs {
		for _, log := range uLogs {
			if task, found := tasks[log.ArticleId]; found {
				if _, found := logs[userId]; !found {
					logs[userId] = &UserLog{
						UserId: userId,
						Words:  make(map[string]float64),
					}
				}
				for w, v := range task.Words {
					logs[userId].Words[w] += v * math.Log1p(float64(log.Ts))
				}
			}
		}
	}
	var val []string
	for userId, l := range logs {
		for w, v := range l.Words {
			val = append(val, fmt.Sprintf("(%d, '%s', %.9f)", userId, db.Escape(w), v))
		}
		if len(val) >= 1000 {
			_, _, err := db.Query(`INSERT INTO tmm.user_reading_kws (user_id, kw, score) VALUES %s ON DUPLICATE KEY UPDATE score=VALUES(score)`, strings.Join(val, ","))
			if err != nil {
				log.Println(err.Error())
			}
			val = []string{}
		}
	}
	if len(val) > 0 {
		_, _, err := db.Query(`INSERT INTO tmm.user_reading_kws (user_id, kw, score) VALUES %s ON DUPLICATE KEY UPDATE score=VALUES(score)`, strings.Join(val, ","))
		if err != nil {
			log.Println(err.Error())
		}
		val = []string{}
	}
}

func (this *Engine) getTasks() {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT
    st.id,
    st.link
FROM tmm.share_tasks AS st
WHERE st.points_left>0 AND st.online_status = 1
ORDER BY st.bonus DESC, st.id DESC LIMIT 2000`)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var (
		idArr   []string
		tasks   []*Task
		taskMap = make(map[uint64]Task)
	)
	for _, row := range rows {
		taskId := row.Uint64(0)
		link := row.Str(1)
		idStr := strings.Replace(link, "https://tmm.tokenmama.io/article/show/", "", -1)
		id, err := strconv.ParseUint(idStr, 10, 64)
		tasks = append(tasks, &Task{TaskId: taskId, ArticleId: id})
		if err != nil || id == 0 {
			log.Println(err.Error())
			continue
		}
		idArr = append(idArr, fmt.Sprintf("%d", id))
	}
	if len(idArr) > 0 {
		rows, _, err := db.Query(`SELECT id, content FROM tmm.articles WHERE id IN (%s)`, strings.Join(idArr, ","))
		if err != nil {
			return
		}
		for _, row := range rows {
			articleId := row.Uint64(0)
			content := row.Str(1)
			topWords, err := this.topWords(content)
			if err != nil {
				continue
			}
			taskMap[articleId] = Task{
				ArticleId: articleId,
				Words:     topWords,
			}
		}
	}
	for _, task := range tasks {
		if t, found := taskMap[task.ArticleId]; found {
			task.Words = t.Words
		}
	}
	this.Lock()
	this.tasks = tasks
	this.Unlock()
}

func (this *Engine) topWords(html string) (map[string]float64, error) {
	reader := bytes.NewBuffer([]byte(html))
	page, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	segments := this.Gse.Segment([]byte(page.Text()))
	var words []string
	for _, seg := range segments {
		w := strings.ToLower(strings.TrimSpace(seg.Token().Text()))
		if len(w) <= 1 || strings.Contains("。？！，、；：“”‘'（）《》〈〉【】『』「」﹃﹄〔〕…—～﹏￥\n", w) {
			continue
		}
		words = append(words, w)
	}
	if len(words) == 0 {
		return nil, nil
	}
	scores := this.weighter.Score(words)
	topWords := make(map[string]float64)
	for idx, item := range scores {
		if idx >= MaxDocWords {
			break
		}
		topWords[item.Key] = item.Val
	}
	return topWords, nil
}
