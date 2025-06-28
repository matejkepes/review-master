package smsgateway

import (
	"log"
	"testing"

	"send_sms/config"
)

func TestSend(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	// resp, sendEmail := Send(config.GatewayAddress, config.GatewayPort, config.GatewayPassword, config.GatewaySocketTimeout, "447123456789", "test")
	resp, sendEmail := Send(config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort, config.Gateways[0].GatewayPassword, config.Gateways[0].GatewaySocketTimeout, "447123456789", "test")
	log.Printf("resp: %s, sendEmail: %t\n", resp, sendEmail)
}

func TestLoginSuccessful1(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	success := LoginSuccessful("{\"method_reply\": \"authentication\", \"reply\": \"ok\",  \"server_password\":\"admin\", \"client_id\":\"id1\"}", config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if !success {
		t.Fatal("gateway login response failed")
	}
}

func TestLoginSuccessful2(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	success := LoginSuccessful("{\"notification\": \"replacing old connection IP - 35.176.91.183\"}", config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if !success {
		t.Fatal("gateway login response failed")
	}
}

func TestLoginSuccessfulNot(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	success := LoginSuccessful("{\"server_password\": \"admin\", \"reply\": \"error\", error_code: \"err-201\", \"client_id\": \"id1\", \"method_reply\": \"authentication\"}", config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if success {
		t.Fatal("gateway login response successful should have failed")
	}
}

func TestSendMsgSuccessful1(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	success := SendMsgSuccessful(
		"{\"client_id\": \"id1\", \"unicode\": \"0\", \"msg\": \"6B656C6C6F\", \"reply\": \"proceeding\", \"number\": \"00447123456789\", \"msg_id\": 3432}",
		config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if !success {
		t.Fatal("gateway send message response failed")
	}
}

func TestSendMsgSuccessful2(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	success := SendMsgSuccessful(
		"{\"number\": \"00447123456789\", \"retry_num\": \"0\", \"unicode\": \"0\", \"sending_time\": \"00:00:02\", \"sms_reference_list\": [4219, 4220], \"sim_data\": \"1,3312,10000,10000,10000\", \"validity_int\": 1, \"pdu_seq\": 2, \"msg\": \"5374696C6C206E6F7420626F6F6B696E67206279204150503F20547261636B2063617220616E642063686F6F736520746F207061792063617368206F7220434152442061707073656E642E6D652F52696465202D202031312F31302F323031392061742031313A30312066726F6D20204348414E5445524C414E4453204352454D41544F5249554D2C204348414E5445524C414E4453204156454E55452C2048554C4C2C20485535344546\", \"reply\": \"ok\", \"pdu_cnt\": 2, \"pdu_cont\": true, \"ccid\": \"8944303412694355136\", \"pdu_mode\": \"yes\", \"validity\": \"1\", \"pdu_num_of_sms\": 2, \"client_id\": \"id1\", \"card_add\": \"21\", \"port_num\": \"4\", \"send_to_sim\": \"21#4\", \"active_sim\": 1, \"queue_type\": \"master\"}",
		config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if !success {
		t.Fatal("gateway send message response failed")
	}
}

func TestSendMsgSuccessfulNot1(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	success := SendMsgSuccessful(
		"{\"client_id\": \"id1\", \"unicode\": \"0\", \"msg\": \"6B656C6C6F\", \"reply\": \"error\", error_code: \"err-200\", \"number\": \"00447123456789\", \"msg_id\": 3432}",
		config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if success {
		t.Fatal("gateway send message response should have failed")
	}
}

func TestSendMsgSuccessfulNot2(t *testing.T) {
	// read config file
	config := config.ReadProperties()

	// This has error_code as quoted (this is not normally the case, facilitate for if it could happen)
	success := SendMsgSuccessful(
		"{\"client_id\": \"id1\", \"unicode\": \"0\", \"msg\": \"6B656C6C6F\", \"reply\": \"error\", \"error_code\": \"err-200\", \"number\": \"00447123456789\", \"msg_id\": 3432}",
		config.Gateways[0].GatewayAddress, config.Gateways[0].GatewayPort)
	log.Printf("success: %t\n", success)
	if success {
		t.Fatal("gateway send message response should have failed")
	}
}
