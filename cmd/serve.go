// Copyright © 2018 Andreas Fritzler <andreas.fritzler@gmail.com>
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
	"log"
	"net/http"
	"os"

	renderer "github.com/afritzler/kube-universe/pkg/renderer"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/afritzler/kube-universe/statik"
)

var port string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a webserver to serve the 3D landscape view",
	Long: `Starts a webserver to serve the 3D landscape view.

By default, the website can be accessed on http://localhost:3000. The JSON representation of
the landscape graph can be found under http://localhost:3000/graph.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&port, "port", "p", "3000", "Port on which the server should listen")
	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))
}

func serve() {
	fmt.Printf("started server on http://localhost:%s\n", port)
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(statikFS))
	http.HandleFunc("/graph", graphResponse)
	http.ListenAndServe(getPort(), nil)
}

func graphResponse(w http.ResponseWriter, r *http.Request) {
	config := os.Getenv("KUBECONFIG")
	if config == "" {
		config = rootCmd.Flag("kubeconfig").Value.String()
	}
	data, err := renderer.GetGraph(config)
	if err != nil {
		fmt.Printf("failed to render landscape graph: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func getPort() string {
	return fmt.Sprintf(":%s", port)
}
