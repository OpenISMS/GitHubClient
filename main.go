package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/joho/godotenv"
	openisms "github.com/openisms/api/bindings/go/io/openisms"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var token string
var org string
var ctx = context.Background()

func main() {

	_ = godotenv.Load()
	//if envErr != nil {
	//	log.Fatalf("Can't load env %v", envErr)
	//}

	token = os.Getenv("TOKEN")
	org = os.Getenv("ORG")
	listen := getEnv("LISTEN", "localhost:2701")

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/audit", getAudit)

	fmt.Printf("Go To http://localhost:2701/audit\n")

	log.Fatal(http.ListenAndServe(listen, mux))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getClient() *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func getAudit(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client := getClient()

	//results := map[string]interface{}{
	//	"myOrgs":  MyOrgs(client, ctx),
	//	"org":     ShowOrg(client, ctx, org),
	//	"members": ListMembers(client, ctx, org),
	//	"repos":   GetRepos(client, ctx, org),
	//	"events":  GetEvents(client, ctx, org),
	//}

	meta := openisms.Meta{
		Created:    timestamppb.Now(),
		Identifier: "random",
	}

	source := openisms.SourceSystem{
		Name:   "GitHub",
		Vendor: "GitHub Inc.",
		Type:   openisms.SourceSystem_source_control_management,
	}

	users := GetUsers(client, ctx, org)

	e := openisms.Event{
		Meta:          &meta,
		Source:        &source,
		Users:         users,
		Devices:       nil,
		Repositories:  nil,
		Certification: nil,
		Pentest:       nil,
		Stats:         nil,
	}

	j, _ := json.MarshalIndent(e, "", "  ")
	_, err := io.WriteString(w, string(j))
	if err != nil {
		log.Fatal(err)
	}
}

func MyOrgs(client *github.Client, ctx context.Context) []map[string]interface{} {

	myOrgs, _, _ := client.Organizations.List(ctx, "", nil)

	var res []map[string]interface{}

	for _, org := range myOrgs {
		info := map[string]interface{}{
			"Name":        org.Login,
			"Description": org.Description,
			"URL":         org.URL,
		}
		res = append(res, info)
	}
	return res
}

func ShowOrg(client *github.Client, ctx context.Context, orgName string) map[string]interface{} {

	org, _, _ := client.Organizations.Get(ctx, orgName)

	info := map[string]interface{}{
		"Name":                          org.Name,
		"Company":                       org.Company,
		"Description":                   org.Description,
		"Email":                         org.Email,
		"BillingEmail":                  org.BillingEmail,
		"Location":                      org.Location,
		"MembersCanCreatePrivateRepos":  org.MembersCanCreatePrivateRepos,
		"MembersCanCreateInternalRepos": org.MembersCanCreateInternalRepos,
		"MembersCanCreatePublicRepos":   org.MembersCanCreatePublicRepos,
		"Plan":                          org.Plan.Name,
		"Seats":                         org.Plan.Seats,
		"FilledSeats":                   org.Plan.FilledSeats,
		"Collaborators":                 org.Plan.Collaborators,
		"URL":                           org.HTMLURL,
	}
	return info
}

func GetUsers(client *github.Client, ctx context.Context, orgName string) []*openisms.User {
	members, _, _ := client.Organizations.ListMembers(ctx, orgName, nil)
	var res []*openisms.User

	for _, user := range members {

		fullUser, _, _ := client.Users.Get(ctx, *user.Login)

		info := openisms.User{
			Person: &openisms.Person{
				Id: strconv.FormatInt(*fullUser.ID, 10),
				Name: &openisms.Name{
					FullName: *fullUser.Name,
					//Role:     *fullUser.RoleName,
				},
				PrimaryEmail: *fullUser.Email,
				//Picture:      fullUser.GravatarID,
				Company: *fullUser.Company,
				Created: timestamppb.New(fullUser.CreatedAt.UTC()),
				Updated: timestamppb.New(fullUser.UpdatedAt.UTC()),
			},
			SecondFactorActive:   nil,
			SecondFactorEnforced: nil,
			Active:               nil,
			//Suspended:            wrapperspb.Bool(fullUser.SuspendedAt.Time.Unix() > 1),
			Disabled:          nil,
			Deleted:           nil,
			Groups:            nil,
			Privileges:        nil,
			Employment:        nil,
			ConnectedAccounts: nil,
			Possessions:       nil,
			Tags:              nil,
		}

		//info := map[string]interface{}{
		//	"Name":        fullUser.Name,
		//	"Email":       fullUser.Email,
		//	"URL":         fullUser.HTMLURL,
		//	"Login":       fullUser.Login,
		//	"AvatarURL":   fullUser.AvatarURL,
		//	"Type":        fullUser.Type,
		//	"CreatedAt":   fullUser.CreatedAt,
		//	"UpdatedAt":   fullUser.UpdatedAt,
		//	"SuspendedAt": fullUser.SuspendedAt,
		//	"Location":    fullUser.Location,
		//	//"TwoFactorAuthentication": fullUser.TwoFactorAuthentication,
		//}
		res = append(res, &info)
	}

	return res

}

func ListMembers(client *github.Client, ctx context.Context, orgName string) []map[string]interface{} {
	members, _, _ := client.Organizations.ListMembers(ctx, orgName, nil)
	var res []map[string]interface{}

	for _, user := range members {

		fullUser, _, _ := client.Users.Get(ctx, *user.Login)

		info := map[string]interface{}{
			"Name":        fullUser.Name,
			"Email":       fullUser.Email,
			"URL":         fullUser.HTMLURL,
			"Login":       fullUser.Login,
			"AvatarURL":   fullUser.AvatarURL,
			"Type":        fullUser.Type,
			"CreatedAt":   fullUser.CreatedAt,
			"UpdatedAt":   fullUser.UpdatedAt,
			"SuspendedAt": fullUser.SuspendedAt,
			"Location":    fullUser.Location,
			//"TwoFactorAuthentication": fullUser.TwoFactorAuthentication,
		}
		res = append(res, info)
	}

	return res
}

func GetEvents(client *github.Client, ctx context.Context, orgName string) []map[string]interface{} {
	activities, _, _ := client.Activity.ListEventsForOrganization(ctx, orgName, nil)

	var res []map[string]interface{}

	for _, act := range activities {

		info := map[string]interface{}{
			"CreatedAt": act.CreatedAt,
			"Actor":     act.Actor.Login,
			"Repo":      act.Repo.Name,
			"Type":      act.Type,
		}
		res = append(res, info)
	}

	return res
}

func GetRepos(client *github.Client, ctx context.Context, orgName string) []map[string]interface{} {
	repos, _, _ := client.Repositories.List(ctx, orgName, nil)

	var res []map[string]interface{}

	for _, repo := range repos {

		// needs repo-scope!
		protection, _, _ := client.Repositories.GetBranchProtection(ctx, *repo.Owner.Login, *repo.Name, *repo.DefaultBranch)

		info := map[string]interface{}{
			"FullName":      repo.FullName,
			"Description":   repo.Description,
			"HTMLURL":       repo.HTMLURL,
			"Language":      repo.Language,
			"Archived":      repo.Archived,
			"Private":       repo.Private,
			"Disabled":      repo.Disabled,
			"Visibility":    repo.Visibility,
			"DefaultBranch": repo.DefaultBranch,
			"IsFork":        repo.Fork,
			"BranchProtection": map[string]interface{}{
				"restrictions":                          getUrlsForBranchRestrictions(protection.Restrictions),
				"enforce_admins":                        protection.EnforceAdmins.Enabled,
				"RequireLinearHistory":                  protection.RequireLinearHistory.Enabled,
				"AllowDeletions":                        protection.AllowDeletions.Enabled,
				"AllowForcePushes":                      protection.AllowForcePushes.Enabled,
				"RequiredConversationResolution":        protection.RequiredConversationResolution.Enabled,
				"RequiredPullRequestReviewsCount":       protection.RequiredPullRequestReviews.RequiredApprovingReviewCount,
				"RequiredPullRequestReviewsBypass":      getUrls(protection.RequiredPullRequestReviews.BypassPullRequestAllowances),
				"RequiredPullRequestReviewsByCodeowner": protection.RequiredPullRequestReviews.RequireCodeOwnerReviews,
				"RequiredStatusChecks":                  protection.RequiredStatusChecks.Checks,
				"RequiredStatusChecksStrict":            protection.RequiredStatusChecks.Strict,
			},
		}
		res = append(res, info)
	}

	return res
}

func getUrls(allowances *github.BypassPullRequestAllowances) []string {
	var res []string
	res = append(res, getUserUrls(allowances.Users)...)
	res = append(res, getTeamUrls(allowances.Teams)...)
	res = append(res, getAppUrls(allowances.Apps)...)
	return res
}

func getUrlsForBranchRestrictions(restrictions *github.BranchRestrictions) []string {
	var res []string
	res = append(res, getUserUrls(restrictions.Users)...)
	res = append(res, getTeamUrls(restrictions.Teams)...)
	res = append(res, getAppUrls(restrictions.Apps)...)
	return res
}

func getUserUrls(users []*github.User) []string {
	var res []string
	for _, u := range users {
		res = append(res, *u.HTMLURL)
	}
	return res
}

func getTeamUrls(teams []*github.Team) []string {
	var res []string
	for _, u := range teams {
		res = append(res, *u.HTMLURL)
	}
	return res
}

func getAppUrls(apps []*github.App) []string {
	var res []string
	for _, u := range apps {
		res = append(res, *u.HTMLURL)
	}
	return res
}
