package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"
	"github.com/gobuffalo/packr"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"typhoon-cli/interfaces"
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

func (u *Utils) GetGoTemplate(object *interfaces.FileObject) (error error, data string)  {
	box := packr.NewBox(object.Path)

	data, err := box.FindString(object.Name)

	if err != nil {
		log.Fatal(err)
	}

	return err, data
}