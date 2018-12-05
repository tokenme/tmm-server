package articleclassifier

import (
	"bytes"
	"fmt"
	"github.com/AlasdairF/Classifier"
	"github.com/PuerkitoBio/goquery"
	//"github.com/davecgh/go-spew/spew"
	"github.com/go-ego/gse"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/tools/articleclassifier/Uint16Uint64"
	"github.com/tokenme/tmm/tools/tfidf"
	"github.com/tokenme/tmm/utils"
	"log"
	"strconv"
	"strings"
)

const MaxDocWords int = 200

type Doc struct {
	Id    uint64
	Cid   uint16
	Words []string
}

type Task struct {
	Id  uint64
	Cid uint16
}

type Classifier struct {
	service    *common.Service
	config     common.Config
	trainer    *classifier.Trainer
	classifier *classifier.Classifier
	Gse        *gse.Segmenter
	weighter   *tfidf.Weighter
}

func NewClassifier(service *common.Service, config common.Config) *Classifier {
	obj := &Classifier{
		service:  service,
		config:   config,
		Gse:      &gse.Segmenter{},
		weighter: tfidf.NewWeighter(),
	}
	obj.Gse.LoadDict(config.GseDict)
	return obj
}

func (this *Classifier) LoadModel() (err error) {
	this.classifier, err = classifier.Load(this.config.ArticleClassifierModel)
	if err != nil {
		return err
	}
	return this.weighter.Load(this.config.ArticleClassifierTFIDF)
}

func (this *Classifier) getCategories() ([][]byte, error) {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT id FROM tmm.article_categories`)
	if err != nil {
		return nil, err
	}
	var categories [][]byte
	for _, row := range rows {
		cat := utils.Uint16ToByte(uint16(row.Uint(0)))
		categories = append(categories, cat)
	}
	return categories, nil
}

func (this *Classifier) Classify(title string, digest string, author string, html string) ([]uint16, error) {
	var content string
	if html != "" {
		reader := bytes.NewBuffer([]byte(html))
		page, err := goquery.NewDocumentFromReader(reader)
		if err == nil {
			content = page.Text()
		}
	}

	txt := fmt.Sprintf("%s %s %s %s", title, digest, author, content)
	segments := this.Gse.Segment([]byte(txt))
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
	var topWords [][]byte
	for idx, item := range scores {
		if idx >= MaxDocWords {
			break
		}
		topWords = append(topWords, []byte(item.Key))
	}
	categoryScores := this.classifier.Classify(topWords)
	sorted := sortUint16Uint64.New(categoryScores)
	sortUint16Uint64.Desc(sorted)
	cats := this.classifier.Categories
	var cids []uint16
	for _, score := range sorted {
		cid := utils.ByteToUint16(cats[score.K])
		//fmt.Println(i, `Category`, cid, `K`, score.K, `Score`, score.V)
		if score.V == 0 {
			continue
		}
		cids = append(cids, cid)
	}
	return cids, nil
}

func (this *Classifier) ClassifyDocs() {
	log.Println("Classifing Docs...")
	db := this.service.Db
	var (
		startId uint64
		limit   uint64 = 1000
	)
	for {
		endId := startId
		log.Println("Reading docs: ", startId)
		rows, _, err := db.Query(`SELECT st.id, st.link, st.title, st.summary FROM tmm.share_tasks AS st WHERE NOT EXISTS (SELECT 1 FROM tmm.share_task_categories AS stc WHERE stc.task_id=st.id LIMIT 1) AND st.creator=0 AND st.id>%d AND st.online_status=1 ORDER BY st.id ASC LIMIT %d`, startId, limit)
		if err != nil {
			log.Println(err.Error())
			break
		}
		var (
			tasks      = make(map[uint64]*Task)
			extraTasks []*Task
			idArr      []string
		)
		for _, row := range rows {
			endId = row.Uint64(0)
			link := row.Str(1)
			idStr := strings.Replace(link, "https://tmm.tokenmama.io/article/show/", "", -1)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil || id == 0 {
				task := &Task{Id: endId}
				cids, _ := this.Classify(row.Str(2), row.Str(3), "", "")
				if len(cids) > 0 {
					task.Cid = cids[0]
					log.Println("Doc: ", link, ", Cid:", task.Cid)
				}
				extraTasks = append(extraTasks, task)
				continue
			}
			tasks[id] = &Task{
				Id: endId,
			}
			idArr = append(idArr, idStr)
		}
		if len(idArr) > 0 {
			rows, _, err := db.Query(`SELECT id, title, digest, author, content FROM tmm.articles WHERE id IN (%s)`, strings.Join(idArr, ","))
			if err != nil {
				continue
			}
			for _, row := range rows {
				if task, found := tasks[row.Uint64(0)]; found {
					cids, _ := this.Classify(row.Str(1), row.Str(2), row.Str(3), row.Str(4))
					if len(cids) > 0 {
						task.Cid = cids[0]
						log.Println("Doc: ", fmt.Sprintf("https://tmm.tokenmama.io/article/show/%d", row.Uint64(0)), ", Cid:", task.Cid)
					}
				}
			}
		}
		var val []string
		for _, task := range tasks {
			if task.Cid == 0 {
				continue
			}
			val = append(val, fmt.Sprintf("(%d, %d, 0)", task.Id, task.Cid))
		}
		for _, task := range extraTasks {
			if task.Cid == 0 {
				continue
			}
			val = append(val, fmt.Sprintf("(%d, %d, 0)", task.Id, task.Cid))
		}
		if len(val) > 0 {
			log.Println("Saving ", len(val), " Classified Docs...")
			_, _, err := db.Query(`INSERT IGNORE INTO tmm.share_task_categories(task_id, cid, is_auto) VALUES %s`, strings.Join(val, ","))
			if err != nil {
				log.Println(err.Error())
			}
		}
		if endId == startId {
			break
		}
		startId = endId
	}
}

func (this *Classifier) Train() error {
	this.trainer = new(classifier.Trainer)
	categories, err := this.getCategories()
	if err != nil {
		return err
	}
	this.trainer.DefineCategories(categories)
	var (
		trainDocs uint
		testDocs  uint
	)
	docs := this.getDocs()
	log.Printf("Total Docs:%d", len(docs))
	for _, doc := range docs {
		var tokens [][]byte
		for _, w := range doc.Words {
			tokens = append(tokens, []byte(w))
		}
		cat := utils.Uint16ToByte(doc.Cid)
		rand := utils.RangeRandUint64(1, 9)
		if rand > 3 {
			trainDocs += 1
			this.trainer.AddTrainingDoc(cat, tokens)
		} else {
			testDocs += 1
			this.trainer.AddTestDoc(cat, tokens)
		}
	}
	verbose := true
	allowance, maxscore, err := this.trainer.Test(verbose)
	if err != nil {
		return err
	}
	log.Printf("Allowance: %v, Maxscore: %v, Train:%d, Test:%d\n", allowance, maxscore, trainDocs, testDocs)
	this.trainer.Create(allowance, maxscore)
	return this.trainer.Save(this.config.ArticleClassifierModel)
}

func (this *Classifier) getDocs() []*Doc {
	log.Println("Loading Docs...")
	db := this.service.Db
	var (
		startId   uint64
		limit     uint64 = 1000
		docs             = make(map[uint64]*Doc)
		extraDocs []*Doc
	)
	for {
		endId := startId
		rows, _, err := db.Query(`SELECT st.id, stc.cid, st.link, st.title, st.summary FROM tmm.share_tasks AS st INNER JOIN tmm.share_task_categories AS stc ON (stc.task_id=st.id) WHERE st.creator=0 AND stc.is_auto=1 AND st.id>%d AND st.online_status=1 ORDER BY st.id ASC LIMIT %d`, startId, limit)
		if err != nil {
			log.Println(err.Error())
			break
		}
		var idArr []string
		for _, row := range rows {
			endId = row.Uint64(0)
			cid := row.Uint(1)
			link := row.Str(2)
			idStr := strings.Replace(link, "https://tmm.tokenmama.io/article/show/", "", -1)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil || id == 0 {
				doc := &Doc{Cid: uint16(cid)}
				txt := fmt.Sprintf("%s %s", row.Str(3), row.Str(4))
				segments := this.Gse.Segment([]byte(txt))
				var words []string
				for _, seg := range segments {
					w := strings.ToLower(strings.TrimSpace(seg.Token().Text()))
					if len(w) <= 1 || strings.Contains("。？！，、；：“”‘'（）《》〈〉【】『』「」﹃﹄〔〕…—～﹏￥\n", w) {
						continue
					}
					words = append(words, w)
				}
				if len(words) == 0 {
					continue
				}
				this.weighter.AddDoc(words)
				doc.Words = words
				extraDocs = append(extraDocs, doc)
				continue
			}
			docs[id] = &Doc{
				Cid: uint16(cid),
				Id:  id,
			}
			idArr = append(idArr, idStr)
		}
		if len(idArr) > 0 {
			rows, _, err := db.Query(`SELECT id, title, digest, author, content FROM tmm.articles WHERE id IN (%s)`, strings.Join(idArr, ","))
			if err != nil {
				continue
			}
			for _, row := range rows {
				if doc, found := docs[row.Uint64(0)]; found {
					//log.Printf("Article: %d, Cat: %d\n", doc.Id, doc.Cid)
					var content string
					reader := bytes.NewBuffer([]byte(row.Str(4)))
					page, err := goquery.NewDocumentFromReader(reader)
					if err == nil {
						content = page.Text()
					}
					txt := fmt.Sprintf("%s %s %s %s", row.Str(1), row.Str(2), row.Str(3), content)
					segments := this.Gse.Segment([]byte(txt))
					var words []string
					for _, seg := range segments {
						w := strings.ToLower(strings.TrimSpace(seg.Token().Text()))
						if len(w) <= 1 || strings.Contains("。？！，、；：“”‘'（）《》〈〉【】『』「」﹃﹄〔〕…—～﹏￥\n", w) {
							continue
						}
						words = append(words, w)
					}
					if len(words) == 0 {
						continue
					}
					this.weighter.AddDoc(words)
					doc.Words = words
				}
			}
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	err := this.weighter.Save(this.config.ArticleClassifierTFIDF)
	if err != nil {
		log.Println(err.Error())
	}
	var documents []*Doc
	for _, doc := range docs {
		if len(doc.Words) == 0 {
			continue
		}
		scores := this.weighter.Score(doc.Words)
		var topWords []string
		for idx, item := range scores {
			if idx >= MaxDocWords {
				break
			}
			topWords = append(topWords, item.Key)
		}
		doc.Words = topWords
		documents = append(documents, doc)
	}

	for _, doc := range extraDocs {
		if len(doc.Words) == 0 {
			continue
		}
		scores := this.weighter.Score(doc.Words)
		var topWords []string
		for idx, item := range scores {
			if idx >= MaxDocWords {
				break
			}
			topWords = append(topWords, item.Key)
		}
		doc.Words = topWords
		documents = append(documents, doc)
	}
	return documents
}
