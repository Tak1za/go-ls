/*
Copyright © 2020 Varun Gupta <varungupta2015135@gmail.com>

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
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ls",
	Short: "List down the content of a directory",
	Long:  `'ls' is a command for Windows that lists down the content of a directory, styled and unstyled can be changed via flags`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}

		list, err := ioutil.ReadDir(cwd)
		if err != nil {
			log.Fatalln(err)
		}

		if ok, _ := cmd.Flags().GetBool("all"); !ok {
			filterHiddenFiles(&list, cwd)
		}

		for _, item := range list {
			if item.IsDir() {
				color.Blue(item.Name())
				continue
			}
			color.White(item.Name())
		}
	},
}

func filterHiddenFiles(list *[]os.FileInfo, cwd string) {
	for i, v := range *list {
		path := cwd + "\\" + v.Name()
		pointer, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			log.Fatalln(err)
		}

		attr, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			log.Fatalln(err)
		}

		if attr&syscall.FILE_ATTRIBUTE_HIDDEN == 2 {
			*list = append((*list)[:i], (*list)[i+1:]...)
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("all", "a", false, "Do not ignore hidden files or directories")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-ls" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-ls")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
