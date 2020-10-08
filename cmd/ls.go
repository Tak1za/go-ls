package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

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
	cTime  time.Time
	perm   string
	size   int64
}

var (
	response          []data
	PRINT_WITH_TIME   bool = false
	PRINT_LONG_FORMAT bool = false
)

func ls(cmd *cobra.Command, args []string) {
	list, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range list {
		var inner data
		inner.name = v.Name()
		inner.isDir = v.IsDir()
		inner.cTime = v.ModTime()
		inner.perm = v.Mode().Perm().String()
		inner.size = v.Size()
		response = append(response, inner)
	}

	if ok, _ := cmd.Flags().GetBool("all"); !ok {
		response = filterDotFiles(response)
	}

	if ok, _ := cmd.Flags().GetBool("author"); ok {
		response = addAuthor(response)
	}

	if ok, _ := cmd.Flags().GetBool("c"); ok {
		response = sortByCTime(response)
		PRINT_WITH_TIME = true
	}

	if ok, _ := cmd.Flags().GetBool("long"); ok {
		response = addAuthor(response)
		response = sortName(response)
		PRINT_LONG_FORMAT = true
	}

	if PRINT_LONG_FORMAT {
		printLongFormat(response)
	} else if PRINT_WITH_TIME {
		printListWithTime(response)
	} else {
		response = sortName(response)
		printList(response)
	}
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

func printLongFormat(list []data) {
	tw := ansiterm.NewTabWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, item := range list {
		if item.isDir {
			color.New(color.FgBlue).FprintfFunc()(
				tw,
				"%v\t%v\t%v\t%v\t%v\n",
				item.perm,
				item.author,
				item.size,
				item.cTime.Local().Format(time.Stamp),
				item.name)
		} else {
			fmt.Fprintf(
				tw,
				"%v\t%v\t%v\t%v\t%v\n",
				item.perm,
				item.author,
				item.size,
				item.cTime.Local().Format(time.Stamp),
				item.name)
		}
	}
	tw.Flush()
}

func printListWithTime(list []data) {
	tw := ansiterm.NewTabWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _, item := range list {
		if item.isDir {
			color.New(color.FgBlue).FprintfFunc()(tw, "%v\t%v\t%v\n", item.name, item.author, item.cTime.Local().Format(time.Stamp))
		} else {
			fmt.Fprintf(tw, "%v\t%v\t%v\n", item.name, item.author, item.cTime.Local().Format(time.Stamp))
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

func sortByCTime(list []data) []data {
	sort.Slice(list, func(i, j int) bool {
		return list[i].cTime.After(list[j].cTime)
	})

	return list
}

func sortName(list []data) []data {
	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].name) < strings.ToLower(list[j].name)
	})

	return list
}
