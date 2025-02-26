/*
Copyright © 2025 NAME HERE yajushsharma12@gmail.com
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/sharmayajush/challenge2015/pkg/degrees"
	"github.com/spf13/cobra"
)

// degreesCmd represents the degrees command
var degreesCmd = &cobra.Command{
	Use:   "degrees",
	Short: "Find the degrees of separation between two actors",
	Long:  `Find the degrees of separation between two actors using a BFS algorithm.`,
	Args:  cobra.ExactArgs(2), // Ensure exactly 2 arguments are provided
	Run: func(cmd *cobra.Command, args []string) {
		start := args[0]
		target := args[1]

		// Call your BFS function
		ans, err := degrees.BfsWithPath(start, target)
		if err != nil {
			log.Fatalf("Error finding path: %v", err)
		}

		// Print the result
		fmt.Printf("Degrees of Separation: %v \n\n", len(ans.Nodes))
		for i, node := range ans.Nodes {
			fmt.Printf("%v. Movie: %s\n", i+1, node.Movie)
			fmt.Printf("%s: %s\n", node.Role1, node.Person1)
			fmt.Printf("%s: %s\n\n", node.Role2, node.Person2)
		}
	},
}

func init() {
	rootCmd.AddCommand(degreesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// degreesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// degreesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
