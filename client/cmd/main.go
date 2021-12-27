package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	addr     string
	port     string
	filename string
	foo      func()
	URL      string
)

type File struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
}

func main() {
	InitConfig()
	addr = viper.GetString("addr")
	port = viper.GetString("port")
	ParseArgs(os.Args)
	URL = "http://" + addr + ":" + port
	foo()
}

func GetFileList() {
	var file []File
	req := URL + "/file/all"
	resp, err := http.Get(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal(body, &file)
	for idx, elem := range file {
		fmt.Printf("%d:\n", idx+1)
		fmt.Printf(" filename: %s\n", elem.Name)
		fmt.Printf(" checksum: %s\n", elem.Checksum)
	}
}

func GetFile() {
	req := URL + "/file/" + filename
	resp, err := http.Get(req)
	if err != nil {
		fmt.Println("NotOK")
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	name := resp.Header.Get("filename")
	body, err := ioutil.ReadAll(resp.Body)
	file, err := os.Create(name)
	if err != nil {
		fmt.Println("NotOK")
		log.Fatalln(err)
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		fmt.Println("NotOK")
		log.Fatalln(err)
	}
	fmt.Println("OK")
}

func UploadFile() {
	req := URL + "/file/new"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	part, err := writer.CreateFormFile("filename", filepath.Base(filename))
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = writer.Close()
	if err != nil {
		log.Fatalln(err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, req, payload)
	if err != nil {
		log.Fatalln(err)
		return
	}

	request.Header.Add("Content-Type", "multipart/form-data")
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("OK")
	} else {
		fmt.Println("NotOK")
	}
}

func UpdateFile() {
	req := URL + "/file/update"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("filename", filepath.Base(filename))
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, req, payload)
	if err != nil {
		log.Fatalln(err)
		return
	}

	request.Header.Add("Content-Type", "multipart/form-data")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("OK")
	} else {
		fmt.Println("NotOK")
	}
}

func DeleteFile() {
	req := URL + "/file/delete/" + filename

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, req, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("OK")
	} else {
		fmt.Println("NotOK")
	}
}

func PrintDoc() {
	docFile, _ := os.Open("doc.txt")
	defer docFile.Close()
	stat, _ := docFile.Stat()
	buffer := make([]byte, stat.Size())
	docFile.Read(buffer)
	fmt.Printf("%s", buffer)
}

func InitConfig() {
	viper.SetConfigFile("conf.yml")
	viper.ReadInConfig()
}

func ParseArgs(args []string) {
	if len(args) < 2 {
		log.Fatalln("action not specified")
	}
	if len(args) == 2 {
		switch args[1] {

		case "getall":
			foo = GetFileList
			return
		case "help":
			foo = PrintDoc
			return
		default:
			foo = PrintDoc
			log.Println("unknown command or no file specified")
			return
		}
	}
	switch args[1] {
	case "getall":
		foo = GetFileList
	case "get":
		foo = GetFile
		filename = args[2]
	case "update":
		foo = UpdateFile
		filename = args[2]
	case "upload":
		foo = UploadFile
		filename = args[2]
	case "delete":
		foo = DeleteFile
		filename = args[2]
	case "help":
		foo = PrintDoc
		return
	default:
		fmt.Println("unknown command")
		foo = PrintDoc
	}

	for idx, arg := range args {
		if idx < 3 {
			continue
		}
		switch args[idx-1] {
		case "addr":
			addr = arg
		case "port":
			port = arg
		}
	}
}
