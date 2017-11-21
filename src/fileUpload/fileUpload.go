//Package fileUpload 文件服务器
package fileUpload

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"utils"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type myhandler struct{}
type home struct {
	Title string
}

type fileStruct struct {
	Title  string
	IPAddr string
}

const (
	templateDir = "./fileUpload/view/"
	uploadDir   = "./fileUpload/upload/"
	cssDir      = "./fileUpload/css/"
)

//FileUpload 上传文件
func FileUpload() {
	server := http.Server{
		Addr:        ":9090",
		Handler:     &myhandler{},
		ReadTimeout: 10 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index
	mux["/upload"] = upload
	mux["/file"] = staticServer
	server.ListenAndServe()
}

func (*myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	if ok, _ := regexp.MatchString("/css/", r.URL.String()); ok {
		http.StripPrefix("/css/", http.FileServer(http.Dir(cssDir))).ServeHTTP(w, r)
	} else {
		http.StripPrefix("/", http.FileServer(http.Dir(uploadDir))).ServeHTTP(w, r)
	}

}

func upload(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		ipAddr := utils.IPV4()
		uploadFile := fileStruct{Title: "上传文件", IPAddr: "http://" + ipAddr + ":9090/upload"}
		t, _ := template.ParseFiles(templateDir + "file.html")
		t.Execute(w, uploadFile)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Fprintf(w, "%v", "上传错误")
			return
		}
		fileext := filepath.Ext(handler.Filename)
		if check(fileext) == false {
			fmt.Fprintf(w, "%v", "不允许的上传类型")
			return
		}
		filename := handler.Filename
		f, _ := os.OpenFile(uploadDir+filename, os.O_CREATE|os.O_WRONLY, 0777)
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Fprintf(w, "%v", "上传失败")
			return
		}
		filedir, _ := filepath.Abs(uploadDir + filename)
		fmt.Fprintf(w, "%v", filename+"上传完成,服务器地址:"+filedir)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	title := home{Title: "首页"}
	t, _ := template.ParseFiles(templateDir + "index.html")
	t.Execute(w, title)
}

func staticServer(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/file", http.FileServer(http.Dir(uploadDir))).ServeHTTP(w, r)
}

func check(name string) bool {
	ext := []string{".exe", ".js"}

	for _, v := range ext {
		if v == name {
			return false
		}
	}
	return true
}
