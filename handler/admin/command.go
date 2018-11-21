package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	. "github.com/tokenme/tmm/handler"
	"github.com/tokenme/tmm/common"
	"github.com/ziutek/mymysql/autorc"
	"net/http"
	"time"
	"html/template"
	"fmt"
	"github.com/shopspring/decimal"
)

type AddArticle struct {
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
}

type Article struct {
	Id          int           `json:"id" form:"id" ,binding:"required" `
	Author      string        `json:"author,omitempty" form:"author"`
	Title       string        `json:"title,omitempty" form:"title"`
	SourceUrl   string        `json:"source_url,omitempty" form:"source_url"`
	Content     template.HTML `json:"content" form:"content" `
	Digest      string        `json:"digest,omitempty" form:"digest"`
	Cover       string        `json:"cover,omitempty" form:"cover"`
	PublishedOn string        `json:"published_on,omitempty" form:"published_on"`
}

func ModifiyArticleHandler(c *gin.Context) {
	db := Service.Db
	if !VerfiyAdmin(c, db) {
		return
	}
	var up = Article{}
	if CheckErr(c.Bind(&up), c) {
		return
	}

	Query := `update tmm.articles as art ,tmm.share_tasks as task 
	set art.author='%s',art.title='%s',art.source_url='%s',
	art.content='%s' ,art.digest = '%s',art.cover='%s',
	art.published_at='%s',task.summary='%s',task.title='%s',task.image='%s' 
	where art.id = %d and task.link = '%s' `
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, up.Id)

	_, _, err := db.Query(Query, db.Escape(up.Author), db.Escape(up.Title), db.Escape(up.SourceUrl),
		db.Escape(string(up.Content)), db.Escape(up.Digest), db.Escape(up.Cover),
		db.Escape(up.PublishedOn), db.Escape(up.Digest), db.Escape(up.Title),
		db.Escape(up.Cover), up.Id, link)

	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": "" })
}

func AddArticleHandler(c *gin.Context) {
	db := Service.Db
	if !VerfiyAdmin(c, db) {
		return
	}
	var article AddArticle
	if CheckErr(c.Bind(&article), c) {
		return
	}

	Query := `INSERT INTO tmm.articles(fileid,author,title,link,source_url,
cover,published_at,digest,content,sortid,inserted_at) 
VALUES (%d,'%s','%s','%s','%s','%s','%s','%s','%s',%d,'%s')`
	_, _, err := db.Query(Query, article.Fileid, db.Escape(article.Author), db.Escape(article.Title),
		db.Escape(article.Link), db.Escape(article.SourceUrl), db.Escape(article.Cover),
		db.Escape(article.PublishedOn), db.Escape(article.Digest), db.Escape(article.Content),
		article.Sortid, db.Escape(time.Now().Format(`2006-01-02`)))
	if CheckErr(err, c) {
		return
	}
	rows, result, err := db.Query(`select id from tmm.articles where link = '%s' `, db.Escape(article.Link))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `select Error `, c) {
		return
	}
	link:=fmt.Sprintf("https://tmm.tokenmama.io/article/show/%d", rows[0].Int(result.Map(`id`)))
	Query = `INSERT INTO tmm.share_tasks (creator, title, summary, link, 
image, points, points_left, bonus, max_viewers) VALUES(0, '%s', '%s', '%s', '%s', 5000, 5000, 10, 10)`
	_, _, err = db.Query(Query, db.Escape(article.Title), db.Escape(article.Digest),
		db.Escape(link), db.Escape(article.Cover))
	if CheckErr(err, c) {
		return
	}
	rows, _, err = db.Query(`select id from tmm.share_tasks where link = '%s'`,db.Escape(link))
	if CheckErr(err, c) {
		return
	}
	if Check(len(rows) == 0, `select Error `, c) {
		return
	}
	_,_,err=db.Query(`INSERT INTO tmm.share_task_categories (task_id,cid,is_auto) VALUES (%d,%d,%d)`,rows[0].Int(0),article.Sortid,0)
	if CheckErr(err,c){
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": "" },
	)
}


func GetArticleHandler(c *gin.Context) {
	db := Service.Db
	if !VerfiyAdmin(c, db) {
		return
	}
	sortid, err := strconv.ParseUint(c.Param("sortid"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	var offset int
	page, _ := strconv.Atoi(c.DefaultQuery(`page`, "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery(`limit`, "10"))
	if page >= 1 {
		offset = (page - 1) * 10
	} else {
		offset = 0
	}
	inCidConstrain := fmt.Sprintf("INNER JOIN tmm.share_task_categories AS stc ON (stc.task_id=st.id AND stc.cid=%d)", sortid)

	Query := `
	SELECT
    st.id,
    st.title,
    st.summary,
    st.link,
    st.image,
    st.max_viewers,
    st.bonus,
    st.points,
    st.points_left,
    st.viewers,
    st.inserted_at,
    st.updated_at,
    st.creator,
    st.online_status
FROM tmm.share_tasks AS st
%s
ORDER BY st.updated_at DESC LIMIT %d OFFSET %d`

	Rows, Result, err := db.Query(Query, inCidConstrain, limit, offset)
	if CheckErr(err, c) {
		return
	}
	if len(Rows) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "没有到数据",
			"data": gin.H{
				"curr_page": page,
				"data":      nil,
			},
		})
		return
	}
	article := []common.ShareTask{}
	for _, Row := range Rows {
		bonus, _ := decimal.NewFromString(Row.Str(Result.Map(`bonus`)))
		points, _ := decimal.NewFromString(Row.Str(Result.Map(`points`)))
		pointsLeft, _ := decimal.NewFromString(Row.Str(Result.Map(`points_left`)))
		article = append(article, common.ShareTask{
			Id:           Row.Uint64(Result.Map(`id`)),
			Title:        Row.Str(Result.Map(`title`)),
			Summary:      Row.Str(Result.Map(`summary`)),
			Link:         Row.Str(Result.Map(`link`)),
			Image:        Row.Str(Result.Map(`image`)),
			MaxViewers:   Row.Uint(Result.Map(`max_viewers`)),
			Bonus:        bonus,
			Points:       points,
			PointsLeft:   pointsLeft,
			Viewers:      Row.Uint(Result.Map(`viewers`)),
			InsertedAt:   Row.ForceLocaltime(Result.Map(`inserted_at`)).Format(time.RFC3339),
			UpdatedAt:    Row.ForceLocaltime(Result.Map(`updated_at`)).Format(time.RFC3339),
			Creator:      Row.Uint64(Result.Map(`creator`)),
			OnlineStatus: int8(Row.Int(Result.Map(`online_status`))),
		})

	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "",
		"data": gin.H{
			"curr_page": page,
			"data":      article,
		},
	})
}

func DeleteArticleHandler(c *gin.Context) {
	db := Service.Db
	if !VerfiyAdmin(c, db) {
		return
	}
	articleId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if CheckErr(err, c) {
		return
	}
	link := fmt.Sprintf(`https://tmm.tokenmama.io/article/show/%d`, articleId)
	_, _, err = db.Query(`delete from tmm.share_tasks where link = '%s'`, link)
	if CheckErr(err, c) {
		return
	}
	_, _, err = db.Query(`delete from tmm.articles where id = %d`, articleId)
	if CheckErr(err, c) {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": "",
	})

}

func VerfiyAdmin(c *gin.Context, db *autorc.Conn) bool {
	userContext, exists := c.Get("USER")
	if Check(!exists, `Need login`, c) {
		return false
	}
	User := userContext.(common.User)
	Query := `select role from user_settings where user_id = %d`
	Rows, _, err := db.Query(Query, User.Id)
	if CheckErr(err, c) {
		return false
	}
	if len(Rows) == 0 {
		c.JSON(http.StatusOK, APIError{Msg: `没有查询到数据`})
		return false
	}
	row := Rows[0]
	if row.Int(0) != 1 {
		return false
	}
	return true
}
