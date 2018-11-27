package admin

import (
	"github.com/gin-gonic/gin"
	. "github.com/tokenme/tmm/handler"
	"fmt"
	"net/http"
	"github.com/tokenme/tmm/common"
)


/*
type Article struct {
	Fileid      int    `json:"fileid" form:"fileid" `
	Author      string `json:"author" form:"author" `
	Title       string `json:"title" form:"title" `
	Link        string `json:"link" form:"link" `
	SourceUrl   string `json:"source_url" form:"source_url"`
	Cover       string `json:"cover" form:"cover" `
	PublishedOn string `json:"published_on" form:"published_on" `
	Digest      string `json:"digest" form:"digest"`
	Content     string `json:"content" form:"content"`
	Sortid      int    `json:"sortid" form:"sortid" `
	Published   int    `json:"published" form:"published"`
	Online      int    `json:"online" form:"online"`
	Id          int    `json:"id"  form:"id"`
}


func ModifiyArticleHandler(c *gin.Context) {
	db := Service.Db
	var up = Article{}
	if CheckErr(c.Bind(&up), c) {
		return
	}
	Query := `update tmm.articles as art ,tmm.share_tasks as task 
	set art.fileid = %d ,art.author='%s',art.title='%s',art.link='%s',art.source_url='%s',
	art.content='%s' ,art.digest = '%s',art.cover='%s',
	art.published_at='%s',art.published = %d ,task.summary='%s',task.title='%s',task.image='%s' 
	where art.id = %d and task.link = '%s' `
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, up.Id)

	_, _, err := db.Query(Query, up.Fileid, db.Escape(up.Author), db.Escape(up.Title),
		db.Escape(up.Link), db.Escape(up.SourceUrl),
		db.Escape(up.Content), db.Escape(up.Digest), db.Escape(up.Cover),
		db.Escape(up.PublishedOn), up.Published, db.Escape(up.Digest), db.Escape(up.Title),
		db.Escape(up.Cover), up.Id, link)

	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": ""})
}

func AddArticleHandler(c *gin.Context) {
	db := Service.Db
	var article Article
	if CheckErr(c.Bind(&article), c) {
		return
	}
	Query := `INSERT INTO tmm.articles(fileid,author,title,link,source_url,
cover,published_at,digest,content,sortid,Published) 
VALUES (%d,'%s','%s','%s','%s','%s','%s','%s','%s',%d,%d)`
	_, result, err := db.Query(Query, article.Fileid, db.Escape(article.Author), db.Escape(article.Title),
		db.Escape(article.Link), db.Escape(article.SourceUrl), db.Escape(article.Cover),
		db.Escape(article.PublishedOn), db.Escape(article.Digest), db.Escape(article.Content),
		article.Sortid, 1)
	if CheckErr(err, c) {
		return
	}
	link := fmt.Sprintf("https://tmm.tokenmama.io/article/show/%d", result.InsertId())

	Query = `INSERT INTO tmm.share_tasks (creator, title, summary, link, 
	image, points, points_left, bonus, max_viewers,online_status) VALUES(0, '%s', '%s', '%s', '%s', 5000, 5000, 10, 10,-1)`

	_, _, err = db.Query(Query, db.Escape(article.Title), db.Escape(article.Digest),
		db.Escape(link), db.Escape(article.Cover))
	if CheckErr(err, c) {
		return
	}
	rows, _, err := db.Query(`select id from tmm.share_tasks where link = '%s'`, db.Escape(link))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `select Error `, c) {
		return
	}
	_, _, err = db.Query(`INSERT INTO tmm.share_task_categories (task_id,cid,is_auto) VALUES (%d,%d,%d)`, rows[0].Int(0), article.Sortid, 0)
	if CheckErr(err, c) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": article},
	)
}

func GetArticleHandler(c *gin.Context) {
	db := Service.Db
	var (
		offset int
		Query  string
		ArticleList []Article
	)
	sortid, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "3"))
	if page >= 1 {
		offset = (page - 1) * limit
	} else {
		offset = 0
	}

	if sortid == 0 {
		Query = fmt.Sprintf(`
	select 	id,fileid,author,
	title,link,source_url,cover,published_at,
	digest,content,sortid,published from tmm.articles order by published_at DESC
	limit %d offset %d
	`, limit, offset)
	} else {
		Query = fmt.Sprintf(`
		select 	id,fileid,author,
		title,link,source_url,cover,published_at,
		digest,content,sortid,published from tmm.articles 
		where sortid = %d  limit %d offset %d`, sortid, limit, offset)
	}

	Rows, Result, err := db.Query(Query)
	if CheckErr(err, c) {
		return
	}
	if len(Rows) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "没有到数据",
			"data": gin.H{
				"curr_page": page,
				"data":      "",
			},
		})
		return
	}

	for _, Row := range Rows {
		Link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, Row.Int(Result.Map(`id`)))
		Query = `select online_status from tmm.share_tasks where link = '%s'`
		a, result, err := db.Query(Query, Link)
		if CheckErr(err, c) {
			return
		}
		article := Article{
			Id:          Row.Int(Result.Map(`id`)),
			Fileid:      Row.Int(Result.Map(`fileid`)),
			Author:      Row.Str(Result.Map(`author`)),
			Link:        Row.Str(Result.Map(`link`)),
			SourceUrl:   Row.Str(Result.Map(`source_url`)),
			Cover:       Row.Str(Result.Map(`cover`)),
			PublishedOn: Row.ForceLocaltime(Result.Map(`published_at`)).Format(time.RFC3339),
			Digest:      Row.Str(Result.Map(`digest`)),
			Content:     Row.Str(Result.Map(`content`)),
			Sortid:      Row.Int(Result.Map(`sortid`)),
			Published:   Row.Int(Result.Map(`published`)),
		}
		if len(a) == 0 {
			continue
		} else {
			article.Online = a[0].Int(result.Map(`online_status`))
		}
		ArticleList = append(ArticleList, article)
	}
	if sortid != 0 {
		Query = fmt.Sprintf(`select count(*) from tmm.articles  where sortid = %d `, sortid)
	} else {
		Query = ` select count(*) from tmm.articles`
	}
	Rows, _, err = db.Query(Query)
	if CheckErr(err, c) {
		return
	}
	count := Rows[0].Int(0)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": gin.H{
			"curr_page": page,
			"data":      ArticleList,
			"amount":    count,
		},
	}, )

}

func DeleteArticleHandler(c *gin.Context) {
	db := Service.Db
	articleId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, articleId)
	_, _, err = db.Query(`delete tmm.share_tasks  ,tmm.share_task_categories  from tmm.share_tasks,tmm.share_task_categories
	 where tmm.share_tasks.link = '%s'
	 and  tmm.share_task_categories.task_id = tmm.share_tasks.id
		`, db.Escape(link))
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`delete from tmm.articles where id = %d`, articleId)
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`delete from `)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": "",
	})

}
*/

type Online struct {
	Id        int  `json:"id" form:"id"`
	IsArticle bool `json:"is_article" from:"is_article"`
	Status    int  `json:"status" form:"status"`
}

func OnlineAndOfflineHandler(c *gin.Context) {
	db := Service.Db
	var On Online
	if CheckErr(c.Bind(&On), c) {
		return
	}
	if On.IsArticle {
		link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, On.Id)
		Query := `update tmm.share_tasks set online_status = %d where link = '%s'`
		_, _, err := db.Query(Query, On.Status, link)
		if CheckErr(err, c) {
			return
		}
	} else {
		_, _, err := db.Query(`update tmm.share_tasks set online_status = %d where id = %d`, On.Status, On.Id)
		if CheckErr(err, c) {
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"data": "",
	})
}


func GetTypeHander(c *gin.Context) {
	db := Service.Db

	Query := `select id,name from tmm.article_categories `
	maps := make(map[int]string)
	Rows, result, err := db.Query(Query)
	if CheckErr(err, c) {
		return
	}
	for _, Row := range Rows {
		maps[Row.Int(result.Map(`id`))] = Row.Str(result.Map(`name`))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": maps},
	)

}

func VerfiyAdminFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := Service.Db
		userContext, exists := c.Get("USER")
		if Check(!exists, `Need login`, c) {
			return
		}
		User := userContext.(common.User)
		Query := `select role from user_settings where user_id = %d`
		Rows, _, err := db.Query(Query, User.Id)
		if len(Rows) == 0 || err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": `401`,
				"msg":  "invalid User",
			})
			return
		}
		row := Rows[0]
		if row.Int(0) != 1 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				`code`: `401`,
				`msg`:  `The User Must Be Admin`,
			})
			return
		}
		c.Next()
		return
	}

}
