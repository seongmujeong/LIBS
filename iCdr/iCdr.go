package iCdr

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Cdr struct {
	prod sarama.SyncProducer
}

var pCdr *Cdr

func StCdr() *Cdr {
	if pCdr == nil {
		pCdr = new(Cdr)
	}
	return pCdr
}

/*
	note
		cloudevent Cdr 스펙에 맞게 Data 만들기
	param
		sId				: ID 정보
		sSource			: Source 정보
		sType			: Type 정보
		sSubject		: SubJect 정보
		mData			: Data 정보
		mExt			: 추가 Key-value 정보
	return
		cloudevents.Event	: cloudevents.Event 구조체
*/

func MakeCdr(sId string, sSource string, sType string, sSubject string, mData map[string]interface{}, mExt map[string]interface{}) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID(sId)
	event.SetSource(sSource)
	event.SetType(sType)
	event.SetSubject(sSubject)
	event.SetTime(time.Now())
	event.SetData(cloudevents.ApplicationJSON, mData)
	for mKey, mVal := range mExt {
		event.SetExtension(mKey, mVal)
	}
	return event
}

/*
	note
		Kafka Producer 등록
	param
		kafkaConn		: Kafka 접속 정보(IP:PORT)
						  ex) 192.168.112.250:9092
	return
		error			: error 정보
*/

func CreateProducer(kafkaConn string) error {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)
	var err error
	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// sync producer
	StCdr().prod, err = sarama.NewSyncProducer([]string{kafkaConn}, config)

	return err
}

/*
	note
		Cdr 전송
	param
		topic			: Kafka Topic 정보
		event			: event 정보 구조체
	return
		cloudevents.Event	: cloudevents.Event 구조체
*/

func PubMsg(topic string, event cloudevents.Event) error {
	// publish sync
	message, _ := event.MarshalJSON()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(message[:])),
	}
	p, o, err := StCdr().prod.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
	return err
}
