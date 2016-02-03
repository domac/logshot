package logsend

import (
	"github.com/Shopify/sarama"
	"github.com/juju/errors"
	"strconv"
	"strings"
	"study2016/logshot/logger"
	"time"
)

var (
	kBrokers     string
	kBatch       int
	kTopic       string
	kBufferTime  int
	kBufferBytes int
)

func init() {
	RegisterNewSender("kafka", InitKafka, NewKafkaSender)
}

type KafkaProducer struct {
	producer sarama.AsyncProducer
	topic    string
}

//构造生产者
func NewKafkaProducer(brokers []string, topic string, bufferTime, bufferBytes, batchSz int) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.Flush.Bytes = bufferBytes
	config.Producer.Flush.Frequency = time.Duration(bufferTime * 1000000)
	config.Producer.Flush.Messages = batchSz

	//设置超时
	config.Net.DialTimeout = time.Second * 15
	config.Net.WriteTimeout = time.Second * 15

	p, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		logger.GetLogger().Errorln(err)
		return nil, err
	}
	k := &KafkaProducer{
		producer: p,
		topic:    topic,
	}
	return k, nil
}

//生产数据写入
func (k *KafkaProducer) Write(s string) {
	k.producer.Input() <- &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(s),
	}
}

//关闭生产者连接
func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}

//output to kafka
type KafkaSender struct {
	sendCh   chan *LogLine
	producer *KafkaProducer
}

//1.初始化配置
//2.监听消息发送通道
func InitKafka(conf map[string]string, sender Sender) error {

	logger.GetLogger().Infoln("init kafka sender")

	//变量初始化
	if val, ok := conf["kafkaBatch"]; ok {
		kBatch, _ = strconv.Atoi(val)
	}
	if val, ok := conf["kafkaBrokers"]; ok {
		kBrokers = val
	}
	if val, ok := conf["kafkaTopic"]; ok {
		kTopic = val
	}
	if val, ok := conf["kafkaBufferTime"]; ok {
		kBufferTime, _ = strconv.Atoi(val)
	}
	if val, ok := conf["kafkaBufferBytes"]; ok {
		kBufferBytes, _ = strconv.Atoi(val)
	}
	//创建kafka的生产者
	myproducer, err := NewKafkaProducer(
		strings.Split(kBrokers, ","),
		kTopic, kBufferTime,
		kBufferBytes, kBatch)

	//如果连接失效,直接返回
	if err != nil {
		return err
	}

	err = sender.SetConfig(myproducer)
	if err != nil {
		return err
	}
	sender.Receive()
	return nil
}

//工厂类,生成本Sender
func NewKafkaSender() Sender {
	sender := &KafkaSender{
		sendCh: make(chan *LogLine, 0),
	}
	return Sender(sender)
}

func (self *KafkaSender) SetConfig(obj interface{}) error {
	switch obj := obj.(type) {
	case *KafkaProducer:
		self.producer = obj
	default:
		return errors.New("kafka setconfig error ")
	}
	return nil
}
func (self *KafkaSender) Name() string {
	return "kafka"
}
func (self *KafkaSender) Stop() error {
	logger.GetLogger().Infoln("kafka sender stop")
	close(self.sendCh)
	return self.producer.Close()
}

func (self *KafkaSender) Receive() {
	go func() {
		//阻塞的方式接收prodChan的消息
		for data := range self.sendCh {
			self.producer.Write(string(data.Line))
		}
	}()
}

func (self *KafkaSender) Send(ll *LogLine) {
	self.sendCh <- ll
}
