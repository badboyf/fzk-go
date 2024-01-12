package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
)

func main() {
	//leetcodeList()
	leetcodeQuestion(true)
	//leetcode()

	//saveContentAndAnswer(LCQuestion{Difficulty: "MEDIUN", FrontendQuestionId: "1029", TranslatedTitle: "两地调度", TitleSlug: "two-city-scheduling"}, false)

	//lcTest()
	//saveInFile(1, 100, &[]LCQuestion{
	//	{ Difficulty: "MEDIUN", FrontendQuestionId: "1", TranslatedTitle: "两数之和", TitleSlug: "two-sum"},
	//	{ Difficulty: "MEDIUN", FrontendQuestionId: "1", TranslatedTitle: "两数之和", TitleSlug: "two-sum"},
	//})
}
func lcTest()  {
}

/* 根据现有的数据，爬取题目和解析 */
func leetcodeQuestion(skipIfExist bool) {
	pth := "lc/file/"
	fileInfoList, err := ioutil.ReadDir(pth)
	if err != nil {
		fmt.Printf("fileInfoList  %v \n", err)
		return
	}
	sleepMilli := time.Duration(3) * time.Second
	for i := range fileInfoList {
		f := fileInfoList[i]
		if f.IsDir() || strings.HasSuffix(f.Name(), ".txt") {
			continue
		}
		completePath := pth + f.Name()
		content, err := ReadAll(completePath)
		if err != nil {
			fmt.Printf("fileInfoList  %v \n", err)
			return
		}

		var m []LCQuestion
		err = json.Unmarshal(content, &m)
		if err != nil {
			fmt.Printf("Unmarshal error  %v \n", err)
			return
		}

		for questioinIdx := range m {
			question := m[questioinIdx]
			fmt.Printf("save content and answer %v \n", question)
			rs := saveContentAndAnswer(question, skipIfExist)
			if rs != "skip" {
				time.Sleep(sleepMilli)
			}
		}
	}
}

/* 保存题的答案和内容 */
func saveContentAndAnswer(question LCQuestion, skipIfExist bool) string {
	result1 := saveContent(question, skipIfExist)
	result2 := saveAnswer(question, skipIfExist)
	if result1 == "skip" && result2 == "skip" {
		return "skip"
	} else if result2 == "fail" || result1 == "fail"{
		return "fail"
	}
	return "success"
}


/* 保存题的内容 */
func saveContent(question LCQuestion, skipIfExist bool) string {
	pthT := "lc/file/content/%v_%v_content.json"
	errPth := fmt.Sprintf("lc/file/content/error_%v_%v_content.json", question.FrontendQuestionId, question.TitleSlug)
	pth := fmt.Sprintf(pthT, question.FrontendQuestionId, question.TitleSlug)
	bodyT := "{\n    \"operationName\": \"questionData\",\n    \"variables\": {\n        \"titleSlug\": \"%v\"\n    },\n    \"query\": \"query questionData($titleSlug: String!) {\\n  question(titleSlug: $titleSlug) {\\n    questionId\\n    questionFrontendId\\n    categoryTitle\\n    boundTopicId\\n    title\\n    titleSlug\\n    content\\n    translatedTitle\\n    translatedContent\\n    isPaidOnly\\n    difficulty\\n    likes\\n    dislikes\\n    isLiked\\n    similarQuestions\\n    contributors {\\n      username\\n      profileUrl\\n      avatarUrl\\n      __typename\\n    }\\n    langToValidPlayground\\n    topicTags {\\n      name\\n      slug\\n      translatedName\\n      __typename\\n    }\\n    companyTagStats\\n    codeSnippets {\\n      lang\\n      langSlug\\n      code\\n      __typename\\n    }\\n    stats\\n    hints\\n    solution {\\n      id\\n      canSeeDetail\\n      __typename\\n    }\\n    status\\n    sampleTestCase\\n    metaData\\n    judgerAvailable\\n    judgeType\\n    mysqlSchemas\\n    enableRunCode\\n    envInfo\\n    book {\\n      id\\n      bookName\\n      pressName\\n      source\\n      shortDescription\\n      fullDescription\\n      bookImgUrl\\n      pressImgUrl\\n      productUrl\\n      __typename\\n    }\\n    isSubscribed\\n    isDailyQuestion\\n    dailyRecordStatus\\n    editorType\\n    ugcQuestionId\\n    style\\n    exampleTestcases\\n    __typename\\n  }\\n}\\n\"\n}"
	body := fmt.Sprintf(bodyT, question.TitleSlug)
	url := "https://leetcode-cn.com/graphql/"

	if skipIfExist && checkFileIsExist(pth) {
		DeleteFile(errPth)
		return "skip"
	}

	reqest, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		fmt.Printf("error  %v \n", err)
		saveInFile(errPth, err.Error())
		return "fail"
	}

	// 增加header
	AddLcReqHeader(reqest)

	//处理返回结果
	client := &http.Client{}
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Printf("client.Do  %v \n", err)
		saveInFile(errPth, err.Error())
		return "fail"
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	if err != nil {
		fmt.Println("json.Unmarshal error \n", err)
		saveInFile(errPth, err.Error())
		return "fail"
	}

	saveInFile(pth, m)
	DeleteFile(errPth);
	return "success"
}
/* 保存题的官方题解 */
func saveAnswer(question LCQuestion, skipIfExist bool) string {
	errPth := fmt.Sprintf("lc/file/content/error_%v_%v_answer.json", question.FrontendQuestionId, question.TitleSlug)
	pthT := "lc/file/content/%v_%v_answer.json"
	pth := fmt.Sprintf(pthT, question.FrontendQuestionId, question.TitleSlug)
	if skipIfExist && checkFileIsExist(pth) {
		DeleteFile(errPth)
		return "skip"
	}

	bodyT := "{\"operationName\":\"solutionDetailArticle\",\"variables\":{\"slug\":\"%v\",\"orderBy\":\"DEFAULT\"},\"query\":\"query solutionDetailArticle($slug: String!, $orderBy: SolutionArticleOrderBy!) {\\n  solutionArticle(slug: $slug, orderBy: $orderBy) {\\n    ...solutionArticle\\n    content\\n    question {\\n      questionTitleSlug\\n      __typename\\n    }\\n    position\\n    next {\\n      slug\\n      title\\n      __typename\\n    }\\n    prev {\\n      slug\\n      title\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\\nfragment solutionArticle on SolutionArticleNode {\\n  rewardEnabled\\n  canEditReward\\n  uuid\\n  title\\n  slug\\n  sunk\\n  chargeType\\n  status\\n  identifier\\n  canEdit\\n  canSee\\n  reactionType\\n  reactionsV2 {\\n    count\\n    reactionType\\n    __typename\\n  }\\n  tags {\\n    name\\n    nameTranslated\\n    slug\\n    tagType\\n    __typename\\n  }\\n  createdAt\\n  thumbnail\\n  author {\\n    username\\n    profile {\\n      userAvatar\\n      userSlug\\n      realName\\n      __typename\\n    }\\n    __typename\\n  }\\n  summary\\n  topic {\\n    id\\n    commentCount\\n    viewCount\\n    __typename\\n  }\\n  byLeetcode\\n  isMyFavorite\\n  isMostPopular\\n  isEditorsPick\\n  hitCount\\n  videosInfo {\\n    videoId\\n    coverUrl\\n    duration\\n    __typename\\n  }\\n  __typename\\n}\\n\"}"
	slug := getSlug(question)
	if slug == "" {
		saveInFile(errPth, "slug null, may vip")
		return "fail"
	}
	body := fmt.Sprintf(bodyT, slug)
	url := "https://leetcode-cn.com/graphql/"

	reqest, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		fmt.Printf("error  %v \n", err)
		saveInFile(errPth, err.Error())
		return "fail"
	}

	// 增加header
	AddLcReqHeader(reqest)

	//处理返回结果
	client := &http.Client{}
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Printf("client.Do  %v \n", err)
		saveInFile(errPth, err.Error())
		return "fail"
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	if err != nil {
		fmt.Println("json.Unmarshal error \n", err)
		saveInFile(errPth, err.Error())
		return "fail"
	}

	saveInFile(pth, m)
	DeleteFile(errPth);
	return "success"
}

/* 获取官方题解slug */
func getSlug(question LCQuestion) string {
	bodyT := "{\"operationName\":\"questionSolutionArticles\",\"variables\":{\"questionSlug\":\"%v\",\"first\":10,\"skip\":0,\"orderBy\":\"DEFAULT\"},\"query\":\"query questionSolutionArticles($questionSlug: String!, $skip: Int, $first: Int, $orderBy: SolutionArticleOrderBy, $userInput: String, $tagSlugs: [String!]) {\\n  questionSolutionArticles(questionSlug: $questionSlug, skip: $skip, first: $first, orderBy: $orderBy, userInput: $userInput, tagSlugs: $tagSlugs) {\\n    totalNum\\n    edges {\\n      node {\\n        ...solutionArticle\\n        __typename\\n      }\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\\nfragment solutionArticle on SolutionArticleNode {\\n  rewardEnabled\\n  canEditReward\\n  uuid\\n  title\\n  slug\\n  sunk\\n  chargeType\\n  status\\n  identifier\\n  canEdit\\n  canSee\\n  reactionType\\n  reactionsV2 {\\n    count\\n    reactionType\\n    __typename\\n  }\\n  tags {\\n    name\\n    nameTranslated\\n    slug\\n    tagType\\n    __typename\\n  }\\n  createdAt\\n  thumbnail\\n  author {\\n    username\\n    profile {\\n      userAvatar\\n      userSlug\\n      realName\\n      __typename\\n    }\\n    __typename\\n  }\\n  summary\\n  topic {\\n    id\\n    commentCount\\n    viewCount\\n    __typename\\n  }\\n  byLeetcode\\n  isMyFavorite\\n  isMostPopular\\n  isEditorsPick\\n  hitCount\\n  videosInfo {\\n    videoId\\n    coverUrl\\n    duration\\n    __typename\\n  }\\n  __typename\\n}\\n\"}"
	body := fmt.Sprintf(bodyT, question.TitleSlug)
	url := "https://leetcode-cn.com/graphql/"

	reqest, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return ""
	}

	// 增加header
	AddLcReqHeader(reqest)

	//处理返回结果
	client := &http.Client{}
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Printf("client.Do  %v \n", err)
		return ""
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	if err != nil {
		fmt.Println("json.Unmarshal error \n", err)
		return ""
	}

	data := m["data"].(map[string]interface{})
	questionSolutionArticles := data["questionSolutionArticles"].(map[string]interface{})
	edges := questionSolutionArticles["edges"].([]interface{})
	if edges == nil || len(edges) == 0 {
		return ""
	}
	slugs := edges[0].(map[string]interface{})
	node := slugs["node"].(map[string]interface{})
	slug := node["slug"].(string)
	return slug
}

func leetcodeList() {
	pageEnd := 20
	pageStart := 1
	pageSize := 100
	for page := pageStart; page <= pageEnd; page++ {
		url := "https://leetcode-cn.com/graphql/"
		//url := "http://localhost:7020/open/inner/test"
		datas := doLeetcodeList(url, page, pageSize)

		fileName := fmt.Sprintf("lc/file/page-%v.txt", page)
		saveInFile(fileName, datas)
	}
}

func doLeetcodeList(url string, page int, pageSize int) *[]LCQuestion {
	client := &http.Client{}
	bodyT := "{\"query\":\"\\n    query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {\\n  problemsetQuestionList(\\n    categorySlug: $categorySlug\\n    limit: $limit\\n    skip: $skip\\n    filters: $filters\\n  ) {\\n    hasMore\\n    total\\n    questions {\\n      acRate\\n      difficulty\\n      freqBar\\n      frontendQuestionId\\n      isFavor\\n      paidOnly\\n      solutionNum\\n      status\\n      title\\n      titleCn\\n      titleSlug\\n      topicTags {\\n        name\\n        nameTranslated\\n        id\\n        slug\\n      }\\n      extra {\\n        hasVideoSolution\\n        topCompanyTags {\\n          imgUrl\\n          slug\\n          numSubscribed\\n        }\\n      }\\n    }\\n  }\\n}\\n    \",\"variables\":{\"categorySlug\":\"algorithms\",\"skip\":%v,\"limit\":%v,\"filters\":{}},\"operationName\":\"problemsetQuestionList\"}"
	body := fmt.Sprintf(bodyT, strconv.Itoa((page - 1) * pageSize), pageSize)
	reqest, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return nil
	}

	//增加header选项
	AddLcReqHeader(reqest)

	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Printf("client.Do  %v \n", err)
		return nil
	}
	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("read error", err)
		return nil
	}
	//bodyStr := string(resBody)
	//fmt.Printf("%v \n", bodyStr)
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	if err != nil {
		fmt.Println("json.Unmarshal error \n", err)
		return nil
	}

	data := m["data"].(map[string]interface{})
	problemsetQuestionList := data["problemsetQuestionList"].(map[string]interface{})
	questions := problemsetQuestionList["questions"].([]interface{})

	qs := []LCQuestion{}
	for i := range questions {
		question := questions[i].(map[string]interface{})
		difficulty := question["difficulty"].(string)
		frontendQuestionId := question["frontendQuestionId"].(string)
		titleCn := question["titleCn"].(string)
		titleSlug := question["titleSlug"].(string)

		q := LCQuestion {
			Difficulty:         difficulty,
			FrontendQuestionId: frontendQuestionId,
			TranslatedTitle:    titleCn,
			TitleSlug:          titleSlug,
		}
		qs = append(qs, q)
		fmt.Printf("%v  %v  %v  %v \n", frontendQuestionId, titleCn, difficulty, titleSlug)
	}
	return &qs
}

func AddLcReqHeader(reqest *http.Request)  {
	//增加header选项
	reqest.Header.Add("Cookie", "__appToken__=; gr_user_id=5b4b3b8f-85b2-4315-9138-ae8c4a5dd2df; aliyungf_tc=f5259218469e4a9a0a693aa4a23d2110044dc942eae44ba5d29f23975d96a7f2; NEW_PROBLEMLIST_PAGE=1; csrftoken=wRHmbzyxXIENx0FzMNBNLv7k4R0zP5mOriCxr0A0ehBH1p2wAolBqW8vgRUDTkG5; a2873925c34ecbd2_gr_last_sent_cs1=elo9uent-agnesi; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMTY2MDIzNCIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImRqYW5nby5jb250cmliLmF1dGguYmFja2VuZHMuTW9kZWxCYWNrZW5kIiwiX2F1dGhfdXNlcl9oYXNoIjoiODQ3MzcxODMzZDg3ZDc2ZGExNzFlZjNhMWE0OWVkNjJmMjNiMjcwNTAwZWY1N2U0ZTM3N2ZjY2VmNDE5NjVmMSIsImlkIjoxNjYwMjM0LCJlbWFpbCI6IiIsInVzZXJuYW1lIjoiZWxvOXVlbnQtYWduZXNpIiwidXNlcl9zbHVnIjoiZWxvOXVlbnQtYWduZXNpIiwiYXZhdGFyIjoiaHR0cHM6Ly9hc3NldHMubGVldGNvZGUtY24uY29tL2FsaXl1bi1sYy11cGxvYWQvZGVmYXVsdF9hdmF0YXIucG5nIiwicGhvbmVfdmVyaWZpZWQiOnRydWUsIl90aW1lc3RhbXAiOjE2MzY5ODEwNzcuMDI0ODcwNCwiZXhwaXJlZF90aW1lXyI6MTYzOTUwODQwMCwibGF0ZXN0X3RpbWVzdGFtcF8iOjE2Mzc5MjM5OTd9")
	reqest.Header.Add("Authority", "leetcode-cn.com")
	reqest.Header.Add("sec-ch-u", "\"Google Chrome\";v=\"95\", \"Chromium\";v=\"95\", \";Not A Brand\";v=\"99\"")
	reqest.Header.Add("Authorization", "" )
	reqest.Header.Add("x-csrftoken", "wRHmbzyxXIENx0FzMNBNLv7k4R0zP5mOriCxr0A0ehBH1p2wAolBqW8vgRUDTkG5")
	reqest.Header.Add("sec-ch-ua-mobile", "?0")
	reqest.Header.Add("Content-Type", "application/json")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36")
	reqest.Header.Add("sec-ch-ua-platform", "macOS")
	reqest.Header.Add("Accept","*/*")
	reqest.Header.Add("Origin", "https://leetcode-cn.com")
	reqest.Header.Add("sec-fetch-site","same-origin")
	reqest.Header.Add("sec-fetch-mode", "cors")
	reqest.Header.Add("sec-fetch-dest", "empty")
	reqest.Header.Add("Referer", "https://leetcode-cn.com/problemset/algorithms/?page=1")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")

}

func leetcode() {
	times := 2
	offset := 0
	for i := 1; i < times; i++ {
		page := strconv.Itoa(offset + i)
		url := "https://leetcode-cn.com/problemset/algorithms/?page=" + page
		doCrawLC(url)
		//createArticleInDb(datas)
	}
}

func doCrawLC(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}
	req.Header.Add("Cookie", "__appToken__=; gr_user_id=5b4b3b8f-85b2-4315-9138-ae8c4a5dd2df; aliyungf_tc=f5259218469e4a9a0a693aa4a23d2110044dc942eae44ba5d29f23975d96a7f2; NEW_PROBLEMLIST_PAGE=1; csrftoken=wRHmbzyxXIENx0FzMNBNLv7k4R0zP5mOriCxr0A0ehBH1p2wAolBqW8vgRUDTkG5; a2873925c34ecbd2_gr_session_id=4c435f12-7d8c-4542-a627-d7634bd760a6; a2873925c34ecbd2_gr_session_id_4c435f12-7d8c-4542-a627-d7634bd760a6=true; a2873925c34ecbd2_gr_last_sent_sid_with_cs1=4c435f12-7d8c-4542-a627-d7634bd760a6; a2873925c34ecbd2_gr_last_sent_cs1=elo9uent-agnesi; a2873925c34ecbd2_gr_cs1=elo9uent-agnesi; LEETCODE_SESSION=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJfYXV0aF91c2VyX2lkIjoiMTY2MDIzNCIsIl9hdXRoX3VzZXJfYmFja2VuZCI6ImRqYW5nby5jb250cmliLmF1dGguYmFja2VuZHMuTW9kZWxCYWNrZW5kIiwiX2F1dGhfdXNlcl9oYXNoIjoiODQ3MzcxODMzZDg3ZDc2ZGExNzFlZjNhMWE0OWVkNjJmMjNiMjcwNTAwZWY1N2U0ZTM3N2ZjY2VmNDE5NjVmMSIsImlkIjoxNjYwMjM0LCJlbWFpbCI6IiIsInVzZXJuYW1lIjoiZWxvOXVlbnQtYWduZXNpIiwidXNlcl9zbHVnIjoiZWxvOXVlbnQtYWduZXNpIiwiYXZhdGFyIjoiaHR0cHM6Ly9hc3NldHMubGVldGNvZGUtY24uY29tL2FsaXl1bi1sYy11cGxvYWQvZGVmYXVsdF9hdmF0YXIucG5nIiwicGhvbmVfdmVyaWZpZWQiOnRydWUsIl90aW1lc3RhbXAiOjE2MzY5ODEwNzcuMDI0ODcwNCwiZXhwaXJlZF90aW1lXyI6MTYzOTUwODQwMCwibGF0ZXN0X3RpbWVzdGFtcF8iOjE2Mzc5MDUxOTJ9.zy40VUvUzBA97yo3EQTuqyqGRy0lZ5G29PS1cv7OyQA")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("lc NewDocumentFromReader error", err)
		return
	}

	//doc.Find("div[role=row][class=odd:bg-overlay-3 dark:odd:bg-dark-overlay-1 even:bg-overlay-1 dark:even:bg-dark-overlay-3]").Each(func(i int, line *goquery.Selection) {
	doc.Find("div[class='odd:bg-overlay-3 dark:odd:bg-dark-overlay-1 even:bg-overlay-1 dark:even:bg-dark-overlay-3']").Each(func(i int, line *goquery.Selection) {
		line.Find("div[role=cell]:div:nth-child(2)").Each(func(i int, name *goquery.Selection) {
			pName := name.Find("a").Text()
			href, _ := name.Find("a").Attr("href")
			fmt.Printf("%v %v \n", pName, href)
		})

		line.Find("div[role=cell]:div:nth-child(5)").Each(func(i int, diffc *goquery.Selection) {
			difficulty := diffc.Find("span").Text()
			fmt.Printf("%v \n", difficulty)
		})
	})

}

func DeleteFile(fileName string)  {
	if !checkFileIsExist(fileName) {
		return
	}
	err := os.Remove(fileName)
	if err != nil {
		fmt.Println("DeleteFile err:", err)
	}
}
func saveInFile(fileName string, qestions interface{})  {
	abs, _ := filepath.Abs(fileName)
	fmt.Printf("write to file %v \n", abs)

	var dstFile *os.File
	var err error
	if !checkFileIsExist(fileName) {
		dstFile, err = os.Create(fileName)
	} else {
		dstFile, err = os.OpenFile(fileName, os.O_RDWR, 0666)
	}
	if err != nil {
		fmt.Println("file err:", err)
		return
	}
	defer dstFile.Close()

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	err = jsonEncoder.Encode(qestions)
	if err != nil {
		fmt.Println("jsonEncoder.Encode err:", err);return
	}
	//fmt.Println(bf.String())

	_, err = io.WriteString(dstFile, bf.String())
	if err != nil {
		fmt.Println("write err:", err)
		return
	}

}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

type LCQuestion struct {
	Difficulty         string `json:"difficulty"`
	FrontendQuestionId string `json:"questionFrontendId"`
	TranslatedTitle    string `json:"translatedTitle"`
	TitleSlug          string `json:"titleSlug"`
}