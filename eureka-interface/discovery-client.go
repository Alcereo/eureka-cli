package discovery

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type Instance struct {
	Id string `json:"instanceId"`
	AppName string `json:"app"`
	Ip string `json:"ipAddr"`
	Port struct{
		Number int `json:"$"`
	} `json:"port"`
	Status string `json:"status"`
}

type Application struct {
	Name string `json:"name"`
	Instances []Instance `json:"instance"`
}

type SingleInstanceApplication struct {
	Instance Instance `json:"instance"`
}

type ApplicationMessage struct {
	Applications []Application `json:"application"`
}

type SingleApplicationMessage struct {
	Application Application `json:"application"`
}

type ApplicationsResponse struct {
	Message ApplicationMessage `json:"applications"`
}

type Client struct {
	EurekaHost string
	EurekaPort int
	URL string
}

func (client *Client) GetInstanceByAppAndId(appName string, instanceId string) []Instance {

	client.URL = fmt.Sprintf(
		"http://%s:%d/eureka/apps/%s/%s",
		client.EurekaHost,
		client.EurekaPort,
		appName,
		instanceId,
	)

	response := client.requestToEureka(true)

	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if len(bytes)==0{

		return []Instance{}
	}else {
		var appResp SingleInstanceApplication

		json.Unmarshal(bytes, &appResp)

		instanceArr := []Instance{}

		instanceArr = append(instanceArr, appResp.Instance)

		return instanceArr
	}

}

func (client *Client) GetInstanceById(instanceId string) []Instance {

	client.URL = fmt.Sprintf(
		"http://%s:%d/eureka/instances/%s",
		client.EurekaHost,
		client.EurekaPort,
		instanceId,
	)

	response := client.requestToEureka(true)

	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	var appResp SingleInstanceApplication

	json.Unmarshal(bytes, &appResp)

	instanceArr := []Instance{}

	instanceArr = append(instanceArr, appResp.Instance)

	return instanceArr
}

func (client *Client) GetInstancesByApp(appName string) []Instance {

	client.URL = fmt.Sprintf(
		"http://%s:%d/eureka/apps/%s",
		client.EurekaHost,
		client.EurekaPort,
		appName,
	)

	response := client.requestToEureka(true)

	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	var appResp SingleApplicationMessage

	json.Unmarshal(bytes, &appResp)

	instanceArr := []Instance{}

	for _, instance := range appResp.Application.Instances {
		instanceArr = append(instanceArr, instance)
	}

	return instanceArr

}

func (client *Client) GetInstances() []Instance {

	client.URL = fmt.Sprintf(
		"http://%s:%d/eureka/apps",
		client.EurekaHost,
		client.EurekaPort,
	)

	response := client.requestToEureka(true)

	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	var appResp ApplicationsResponse

	json.Unmarshal(bytes, &appResp)

	instanceArr := []Instance{}

	for _, application := range appResp.Message.Applications {
		for _, instance := range application.Instances {
			instanceArr = append(instanceArr, instance)
		}
	}

	return instanceArr
}


func (client Client) requestToEureka(mayNotFound bool) *http.Response {

	httpClient := &http.Client{}
	url := client.URL

	request, err := http.NewRequest(
		"GET",
		url,
		nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)


	if err != nil {
		if (response == nil){

			url = fmt.Sprintf(
				"http://%s:%d/",
				client.EurekaHost,
				client.EurekaPort,
			)

			fmt.Fprintln(os.Stderr, "   Cant connect to Eureka server at: ",url,", check connection\n" +
					"   or enter other address by: -u -p global flags: \n" +
					"       eureka-cli -u $HOST -p $PORT command")

			os.Exit(1)
		}
	}

	if response.StatusCode == 404 && mayNotFound {
		return response
	}

	if response.StatusCode != 200 {
		log.Fatalln("Request to ",url," return code: ", response.StatusCode)
	}


	return response
}