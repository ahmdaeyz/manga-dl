// Copyright © 2018 ahmdaeyz <ahmedalarabe5@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"os/user"

	"github.com/ahmdaeyz/manga-dl/manga"
	"github.com/spf13/cobra"
)

// dlChCmd represents the dlCh command
var dlChCmd = &cobra.Command{
	Use:   "dlCh",
	Short: "Download a specific chapter",
	Run: func(cmd *cobra.Command, args []string) {
		chapters := manga.GetChapters(args[0])
		cbz, err := cmd.Flags().GetBool("cbz")
		chNum, err := cmd.Flags().GetFloat64("ch-num")
		dirPath, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			log.Fatal(err)
		}
		for _, chapter := range chapters {
			if chapter.ChapterNum == chNum {
				chapter.DownloadByChapter(dirPath, cbz)
				break
			}
		}

	},
}

func init() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dlChCmd.Flags().BoolP("cbz", "z", false, "compress to comic book zip")
	dlChCmd.Flags().StringP("output-dir", "o", user.HomeDir, "Specifies the output directory.")
	dlChCmd.Flags().Float64P("ch-num", "n", 1, "chapter number")
	rootCmd.AddCommand(dlChCmd)

}
