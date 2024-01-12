package business

import (
	"bufio"
	"fmt"
	"fzk-test-to/util"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const file_list_path = "/Users/fengzhikui/self/image/二次女-去重.txt"
const img_path = "/Users/fengzhikui/self/image/二次女10/"
const db_path = "/Users/fengzhikui/data/fzknotebook/fzk-test-go/zfile/tu_img/db.txt"

func TLProcess() {
	content := util.File.ReadFilePanic(file_list_path)
	db := util.DBStr.Init(db_path)
	ids := db.Content

	split := strings.Split(content, "\n")
	for i, item := range split {
		if util.Util.Contains(item, ids) {
			continue
		}
		if strings.TrimSpace(item) == "" {
			continue
		}
		downloadPath := strings.ReplaceAll(item, "thumb.300_0.", "")
		err := DoTLDownloadImg(downloadPath)
		fmt.Printf("\n%v download %s result=%v\n", i, item, err)
		if err != nil {
			continue
		}

		err = db.StoreAndFlushOne(item)
		if err != nil {
			fmt.Printf("StoreAndFlushOne %s\n", err)
		} else {
			ids = append(ids, item)
		}

		//time.Sleep(time.Second * 1)
	}
	fmt.Printf("%s", split[0])
}

func DoTLDownloadImg(imgUrl string) error {
	fileName := path.Base(imgUrl)

	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred!")
		return err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(img_path + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d %s", written, img_path + fileName)

	return nil
}

type TLContent struct {
	Id           string `json:"id"`
	DownloadPath string `json:"download_path"`
	Success      bool   `json:"success"`
}
