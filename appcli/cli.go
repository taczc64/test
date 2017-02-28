package cliTest

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
)

func TestCLi() {
	os.Args = []string{"greet", "--name", "tac"}
	app := cli.NewApp()
	app.Name = "greet"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "waihao", Value: "bob", Usage: "a name to say"},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Printf("hello %v\n", c.String("waihao"))
		return nil
	}
	app.UsageText = "app [first_arg] [second_arg]"
	app.Author = "GavinTang"
	app.Email = "gavin.tang@btcc.co"
	app.Authors = []cli.Author{{Name: "Oliver allen", Email: "oliver@shop.com"}}
	app.Run(os.Args)
}
