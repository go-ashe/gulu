package mqttx

import (
	"context"
	"log"
	"slices"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ashe/gulu/logx"
)

type Config struct {
	Broker    string `json:"broker" yaml:"broker"`
	ClientID  string `json:"client_id" yaml:"client_id"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	QoS       byte   `json:"qos" yaml:"qos"`
	Clean     bool   `json:"clean" yaml:"clean"`
	KeepAlive int    `json:"keep_alive" yaml:"keep_alive"`
	Timeout   int    `json:"timeout" yaml:"timeout"`

	// PublishLogTopics 指定打印发送日志的 topic 最后一段，按名称精确匹配。
	// 例如配置 command 可匹配 device/{serialNumber}/command。
	// nil 或空列表关闭发送日志，发送失败日志也受此配置影响。
	PublishLogTopics []string `json:"publish_log_topics" yaml:"publish_log_topics"`
}

type Server struct {
	config   *Config
	client   mqtt.Client
	subs     map[string]mqtt.MessageHandler // topic -> 回调
	isClosed bool
}

func NewServer(cfg *Config) (c *Server, err error) {
	c = &Server{
		config: cfg,
		subs:   make(map[string]mqtt.MessageHandler),
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.config.Broker)
	opts.SetClientID(c.config.ClientID)
	opts.SetUsername(c.config.Username)
	opts.SetPassword(c.config.Password)
	opts.SetCleanSession(c.config.Clean)
	opts.SetKeepAlive(time.Duration(c.config.KeepAlive) * time.Second)
	opts.SetConnectTimeout(time.Duration(c.config.Timeout) * time.Second)
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(client mqtt.Client) {
		log.Printf("[MQTT] Connected to %s", c.config.Broker)
		c.resubscribeAll()
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v", err)
	}

	c.client = mqtt.NewClient(opts)

	err = c.Start(context.Background())
	return c, err
}

// resubscribeAll 重连后恢复订阅
func (s *Server) resubscribeAll() {
	for topic, cb := range s.subs {
		if token := s.client.Subscribe(topic, s.config.QoS, cb); token.Wait() && token.Error() != nil {
			log.Printf("[MQTT] Failed to resubscribe %s: %v", topic, token.Error())
		} else {
			log.Printf("[MQTT] Resubscribed %s", topic)
		}
	}
}

// Publish 发布消息
func (s *Server) Publish(topic string, payload any) (err error) {
	defer func() {
		topicName := topic[strings.LastIndexByte(topic, '/')+1:]
		if slices.Contains(s.config.PublishLogTopics, topicName) {
			logx.Info("===[MQTT] 发送消息", "topic", topic, "data", payload, "error", err)
		}
	}()

	if token := s.client.Publish(topic, s.config.QoS, false, payload); token.Wait() && token.Error() != nil {
		err = token.Error()
	}
	return
}

// Subscribe 单个订阅
func (s *Server) Subscribe(topic string, callback mqtt.MessageHandler) error {
	if token := s.client.Subscribe(topic, s.config.QoS, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	s.subs[topic] = callback
	return nil
}

// Start 启动服务器
func (s *Server) Start(_ context.Context) error {
	if token := s.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Server) Stop(_ context.Context) {
	if !s.isClosed {
		s.isClosed = true
		if s.client != nil && s.client.IsConnected() {
			s.client.Disconnect(250)
			log.Println("[MQTT] Disconnected")
		}
	}
}
