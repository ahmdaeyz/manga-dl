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
	"strings"

	"github.com/ahmdaeyz/manga-dl/manga"
	"github.com/spf13/cobra"
)

// dlVolCmd represents the dlVol command
var dlVolCmd = &cobra.Command{
	Use:   "dlVol",
	Short: "Download a specific volume",
	Long:  `NOTE: for vTBD just enter what is there\nyou may download the whole manga\nwith "dl" if the volume name is not available'`,
	Run: func(cmd *cobra.Command, args []string) {
		mng := manga.GetManga(args[0])
		volumes := mng.GetVolumes()
		dirPath, err := cmd.Flags().GetString("output-dir")
		cbz, err := cmd.Flags().GetBool("cbz")
		volNum, err := cmd.Flags().GetString("vol-num")
		if err != nil {
			log.Fatal(err)
		}
		for _, volume := range volumes {
			if Num := strings.Replace(volume.VolNum, "Volume ", "", -1); Num == volNum {
				volume.DownloadByVolume(dirPath, cbz)
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
	dlVolCmd.Flags().BoolP("cbz", "z", false, "compress to comic book zip")
	dlVolCmd.Flags().StringP("output-dir", "o", user.HomeDir, "Specifies the output directory.")
	dlVolCmd.Flags().StringP("vol-num", "n", "01", "volume number")
	rootCmd.AddCommand(dlVolCmd)

}
