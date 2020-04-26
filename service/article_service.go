package service

import (
	"encoding/json"
	"feelingliu/modles"
	"feelingliu/tools"
	"feelingliu/utils"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

type Article struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required,max=32"`
	Content     string `json:"content" db:"content" binding:"required"`
	Html        string `json:"html" db:"html" binding:"required"`
	TagID       int    `json:"tag_id" binding:"required"`
	CreatedTime string `json:"created_time" db:"created_time"`
	UpdatedTime string `json:"updated_time" db:"updated_time"`
	Status      string `json:"status" db:"status" binding:"required"`
}

type Articles struct {
	Items []Article `json:"items"`
	Total int       `json:"total"`
}

type ArticleDetail struct {
	A     Article `json:"article"`
	Tags  []Tag   `json:"tags"`
	Views int     `json:"views"`
}

type Options struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Search bool   `json:"search"`
	Admin  bool   `json:"admin"`
	C      string `json:"c"` // category
	T      string `json:"t"` // tag
	Q      string `json:"q"` // 搜索的关键字
}

type Option func(*Options)

var defaultOptions = Options{
	Limit:  10,
	Page:   1,
	C:      "",
	T:      "",
	Q:      "",
	Search: false, // 搜索文章结果不进行缓存
	Admin:  false, // 是否是admin页面请求，如果不是，文章就不包括草稿文章
}

func newOptions(opts ...Option) Options {
	// 初始化默认值
	opt := defaultOptions

	for _, o := range opts {
		o(&opt) // 依次调用opts函数列表中的函数，为服务选项（opt变量）赋值
	}

	return opt
}

func GetArticlesByTag(opts ...Option) (data Articles, err error) {
	options := newOptions(opts...)
	//baseSql := "SELECT %s FROM article a  INNER JOIN blog_tag_article ta ON a.id=ta.article_id INNER JOIN tag t ON ta.tag_id=t.id AND t.tag_name=" + "'" + options.T + "'" + ""
	baseSql := "SELECT %s FROM article a INNER JOIN tag t ON a.tag_id=t.id AND t.tag_name=" + "'" + options.T + "'" + ""
	data, err = genArticles(baseSql, opts...)
	return
}

func genArticles(baseSql string, opts ...Option) (data Articles, err error) {
	options := newOptions(opts...)
	key := articleCacheKey(options)
	if !options.Search {
		cacheData, e := getArticleCache(key)
		fmt.Println(cacheData)
		if e != redis.Nil {
			utils.WriteErrorLog(fmt.Sprintf("[ %s ] 读取缓存失败, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), e))
		}
		if cacheData.Total != 0 {
			return cacheData, nil
		}
	}

	articles := make([]Article, 0)

	var f string
	if !options.Admin {
		f = " WHERE a.status='published'"
	}
	offset := (options.Page - 1) * options.Limit
	selectSql := fmt.Sprintf(baseSql, "a.id, a.title, a.created_time, a.updated_time, a.status") + f + fmt.Sprintf(" ORDER BY a.id DESC limit %d offset %d", options.Limit, offset)
	if db := modles.DB.Raw(selectSql).Scan(&articles); db.Error != nil {
		return
	}

	//var total int
	//if findtotal := modles.DB.Model(Article{}).Where("status = ?, ", "published").Count(&total); findtotal.Error != nil {
	//	return
	//}
	//fmt.Println("total : ", total)
	//data.Total = total
	data.Total = len(articles)
	data.Items = articles

	if !options.Search {
		if e := setArticleCache(key, data); e != nil {
			utils.WriteErrorLog(fmt.Sprintf("[ %s ] 写入缓存失败, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), e))
		}
	}

	return
}

func articleCacheKey(opts Options) string {
	if opts.Admin {
		return fmt.Sprintf("article_%d_%d_%s_%s_%s", opts.Limit, opts.Page, "admin", opts.C, opts.T)
	} else {
		return fmt.Sprintf("article_%d_%d_%s_%s", opts.Limit, opts.Page, opts.C, opts.T)
	}
}

func getArticleCache(key string) (a Articles, err error) {
	data, e := tools.GetKey(key)
	if e != nil {
		fmt.Println(err)
		return a, e
	}

	if e := json.Unmarshal([]byte(data), &a); e != nil {
		return a, e
	}
	return a, nil
}

func SetLimitPage(limit, page string) Option {
	return func(o *Options) {
		if limit != "" && page != "" {
			p, _ := strconv.Atoi(page)
			l, _ := strconv.Atoi(limit)
			o.Limit = l
			o.Page = p
		}
	}
}

func SetAdmin(admin string) Option {
	return func(o *Options) {
		if admin != "" {
			o.Admin = true
		}
	}
}

func SetTag(t string) Option {
	return func(o *Options) {
		o.T = t
	}
}

func SearchArticle(key, status string, opts ...Option) (data Articles, err error) {
	var baseSql string
	if status == "" {
		baseSql = `SELECT %s FROM article a WHERE a.title LIKE '%%` + key + `%%'`
	} else {
		baseSql = `SELECT %s FROM article a WHERE a.title LIKE '%%` + key + `%%' AND a.status='` + status + `'`
	}

	data, err = genArticles(baseSql, opts...)
	return
}

//func SearchFromES(opts ...Option) (articles Articles, e error) {
//
//	const searchMatch = `{"query" : {
//    "multi_match": {
//      "fields":  [ "content", "title" ],
//      "query":     "%s",
//      "fuzziness": "AUTO"
//    }
//} }`
//	var (
//		r     map[string]interface{}
//		items []Article
//		total int
//	)
//
//	options := newOptions(opts...)
//	offset := (options.Page - 1) * options.Limit
//
//	res, err := es.Search(
//		es.Search.WithContext(context.Background()),
//		es.Search.WithIndex(utils.ESInfo.Index),
//		es.Search.WithBody(strings.NewReader(fmt.Sprintf(searchMatch, options.Q))),
//		es.Search.WithTrackTotalHits(true),
//		es.Search.WithSize(options.Limit),
//		es.Search.WithFrom(offset),
//	)
//	if err != nil {
//		return articles, fmt.Errorf("error getting response: %s", err)
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		var e map[string]interface{}
//		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
//			return articles, fmt.Errorf("error parsing the response body: %s", err)
//		} else {
//			return articles, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
//		}
//	}
//
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		return articles, fmt.Errorf("error parsing the response body: %s", err)
//	}
//
//	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
//		source, _ := json.Marshal(hit.(map[string]interface{})["_source"])
//		var a Article
//		if err := json.Unmarshal(source, &a); err != nil {
//			return articles, fmt.Errorf("error parsing the response body: %s", err)
//		}
//		items = append(items, a)
//	}
//	total = int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
//
//	articles.Items = items
//	articles.Total = total
//	return articles, nil
//}

func (a Article) GetAll(opts ...Option) (data Articles, err error) {
	baseSql := "select %s from article a"
	data, err = genArticles(baseSql, opts...)
	return
}

func SetSearch(search bool) Option {
	return func(o *Options) {
		o.Search = search
		o.Page = defaultOptions.Page // 如果不是在第一页执行的搜索，比如：page=3，有可能会搜不到数据，必须从第一页开始搜索
	}
}

func setArticleCache(key string, value Articles) error {
	//fmt.Println(value.Items)
	//fmt.Println(value.Total)
	marshal, _ := json.Marshal(value)
	//fmt.Println(string(marshal))
	//fmt.Println(key)
	e := tools.SetKey(key, marshal, tools.SetTimeout(true))
	return e
}

func (a Article) GetOne(opts ...Option) (ArticleDetail, error) {
	options := newOptions(opts...)
	var one Article
	if db := modles.DB.Where("id = ?", a.ID).Find(&one); db.Error != nil {
		return ArticleDetail{}, db.Error
	}

	tags, _ := GetTagsByArticleID(a.ID)

	viewKey := one.ViewKey()
	n, err := getViews(viewKey)
	fmt.Println("n", n)
	if err != nil {
		utils.WriteErrorLog(fmt.Sprintf("[ %s ] 获取阅读量失败, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), err))
	}

	if !options.Admin {
		if e := addView(viewKey); e != nil {
			utils.WriteErrorLog(fmt.Sprintf("[ %s ] 添加阅读量失败, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), e))
		}
	}
	return ArticleDetail{one, tags, n}, nil
}

func addView(key string) error {
	e := tools.INCRKey(key)
	return e
}

func GetTagsByArticleID(articleID int) ([]Tag, error) {
	var t []Tag
	sql := "SELECT t.* FROM tag t RIGHT JOIN blog_tag_article ta ON t.id=ta.tag_id WHERE ta.article_id=" + "'" + strconv.Itoa(articleID) + "'" + ""
	if db := modles.DB.Raw(sql).Scan(&t); db.Error != nil {
		return nil, db.Error
	}
	return t, nil
}

func (a Article) ViewKey() string {
	viewKey := a.Title + ":view"
	return viewKey
}

func getViews(key string) (n int, err error) {
	data, e := tools.GetKey(key)
	if e != nil {
		return n, e
	}

	if e := json.Unmarshal([]byte(data), &n); e != nil {
		return n, e
	}
	return n, nil
}

func (a *Article) Create() (Article, error) {
	createTime := time.Now().Format(modles.AppInfo.TimeFormat)

	var article Article = Article{
		Title:       a.Title,
		Content:     a.Content,
		Html:        a.Html,
		TagID:       a.TagID,
		CreatedTime: createTime,
		Status:      a.Status,
	}
	db := modles.DB.Create(&article)
	if db.Error != nil {
		return Article{}, db.Error
	}

	//if article.Status == "published" {
	//	//	if e := article.IndexBlog(); e != nil {
	//	//		utils.WriteErrorLog(fmt.Sprintf("[ %s ] 存入elastic出错, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), e))
	//	//	}
	//	//}
	return article, nil
}

//func (a Article) IndexBlog() error {
//	req := esapi.IndexRequest{
//		Index:      utils.ESInfo.Index,
//		DocumentID: strconv.Itoa(a.ID),
//		Body:       esutil.NewJSONReader(a),
//		Refresh:    "true",
//	}
//
//	res, err := req.Do(context.Background(), es)
//	if err != nil {
//		return err
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		var e map[string]interface{}
//		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
//			return err
//		}
//		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
//	}
//	return nil
//}

func (a *Article) Delete() error {
	var article = Article{ID: a.ID}
	db := modles.DB.Delete(&article)
	if db.Error != nil {
		return db.Error
	}
	// 删除阅读量
	viewKey := a.ViewKey()
	if e := tools.DelKey(viewKey); e != nil {
		utils.WriteErrorLog(fmt.Sprintf("[ %s ] 删除阅读量失败, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), e))
	}
	// 从ES中删除
	//if e := a.DeleteFromES(); e != nil {
	//	utils.WriteErrorLog(fmt.Sprintf("[ %s ] 从elastic中删除出错, %v\n", time.Now().Format(utils.AppInfo.TimeFormat), e))
	//}
	return nil
}

func (a *Article) Edit() error {
	updateTime := time.Now().Format(modles.AppInfo.TimeFormat)
	var newarticle = Article{}
	modles.DB.Where("id = ?", a.ID).Find(&newarticle)

	newarticle.ID = a.ID
	newarticle.Title = a.Title
	newarticle.Content = a.Content
	newarticle.Html = a.Html
	newarticle.TagID = a.TagID
	newarticle.UpdatedTime = updateTime
	newarticle.Status = a.Status
	save := modles.DB.Save(&newarticle)
	if save.Error != nil {
		return save.Error
	}

	//if len(a.TagID) > 0 {
	//	for _, tagID := range a.TagID {
	//		_, e := db.Exec("insert into blog_tag_article (tag_id, article_id) values (?, ?)", tagID, a.ID)
	//		if e != nil {
	//			return e
	//		}
	//	}
	//}
	//if a.Status == "published" {
	//	if e := a.IndexBlog(); e != nil {
	//		utils.WriteErrorLog(fmt.Sprintf("[ %s ] 从elastic更新出错, %v\n", time.Now().Format(utils.AppInfo.TimeFormat), e))
	//	}
	//} else {
	//	if e := a.DeleteFromES(); e != nil {
	//		utils.WriteErrorLog(fmt.Sprintf("[ %s ] 从elastic删除出错, %v\n", time.Now().Format(utils.AppInfo.TimeFormat), e))
	//	}
	//}
	return nil
}
