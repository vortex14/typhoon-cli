package gitlab

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/mitchellh/mapstructure"
	"github.com/xanzy/go-gitlab"
	"net/url"
	"os"
	"strings"
	"typhoon-cli/src/interfaces"
)

type Server struct {
	Cluster interfaces.Cluster
}

func (s *Server) getGitlabProjects(gitlabClient *gitlab.Client, page int) []*gitlab.Project {
	projects, _, _ := gitlabClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
		ListOptions:              gitlab.ListOptions{
			PerPage: 100,
			Page: page,
		},
	})
	return projects

}

func (s *Server) GetAllProjectsList() []*interfaces.GitlabProject {
	settings := s.Cluster.GetEnvSettings()
	color.Green("Sync gitlab projects. waiting for %s", settings.Gitlab)
	var scrapedProjects []*interfaces.GitlabProject
	gitlabClient, _ := s.GetClient()

	count := 10
	bar := pb.StartNew(count)
	bar.SetMaxWidth(100)
	for i := 1; i <= count; i++ {

		description := fmt.Sprintf("scan gitlab page: %d", i)

		tmpl := `{{string . "title"}} - {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}}  {{percent .}} {{etime .}}`

		bar.SetTemplateString(tmpl)

		bar.Set("title", description)

		projects := s.getGitlabProjects(gitlabClient, i)
		for _, project := range projects {
			scrapedProjects = append(scrapedProjects, &interfaces.GitlabProject{
				Name: project.Name,
				Git: project.WebURL + ".git",
				Id: project.ID,
			})
		}

		bar.Increment()
	}
	bar.Finish()

	return scrapedProjects
}


func (s *Server) SyncGitlabProjects()  {

	scrapedAllProjects := s.GetAllProjectsList()
	clusterProjects := s.Cluster.GetProjects()

	foundCount := 0

	for _, project := range clusterProjects {
		for _, gitLabProject := range scrapedAllProjects {
			if gitLabProject.Git == project.Git {
				foundCount += 1
				project.GitlabId = gitLabProject.Id
			}
		}

	}
	s.Cluster.SaveConfig()
	color.Green("A total of %d projects were found on gitlab. Found %d out of %d projects for this cluster", len(scrapedAllProjects), foundCount, len(clusterProjects))
}

func (s *Server) GetPipelineHistory(client *gitlab.Client, GitlabId int)  {
	pipelines, _, _ := client.Pipelines.ListProjectPipelines(GitlabId, &gitlab.ListProjectPipelinesOptions{
		ListOptions:   gitlab.ListOptions{},
	}, func(request *retryablehttp.Request) error {
		return nil
	})


	for _, pipeline := range pipelines {
		color.Green("%s", pipeline.String())
	}
}

func (s *Server) GetClient() (*gitlab.Client, error) {
	settings := s.Cluster.GetEnvSettings()
	gitlabClient, err := gitlab.NewClient(settings.GitlabToken, gitlab.WithBaseURL(settings.Gitlab))
	return gitlabClient, err
}

func pathEscape(s string) string {
	return strings.Replace(url.PathEscape(s), ".", "%2E", -1)
}

func (s *Server) GetVariables() []*gitlab.PipelineVariable {
	meta := s.Cluster.GetMeta()

	var variables []*gitlab.PipelineVariable

	for m := range meta {
		row := meta[m]

		if m == "variables" {

			err := mapstructure.Decode(row, &variables)
			if err != nil {
				color.Red("%s", err)
				return nil
			}

		}
	}

	for _, variable := range variables {
		variable.VariableType = "file"
	}

	return variables
}

func (s *Server) HistoryPipelines()  {
	//s.GetPipelineHistory(gitlabClient, project.GitlabId)
}

func (s *Server) Deploy()  {

	gitlabClient, _ := s.GetClient()

	clusterProjects := s.Cluster.GetProjects()
	countGitlabIds := 0
	for _, project := range clusterProjects {
		if project.GitlabId > 0 {
			countGitlabIds += 1
		}
	}
	if countGitlabIds == 0 {
		color.Red("Not found Gitlab ids into %s cluster %s", s.Cluster.GetName(), s.Cluster.GetClusterConfigPath())
		os.Exit(1)
	}

	variables := s.GetVariables()

	for _, project := range clusterProjects {
		color.Green("branch: %s , %s id: %d",project.Branch, project.Name, project.GitlabId)

		pipeline, response, errorPipeline := gitlabClient.Pipelines.CreatePipeline(project.GitlabId, &gitlab.CreatePipelineOptions{
			Ref:       &project.Branch,
			Variables: variables,
		}, func(request *retryablehttp.Request) error {
			return nil
		})

		if errorPipeline != nil {
			color.Red("%s", errorPipeline)
			color.Red("%+v", response)
			os.Exit(1)
		}
		color.Yellow("created pipeline: %s %s for %s", project.Name, pipeline.WebURL, project.Git)

		//return
	}
}
