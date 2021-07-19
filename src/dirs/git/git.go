package git

import (
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"os"
	"typhoon-cli/src/environment"
)

type Git struct {
	settings *environment.Settings
	repo *git.Repository
	Path string
}

func (g *Git) TestGetDir(projectName string) error {
	if _, err := os.Stat(g.settings.Projects + "/" + projectName + "/.git"); os.IsNotExist(err) {
		return err
	}
	return nil
}

func (g *Git) Commit()  {
	g.LoadRepo()

}

func (g *Git) Push()  {
	g.LoadRepo()
}

func (g Git) AddAll()  {
	g.LoadRepo()

}

func (g *Git) AddAndCommit(remote string)  {
	g.LoadRepo()

}

func (g *Git) LoadRepo() {
	if g.repo == nil {
		repo, err := git.PlainOpen(g.Path)
		if err != nil {
			color.Red("%s", err.Error())
			os.Exit(1)
		}
		g.repo = repo
	}
}

