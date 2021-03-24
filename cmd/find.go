/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/spf13/cobra"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Search within specified index",
	Long: `This command searches within the selected index
	and returns the books in which the query was found`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Find(args)
	},
}

func init() {
	findCmd.Flags().String("index", "", "location of existing index to be used")
	findCmd.Flags().String("query", "", "query to be evaluated")
	rootCmd.AddCommand(findCmd)
}

func Find(args []string) error {
	index := args[0]
	query := args[1]

	ix, err := bleve.Open(index)
	if err != nil {
		return err
	}
	res, err := ix.Search(bleve.NewSearchRequest(bleve.NewMatchQuery(query)))
	if err != nil {
		return err
	}
	fmt.Println(res.String())
	return nil
}
