package fsitem

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

type FsItem struct {
	path             *gio.File
	action           string
	pathInfo         *gio.FileInfo
	emblems          []string
	numEmblemPattern string
	numEmblemFound   bool
}

func (fi *FsItem) popEmblem(index int) {
	fi.emblems = append(fi.emblems[:index], fi.emblems[index+1:]...)
}

func (fi *FsItem) changeNum(num int) int {
	if fi.isIncreaseAction() {
		if num == 19 {
			return -1
		}
		return num + 1
	} else {
		if num <= 1 {
			return -1
		}
		return num - 1
	}
}

func (fi *FsItem) matchesPattern(pattern, str string) bool {
	matched, _ := regexp.MatchString(`^emblem-num-[0-9]+-symbolic$`, str)
	return matched
}

func (fi *FsItem) isIncreaseAction() bool {
	return fi.action == "inc"
}

func (fi *FsItem) setOne() {
	one := "emblem-num-" + "1" + "-symbolic"
	fi.emblems = append(fi.emblems, one)
	fi.setEmblems()
}

func (fi *FsItem) setEmblems() {
	fi.pathInfo.SetAttributeStringv("metadata::emblems", fi.emblems)
	fi.path.SetAttributesFromInfo(context.Background(), fi.pathInfo, gio.FileQueryInfoNone)
}

func (fi *FsItem) executeNumAction() {
	for i, emblem := range fi.emblems {
		if fi.matchesPattern(fi.numEmblemPattern, emblem) {
			fi.numEmblemFound = true
			num, err := strconv.Atoi(strings.Split(emblem, "-")[2])
			if err != nil {
				log.Fatalln(err.Error())
			}
			num = fi.changeNum(num)
			if num == -1 {
				fi.popEmblem(i)
				fi.setEmblems()
				return
			} else {
				emblem = "emblem-num-" + strconv.Itoa(num) + "-symbolic"
				fi.emblems[i] = emblem
				fi.setEmblems()
				return
			}
		}
	}

	if !fi.numEmblemFound && fi.action != "dec" {
		fi.setOne()
		return
	}
}

func (fi *FsItem) executeToggleAction() {
	var emblem string
	if fi.action == "clock" {
		emblem = "emblem-urgent"
	} else {
		emblem = "vcs-normal"
	}

	for i := range fi.emblems {
		if fi.emblems[i] == emblem {
			// remove emblem
			fi.popEmblem(i)
			fi.setEmblems()
			return
		}
	}
	// if the emblem is not there, add it
	fi.emblems = append(fi.emblems, emblem)
	fi.setEmblems()
}

func (fi *FsItem) executeClearAction() {
	fi.emblems = []string{}
	fi.setEmblems()
}

func (fi *FsItem) ExecuteAction() {
	switch fi.action {
	case "inc", "dec":
		fi.executeNumAction()
	case "clear":
		fi.executeClearAction()
	case "clock", "check":
		fi.executeToggleAction()
	}
}

func NewFsItem(path, action string) *FsItem {
	fi := &FsItem{}
	fi.path = gio.NewFileForPath(path)
	fi.action = action
	var err error
	fi.pathInfo, err = fi.path.QueryInfo(context.Background(), "metadata::emblems", gio.FileQueryInfoNone)
	if err != nil {
		log.Fatalln(err.Error())
	}

	fi.emblems = fi.pathInfo.AttributeStringv("metadata::emblems")
	fi.numEmblemPattern = `^emblem-num-[0-9]+-symbolic$`
	fi.numEmblemFound = false
	return fi
}
