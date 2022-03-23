package iCdr

import (
	"fmt"
	"os"
	"testing"
)

func TestSendCdr(t *testing.T) {

	err := CreateProducer("192.168.112.250:9092")
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	err = PubMsg("test", MakeCdr("sId",
		"sSource",
		"sType",
		"sSubject",
		map[string]interface{}{
			"id":      "abcd",
			"message": "Hello, test!",
		},
		map[string]interface{}{
			"ext": "ext",
			"int": 123,
		},
	))

	if err != nil {
		t.Error("Send Cdr Fail")
	}
}
