package ftpConnection

import (
	"log"
	"testing"
)

func TestFtpConnection_Connect(t *testing.T) {
	var fail bool
	client := NewFtp("h4.netangels.ru", "21", "c8555_testakul_atonko_ru", "twBDXHqa9mz7eblE")
	err := client.Connect()
	if err != nil {
		log.Println(err)
		fail = true
	}

	files, err := client.GetRemoteTxtFileNames("www")
	if err != nil {
		log.Println(err)
		fail = true
	}
	if len(files) != 5 {
		log.Println(err)
		fail = true
	}

	client.Disconnect()
	if fail {
		t.Fail()
	}
}
