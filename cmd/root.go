/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
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
			fmt.Println(Red+"username is required"+Reset)
			os.Exit(1)
		}
		username := args[0]
		if username == "" {
			fmt.Println(Red+"username is required"+Reset)
			os.Exit(1)
		}
		listAllUserRepos(username)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hub-stats.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listAllUserRepos(username string) {
	reposUrl := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	fmt.Println(Cyan+"==== fetching from =====", reposUrl, "===="+Reset)

	resp, err := http.Get(reposUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("StatusCode:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	switch resp.StatusCode {
	case 200:
		fmt.Println(Green + "==== success =====" + Reset)
		var response GithubRepos
		err = json.Unmarshal(body, &response)
		if err != nil {
			panic(err)
		}
		for _, repo := range response {
			fmt.Println(repo.Name)
		}

	case 404:
		var response GithubResponseError
		err = json.Unmarshal(body, &response)
		if err != nil {
			panic(err)
		}
		fmt.Println(response.Message)
		fmt.Println("Learn more about this error :", response.DocumentationURL)

	default:
		fmt.Println("StatusCode:", resp.StatusCode)
		fmt.Println("body:", string(body))

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
