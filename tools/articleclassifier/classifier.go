package articleclassifier

import (
	"github.com/AlasdairF/Classifier"
	"github.com/go-ego/gse"
	"github.com/tokenme/tmm/common"
	"github.com/tokenme/tmm/utils"
	"log"
	"strconv"
	"strings"
)

type Doc struct {
	Id    uint64
	Cid   uint16
	Words []string
}

type Classifier struct {
	service *common.Service
	config  common.Config
	trainer *classifier.Trainer
	Gse     *gse.Segmenter
}

func NewClassifier(service *common.Service, config common.Config) *Classifier {
	classifier := &Classifier{
		service: service,
		config:  config,
		trainer: new(classifier.Trainer),
		Gse:     &gse.Segmenter{},
	}
	classifier.Gse.LoadDict(config.GseDict)
	classifier.getCategories()
	return classifier
}

func (this *Classifier) getCategories() error {
	db := this.service.Db
	rows, _, err := db.Query(`SELECT id FROM tmm.article_categories`)
	if err != nil {
		return err
	}
	var categories [][]byte
	for _, row := range rows {
		cat := utils.Uint16ToByte(uint16(row.Uint(0)))
		categories = append(categories, cat)
	}
	this.trainer.DefineCategories(categories)
	return nil
}

func (this *Classifier) Train() error {
	db := this.service.Db
	var (
		startId uint64
		limit   uint64 = 1000
	)
	for {
		endId := startId
		rows, _, err := db.Query(`SELECT st.id, stc.cid, st.link FROM tmm.share_tasks AS st INNER JOIN tmm.share_task_categories AS stc ON (stc.task_id=st.id) WHERE st.creator=0 AND stc.is_auto=0 ORDER BY st.id ASC LIMIT %d, %d`, startId, limit)
		if err != nil {
			log.Println(err.Error())
			break
		}
		var (
			docs  = make(map[uint64]*Doc, len(rows))
			idArr []string
		)
		for _, row := range rows {
			endId = row.Uint64(0)
			cid := row.Uint(1)
			link := row.Str(2)
			idStr := strings.Replace(link, "https://tmm.tokenmama.io/article/show/", "", -1)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil || id == 0 {
				log.Println(err.Error())
				continue
			}
			docs[id] = &Doc{
				Cid: uint16(cid),
				Id:  id,
			}
			idArr = append(idArr, idStr)
		}
		if len(idArr) > 0 {
			rows, _, err := db.Query(`SELECT id, content FROM tmm.articles WHERE id IN (%s)`, strings.Join(idArr, ","))
			if err != nil {
				continue
			}
			for _, row := range rows {
				if doc, found := docs[row.Uint64(0)]; found {
					log.Printf("Article: %d, Cat: %d\n", doc.Id, doc.Cid)
					content := row.Str(1)
					segments := this.Gse.Segment([]byte(content))
					var tokens [][]byte
					for _, seg := range segments {
						w := strings.ToLower(seg.Token().Text())
						tokens = append(tokens, []byte(w))
					}
					cat := utils.Uint16ToByte(doc.Cid)
					rand := utils.RangeRandUint64(1, 9)
					if rand > 3 {
						this.trainer.AddTrainingDoc(cat, tokens)
					} else {
						this.trainer.AddTestDoc(cat, tokens)
					}

				}
			}
		}
		if endId == startId {
			break
		}
		startId = endId
	}
	verbose := true
	allowance, maxscore, err := this.trainer.Test(verbose)
	if err != nil {
		return err
	}
	log.Printf("Allowance: %v, Maxscore: %v\n", allowance, maxscore)
	this.trainer.Create(allowance, maxscore)
	return this.trainer.Save(this.config.ArticleClassifierModel)
}
