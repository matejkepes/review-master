package email

import (
	"container/ring"
	"testing"
	"time"
)

// EmailLastSent - used to check when email was last sent so do not send too frquently.
// Initialise email last sent to a long time ago.
var EmailLastSent = time.Now().AddDate(-10, 0, 0)

// FailoverEmailLastSent - used to check when email was last sent so do not send too frquently for failover.
// Initialise email last sent to a long time ago.
var FailoverEmailLastSent = time.Now().AddDate(-10, 0, 0)

var sendSmsLastErrorsLen = 20

// SendSmsLastErrors - list of times of last send SMS errors which should initiate an email
var SendSmsLastErrors = ring.New(sendSmsLastErrorsLen)

// FailoverSendSmsLastErrors - list of times of last send SMS errors which should initiate an email for failover
var FailoverSendSmsLastErrors = ring.New(sendSmsLastErrorsLen)

func init() {
	// initialise last send SMS errors list to a time before check period
	for i := 0; i < sendSmsLastErrorsLen; i++ {
		SendSmsLastErrors.Value = time.Now().AddDate(-10, 0, 0)
		SendSmsLastErrors = SendSmsLastErrors.Next()
		FailoverSendSmsLastErrors.Value = time.Now().AddDate(-10, 0, 0)
		FailoverSendSmsLastErrors = FailoverSendSmsLastErrors.Next()
	}
}

func TestSend1(t *testing.T) {
	success := Send(TestSmtpServer, TestSmtpServerPort, TestPassword, TestFrom, TestTo, "test", "testing")
	if !success {
		t.Fatal("Error sending email")
	}
}

func TestSend2(t *testing.T) {
	success := Send(TestSmtpServer, TestSmtpServerPort, "wrong", TestFrom, TestTo, "test", "testing")
	if success {
		t.Fatal("Error should not have sent email")
	}
}

func TestSend3(t *testing.T) {
	success := Send(TestSmtpServer, TestSmtpServerPort, "wrong", TestFrom, "someone@somewhere.com, some@else.co.uk,som@what.com  ,  who@where.com", "test", "testing")
	if success {
		t.Fatal("Error should not have sent email")
	}
}

func TestSend4(t *testing.T) {
	success := Send(TestSmtpServer, TestSmtpServerPort, TestPassword, TestFrom, "", "test", "testing")
	if !success {
		t.Fatal("Error should return sending email successful even though nothing was sent because no recipients")
	}
}

func TestSend5(t *testing.T) {
	success := Send(TestSmtpServer, TestSmtpServerPort, TestPassword, TestFrom, ", , ,,,,", "test", "testing")
	if !success {
		t.Fatal("Error should return sending email successful even though nothing was sent because no recipients")
	}
}

func TestCheckSend1(t *testing.T) {
	// Should NOT send
	// shared.SendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	// success := CheckSend(shared.SendSmsLastErrors, &shared.EmailLastSent)
	success := CheckSend(SendSmsLastErrors, &EmailLastSent)
	// fmt.Println()
	// shared.SendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	if success {
		t.Fatal("error check should return false")
	}
}

func TestFailoverCheckSend1(t *testing.T) {
	// Should NOT send
	// shared.FailoverSendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	// success := CheckSend(shared.SendSmsLastErrors, &shared.EmailLastSent)
	success := CheckSend(SendSmsLastErrors, &EmailLastSent)
	// fmt.Println()
	// shared.FailoverSendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	if success {
		t.Fatal("error check should return false")
	}
}

func TestCheckSend2(t *testing.T) {
	// Should send
	// initialise last send SMS errors list to a few minutes ago
	oldest := time.Now().Add(-time.Minute * 1)
	// for i := 0; i < shared.SendSmsLastErrors.Len(); i++ {
	// 	oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.SendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.SendSmsLastErrors = shared.SendSmsLastErrors.Next()
	// }
	for i := 0; i < SendSmsLastErrors.Len(); i++ {
		oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
		SendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
		SendSmsLastErrors = SendSmsLastErrors.Next()
	}
	// shared.SendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	// success := CheckSend(shared.SendSmsLastErrors, &shared.EmailLastSent)
	success := CheckSend(SendSmsLastErrors, &EmailLastSent)
	// fmt.Println()
	// shared.SendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	if !success {
		t.Fatal("error check should return true")
	}
	// check that oldest does not exist (should have been replaced)
	// for i := 0; i < shared.SendSmsLastErrors.Len(); i++ {
	// 	if oldest.Equal(shared.SendSmsLastErrors.Value.(time.Time)) {
	// 		t.Fatal("oldest time has not been replaced in list")
	// 	}
	// 	shared.SendSmsLastErrors = shared.SendSmsLastErrors.Next()
	// }
	for i := 0; i < SendSmsLastErrors.Len(); i++ {
		if oldest.Equal(SendSmsLastErrors.Value.(time.Time)) {
			t.Fatal("oldest time has not been replaced in list")
		}
		SendSmsLastErrors = SendSmsLastErrors.Next()
	}
}

func TestFailoverCheckSend2(t *testing.T) {
	// Should send
	// initialise last send SMS errors list to a few minutes ago
	oldest := time.Now().Add(-time.Minute * 1)
	// for i := 0; i < shared.FailoverSendSmsLastErrors.Len(); i++ {
	// 	oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.FailoverSendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.FailoverSendSmsLastErrors = shared.FailoverSendSmsLastErrors.Next()
	// }
	for i := 0; i < FailoverSendSmsLastErrors.Len(); i++ {
		oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
		FailoverSendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
		FailoverSendSmsLastErrors = FailoverSendSmsLastErrors.Next()
	}
	// shared.FailoverSendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	success := CheckSend(FailoverSendSmsLastErrors, &FailoverEmailLastSent)
	// fmt.Println()
	// shared.FailoverSendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	if !success {
		t.Fatal("error check should return true")
	}
	// check that oldest does not exist (should have been replaced)
	// for i := 0; i < shared.FailoverSendSmsLastErrors.Len(); i++ {
	// 	if oldest.Equal(shared.FailoverSendSmsLastErrors.Value.(time.Time)) {
	// 		t.Fatal("oldest time has not been replaced in list")
	// 	}
	// 	shared.FailoverSendSmsLastErrors = shared.FailoverSendSmsLastErrors.Next()
	// }
	for i := 0; i < FailoverSendSmsLastErrors.Len(); i++ {
		if oldest.Equal(FailoverSendSmsLastErrors.Value.(time.Time)) {
			t.Fatal("oldest time has not been replaced in list")
		}
		FailoverSendSmsLastErrors = FailoverSendSmsLastErrors.Next()
	}
}

func TestCheckSend3(t *testing.T) {
	// Should NOT send
	// initialise last send SMS errors list to a few minute ago
	oldest := time.Now().Add(-time.Minute * 1)
	// for i := 0; i < shared.SendSmsLastErrors.Len(); i++ {
	// 	oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.SendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.SendSmsLastErrors = shared.SendSmsLastErrors.Next()
	// }
	for i := 0; i < SendSmsLastErrors.Len(); i++ {
		oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
		SendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
		SendSmsLastErrors = SendSmsLastErrors.Next()
	}
	// initialise email last sent to a minute ago
	// shared.EmailLastSent = time.Now().Add(time.Minute)
	EmailLastSent = time.Now().Add(time.Minute)
	// shared.SendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	// success := CheckSend(shared.SendSmsLastErrors, &shared.EmailLastSent)
	success := CheckSend(SendSmsLastErrors, &EmailLastSent)
	// fmt.Println()
	// shared.SendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	if success {
		t.Fatal("error check should return true")
	}
	// check that oldest does not exist (should have been replaced)
	// for i := 0; i < shared.SendSmsLastErrors.Len(); i++ {
	// 	if oldest.Equal(shared.SendSmsLastErrors.Value.(time.Time)) {
	// 		t.Fatal("oldest time has not been replaced in list")
	// 	}
	// 	shared.SendSmsLastErrors = shared.SendSmsLastErrors.Next()
	// }
	for i := 0; i < SendSmsLastErrors.Len(); i++ {
		if oldest.Equal(SendSmsLastErrors.Value.(time.Time)) {
			t.Fatal("oldest time has not been replaced in list")
		}
		SendSmsLastErrors = SendSmsLastErrors.Next()
	}
}

func TestFailoverCheckSend3(t *testing.T) {
	// Should NOT send
	// initialise last send SMS errors list to a few minute ago
	oldest := time.Now().Add(-time.Minute * 1)
	// for i := 0; i < shared.FailoverSendSmsLastErrors.Len(); i++ {
	// 	oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.FailoverSendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
	// 	shared.FailoverSendSmsLastErrors = shared.FailoverSendSmsLastErrors.Next()
	// }
	for i := 0; i < FailoverSendSmsLastErrors.Len(); i++ {
		oldest = time.Now().Add(-time.Minute * time.Duration(i+1))
		FailoverSendSmsLastErrors.Value = time.Now().Add(-time.Minute * time.Duration(i+1))
		FailoverSendSmsLastErrors = FailoverSendSmsLastErrors.Next()
	}
	// initialise email last sent to a minute ago
	// shared.FailoverEmailLastSent = time.Now().Add(time.Minute)
	FailoverEmailLastSent = time.Now().Add(time.Minute)
	// shared.FailoverSendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	// success := CheckSend(shared.FailoverSendSmsLastErrors, &shared.FailoverEmailLastSent)
	success := CheckSend(FailoverSendSmsLastErrors, &FailoverEmailLastSent)
	// fmt.Println()
	// shared.FailoverSendSmsLastErrors.Do(func(p interface{}) {
	// 	fmt.Println(p)
	// })
	if success {
		t.Fatal("error check should return true")
	}
	// check that oldest does not exist (should have been replaced)
	// for i := 0; i < shared.FailoverSendSmsLastErrors.Len(); i++ {
	// 	if oldest.Equal(shared.FailoverSendSmsLastErrors.Value.(time.Time)) {
	// 		t.Fatal("oldest time has not been replaced in list")
	// 	}
	// 	shared.FailoverSendSmsLastErrors = shared.FailoverSendSmsLastErrors.Next()
	// }
	for i := 0; i < FailoverSendSmsLastErrors.Len(); i++ {
		if oldest.Equal(FailoverSendSmsLastErrors.Value.(time.Time)) {
			t.Fatal("oldest time has not been replaced in list")
		}
		FailoverSendSmsLastErrors = FailoverSendSmsLastErrors.Next()
	}
}
