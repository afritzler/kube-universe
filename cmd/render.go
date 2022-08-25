// Copyright © 2018 Andreas Fritzler
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
	"fmt"
	"os"

	renderer "github.com/afritzler/kube-universe/pkg/renderer"
	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Renders the cluster graph as JSON",
	Long:  `Renders the cluster graph into JSON format and prints it out to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		render()
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}

func render() {
	kubeconfig := rootCmd.Flag("kubeconfig").Value.String()
	data, err := renderer.GetGraph(kubeconfig)
	if err != nil {
		fmt.Printf("failed to render cluster graph: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", data)
}
