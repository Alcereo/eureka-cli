package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"github.com/alcereo/eureka-cli/eureka-interface"
	"text/template"
)

const infoUrlAppNameRequiredCode = 19
const infoUrlIdRequiredCode = 18

var instancesListTemplate = "" +
					"{{- printf \"%-20.20s\" \"APP NAME\" }}{{- printf \"%-10.10s\" \"STATUS\"  }}{{- printf \"%-18.18s\" \"ID\"  }}{{- printf \"%-18.18s\" \"IP ADDRESS\"  }}{{- printf \"%-18.18s\" \"PORT\"  }} \n" +
	"{{range .Items}}{{- printf \"%-20.20s\"  .AppName    }}{{- printf \"%-10.10s\" .Status     }}{{- printf \"%-18.18s\"  .Id    }}{{- printf \"%-18.18s\"  .Ip            }}{{- printf \"%d\"  .Port.Number   }} \n{{end}}"

func main() {

	var eurekaHost string
	var eurekaPort int

	app := cli.NewApp()
	app.Name = "eureka-cli"
	app.Version = "0.1"
	app.Usage = "A command-line interface to perform with netflix eureka"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "u, host",
			Value:       "localhost",
			Usage:       "IP host adrees of eureka server",
			EnvVar:      "EUREKA_SERVER_HOST",
			Destination: &eurekaHost,
		},
		cli.IntFlag{
			Name:        "p, port",
			Value:       8761,
			Usage:       "Port of eureka server",
			EnvVar:      "EUREKA_SERVER_PORT",
			Destination: &eurekaPort,
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "info",
			Description: "Query info about instances \n\n" +
				"Sample:\n    eureka info -a $APP_NAME -i $INSTANCE_ID",
			Usage: "Query info about instances",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "a, app-name",
					Usage:  "Application name registered in Eureka server",
					EnvVar: "INFO_SPRING_APPLICATION_NAME",
				},
				cli.StringFlag{
					Name:   "i, id",
					Usage:  "Instance ID registered in Eureka server",
					EnvVar: "INFO_EUREKA_INSTANCE_INSTANCE_ID",
				},
			},
			Action: func(context *cli.Context) {

				appName := context.String("app-name")
				instanceId := context.String("id")

				client := discovery.Client{
					EurekaHost: eurekaHost,
					EurekaPort: eurekaPort,
				}

				var instances []discovery.Instance

				switch {
				case appName == "" && instanceId == "":
					{
						instances = client.GetInstances()
					}
				case appName == "":
					{
						instances = client.GetInstanceById(instanceId)
					}
				case instanceId == "":
					{
						instances = client.GetInstancesByApp(appName)
					}
				default:
					{
						instances = client.GetInstanceByAppAndId(appName, instanceId)
					}
				}

				t,_ := template.New("instances").Parse(instancesListTemplate)

				t.Execute(os.Stdout,
					struct {
						Items []discovery.Instance
						}{Items:instances})

			},

			// GET URL
			Subcommands: cli.Commands{
				cli.Command{
					Name:        "url",
					Description: "Get url of concrete instance",
					Usage:       "Get url of concrete instance",
					ArgsUsage:   "$APP_NAME $INSTANCECE_ID",
					Action: func(context *cli.Context) error {

						if context.Args().Get(0) == "" {
							cli.ShowCommandHelpAndExit(context, "url", infoUrlAppNameRequiredCode)
						}

						if context.Args().Get(1) == "" {
							cli.ShowCommandHelpAndExit(context, "url", infoUrlIdRequiredCode)
						}

						//TODO: EurekaService.getInstance() .Host + . Port
						fmt.Println("APP:", context.String("app-name"))
						fmt.Println("Some info url")

						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
