package main

import (
	"github.com/behnamgolds/embleman-go/internal/fsitem"
	"github.com/behnamgolds/embleman-go/internal/utils"
)

func main() {
	args := utils.ParseCmdArgs()
	if args[0] == "clear" {
		for _, path := range args[1:] {
			fi := fsitem.NewFsItem(path, "clear")
			fi.ExecuteAction()
		}
	} else {
		fi := fsitem.NewFsItem(args[1], args[0])
		fi.ExecuteAction()
	}

	utils.Refresh()
}
