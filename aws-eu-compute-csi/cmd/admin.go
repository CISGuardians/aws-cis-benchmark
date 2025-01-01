/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"sync"

	// "github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/spf13/cobra"

	// "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

var profiles []string

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
		// Check if the profile flag is empty
		if len(profiles) == 0 {
			log.Fatalf("--profile option is required. Please specify an AWS profile or profiles with comma separated.")
		}
		checkProfilesInParallel(profiles)
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
	// Make the --profile flag required
	adminCmd.Flags().StringSliceVarP(&profiles, "profiles", "p", []string{}, "Comma-separated list AWS profile to use (required) eg: --profiles default,dev,test")
	adminCmd.MarkFlagRequired("profiles")
}


func checkProfilesInParallel(profiles []string) {
	var wg sync.WaitGroup
	results := make(chan string, len(profiles)) // Channel to collect results

	for _, prof := range profiles {
		wg.Add(1)
		go func(prof string) {
			defer wg.Done()
			prof = strings.TrimSpace(prof)

			// Attempt to create a session with the given profile
			sess, sess_err := session.NewSessionWithOptions(session.Options{

				Profile: prof, // Specify the AWS profile
			})
			if sess_err != nil {
			
				// log.Printf("Failed to create session: %v", sess_err)
				// log.Printf("Skipping profile %s: Failed to create session: %v", prof, sess_err)
				results <- fmt.Sprintf("%-15s   | Session creation failed", prof)
				return // Exit the goroutine early
			}
				
			svc := iam.New(sess)

			// List IAM users
			result, err := svc.ListUsers(&iam.ListUsersInput{})
			if err != nil {
				
				// log.Printf("Skipping profile %s: Unable to list users: %v", prof, err)
				results <- fmt.Sprintf("%-15s   | Failed to list users", prof)
				return // Exit the goroutine early
			}
			
			for _, user := range result.Users {
				// Get policies attached to each user
				policies, err := svc.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
					UserName: user.UserName,
				})
				if err != nil {
					log.Printf("Skipping user %s: Unable to list policies: %v", *user.UserName, err)
					continue // Skip this user and move to the next one
				}
		
				// Check for AdministratorAccess policy
				for _, policy := range policies.AttachedPolicies {
					if *policy.PolicyName == "AdministratorAccess" {
						// fmt.Printf("User %s has administrative access.\n", *user.UserName)
						results <- fmt.Sprintf("%-15s   | This profile has administrative access", prof)
					}
				}
			}
			
	
		}(prof)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(results)

	// Print the results in matrix format
	fmt.Println("Profile           | Status")
	fmt.Println("----------------------------")
	for result := range results {
		fmt.Println(result)
	}
}