package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hectane/go-acl/api"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
)

type data struct {
	name   string
	author string
	isDir  bool
}

var response []data

func ls(cmd *cobra.Command, args []string) {
	list, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range list {
		var inner data
		inner.name = v.Name()
		inner.isDir = v.IsDir()
		response = append(response, inner)
	}

	if ok, _ := cmd.Flags().GetBool("all"); !ok {
		response = filterDotFiles(response)
	}

	if ok, _ := cmd.Flags().GetBool("author"); ok {
		response = addAuthor(response)
	}

	printList(response)
}

func printList(list []data) {
	tw := ansiterm.NewTabWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, item := range list {
		if item.isDir {
			color.New(color.FgBlue).FprintfFunc()(tw, "%v\t%v\n", item.name, item.author)
		} else {
			fmt.Fprintf(tw, "%v\t%v\n", item.name, item.author)
		}
	}
	tw.Flush()
}

func filterDotFiles(list []data) []data {
	var filteredList []data
	for _, v := range list {
		if strings.Index(v.name, ".") != 0 {
			filteredList = append(filteredList, v)
		}
	}

	return filteredList
}

func addAuthor(list []data) []data {
	cwd, _ := os.Getwd()
	var filteredList []data
	for _, v := range list {
		path := cwd + "\\" + v.name
		var owner *windows.SID
		err := api.GetNamedSecurityInfo(
			path,
			api.SE_FILE_OBJECT,
			api.OWNER_SECURITY_INFORMATION,
			&owner,
			nil,
			nil,
			nil,
			nil)
		if err != nil {
			log.Fatalln(err)
		}
		systemName, _ := os.Hostname()
		account, _, _, _ := owner.LookupAccount(systemName)
		v.author = account
		filteredList = append(filteredList, v)
	}

	return filteredList
}
