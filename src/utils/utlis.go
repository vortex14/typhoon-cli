package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	"path/filepath"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"typhoon-cli/src/interfaces"
)

type Utils struct {}

func (u *Utils) GoRunTemplate(goTemplate *interfaces.GoTemplate) bool {
	tmpl, _ := template.New("new").Parse(goTemplate.Source)
	status := true
	f, err := os.Create(goTemplate.ExportPath)
	if err != nil {
		log.Println("create file: ", err)
		status = false
	}

	err = tmpl.Execute(f, &goTemplate.Data)
	if err != nil {
		log.Print("execute: ", err)
		status = false
	}
	f.Close()

	return status
}

func (u *Utils) ParseLog(object *interfaces.FileObject) error {
	currentPath, _ := os.Getwd()
	logPath := fmt.Sprintf("%s/%s", currentPath, object.Path)
	dat, err := ioutil.ReadFile(logPath)

	color.Red("Log file path: %s", logPath)
	if err != nil {

		color.Red("Log file not found")
		os.Exit(0)


	}

	logDataMap := logfmt.NewDecoder(strings.NewReader(string(dat)))
	for logDataMap.ScanRecord() {
		for logDataMap.ScanKeyval() {
			color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
		}
	}



	return nil
}

func (u *Utils) CopyDir(name string, object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	errC := os.MkdirAll(name, 0755)
	if errC != nil{
		color.Red("%s",errC)
	}
	box.Walk(func(s string, file packd.File) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(s)
		//color.Yellow("data : %s, info: %s",s, pathDir)
		if pathDir != "." {
			errD := os.Mkdir(name+"/"+pathDir, 0755)
			if errD != nil{
				//color.Red("%s",errD)
			}

		}

		f, err := os.Create(name+"/"+s)
		if err != nil {
			log.Println("create err", err)
		}
		f.WriteString(file.String())

		return nil
	})
	return nil
}

func (u *Utils) CopyDirAndReplaceLabel(name string, label *interfaces.ReplaceLabel ,object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	errC := os.MkdirAll(name, 0755)
	if errC != nil{
		color.Red("%s",errC)
	}
	box.Walk(func(s string, file packd.File) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(s)
		//color.Yellow("data : %s, info: %s",s, pathDir)
		if pathDir != "." {
			errD := os.Mkdir(name+"/"+pathDir, 0755)
			if errD != nil{
				//color.Red("%s",errD)
			}

		}

		f, err := os.Create(name+"/"+s)
		if err != nil {
			log.Println("create err", err)
		}
		data := file.String()
		data = strings.ReplaceAll(data, label.Label, label.Value)
		f.WriteString(data)

		return nil
	})
	return nil
}

func (u *Utils) GetGoTemplate(object *interfaces.FileObject) (error error, data string)  {
	box := packr.NewBox(object.Path)

	data, err := box.FindString(object.Name)

	if err != nil {
		log.Fatal(err)
	}

	return err, data
}

func (u *Utils) CheckSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}