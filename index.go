package main

import (
	"context"
	"embed"
	"github.com/google/go-github/v47/github"
	"html/template"
	"log"
	"net/http"
)

type IndexPageData struct {
	TokenAvailable bool
	OrgsAvailable  []GitHubOrg
	OrgSelected    string
	Octocat        string
}

type GitHubOrg struct {
	Name        string
	Description string
	URL         string
}

//go:embed templates/*
var templates embed.FS

func getRoot(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFS(templates, "templates/*")

	client := getClient()
	ctx := context.Background()
	octo, _, _ := client.Octocat(ctx, "OpenISMS")

	data := IndexPageData{
		TokenAvailable: len(token) > 0,
		OrgsAvailable:  ListGitHubOrgs(client, ctx),
		OrgSelected:    org,
		Octocat:        octo,
	}

	err := tmpl.ExecuteTemplate(w, "index.gohtml", data)
	if err != nil {
		log.Fatal(err)
	}
}

func ListGitHubOrgs(client *github.Client, ctx context.Context) []GitHubOrg {

	myOrgs, _, _ := client.Organizations.List(ctx, "", nil)

	var res []GitHubOrg

	for _, org := range myOrgs {
		res = append(res, GitHubOrg{
			Name:        *org.Login,
			Description: *org.Description,
			URL:         *org.URL,
		})
	}
	return res
}
