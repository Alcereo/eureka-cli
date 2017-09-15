package discovery

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
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
}

func (client Client) GetInstanceByAppAndId(appName string, instanceId string) []Instance {

	httpClient := &http.Client{}

	url := fmt.Sprintf(
		"http://%s:%d/eureka/apps/%s/%s",
		client.EurekaHost,
		client.EurekaPort,
		appName,
		instanceId,
	)

	request, err := http.NewRequest(
		"GET",
		url,
		nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)

	if err != nil {
		log.Fatal(err)
		return nil
	}

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

func (client Client) GetInstanceById(instanceId string) []Instance {
	httpClient := &http.Client{}

	url := fmt.Sprintf(
		"http://%s:%d/eureka/instances/%s",
		client.EurekaHost,
		client.EurekaPort,
		instanceId,
	)

	request, err := http.NewRequest(
		"GET",
		url,
		nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	bytes, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	var appResp SingleInstanceApplication

	json.Unmarshal(bytes, &appResp)

	instanceArr := []Instance{}

	instanceArr = append(instanceArr, appResp.Instance)

	return instanceArr
}

func (client Client) GetInstancesByApp(appName string) []Instance {

	httpClient := &http.Client{}

	url := fmt.Sprintf(
		"http://%s:%d/eureka/apps/%s",
		client.EurekaHost,
		client.EurekaPort,
		appName,
	)

	request, err := http.NewRequest(
		"GET",
		url,
		nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)

	if err != nil {
		log.Fatal(err)
		return nil
	}

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

func (client Client) GetInstances() []Instance {

	httpClient := &http.Client{}

	url := fmt.Sprintf(
		"http://%s:%d/eureka/apps",
		client.EurekaHost,
		client.EurekaPort,
	)

	request, err := http.NewRequest(
		"GET",
		url,
		nil)

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)

	if err != nil {
		log.Fatal(err)
		return nil
	}

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