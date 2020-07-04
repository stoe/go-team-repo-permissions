package cmd

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	graphql "github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type team struct {
	Slug string
	Name string
}

type repo struct {
	Node struct {
		NameWithOwner string
	}
	Permission string
}

type rl struct {
	Cost      int
	Limit     int64
	NodeCount int64
	Remaining int64
	ResetAt   time.Time
}

var (
	token        string
	organization string
	writeCSV     bool

	client *graphql.Client
	ctx    = context.Background()

	rootCmd = &cobra.Command{
		Use:     "ghec-team-repo-permissions",
		Short:   "Get repository permissions for your organization teams",
		RunE:    getTeamRepoPermissions,
		Version: "0.0.2",
	}

	teamQuery struct {
		Organization struct {
			Teams struct {
				PageInfo struct {
					HasNextPage bool
					EndCursor   graphql.String
				}
				Nodes []team
			} `graphql:"teams(first: 100, after: $page)"`
		} `graphql:"organization(login: $org)"`
		RateLimit rl
	}

	repoQuery struct {
		Organization struct {
			Team struct {
				Repositories struct {
					Edges    []repo
					PageInfo struct {
						HasNextPage bool
						EndCursor   graphql.String
					}
				} `graphql:"repositories(first: 100, after: $page)"`
			} `graphql:"team(slug: $team)"`
		} `graphql:"organization(login: $org)"`
		RateLimit rl
	}
)

func getTeamRepoPermissions(cmd *cobra.Command, args []string) error {
	if token == "" {
		return errors.New("github.com personal access token (PAT) required")
	}

	if organization == "" {
		return errors.New("github.com organization required")
	}

	var teams []team

	variables := map[string]interface{}{
		"org":  graphql.String(organization),
		"page": (*graphql.String)(nil),
	}

	for {
		if err := client.Query(ctx, &teamQuery, variables); err != nil {
			panic(err)
		}

		teams = append(teams, teamQuery.Organization.Teams.Nodes...)

		if teamQuery.RateLimit.Remaining < 500 {
			throttle(teamQuery.RateLimit)
		}

		// break on last page
		if !teamQuery.Organization.Teams.PageInfo.HasNextPage {
			break
		}

		variables["page"] = *&teamQuery.Organization.Teams.PageInfo.EndCursor
	}

	var writer *csv.Writer
	var fp string

	if writeCSV {
		fp = "./team-repo-permissions.csv"
		os.Remove(fp)

		file, err := os.Create(fp)
		defer file.Close()

		if err != nil {
			return err
		}

		writer = csv.NewWriter(file)
	} else {
		writer = csv.NewWriter(os.Stdout)
	}

	header := []string{"team", "repository", "permission"}
	writer.Write(header)

	for _, t := range teams {
		variables := map[string]interface{}{
			"org":  graphql.String(organization),
			"team": graphql.String(t.Slug),
			"page": (*graphql.String)(nil),
		}

		for {
			if err := client.Query(ctx, &repoQuery, variables); err != nil {
				panic(err)
			}

			for _, r := range repoQuery.Organization.Team.Repositories.Edges {
				writer.Write([]string{
					t.Slug,
					r.Node.NameWithOwner,
					r.Permission,
				})
				writer.Flush()
			}

			if repoQuery.RateLimit.Remaining < 500 {
				throttle(repoQuery.RateLimit)
			}

			// break on last page
			if !repoQuery.Organization.Team.Repositories.PageInfo.HasNextPage {
				break
			}

			variables["page"] = *&repoQuery.Organization.Team.Repositories.PageInfo.EndCursor
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	if writeCSV {
		fmt.Println("output saved to", fp)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initClient)

	rootCmd.PersistentFlags().StringVarP(
		&token, "token", "t", "GITHUB_TOKEN",
		"github.com personal access token",
	)
	token = os.Getenv("GITHUB_TOKEN")

	rootCmd.PersistentFlags().StringVarP(
		&organization, "org", "o", "",
		"github.com organization",
	)

	rootCmd.PersistentFlags().BoolVar(
		&writeCSV, "csv", false,
		"save to CSV file (team-repo-permissions.csv) in current directory",
	)
}

// initClient creates the github.com GraphQL client
func initClient() {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(ctx, src)
	client = graphql.NewEnterpriseClient(
		"https://api.github.com/graphql", httpClient,
	)
}

func throttle(rl rl) {
	time.Sleep(3600 * time.Millisecond)
}
