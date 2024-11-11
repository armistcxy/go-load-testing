package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/armistcxy/go-load-testing/internal/attacker"
	"github.com/urfave/cli/v2"
)

func main() {
	attackCommand := &cli.Command{
		Name:    "attack",
		Aliases: []string{},
		Usage:   "Start load testing",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Link to config file (json format)",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "frequency",
				Aliases:  []string{"f", "freq"},
				Usage:    "Attack frequency",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "per",
				Usage:    "Interval that frequency is applied (often 1 second)",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "duration",
				Aliases:  []string{"d", "dur"},
				Usage:    "Duration (how long this attack phase will run)",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			path := ctx.String("path")
			freq := ctx.Int("frequency")
			per := ctx.Int("per")
			duration := ctx.Int("duration")

			fig, err := attacker.RetrieveFigure(path)
			if err != nil {
				return fmt.Errorf("failed to retrive figure: %s", err)
			}

			builder := attacker.NewCustomTargetBuilder(fig)
			tgt, err := builder.BuildCustomTargeter()
			if err != nil {
				return fmt.Errorf("failed to build targeter: %s", err)
			}

			atk := attacker.NewAttacker(tgt, "results.bin")
			atk.Attack(freq, time.Duration(per)*time.Second, time.Duration(duration)*time.Second)

			return nil
		},
	}

	plotCommand := &cli.Command{
		Name:    "plot",
		Aliases: []string{"display"},
		Usage:   "Plot/Display a graph of request results",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Link to result file (binary format)",
				Required: false, // default is "results.bin"
			},
		},
		Action: func(ctx *cli.Context) error {
			binaryResultPath := ctx.String("path")
			if binaryResultPath == "" {
				binaryResultPath = "results.bin"
			}
			createHTMLFileCmd := exec.Command("sh", "-c", fmt.Sprintf("cat %s | vegeta plot > plot.html", binaryResultPath))
			if err := createHTMLFileCmd.Run(); err != nil {
				return err
			}

			return openBrowser("plot.html")
		},
	}

	figureCommand := &cli.Command{
		Name:    "figure",
		Aliases: []string{"fig"},
		Usage:   "Set up path for figure of HTTP request",
	}

	app := &cli.App{
		Name:  "goltest",
		Usage: "Load testing with custom payload",
		Commands: []*cli.Command{
			attackCommand,
			plotCommand,
			figureCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler")
	case "darwin": // macOS
		cmd = "open"
	default:
		return fmt.Errorf("unsupported platform")
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
