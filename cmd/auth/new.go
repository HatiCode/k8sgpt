/*
Copyright 2023 The K8sGPT Authors.
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

package auth

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/k8sgpt-ai/k8sgpt/pkg/ai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Configure new provider",
	Long:  "The new command allows to configure a new backend AI provider",
	Run: func(cmd *cobra.Command, args []string) {

		// get ai configuration
		err := viper.UnmarshalKey("ai", &configAI)
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		// search for provider with same name
		providerIndex := -1
		for i, provider := range configAI.Providers {
			if backend == provider.Name {
				providerIndex = i
				break
			}
		}

		// check if backend is not empty
		if backend == "" {
			color.Red("Error: Backend AI cannot be empty.")
			os.Exit(1)
		}

		// check if model is not empty
		if model == "" {
			color.Red("Error: Model cannot be empty.")
			os.Exit(1)
		}

		if ai.NeedPassword(backend) && password == "" {
			fmt.Printf("Enter %s Key: ", backend)
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				color.Red("Error reading %s Key from stdin: %s", backend,
					err.Error())
				os.Exit(1)
			}
			password = strings.TrimSpace(string(bytePassword))
		}

		// create new provider object
		newProvider := ai.AIProvider{
			Name:     backend,
			Model:    model,
			Password: password,
			BaseURL:  baseURL,
		}

		if providerIndex == -1 {
			// provider with same name does not exist, add new provider to list
			configAI.Providers = append(configAI.Providers, newProvider)
			viper.Set("ai", configAI)
			if err := viper.WriteConfig(); err != nil {
				color.Red("Error writing config file: %s", err.Error())
				os.Exit(1)
			}
			color.Green("%s added to the AI backend provider list", backend)
		} else {
			// provider with same name exists, update provider info
			color.Yellow("Provider with same name already exists, use update command to modify an existing provider configuration")
		}
	},
}

func init() {
	// add flag for backend
	newCmd.Flags().StringVarP(&backend, "backend", "b", "openai", "Backend AI provider")
	// add flag for model
	newCmd.Flags().StringVarP(&model, "model", "m", "gpt-3.5-turbo", "Backend AI model")
	// add flag for password
	newCmd.Flags().StringVarP(&password, "password", "p", "", "Backend AI password")
	// add flag for url
	newCmd.Flags().StringVarP(&baseURL, "baseurl", "u", "", "URL AI provider, (e.g `http://localhost:8080/v1`)")
}
