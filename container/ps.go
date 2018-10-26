package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/chengxiangwang/libcontainer/config"
	"github.com/chengxiangwang/libcontainer/utils"
	"text/tabwriter"
)

func PrintShow(showAll bool) error {
	writer := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	printHeader(writer)
	_, err := os.Stat(config.CONTAINER_ROOT)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(config.CONTAINER_ROOT, os.ModeDir); err != nil {
			return err
		}
	}
	rootDir, err := ioutil.ReadDir(config.CONTAINER_ROOT)
	if err != nil {
		return err
	}
	for _, chDir := range rootDir {
		if chDir.IsDir() {
			configFile := path.Join(config.CONTAINER_ROOT, chDir.Name(), "config.json")
			if !utils.PathExists(configFile) {
				continue
			}
			data, _ := ioutil.ReadFile(configFile)
			c, err := JsonAsContainer(data)
			if err != nil || c == nil {
				continue
			}
			printContainer(c, writer)
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

func printHeader(writer *tabwriter.Writer) {
	fmt.Fprint(writer, "ID\tNAME\tIMAGEID\tCOMMAND\tCREATED\tSTATUS\n")
}

func printContainer(c *Container, writer *tabwriter.Writer) {
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n", c.ID.String(), c.Name, c.ImageID, c.Args[0], c.Created, c.Status.Name)
}
