package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/urfave/cli"
)

type Application struct {
	Name string `json:"name"`
}

type ApplicationMessage struct {
	Applications []Application `json:"application"`
}

type ApplicationsResponse struct {
	Message ApplicationMessage `json:"applications"`
}

const infoUrlAppNameRequiredCode = 19
const infoUrlIdRequiredCode = 18

func main() {

	var eurekaHost string
	var eurekaPort int

	app := cli.NewApp()
	app.Name = "eureka-cli"
	app.Version = "0.1"
	app.Usage = "A command-line interface to perform with netflix eureka"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "u, host",
			Value: "localhost",
			Usage: "IP host adrees of eureka server",
			EnvVar: "EUREKA_SERVER_HOST",
			Destination: &eurekaHost,
		},
		cli.IntFlag{
			Name: "p, port",
			Value: 8761,
			Usage: "Port of eureka server",
			EnvVar: "EUREKA_SERVER_PORT",
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
					Name: "a, app-name",
					Usage: "Application name registered in Eureka server",
					EnvVar: "INFO_SPRING_APPLICATION_NAME",
				},
				cli.StringFlag{
					Name: "i, id",
					Usage: "Instance ID registered in Eureka server",
					EnvVar: "INFO_EUREKA_INSTANCE_INSTANCE_ID",
				},
			},
			Action:func(context *cli.Context) {

				//TODO: EurekaService.GetInstances()
				fmt.Println("Some fuck! App: ", context.String("app-name"))

			},

			// GET URL
			Subcommands: cli.Commands{
				cli.Command{
					Name: "url",
					Description: "Get url of concrete instance",
					Usage: "Get url of concrete instance",
					ArgsUsage: "$APP_NAME $INSTANCECE_ID",
					Action:func(context *cli.Context) error {

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

	fmt.Println("Host: ",eurekaHost)
	fmt.Println("Port: ",eurekaPort)

	return
	// test

	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://localhost:8761/eureka/apps", nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}else {

		bytes, _ := ioutil.ReadAll(response.Body)
		response.Body.Close()

		var appResp ApplicationsResponse

		json.Unmarshal(bytes, &appResp)

		for _, els := range appResp.Message.Applications{

			fmt.Println(els.Name)
		}

	}

}