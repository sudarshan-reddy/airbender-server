package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sudarshan-reddy/mqtt/mq"
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
	mqttClient, err := mq.NewClient(cfg.MQTTClient, cfg.MQTTURL, cfg.MQTTTopic, false)
	failOnError(err, "failed to load client")
	defer mqttClient.Close()
	serv := &server{client: mqttClient}
	responseCh := make(chan string)
	defer close(responseCh)
	go serv.startSubscribing(responseCh)
	router := httprouter.New()
	router.GET(apiVersion+"/currentStatus", currentStatus)
	http.ListenAndServe(cfg.ListenAddr, router)
}

type server struct {
	client *mq.Client
}

func (s *server) startSubscribing(response chan string) {
	for writer := range s.client.Subscribe() {
		writer.WritePayload(s)
	}
}

func (s *server) Write(p []byte) (n int, err error) {
	lastResponse = p
	return 0, nil
}

func currentStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(lastResponse)
}
