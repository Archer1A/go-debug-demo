package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func sayHelloName(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm() // 解析参数，默认是不会解析的
	fmt.Println(r.Form) // 这些信息是输出入服务器的打印信息
	fmt.Println("path " ,r.URL.Path)
	fmt.Println("scheme " , r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k,v := range r.Form{
		fmt.Println("key:" ,k)
		fmt.Println("value",strings.Join(v,""))
	}
	fmt.Fprintf(w,"hello vic")

}

func login(w http.ResponseWriter,r *http.Request)  {
	fmt.Println("method " ,r.Method)
	if r.Method == "GET" {
		t,_ := template.ParseFiles("login.html")
		log.Println(t.Execute(w, nil))
	} else{
		r.ParseForm() // 如果注释这句不会有输出 默认情况下，Handler里面是不会自动解析form的，必须显式的调用r.ParseForm()后，你才能对这个表单数据进行操作
		//请求的是登录数据，那么执行登录的逻辑判断
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}



}
func upload(w http.ResponseWriter,r *http.Request)  {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h,strconv.FormatInt(crutime,10))
		token := fmt.Sprintf("/x",h.Sum(nil))

		t,_ := template.ParseFiles("upload.html")
		fmt.Println(token)
		t.Execute(w,nil)

	} else{
		r.ParseMultipartForm(32 << 20)
		file,handler,err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

	}

}

func main()  {
	http.HandleFunc("/v2",sayHelloName)//设置访问路由
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9090",nil)// 设置监听端口
	if err != nil {
		log.Fatal("ListenAndServe: " ,err)
	}

}

