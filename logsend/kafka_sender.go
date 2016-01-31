package logsend

import (
	"github.com/Shopify/sarama"
	"strconv"
	"strings"
	"time"
	"fmt"
)

var (
	prodChan     = make(chan *LogLine, 0)
	kBrokers     string
	kBatch       int
	kTopic       string
	kBufferTime  int
	kBufferBytes int
	cCapacity    int

	myproducer *KafkaProducer
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
	p, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		fmt.Println("kafka connect fialure ....")
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
	sendCh chan *LogLine
}

//1.初始化配置
//2.监听消息发送通道
func InitKafka(conf map[string]string) error {
	Conf.Logger.Printf("kafka sender setting conifig ...")
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
	var err error
	myproducer, err = NewKafkaProducer(
		strings.Split(kBrokers, ","),
		kTopic, kBufferTime,
		kBufferBytes, kBatch)

	//如果连接失效,直接返回
	if err != nil {
		return err
	}

	go func() {
		//阻塞的方式接收prodChan的消息
		for data := range prodChan {
			myproducer.Write(string(data.Line))
		}
	}()
	return nil
}

//工厂类,生成本Sender
func NewKafkaSender() Sender {
	sender := &KafkaSender{
		sendCh: prodChan,
	}
	return Sender(sender)
}

func (self *KafkaSender) Send(ll *LogLine) {
	prodChan <- ll
}
func (self *KafkaSender) SetConfig(iniConfig map[string]string) error {
	return nil
}
func (self *KafkaSender) Name() string {
	return "kafka"
}
func (self *KafkaSender) Stop() error {
	return myproducer.Close()
}
