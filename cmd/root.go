/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hub-stats",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Hello World!")
		cmd.Printf("args: %v\n", args)
		if len(args) == 0 {
			fmt.Println(Red + "username is required" + Reset)
			os.Exit(1)
		}
		username := args[0]
		if username == "" {
			fmt.Println(Red + "username is required" + Reset)
			os.Exit(1)
		}
		listAllUserRepos(username, 1)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var totalRepositoryList GithubRepos

func listAllUserRepos(username string, page int) {
	fmt.Println(Cyan+"==== fetching page =====", page, "===="+Reset)
	fmt.Println(Cyan+"==== total fetched =====", len(totalRepositoryList), "===="+Reset)
	if totalRepositoryList == nil || len(totalRepositoryList)%100 == 0 && page < 10 {
		repositoryList, err := listUserRepos(username, page)
		if err != nil {
			fmt.Println(err)
		} else {
			totalRepositoryList = append(totalRepositoryList, repositoryList...)
		}
		listAllUserRepos(username, page+1)
	}

}
func listUserRepos(username string, page int) (GithubRepos, error) {
	var repositoryList GithubRepos
	reposUrl := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", username, page)
	fmt.Println(Cyan+"==== fetching from =====", reposUrl, "===="+Reset)

	resp, err := http.Get(reposUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	switch resp.StatusCode {
	case 200:
		fmt.Println(Green + "==== success =====" + Reset)

		err = json.Unmarshal(body, &repositoryList)
		if err != nil {
			panic(err)
		}

		return repositoryList, nil

	case 404:
		var response GithubResponseError
		err = json.Unmarshal(body, &response)
		if err != nil {
			panic(err)
		}
		errMessage := fmt.Sprintf("%s Learn more about this error : %s StatusCode: %d", response.Message, response.DocumentationURL, resp.StatusCode)
		log.Fatal(errMessage)
		return nil, errors.New(errMessage)
	default:

		fmt.Println("body:", string(body))
		errMessage := fmt.Sprintf("Something went wrong : %s StatusCode: %d", string(body), resp.StatusCode)
		log.Fatal(errMessage)
		return nil, errors.New(errMessage)

	}
}

type GithubRepos []struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

type GithubResponseError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}
