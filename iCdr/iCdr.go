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

func SendCdr() bool {
	bRet := false
	err := CreateProducer("192.168.112.250:9092")
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	err = PubMsg("test", MakeCdr("f1d979bc-36f7-4e34-b965-fb5f1d21a044",
		"ipron-cloud/flow-service",
		"bridgetec.ipron.tracking.event.cdr",
		"flow/61276d4c7790a01da4bac0fd",
		map[string]interface{}{
			"id":      "abcd",
			"message": "Hello, test!",
		},
		map[string]interface{}{
			"eventSubject": "call/61276d4c7790a01da4bac0fd",
		},
	))
	if err == nil {
		bRet = true
	}

	return bRet
}

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

func PubMsg(topic string, e cloudevents.Event) error {
	// publish sync
	message, _ := e.MarshalJSON()
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
