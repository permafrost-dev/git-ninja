/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/permafrost-dev/git-ninja/lib/integrations/jira"
	"github.com/spf13/cobra"
)

var jiraissuesCmd = &cobra.Command{
	Use:    "jira:issues",
	Short:  "List open JIRA Ticket IDs",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Getenv("JIRA_SUBDOMAIN")) == 0 || len(os.Getenv("JIRA_EMAIL_ADDRESS")) == 0 || len(os.Getenv("JIRA_API_TOKEN")) == 0 {
			fmt.Println("Error: JIRA_SUBDOMAIN, JIRA_EMAIL_ADDRESS and JIRA_API_TOKEN environment variables must be set.")
			return
		}

		ids := jira.GetJiraTicketIDs(os.Getenv("JIRA_SUBDOMAIN"), os.Getenv("JIRA_EMAIL_ADDRESS"))

		if len(ids) > 0 {
			fmt.Println("Open JIRA Tickets:")
			for _, id := range ids {
				fmt.Println(id)
			}
		}

		if len(ids) == 0 {
			fmt.Println("No open JIRA tickets.")
		}
	},
}

func init() {
	rootCmd.AddCommand(jiraissuesCmd)
}
