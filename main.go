package main

import (
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/julienschmidt/httprouter"
	"github.com/sudarshan-reddy/airbender-server/mq"
)

const (
	apiVersion = "/v1"
)

var lastResponse []byte

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}

func main() {
	cfg, err := loadConfigs()
	failOnError(err, "failed to load configs")
	mqttClient, err := mq.NewClient(cfg.MQTTClient, cfg.MQTTURL, cfg.MQTTTopic)
	failOnError(err, "failed to load client")
	defer mqttClient.Close()

	serv := &server{client: mqttClient}
	responseCh := make(chan string)
	defer close(responseCh)
	serv.startSubscribing(responseCh)
	go serv.handleResponses(responseCh)
	router := httprouter.New()
	router.GET(apiVersion+"/currentStatus", currentStatus)
	http.ListenAndServe(":8080", router)
}

type server struct {
	client *mq.Client
}

func (s *server) startSubscribing(response chan string) {
	s.client.Subscribe(func(client mqtt.Client, message mqtt.Message) {
		resp := message.Payload()
		if resp != nil {
			response <- string(resp)
		}
	})
}

func (s *server) handleResponses(responseCh chan string) {
	for response := range responseCh {
		//Write to Database here
		lastResponse = []byte(response)
		fmt.Println(response)
	}
}

func currentStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Write(lastResponse)
}
