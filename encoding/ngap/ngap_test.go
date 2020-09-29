package ngap

import (
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"
)

// send message
var TestNGSetupRequest string = "00150028000003001b00080002f839000000040066001000000000010002f839000010080102030015400100"
var TestInitialUEMessage string = "000f40470000040055000200000026001d1c7e004179000d0102f8392143000010325476981001202e0480a000000079000f4002f839000004001002f839000001005a4001180070400100"
var TestULAuthenticationResponse string = "002e403c000004000a0002000100550002000000260016157e00572d10803adcacc364fc000bdc0f65e324eaa10079000f4002f839000004001002f839000001"
var TestULSecurityModeComplete string = "002e405c000004000a0002000100550002000000260036357e04a860200b007e005e7700090500000001000001f171001c7e004179000d0102f8392143000010325476981001202e0480a000000079000f4002f839000004001002f839000001"
var TestInitialContextSetupResponse string = "200e000f000002000a00020001005500020000"
var TestULRegistrationComplete string = "002e4031000004000a000200010055000200000026000b0a7e042cbd08cf017e00430079000f4002f839000004001002f839000001"
var TestPDUSessionResourceSetupResponse string = "201d0024000003000a00020001005500020000004b40110000010d0003e0c0a80103000003e70000"

// receive message
var TestNGSetupResponse string = "20150031000004000100050100414d4600600008000002f839cafe0000564001ff005000100002f839000110080102031008112233"
var TestDLAuthenticationRequest string = "0004403e000003000a000200010055000200000026002b2a7e00560002000021fc64081953bb33c0682edf1690b25821201094bbaf40940a8000c6a72c4efbaf0337"
var TestDLSecurityModeCommand string = "00044029000003000a0002000100550002000000260016157e036c2b24e2007e005d02000480a00000e1360100"
var TestInitialContextSetupRequest string = "000e0080a7000009000a00020001005500020000001c00070002f839cafe000000000a2201010203100811223300770009000004000000000000005e002013663ab7286c9a6af7cba0b1fd9e6ed48045d4356d46ff3944c81c63324fd803002440040002f839002240080000000100ffff0100264036357e02930d75cf017e0242010177000b0202f839cafe000000000154070002f839000001150a040101020304011122335e010616012c"
var TestInitialContextSetupRequest2 string = "000e0080f500000b000a00020001005500020000006e0008080f4240200f4240001c00070002f839cafe000047002a000001402001020321000003008b000a01f07f00000800000001008600010000880007000000000938000000000a2201010203100811223300770009000000100000000000005e0020473007e30d4d0d77a7073e5b43b909562b7a8c461fc7ef0b73ab4026edbb91aa002440040002f839002240080000000100ffff010026404a497e02809e40eb027e006801003a2e0101c211000901000631310101ff00060103e80103e859322905013c3c0001220401010203790006002041010109250908696e7465726e65741201"
var TestDLPDUSessionEstablishmentAccept string = "001d006d000003000a00020001005500020000004a005a0040012f7e0222994e9f027e00680100202e0100c21100090100063131010100000601e80301e80359322905013c3c00011201402001020321000003008b000a01f0c0a801120000000100860001000088000700000000093800"

func recvfromNW(gnb *GNB, msg string) {
	in, _ := hex.DecodeString(msg)
	gnb.Decode(&in)
	fmt.Printf("")
}

func TestMakePDUSessionResourceSetupResponse(t *testing.T) {
	log.SetPrefix("[NGAP TEST]")

	pattern := []struct {
		in_str string
	}{
		{TestNGSetupResponse},
		{TestDLAuthenticationRequest},
		{TestDLSecurityModeCommand},
		{TestInitialContextSetupRequest},
		{TestDLPDUSessionEstablishmentAccept},
	}

	gnb := NewNGAP("ngap_test.json")
	gnb.UE.PowerON()

	for _, p := range pattern {
		recvfromNW(gnb, p.in_str)
	}

	gnb.SetDebugLevel(1)
	gnb.UE.SetDebugLevel(1)
	v := gnb.MakePDUSessionResourceSetupResponse()
	expect_str := TestPDUSessionResourceSetupResponse
	expect, _ := hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		//fmt.Printf("PDUSessionResourceSetupResponse test not implemented yet.\n")
		t.Errorf("PDUSessionResourceSetupResponse\nexpect: %x\nactual: %x", expect, v)
	}

}

func TestInitialContestSetupResponse(t *testing.T) {
	gnb := NewNGAP("ngap_test.json")
	gnb.UE.PowerON()

	recvfromNW(gnb, TestDLAuthenticationRequest)
	v := gnb.MakeInitialContextSetupResponse()
	expect_str := TestInitialContextSetupResponse
	expect, _ := hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		t.Errorf("InitialContextSetupResponse\nexpect: %x\nactual: %x", expect, v)
	}
}

func TestMakeInitialUEMessage(t *testing.T) {
	gnb := NewNGAP("ngap_test.json")
	gnb.UE.PowerON()

	pdu := gnb.UE.MakeRegistrationRequest()
	gnb.RecvfromUE(&pdu)
	v := gnb.MakeInitialUEMessage()
	expect_str := TestInitialUEMessage
	expect, _ := hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		t.Errorf("InitialUEMessage\nexpect: %x\nactual: %x", expect, v)
	}
}

func TestMakeUplinkNASTransport(t *testing.T) {
	gnb := NewNGAP("ngap_test.json")
	gnb.UE.PowerON()

	var expect_str string
	var expect []byte

	recvfromNW(gnb, TestDLAuthenticationRequest)
	pdu := gnb.UE.MakeAuthenticationResponse()
	gnb.RecvfromUE(&pdu)

	v := gnb.MakeUplinkNASTransport()
	expect_str = TestULAuthenticationResponse
	expect, _ = hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		t.Errorf("UplinkNASTransport1\nexpect: %x\nactual: %x", expect, v)
	}

	recvfromNW(gnb, TestDLSecurityModeCommand)
	pdu = gnb.UE.MakeSecurityModeComplete()
	gnb.RecvfromUE(&pdu)

	v = gnb.MakeUplinkNASTransport()
	expect_str = TestULSecurityModeComplete
	expect, _ = hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		t.Errorf("UplinkNASTransport2\nexpect: %x\nactual: %x", expect, v)
	}

	recvfromNW(gnb, TestInitialContextSetupRequest)
	pdu = gnb.UE.MakeRegistrationComplete()
	gnb.RecvfromUE(&pdu)

	v = gnb.MakeUplinkNASTransport()
	expect_str = TestULRegistrationComplete
	expect, _ = hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		t.Errorf("UplinkNASTransport3\nexpect: %x\nactual: %x", expect, v)
	}
}

func TestMakeNGSetupRequest(t *testing.T) {
	gnb := NewNGAP("ngap_test.json")
	v := gnb.MakeNGSetupRequest()
	expect_str := TestNGSetupRequest
	expect, _ := hex.DecodeString(expect_str)
	if reflect.DeepEqual(expect, v) == false {
		t.Errorf("NGSetupRequest\nexpect: %x\nactual: %x", expect, v)
	}
}

func TestDecode(t *testing.T) {

	pattern := []struct {
		in_str string
		desc   string
	}{
		{TestNGSetupResponse,
			"NG Setup Response"},
		{TestDLAuthenticationRequest,
			"DL Authentication Request"},
		{TestDLSecurityModeCommand,
			"DL Security Mode Command"},
		{TestInitialContextSetupRequest,
			"Initial Context Setup Request"},
		{TestInitialContextSetupRequest2,
			"Initial Context Setup Request #2"},
		{TestDLPDUSessionEstablishmentAccept,
			"PDU Session Establishment Accept"},
	}

	gnb := NewNGAP("ngap_test.json")

	for _, p := range pattern {
		fmt.Printf("---------- test decode: %s\n", p.desc)

		gnb.SetDebugLevel(1)
		gnb.UE.SetDebugLevel(1)
		recvfromNW(gnb, p.in_str)

		if gnb.DecodeError != nil {
			t.Errorf("%s: %v", p.desc, gnb.DecodeError)
		}
	}
}
