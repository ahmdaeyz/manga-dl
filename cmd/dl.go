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
		"github.com/spf13/cobra"
		"manga-dl/manga"
		"fmt"
		"strings"
		"log"
	"os/user"
)

// dlCmd represents the dl command
var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "Download the whole manga",
	Run: func(cmd *cobra.Command, args []string) {
		mng := manga.GetManga(args[0])
		fmt.Println("Now Downloading Stay Tuned")
		volumes := mng.GetVolumes()
		mangaName := strings.Replace(args[0], "http://fanfox.net/manga/", "", -1)
		mangaName = strings.Replace(mangaName, "/", "", -1)
		dirPath, err := cmd.Flags().GetString("output-dir")
		cbz, err := cmd.Flags().GetBool("cbz")
		byVolume,err:=cmd.Flags().GetBool("by-vol")
		byChapter,err:=cmd.Flags().GetBool("by-ch")
		if err != nil {
			log.Fatal(err)
		}
		manga.CreateDirIfNotExist(dirPath + "/" + mangaName)
		if byVolume{
			for i:=len(volumes)-1;i>=0 ;i-- {
				volumes[i].DownloadByVolume(dirPath+"/"+mangaName, cbz)
			}
		}else if byChapter{
			for i:=len(volumes)-1;i>=0 ;i--{
				for _,chapter:= range volumes[i].GetChapters(){
					chapter.DownloadByChapter(dirPath + "/" + mangaName,cbz)
				}
			}
		}

	},
}

func init() {
	user,err:=user.Current()
	if err!=nil{
		log.Fatal(err)
	}
	dlCmd.Flags().BoolP("cbz","z",false,"Compress the chapters or volumes to comic book zip\nthe output is a file for each volume or chapter.")
	dlCmd.Flags().StringP("output-dir","o",user.HomeDir,"Specifies the output directory.")
	dlCmd.Flags().BoolP("by-vol","v",false,"Download by volume")
	dlCmd.Flags().BoolP("by-ch","c",true,"Download by chapter")
	rootCmd.AddCommand(dlCmd)

}
