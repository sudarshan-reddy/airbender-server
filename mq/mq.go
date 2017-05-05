package mq

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client is a very light wrapper on mqtt
type Client struct {
	Client mqtt.Client
	Topic  string
}

// NewClient returns a new instance of Client
func NewClient(clientID, raw, topic string) (*Client, error) {
	uri, _ := url.Parse(raw)
	server := (fmt.Sprintf("tcp://%s", uri.Host))
	username := uri.User.Username()
	password, _ := uri.User.Password()

	connOpts := mqtt.NewClientOptions().AddBroker(server).SetClientID(clientID).SetCleanSession(true)
	connOpts.SetUsername(username)
	connOpts.SetPassword(password)
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)
	connOpts.SetMaxReconnectInterval(1 * time.Second)
	connOpts.SetKeepAlive(30 * time.Second)
	client := &Client{Client: mqtt.NewClient(connOpts), Topic: topic}

	if token := client.Client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

// Subscribe subscribes message to the predefined topic
func (m *Client) Subscribe(onMessageReceived func(mqtt.Client, mqtt.Message)) error {
	fmt.Println("subscribe", m.Topic)
	if token := m.Client.Subscribe(m.Topic, byte(0), onMessageReceived); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Close disconnects
func (m *Client) Close() {
	m.Client.Disconnect(250)
}
