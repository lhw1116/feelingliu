package service

import (
	"encoding/json"
	"errors"
	"feelingliu/modles"
	"feelingliu/tools"
	"feelingliu/utils"
	"fmt"
	"strconv"
	"time"
)

type Article struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required,max=32"`
	Content     string `json:"content" db:"content" binding:"required"`
	Html        string `json:"html" db:"html" binding:"required"`
	CategoryID  int    `json:"category_id" db:"category_id" binding:"required"`
	TagID       []int  `json:"tag_id" binding:"required"`
	CreatedTime string `json:"created_time" db:"created_time"`
	UpdatedTime string `json:"updated_time" db:"updated_time"`
	Status      string `json:"status" db:"status" binding:"required"`
}

type Articles struct {
	Items []Article `json:"items"`
	Total int       `json:"total"`
}

type ArticleDetail struct {
	A     Article  `json:"article"`
	Tags  []Tag    `json:"tags"`
	Views uint8    `json:"views"`
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
	baseSql := "SELECT %s FROM blog_article a  INNER JOIN blog_tag_article ta ON a.id=ta.article_id INNER JOIN blog_tag t ON ta.tag_id=t.id AND t.tag_name=" + "'" + options.T + "'" + ""
	data, err = genArticles(baseSql, opts...)
	return
}

func genArticles(baseSql string, opts ...Option) (data Articles, err error) {
	options := newOptions(opts...)
	key := articleCacheKey(options)
	if !options.Search {
		cacheData, e := getArticleCache(key)
		if e != nil {
			utils.WriteErrorLog(fmt.Sprintf("[ %s ] 读取缓存失败, %v\n", time.Now().Format(modles.AppInfo.TimeFormat), e))
		}
		if cacheData.Total != 0 {
			return cacheData, nil
		}
	}

	articles := make([]Article, 0)

	var f string
	if !options.Admin {
		f = " WHERE a.`status`='published'"
	}
	offset := (options.Page - 1) * options.Limit
	selectSql := fmt.Sprintf(baseSql, "a.id, a.title, a.created_time, a.updated_time, a.`status`") + f + fmt.Sprintf(" ORDER BY a.id DESC limit %d offset %d", options.Limit, offset)

	if err = db.Select(&articles, selectSql); err != nil {
		return
	}

	var total int
	if err = db.Get(&total, fmt.Sprintf(baseSql, "count(1)")+f); err != nil {
		return
	}

	data.Total = total
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
	if e != nil || data == nil {
		return a, e
	}

	v, ok := data.([]uint8)
	if ok {
		if e := json.Unmarshal([]byte(v[:]), &a); e != nil {
			return a, e
		}
		return a, nil
	} else {
		return a, errors.New("返回数据类型有误，json无法解析")
	}
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
		baseSql = `SELECT %s FROM blog_article a WHERE a.title LIKE '%%` + key + `%%'`
	} else {
		baseSql = `SELECT %s FROM blog_article a WHERE a.title LIKE '%%` + key + `%%' AND a.status='` + status + `'`
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
	baseSql := "select %s from blog_article a"
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
	marshal, _ := json.Marshal(value)
	e := tools.SetKey(key, marshal, tools.SetTimeout(true))
	return e
}

func (a Article) GetOne(opts ...Option) (ArticleDetail, error) {
	options := newOptions(opts...)
	var one Article
	if e := db.Get(&one, "select * from blog_article where id=?", a.ID); e != nil {
		return ArticleDetail{}, e
	}
	tags, _ := GetTagsByArticleID(a.ID)

	viewKey := one.ViewKey()
	n, err := getViews(viewKey)
	if err != nil {
		utils.WriteErrorLog(fmt.Sprintf("[ %s ] 获取阅读量失败, %v\n", time.Now().Format(utils.AppInfo.TimeFormat), err))
	}

	if !options.Admin {
		if e := addView(viewKey); e != nil {
			utils.WriteErrorLog(fmt.Sprintf("[ %s ] 添加阅读量失败, %v\n", time.Now().Format(utils.AppInfo.TimeFormat), e))
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
	if e := db.Select(&t, "SELECT t.* FROM blog_tag t RIGHT JOIN blog_tag_article ta ON t.id=ta.tag_id WHERE ta.article_id=?", articleID); e != nil {
		return nil, e
	}
	return t, nil
}

func (a Article) ViewKey() string {
	viewKey := a.Title + ":view"
	return viewKey
}

func getViews(key string) (n uint8, err error) {
	data, e := tools.GetKey(key)
	if e != nil || data == nil {
		return n, e
	}

	v, ok := data.([]uint8)
	if ok {
		if e := json.Unmarshal([]byte(v[:]), &n); e != nil {
			return n, e
		}
		return n, nil
	} else {
		return n, errors.New("返回数据类型有误，json无法解析")
	}
}