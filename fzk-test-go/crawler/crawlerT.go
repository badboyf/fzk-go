package main

import (
	"log"
	"strings"

	"fmt"

	"regexp"

	"context"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const DB = "golangcc"
const author = "zhikui.feng"
const TABLE_ARTICLE = "article"
const TABLE_ARTICLE_CRAWL = "articleCrawler"
const MONGO_URL = "mongodb://192.168.50.157:27017"

var client *mongo.Client

func main() {
	client = getConn()
	defer close(client)

	//zh()
	//crawlSF()

	zhDetail()
}

func zhDetail() {
	url := "https://zhuanlan.zhihu.com/p/187505981"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("sf NewDocumentFromReader error", err)
		return
	}

	html, e := doc.Find(".Post-RichText").Html()
	if e != nil {
		fmt.Printf("get html error")
		return
	}
	//fmt.Printf("html================================================ \n%s \n", html)

	md := Convert(html)
	fmt.Printf("md================================================ \n%s \n", md)
}

func crawlSF() {
	end := 3
	page := 2
	for i := page; i < end; i++ {
		url := "https://segmentfault.com/t/golang/blogs?page=" + strconv.Itoa(page)
		doCrawSF(url)
	}
}

func doCrawSF(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}

	//body, err := ioutil.ReadAll(resp.Body)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("sf NewDocumentFromReader error", err)
		return
	}

	doc.Find("section").Each(func(i int, s *goquery.Selection) {
		id, exist := s.Find("li[data-id]").Attr("data-id")
		if !exist {
			fmt.Printf("sf not exist")
			return
		}
		desc := s.Find("p").Text()
		title, content, success := doCrawSFById(id)
		if !success {
			fmt.Printf("doCrawSFById false %v \n", id)
			return
		}
		fmt.Printf("create %v \n", id)
		doCreateArticleInDb(&CreateArticleVo{ID: "sf_" + id, Title: title, Html: content, Desc: desc})
	})
}

func doCrawSFById(id string) (string, string, bool) {
	url := "https://segmentfault.com/a/" + id
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error  %v \n", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("NewDocumentFromReader error", err)
		return "", "", false
	}

	title := ""
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})

	content := ""
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		content, err = s.Html()
		if err != nil {
			fmt.Printf("article error %v", err)
			return
		}
	})
	return title, content, true
}

func zh() {
	times := 1
	offset := 10
	for i := 0; i < times; i++ {
		page := strconv.Itoa(offset*i) + ".0000"
		url := "https://www.zhihu.com/api/v4/topics/19625982/feeds/top_activity?include=data%5B%3F%28target.type%3Dtopic_sticky_module%29%5D.target.data%5B%3F%28target.type%3Danswer%29%5D.target.content%2Crelationship.is_authorized%2Cis_author%2Cvoting%2Cis_thanked%2Cis_nothelp%3Bdata%5B%3F%28target.type%3Dtopic_sticky_module%29%5D.target.data%5B%3F%28target.type%3Danswer%29%5D.target.is_normal%2Ccomment_count%2Cvoteup_count%2Ccontent%2Crelevant_info%2Cexcerpt.author.badge%5B%3F%28type%3Dbest_answerer%29%5D.topics%3Bdata%5B%3F%28target.type%3Dtopic_sticky_module%29%5D.target.data%5B%3F%28target.type%3Darticle%29%5D.target.content%2Cvoteup_count%2Ccomment_count%2Cvoting%2Cauthor.badge%5B%3F%28type%3Dbest_answerer%29%5D.topics%3Bdata%5B%3F%28target.type%3Dtopic_sticky_module%29%5D.target.data%5B%3F%28target.type%3Dpeople%29%5D.target.answer_count%2Carticles_count%2Cgender%2Cfollower_count%2Cis_followed%2Cis_following%2Cbadge%5B%3F%28type%3Dbest_answerer%29%5D.topics%3Bdata%5B%3F%28target.type%3Danswer%29%5D.target.annotation_detail%2Ccontent%2Chermes_label%2Cis_labeled%2Crelationship.is_authorized%2Cis_author%2Cvoting%2Cis_thanked%2Cis_nothelp%2Canswer_type%3Bdata%5B%3F%28target.type%3Danswer%29%5D.target.author.badge%5B%3F%28type%3Dbest_answerer%29%5D.topics%3Bdata%5B%3F%28target.type%3Danswer%29%5D.target.paid_info%3Bdata%5B%3F%28target.type%3Darticle%29%5D.target.annotation_detail%2Ccontent%2Chermes_label%2Cis_labeled%2Cauthor.badge%5B%3F%28type%3Dbest_answerer%29%5D.topics%3Bdata%5B%3F%28target.type%3Dquestion%29%5D.target.annotation_detail%2Ccomment_count%3B&limit=10&after_id=" + page
		datas := doCrawZh(url)
		createArticleInDb(datas)
	}
}

func doCreateArticleInDb(v *CreateArticleVo) {
	articles := make([]interface{}, 0, 1)
	articleIds := make([]string, 0, 1)
	newArticle := Article{
		ID:         v.ID,
		Author:     author,
		Tags:       []ArticleTag{},
		Categories: []ArticleTag{},
		Title:      v.Title,
		Html:       v.Html,
		Desc:       v.Desc,
		IsTop:      false,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	articleIds = append(articleIds, newArticle.ID)
	articles = append(articles, newArticle)
	deleteByIds(&articleIds, TABLE_ARTICLE, client)
	insertBatch(&articles, TABLE_ARTICLE, client)
}

func createArticleInDb(datas *[]ZHData) {
	var articleCrawlers = make([]interface{}, 0, 1)
	var articleCrawlerIds = make([]string, 0, 1)
	for _, v := range *datas {
		articleCrawlers = append(articleCrawlers, v.Target)
		articleCrawlerIds = append(articleCrawlerIds, strconv.Itoa(v.Target.Id))
	}
	deleteByIds(&articleCrawlerIds, TABLE_ARTICLE_CRAWL, client)
	insertBatch(&articleCrawlers, TABLE_ARTICLE_CRAWL, client)

	articles := make([]interface{}, 0, 1)
	articleIds := make([]string, 0, 1)
	for _, v := range *datas {
		newArticle := Article{
			ID:         "zh_" + strconv.Itoa(v.Target.Id),
			Author:     author,
			Tags:       []ArticleTag{},
			Categories: []ArticleTag{},
			Title:      v.Target.Title,
			Html:       v.Target.Content,
			Desc:       v.Target.Excerpt,
			IsTop:      false,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		articleIds = append(articleIds, newArticle.ID)
		articles = append(articles, newArticle)
	}
	deleteByIds(&articleIds, TABLE_ARTICLE, client)
	insertBatch(&articles, TABLE_ARTICLE, client)
}

func doCrawZh(url string) *[]ZHData {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error  %v \n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error", err)
		return nil
	}
	//fmt.Println(string(body))
	model := &ZHModel{}
	_ = json.Unmarshal(body, model)
	for _, v := range model.Data {
		fmt.Println(v.Target.Title)
		//fmt.Println(v.Target.Content)
	}
	return &model.Data
}

func deleteByIds(ids *[]string, table string, client *mongo.Client) {
	d := client.Database(DB)
	var c = d.Collection(table)
	filter := bson.M{"_id": bson.M{"$in": ids}}
	many, err := c.DeleteMany(context.TODO(), filter)
	if err != nil {
		fmt.Printf("insert error: %v \n", err)
	}
	fmt.Printf("删除结果：%v \n", many)
}
func insertBatch(v interface{}, table string, client *mongo.Client) {
	var realV = v.(*[]interface{})
	d := client.Database(DB)
	var c = d.Collection(table)
	many, err := c.InsertMany(context.TODO(), *realV)
	if err != nil {
		fmt.Printf("insert error: %v \n", err)
	}
	fmt.Printf("插入结果：%v \n", many)
}

func getConn() *mongo.Client {
	param := fmt.Sprintf(MONGO_URL)
	clientOptions := options.Client().ApplyURI(param)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	return client
}

func close(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

type ZHModel struct {
	Data    []ZHData      `json:"data"`
	Paging  ZhModelPaging `json:"paging"`
	Version string        `json:"version"`
}

type ZhModelPaging struct {
	IsEnd    bool   `json:"is_end"`
	IsStart  bool   `json:"is_start"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Totals   int    `json:"totals"`
}

type ZHData struct {
	TargetDescription string       `json:"target_description"`
	Type              string       `json:"type"`
	Target            ZHDataTarget `json:"target"`
}

type ZHDataTarget struct {
	Author          ZHDataTargetAuthor `json:"author"`
	CommentCount    int                `json:"comment_count"`
	Content         string             `json:"content"`
	Created         int                `json:"created"`
	Excerpt         string             `json:"excerpt"`
	ExcerptTitle    string             `json:"excerpt_title"`
	Id              int                `json:"id" bson:"_id"`
	ImageUrl        string             `json:"image_url"`
	IsLabeled       bool               `json:"is_labeled"`
	Title           string             `json:"title"`
	TopicThumbnails []string           `json:"topic_thumbnails"`
	Type            string             `json:"type"`
	Updated         int                `json:"updated"`
	Url             string             `json:"url"`
	VoteupCount     int                `json:"voteup_count"`
}
type ZHDataTargetAuthor struct {
	AvatarUrl         string                 `json:"avatar_url"`
	AvatarUrlTemplate string                 `json:"avatar_url_template"`
	Badge             []interface{}          `json:"badge"`
	BadgeV2           map[string]interface{} `json:"badge_v2"`
	EduMemberTag      map[string]interface{} `json:"edu_member_tag"`
	Gender            int                    `json:"gender"`
	Headline          string                 `json:"headline"`
	Id                int                    `json:"id"`
	IsAdvertiser      bool                   `json:"is_advertiser"`
	IsOrg             bool                   `json:"is_org"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type"`
	Url               string                 `json:"url"`
	UrlToken          string                 `json:"url_token"`
	UserType          string                 `json:"user_type"`
}

type Article struct {
	ID           string       `json:"id" bson:"_id,omitempty"`
	Author       string       `json:"author"`
	Tags         []ArticleTag `json:"tags"`       // 用户自定义分类
	Categories   []ArticleTag `json:"categories"` // 系统分类
	Title        string       `json:"title"`
	Desc         string       `json:"desc"`
	Html         string       `json:"html"`
	Markdown     string       `json:"markdown"`
	IsTop        bool         `json:"isTop" bson:"isTop"`                 // 置顶
	LookCount    int          `json:"lookCount" bson:"lookCount"`         // 浏览数
	LikeCount    int          `json:"likeCount" bson:"likeCount"`         // 点赞数
	CollectCount int          `json:"collect_count" bson:"collect_count"` // 收藏数
	CommentCount int          `json:"comment_count" bson:"comment_count"` // 评论数
	PublishTime  *time.Time   `json:"publish_time" bson:"publish_time"`
	Status       string       `json:"status"` // published, not_published
	CreateTime   time.Time    `json:"createTime" bson:"createTime"`
	UpdateTime   time.Time    `json:"updateTime" bson:"updateTime"`
}
type ArticleTag struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type CreateArticleVo struct {
	ID     string
	Author string
	Title  string
	Html   string
	Desc   string
}

var closeTag = map[string]string{
	"del":      "~~",
	"b":        "**",
	"strong":   "**",
	"i":        "_",
	"em":       "_",
	"dfn":      "_",
	"var":      "_",
	"cite":     "_",
	"br":       "\n",
	"strotyle": "",
}

var blockTag = []string{
	"address", "div", "figure", "p", "figcaption",
	"article", "aside", "nav", "footer", "fieldset", "menu",
	"header", "section", "center", "frameset", "details", "summary",
}

var nextlineTag = []string{
	"pre", "blockquote", "table",
}

//convert html to markdown
//将html转成markdown
func Convert(htmlstr string) (md string) {
	var maps map[string]string
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(htmlstr))
	doc = trimAttr(doc)
	doc, maps = compress(doc)
	doc = handleSpan(doc)      //<span>
	doc = handleNextLine(doc)  //<div>...
	doc = handleBlockTag(doc)  //<div>...
	doc = handleA(doc)         //<a>
	doc = handleImg(doc)       //<img>
	doc = handleHead(doc)      //h1~h6
	doc = handleClosedTag(doc) //<strong>、<i>、eg..
	doc = handleHr(doc)        //<hr>
	doc = handleLi(doc)        //<li>
	md, _ = doc.Find("body").Html()
	md = depress(md, maps)
	return
}

// 解压，释放code和pre
func depress(md string, maps map[string]string) string {
	// 先替换pre，再替换code，因为有的code在pre标签里面
	for key, val := range maps {
		if strings.HasPrefix(key, "{$blockquote") {
			md = strings.Replace(md, key, "\n\r"+val+"\n\r", -1)
		}
	}

	for key, val := range maps {
		if strings.HasPrefix(key, "{$pre") {
			md = strings.Replace(md, key, "\n\r"+val+"\n\r", -1)
		}
	}

	for key, val := range maps {
		if strings.HasPrefix(key, "{$code") || strings.HasPrefix(key, "{$textarea") {
			md = strings.Replace(md, key, val, -1)
		}
	}

	if doc, err := goquery.NewDocumentFromReader(strings.NewReader(md)); err == nil {
		doc = trimAttr(doc)
		md, _ = doc.Find("body").Html()
		md = strings.Replace(md, "<span>", "", -1)
		md = strings.Replace(md, "</span>", "", -1)
	}

	md = strings.ReplaceAll(md, " = ", "=")

	return md
}

// trip attr
func trimAttr(doc *goquery.Document) *goquery.Document {
	attrs := []string{
		"border" /*"colspan", "rowspan",*/, "style", "cellspacing",
		"cellpadding", "bgcolor", "width", "align", "frame", "id", "class",
	}
	elements := []string{
		"table", "thead", "tbody", "tr", "td", "th", "h1", "h2", "h3", "h4", "img",
		"h5", "h6", "i", "em", "strong", "span", "br", "hr", "ul", "li", "ol",
	}
	elements = append(elements, blockTag...)
	elements = append(elements, nextlineTag...)
	for _, tag := range elements {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			for _, attr := range attrs {
				selection.RemoveAttr(attr)
			}
		})
	}
	return doc
}

//压缩html
func compress(doc *goquery.Document) (*goquery.Document, map[string]string) {
	//blockquote、pre、code，并替换 span 为空

	var maps = make(map[string]string)

	if ele := doc.Find("textarea"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$textarea%v}", i)
			cont := "<textarea>" + getInnerHtml(selection) + "</textarea>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}

	if ele := doc.Find("code"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$code%v}", i)
			cont := "<code>" + getInnerHtml(selection) + "</code>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}

	if ele := doc.Find("pre"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$pre%v}", i)
			cont := "<pre>" + getInnerHtml(selection) + "</pre>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}

	if ele := doc.Find("blockquote"); len(ele.Nodes) > 0 {
		ele.Each(func(i int, selection *goquery.Selection) {
			key := fmt.Sprintf("{$blockquote%v}", i)
			cont := "<blockquote>" + getInnerHtml(selection) + "</blockquote>"
			selection.BeforeHtml(key)
			selection.Remove()
			maps[key] = cont
		})
	}

	htmlstr, _ := doc.Html()
	htmlstr = strings.Replace(htmlstr, "=", " = ", -1)
	htmlstr = strings.Replace(htmlstr, "\n", "", -1)
	htmlstr = strings.Replace(htmlstr, "\r", "", -1)
	htmlstr = strings.Replace(htmlstr, "\t", "", -1)
	htmlstr = strings.Replace(htmlstr, "<dl", "<ul", -1)
	htmlstr = strings.Replace(htmlstr, "</dl", "</ul", -1)
	htmlstr = strings.Replace(htmlstr, "<dt", "<li", -1)
	htmlstr = strings.Replace(htmlstr, "</dt", "</li", -1)
	htmlstr = strings.Replace(htmlstr, "<dd", "<li", -1)
	htmlstr = strings.Replace(htmlstr, "</dd", "</li", -1)
	//正则匹配，把“>”和“<”直接的空格全部去掉
	//去除标签之间的空格，如果是存在代码预览的页面，不要替换空格，否则预览的代码会错乱
	r, _ := regexp.Compile(">\\s{1,}<")
	htmlstr = r.ReplaceAllString(htmlstr, "> <")
	//多个空格替换成一个空格
	r2, _ := regexp.Compile("\\s{1,}")
	htmlstr = r2.ReplaceAllString(htmlstr, " ")
	doc, _ = goquery.NewDocumentFromReader(strings.NewReader(htmlstr))
	return doc, maps
}

func handleBlockTag(doc *goquery.Document) *goquery.Document {
	for _, tag := range blockTag {
		hasTag := true
		for hasTag {
			if tagEle := doc.Find(tag); len(tagEle.Nodes) > 0 {
				tagEle.Each(func(i int, selection *goquery.Selection) {
					selection.BeforeHtml("\n" + getInnerHtml(selection) + "\n")
					selection.Remove()
				})
			} else {
				hasTag = false
			}
		}
	}
	return doc
}

//func handleBlockquote(doc *goquery.Document) *goquery.Document {
//	if tagEle := doc.Find("blockquote"); len(tagEle.Nodes) > 0 {
//		tagEle.Each(func(i int, selection *goquery.Selection) {
//			cont := getInnerHtml(selection)
//			cont = strings.Replace(cont, "\r", "", -1)
//			cont = strings.Replace(cont, "\n", "", -1)
//			selection.BeforeHtml("\r\n<blockquote>" + cont + "\n</blockquote>\n")
//			selection.Remove()
//		})
//	}
//
//	doc.Find("code").Each(func(i int, selection *goquery.Selection) {
//		fmt.Println(selection.Html())
//	})
//
//	return doc
//}

//[ok]handle tag <a>
func handleA(doc *goquery.Document) *goquery.Document {
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		if href, ok := selection.Attr("href"); ok {
			if cont, err := selection.Html(); err == nil {
				md := fmt.Sprintf(`[%v](%v)`, cont, href)
				selection.BeforeHtml(md)
				selection.Remove()
			}
		}
	})
	return doc
}

func handleSpan(doc *goquery.Document) *goquery.Document {
	doc.Find("span").Each(func(i int, selection *goquery.Selection) {
		//if href, ok := selection.Attr("href"); ok {
		txt := selection.Text()
		md := fmt.Sprintf(`%v`, txt)
		selection.BeforeHtml(md)
		selection.Remove()
		//}
	})
	return doc
}

//[ok]handle tag ul、ol、li
//处理步骤：
//1、先给每个li标签里面的内容加上"- "或者"\t- "
//2、提取li内容
func handleLi(doc *goquery.Document) *goquery.Document {
	var tags = []string{"ol", "ul", "li"}
	doc.Find("li").Each(func(i int, selection *goquery.Selection) {
		l := len(selection.ParentsFiltered("li").Nodes)
		tab := strings.Join(make([]string, l+2), "{$@$space}")
		selection.PrependHtml("\r$@$" + tab)
	})
	for _, tag := range tags {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			if tag == "ul" {
				selection.BeforeHtml("\n" + selection.Text() + "\n")
			} else {
				selection.BeforeHtml(selection.Text())
			}
			selection.Remove()
		})
	}
	htmlstr, _ := doc.Find("body").Html()
	for i := 10; i > 0; i-- {
		oldTab := "$@$" + strings.Join(make([]string, i), "{$@$space}")
		newTab := strings.Join(make([]string, i-1), "  ") + "- "
		htmlstr = strings.Replace(htmlstr, oldTab, newTab, -1)
	}
	doc, _ = goquery.NewDocumentFromReader(strings.NewReader(htmlstr))
	return doc
}

//[ok]handle tag <hr/>
func handleHr(doc *goquery.Document) *goquery.Document {
	doc.Find("hr").Each(func(i int, selection *goquery.Selection) {
		selection.BeforeHtml("\n- - -\n")
		selection.Remove()
	})
	return doc
}

//[ok]handle tag <img/>
func handleImg(doc *goquery.Document) *goquery.Document {
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		if src, ok := selection.Attr("src"); ok {
			alt := ""
			if val, ok := selection.Attr("alt"); ok {
				alt = val
			}
			md := fmt.Sprintf(`![%v](%v)`, alt, src)
			selection.BeforeHtml(md)
			selection.Remove()
		}
	})
	return doc
}

//[ok]handle tag h1~h6
func handleHead(doc *goquery.Document) *goquery.Document {
	heads := map[string]string{
		"title": "# ",
		"h1":    "# ",
		"h2":    "## ",
		"h3":    "### ",
		"h4":    "#### ",
		"h5":    "##### ",
		"h6":    "###### ",
	}
	for tag, replace := range heads {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			text, _ := selection.Html()
			selection.BeforeHtml("\n\r" + replace + text + "\n\r")
			selection.Remove()
		})
	}
	return doc
}

func handleClosedTag(doc *goquery.Document) *goquery.Document {
	for tag, close := range closeTag {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			text, e := selection.Html()
			if nil != e {
				log.Println(e)
			}
			if strings.TrimSpace(text) != "" {
				selection.BeforeHtml(close + text + close)
			}

			if tag == "br" {
				selection.BeforeHtml(close + close)
			}

			selection.Remove()
		})
	}
	return doc
}

func handleNextLine(doc *goquery.Document) *goquery.Document {
	for _, tag := range nextlineTag {
		doc.Find(tag).Each(func(i int, selection *goquery.Selection) {
			selection.BeforeHtml("\n\n")
			selection.AfterHtml("\n\n")
		})
	}
	return doc
}

func getInnerHtml(selection *goquery.Selection) (html string) {
	html, _ = selection.Html()
	return
}
