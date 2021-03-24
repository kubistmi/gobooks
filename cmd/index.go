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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/ledongthuc/pdf"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Create new index",
	Long: `This command creates new index at selected location,
walks through the selected folder and subfolders,
read all PDF files and add them into the index.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return IndexFolder(args)
	},
}

func init() {
	indexCmd.Flags().String("index", "", "location where the new index should be created")
	indexCmd.Flags().String("folder", "", "folder location that should be analyzed")
	rootCmd.AddCommand(indexCmd)
}

func IndexFolder(args []string) error {
	index := args[0]
	folder := args[1]

	mapping := bleve.NewIndexMapping()
	ix, err := bleve.New(index, mapping)
	if err != nil {
		return nil
	}
	err = filepath.Walk(folder, newWalkFunc(&ix))
	if err != nil {
		return err
	}
	return nil
}

func newWalkFunc(ix *bleve.Index) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".pdf" {
			return nil
		}

		text, err := GetPdf(path)
		if err != nil {
			return err
		}

		return (*ix).Index(path, text)
	}
}

func GetPdf(path string) (string, error) {

	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	t, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	b := new(strings.Builder)
	_, err = io.Copy(b, t)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
