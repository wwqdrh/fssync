package server

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/wwqdrh/logger"
)

// source: ./testdata or testdata
func ListDirFile(source string, prefix bool) ([]string, error) {
	source = strings.TrimLeft(source, "./")
	dirStack := []string{source}

	res := []string{}
	for len(dirStack) > 0 {
		cur := dirStack[0]
		dirStack = dirStack[1:]

		files, err := ioutil.ReadDir(cur)
		if err != nil {
			logger.DefaultLogger.Warn(cur + " 不是文件夹")
			continue
		}

		for _, item := range files {
			if item.IsDir() {
				dirStack = append(dirStack, path.Join(cur, item.Name()))
			} else {
				res = append(res, path.Join(cur, item.Name()))
			}
		}
	}

	if !prefix {
		for i := 0; i < len(res); i++ {
			cur := strings.TrimPrefix(res[i], source)
			cur = strings.TrimPrefix(cur, "/")
			res[i] = cur
		}
	}

	return res, nil
}
