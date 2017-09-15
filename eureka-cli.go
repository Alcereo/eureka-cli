package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"github.com/alcereo/eureka-cli/eureka-interface"
	"text/template"
	"time"
)

const timeoutErrorCode = 1
const infoUrlInstanceNotFoundCode = 2
const idEmptyErrorCode = 3
const appNameEmptyErrorCode = 4


var instancesListTemplate = "" +
					"{{- printf \"%-20.20s\" \"APP NAME\" }}{{- printf \"%-10.10s\" \"STATUS\"  }}{{- printf \"%-35.35s\" \"ID\"  }}{{- printf \"%-18.18s\" \"IP ADDRESS\"  }}{{- printf \"%-18.18s\" \"PORT\"  }} \n" +
	"{{range .Items}}{{- printf \"%-20.20s\"  .AppName    }}{{- printf \"%-10.10s\" .Status     }}{{- printf \"%-35.35s\"  .Id    }}{{- printf \"%-18.18s\"  .Ip            }}{{- printf \"%d\"  .Port.Number   }} \n{{end}}"


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
					ArgsUsage:   "$APP_NAME $INSTANCE_ID",
					Action: func(context *cli.Context) error {

						appName := context.Args().Get(0)
						instanceId := context.Args().Get(1)

						if appName == "" {
							cli.ShowCommandHelpAndExit(context, "url", appNameEmptyErrorCode)
						}

						if instanceId == "" {
							cli.ShowCommandHelpAndExit(context, "url", idEmptyErrorCode)
						}

						client := discovery.Client{
							EurekaHost: eurekaHost,
							EurekaPort: eurekaPort,
						}

						instances := client.GetInstanceByAppAndId(appName, instanceId)

						if len(instances) == 0 {
							return cli.NewExitError(
								fmt.Sprintf(
									"Instance with App name: \"%s\", and Id: \"%s\" not found",
									appName,
									instanceId,
								),
								infoUrlInstanceNotFoundCode,
							)
						}

						instance := instances[0]

						fmt.Printf(
							"http://%s:%d",
							instance.Ip,
							instance.Port.Number,
						)

						return nil
					},
				},
			},
		},
		{
			Name: "wait",
			Description: "Wait for UP instance status",
			Usage: "Wait for UP instance status",
			ArgsUsage:   "$APP_NAME $INSTANCE_ID",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name: "t, time",
					Usage: "Time in seconds to wait for",
					Value: 30,
					EnvVar: "EUREKA_WAIT_TIME",
				},
			},
			Action: func(context *cli.Context) error {

				appName := context.Args().Get(0)
				instanceId := context.Args().Get(1)

				if appName == "" {
					cli.ShowCommandHelpAndExit(context, "wait", appNameEmptyErrorCode)
				}

				if instanceId == "" {
					cli.ShowCommandHelpAndExit(context, "wait", idEmptyErrorCode)
				}

				fmt.Printf(
					"Wait for instanceID: \"%s\" app name: \"%s\"...\n",
					instanceId,
					appName,
				)

				timeout := time.After(time.Second * time.Duration(context.Int("time")))
				success := make(chan discovery.Instance)
				start := time.Now()

				client := discovery.Client{
					EurekaHost: eurekaHost,
					EurekaPort: eurekaPort,
				}

				go func() {
					for {
						instances := client.GetInstanceByAppAndId(appName, instanceId)

						if len(instances) != 0 {
							instance := instances[0]

							if instance.Status == "UP" {
								success <- instance
							}
						}

						time.Sleep(1 * time.Second)
					}
				}()

				select {
				case <-timeout:{
					return cli.NewExitError("Wait timout exit", timeoutErrorCode)
				}
				case instance := <-success:{

					fmt.Println("It took: ", time.Now().Sub(start))

					t,_ := template.New("instances").Parse(instancesListTemplate)

					t.Execute(os.Stdout,
							struct {
							Items []discovery.Instance
						}{
								Items:append(make([]discovery.Instance, 0), instance),
							})

					return nil
				}

				}
			},
		},
	}

	app.Run(os.Args)

}
