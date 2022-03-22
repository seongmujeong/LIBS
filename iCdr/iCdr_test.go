package iCdr

import (
	"testing"
)

func TestSendCdr(t *testing.T) {
	bRet := SendCdr()

	if bRet != true {
		t.Error("Send Cdr Fail")
	}
}
