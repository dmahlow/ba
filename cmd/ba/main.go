package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"
	bettererrors "github.com/xtuc/better-errors"

	"github.com/bytearena/core/common/utils"

	"github.com/bytearena/ba/subcommand/build"
	"github.com/bytearena/ba/subcommand/generate"
	mapcmd "github.com/bytearena/ba/subcommand/map"
	"github.com/bytearena/ba/subcommand/train"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := makeapp()
	app.Version = utils.GetVersion()
	app.Run(os.Args)
}

func makeapp() *cli.App {
	app := cli.NewApp()
	app.Description = "Byte Arena cli tool"
	app.Name = "Byte Arena cli tool"
	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Build an agent",
			BashComplete: func(c *cli.Context) {
				completion, err := build.BashComplete(c.Args().Get(0))

				if err != nil {
					commandFailWith("build", false, c, err)
				}

				fmt.Fprintln(c.App.Writer, completion)
			},
			Action: func(c *cli.Context) error {
				showUsage, err := build.Main(c.Args().Get(0))

				if err != nil {
					commandFailWith("build", showUsage, c, err)
				}

				return nil
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "Generate a boilerplate agent",
			BashComplete: func(c *cli.Context) {
				completion, err := generate.BashComplete()

				if err != nil {
					commandFailWith("generate", false, c, err)
				}

				fmt.Fprintln(c.App.Writer, completion)
			},
			Action: func(c *cli.Context) error {
				showUsage, err := generate.Main(c.Args().Get(0))

				if err != nil {
					commandFailWith("generate", showUsage, c, err)
				}

				return nil
			},
		},
		{
			Name:    "train",
			Aliases: []string{"t"},
			Usage:   "Train your agent",
			Flags: []cli.Flag{
				cli.IntFlag{Name: "tps", Value: 20, Usage: "Number of ticks per second"},
				cli.StringFlag{Name: "host", Value: "", Usage: "IP serving the trainer; required"},
				cli.StringSliceFlag{Name: "agent", Usage: "Agent images"},
				cli.IntFlag{Name: "port", Value: 8080, Usage: "Port serving the trainer"},
				cli.StringFlag{Name: "viz-host", Value: "127.0.0.1", Usage: "Specify a host for the visualization server"},
				cli.StringFlag{Name: "record-file", Value: "", Usage: "Destination file for recording the game"},
				cli.StringFlag{Name: "map", Value: "hexagon", Usage: "Name of the map used by the trainer"},
				cli.BoolFlag{Name: "no-browser", Usage: "Disable automatic browser opening at start"},
				cli.BoolFlag{Name: "debug", Usage: "Enable debug logging"},
				cli.BoolFlag{Name: "quiet", Usage: "Decrease verbosity of the output"},
				cli.BoolFlag{Name: "profile", Usage: "Enable execution profiling"},
				cli.BoolFlag{Name: "dump-raw-comm", Usage: "Dump all the communication between the agent and the server"},
				cli.IntFlag{Name: "duration", Usage: "If set, game will stop after this durarion (in seconds)"},
			},
			Action: func(c *cli.Context) error {
				tps := c.Int("tps")
				host := c.String("host")
				agents := c.StringSlice("agent")
				port := c.Int("port")
				vizhost := c.String("viz-host")
				recordFile := c.String("record-file")
				mapName := c.String("map")
				nobrowser := c.Bool("no-browser")
				isDebug := c.Bool("debug")
				isQuiet := c.Bool("quiet")
				shouldProfile := c.Bool("profile")
				dumpRaw := c.Bool("dump-raw-comm")
				duration := c.Int("duration")

				showUsage, err := train.TrainAction(
					tps,
					host,
					port,
					vizhost,
					nobrowser,
					recordFile,
					agents,
					isDebug,
					isQuiet,
					mapName,
					shouldProfile,
					dumpRaw,
					duration,
				)

				if err != nil {
					commandFailWith("train", showUsage, c, err)
				}

				return nil
			},
		},
		{
			Name:    "map",
			Aliases: []string{},
			Usage:   "Operations on map packs",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "debug", Usage: "Enable debug logging"},
			},
			Subcommands: []cli.Command{
				{
					Name:  "update",
					Usage: "Fetch the trainer maps if needed",
					Action: func(c *cli.Context) error {
						isDebug := c.Bool("debug")

						debug := func(str string) {}

						if isDebug {
							debug = func(str string) {
								fmt.Println(str)
							}
						}

						mapcmd.MapUpdateAction(debug)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "List the trainer maps locally available",
					Action: func(c *cli.Context) error {
						mapcmd.MapListAction()
						return nil
					},
				},
			},
		},
	}

	return app
}

func commandFailWith(name string, showUsage bool, c *cli.Context, err error) {
	berror := bettererrors.
		New("Failed to execute command").
		SetContext("command", name).
		With(err)

	if showUsage {
		cli.ShowCommandHelp(c, name)
	}

	utils.FailWith(berror)
}
