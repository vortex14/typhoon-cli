package v1_1

import (
	"github.com/gobuffalo/packr"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Import struct {
	Old string
	New string
}

func MigrateComponents() {
	if _, err := os.Stat("project"); os.IsNotExist(err) {
		_ = os.Mkdir("project", 755)
	}

	if _, err := os.Stat("fetcher"); !os.IsNotExist(err) {
		_ = os.Rename("fetcher", "project/fetcher")
	}

	if _, err := os.Stat("processor"); !os.IsNotExist(err) {
		_ = os.Rename("processor", "project/processor")
	}

	if _, err := os.Stat("donor"); !os.IsNotExist(err) {
		_ = os.Rename("donor", "project/donor")
	}

	if _, err := os.Stat("result_transporter"); !os.IsNotExist(err) {
		_ = os.Rename("result_transporter", "project/result_transporter")
	}

	if _, err := os.Stat("scheduler"); !os.IsNotExist(err) {
		_ = os.Rename("scheduler", "project/scheduler")
	}
}

func GetComponentTemplate() (error error, data string)  {
	box := packr.NewBox(".")

	data, err := box.FindString("component.gopy")

	if err != nil {
		log.Fatal(err)
	}

	return err, data
}

func GetConfigTemplate() (error error, data string) {
	box := packr.NewBox(".")

	data, err := box.FindString("config.goyaml")

	if err != nil {
		log.Fatal(err)
	}

	return err, data
}


func VisitAndReplace(path string, fi os.FileInfo, err error) error {

	var newImports [] *Import

	newImports = append(newImports, &Import{
		Old: "from executions.strategies",
		New: "from typhoon.components.fetcher.executions.strategies",
	},
	&Import{
		Old: "from responses.strategies",
		New: "from typhoon.components.fetcher.responses.strategies",
	},
	&Import{
		Old: "from executions.strategies.base_pre_fetch",
		New: "from typhoon.components.fetcher.executions.strategies.base_pre_fetch",
	},
	&Import{
		Old: "from executable.pipeline_group",
		New: "from typhoon.components.processor.executable.pipeline_group",
	},
	&Import{
		Old: "from executable.base_handler",
		New: "from typhoon.components.processor.executable.base_handler",
	},
	&Import{
		Old: "from project.executable",
		New: "from project.processor.executable",
	},
	&Import{
		Old: "from project.resource",
		New: "from project.processor.resource",
	},
	&Import{
		Old: "from executable.text_pipelines.base_pipeline import BasePipeline",
		New: "from typhoon.components.processor.executable.text_pipelines.base_pipeline import BasePipeline",
	},
	&Import{
		Old: "from project.consumers",
		New: "from project.result_transporter.consumers",
	},
	&Import{
		Old: "from executions.base_consumer",
		New: "from typhoon.components.result_transporter.executions.base_consumer",
	},
	&Import{
		Old: "from project.result_transporter.consumers import exceptions",
		New: "from typhoon.extensions.result_transporter import exceptions",
	},
	&Import{
		Old: "from project.v1",
		New: "from project.donor.v1",
	},
	)

	if err != nil {
		return err
	}

	if !!fi.IsDir() {
		return nil //
	}

	if strings.Contains(path, ".pyc") || strings.Contains(path, ".json"){
		return nil
	}

	//if !strings.Contains(path, "first_mongo_callback.py") {
	//	return nil
	//}

	data, err := ioutil.ReadFile(path)
	newContents := string(data)
	if err != nil {
		panic(err)
	}

	for _, v := range newImports {
		newContents = strings.Replace(newContents, v.Old, v.New, -1)
	}
	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}

	return nil
}