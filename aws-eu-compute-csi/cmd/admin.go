/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"log"

	// "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// adminCmd represents the admin command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("admin called")
	},
}

func init() {
	rootCmd.AddCommand(adminCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adminCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adminCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	admin()
}

func admin() {
	// Create a new session using default AWS credentials
	sess := session.Must(session.NewSession())
	svc := iam.New(sess)

	// List IAM users
	result, err := svc.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		log.Fatalf("Unable to list users: %v", err)
	}

	for _, user := range result.Users {
		// Get policies attached to each user
		policies, err := svc.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
			UserName: user.UserName,
		})
		if err != nil {
			log.Fatalf("Unable to list policies for user %s: %v", *user.UserName, err)
		}

		// Check for AdministratorAccess policy
		for _, policy := range policies.AttachedPolicies {
			if *policy.PolicyName == "AdministratorAccess" {
				fmt.Printf("User %s has administrative access.\n", *user.UserName)
			}
		}
	}
}
