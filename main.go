package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/chengxiangwang/libcontainer/container"
	"github.com/urfave/cli"
	"os"
)

var runCommand = cli.Command{
	Name:  "run",
	Usage: `goC run -ti [command]`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "cpu",
			Usage: "set cpu share",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "set cpu core",
		},
		cli.StringFlag{
			Name:  "memory",
			Usage: "set memory",
		},
		cli.StringSliceFlag{
			Name:  "v",
			Usage: "create volumn",
		},
	},
	Action: func(context *cli.Context) error {
		if (len(context.Args())) < 1 {
			return fmt.Errorf("Missing container command")
		}

		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		name := context.String("name")
		tty := (context.Bool("ti") || context.Bool("it"))
		volumns := context.StringSlice("v")
		limits := map[string]string{
			"cpu":    context.String("cpu"),
			"cpuset": context.String("cpuset"),
			"memory": context.String("memory"),
		}
		//RunContainer(name string, tty bool, initArgs []string, limits map[string]string, volumns []string)
		err := container.RunContainer(name, tty, cmdArray, limits, volumns)
		if err != nil {
			return err
		}
		return nil
	},
}
var psCommand = cli.Command{
	Name:  "ps",
	Usage: "libcontainer ps",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "a",
			Usage: "show all container",
		},
	},
	Action: func(context *cli.Context) error {
		showAll := context.Bool("b")
		if err := container.PrintShow(showAll); err != nil {
			return err
		}
		return nil
	},
}

var logCommand = cli.Command{
	Name:  "log",
	Usage: "libcontainer log container_id",
	Flags: []cli.Flag{},
	Action: func(context *cli.Context) error {
		if (len(context.Args())) < 1 {
			return fmt.Errorf("Missing container id")
		}
		containerId := context.Args()[0]
		container.PrintContainerLog(containerId)
		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "goC"
	app.Usage = "goC run -it /bin/bash"
	app.Commands = []cli.Command{
		runCommand,
		psCommand,
		logCommand,
	}
	app.Before = func(context *cli.Context) error {
		//log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
