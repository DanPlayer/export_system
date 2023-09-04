package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var Client = NewRabbitMQ()

// RabbitMQ rabbitMQ结构体
type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	consumer map[string]<-chan amqp.Delivery
}

type Options struct {
	UserName string
	Password string
	Host     string
	Vhost    string
}

// NewRabbitMQ 创建简单模式下RabbitMQ实例
func NewRabbitMQ() *RabbitMQ {
	options := Options{
		UserName: "guest",
		Password: "guest",
		Host:     "127.0.0.1",
		Vhost:    "/",
	}
	// 创建RabbitMQ实例
	rabbitmq := &RabbitMQ{}
	var err error
	// 获取connection
	rabbitmq.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/%s", options.UserName, options.Password, options.Host, options.Vhost))
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	// 获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")

	rabbitmq.consumer = make(map[string]<-chan amqp.Delivery, 0)

	return rabbitmq
}

// Create 创建队列
func (r *RabbitMQ) Create(key string) error {
	exchange := fmt.Sprintf("%s-exchange", key)
	queue := fmt.Sprintf("%s-queue", key)
	bindKey := fmt.Sprintf("%s-bindkey", key)

	// 申请交换机
	err := r.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 申请队列
	q, err := r.channel.QueueDeclare(
		queue,
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// 绑定队列与交换机
	err = r.channel.QueueBind(
		q.Name,
		bindKey,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	// 创建消费者
	consumer, err := r.channel.Consume(
		queue, // queue
		// 用来区分多个消费者
		"", // consumer
		// 是否自动应答
		true, // auto-ack
		// 是否独有
		false, // exclusive
		// 设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		// 列是否阻塞
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println(err)
	}

	r.consumer[key] = consumer

	return nil
}

func (r *RabbitMQ) Push(key, data string) error {
	exchange := fmt.Sprintf("%s-exchange", key)
	bindKey := fmt.Sprintf("%s-bindkey", key)

	err := r.channel.Publish(
		exchange,
		bindKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	return err
}

// Pop 消费队列
func (r *RabbitMQ) Pop(key string) <-chan string {
	//接收消息
	get := <-r.consumer[key]
	item := string(get.Body)

	list := make(chan string)
	go func() {
		list <- item
	}()
	return list
}

// Destroy 断开channel 和 connection
func (r *RabbitMQ) Destroy(key string) error {
	exchange := fmt.Sprintf("%s-exchange", key)
	queue := fmt.Sprintf("%s-queue", key)

	_, err := r.channel.QueueDelete(queue, true, false, true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = r.channel.ExchangeDelete(exchange, true, false)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RabbitMQ) Close() error {
	err := r.channel.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return r.conn.Close()
}

// failOnErr 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}
