package stompngo

import (
	"net"
)

// ack_test BEGIN

type (
	terrData struct {
		proto   string
		headers Headers
		errval  Error
	}
)

var (
	terrList = []terrData{
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/a"},
			EREQMIDACK},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/a"},
			EREQSUBACK},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/a", HK_SUBSCRIPTION, "sub11a"},
			EREQMIDACK},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/a"},
			EREQIDACK},
	}
)

const (
// None at present.
)

// ack_test END

// codec_test BEGIN

type (
	//
	testdata struct {
		encoded string
		decoded string
	}
)

var (
	//
	tdList = []testdata{
		{"stringa", "stringa"},
		{"stringb", "stringb"},
		{"stringc", "stringc"},
		{"stringd", "stringd"},
		{"stringe", "stringe"},
		{"stringf", "stringf"},
		{"stringg", "stringg"},
		{"stringh", "stringh"},
		{"\\\\", "\\"},
		{"\\n", "\n"},
		{"\\c", ":"},
		{"\\\\\\n\\c", "\\\n:"},
		{"\\c\\n\\\\", ":\n\\"},
		{"\\\\\\c", "\\:"},
		{"c\\cc", "c:c"},
		{"n\\nn", "n\nn"},
	}
)

const (
// None at present.
)

// codec_test END

// connbv_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// connbv_test END

// conndisc_test BEGIN

type (
	frameData struct {
		data string
		resp error
	}
	verData struct {
		ch Headers // Client headers
		sh Headers // Server headers
		e  error   // Expected error
	}
)

var (
	frames = []frameData{ // Many are possible but very unlikely
		{"EBADFRM", EBADFRM},
		{"EUNKFRM\n\n\x00", EUNKFRM},
		{"ERROR\n\n\x00", nil},
		{"ERROR\n\x00", EBADFRM},
		{"ERROR\n\n", EBADFRM},
		{"ERROR\nbadconhdr\n\n\x00", EUNKHDR},
		{"ERROR\nbadcon:badmsg\n\n\x00", nil},
		{"ERROR\nbadcon:badmsg\n\nbad message\x00", nil},
		{"CONNECTED\n\n\x00", nil},
		{"CONNECTED\n\nconnbody\x00", EBDYDATA},
		{"CONNECTED\n\nconnbadbody", EBDYDATA},
		{"CONNECTED\nk1:v1\nk2:v2\n\nconnbody\x00", EBDYDATA},
		{"CONNECTED\nk1:v1\nk2:v2\n\nconnbody", EBDYDATA},
		{"CONNECTED\nk1:v1\nk2:v2\n\n\x00", nil},
	}
	verChecks = []verData{
		{Headers{HK_ACCEPT_VERSION, SPL_11}, Headers{HK_VERSION, SPL_11}, nil},
		{Headers{}, Headers{}, nil},
		{Headers{HK_ACCEPT_VERSION, "1.0,1.1,1.2"}, Headers{HK_VERSION, SPL_12}, nil},
		{Headers{HK_ACCEPT_VERSION, "1.3"}, Headers{HK_VERSION, "1.3"}, EBADVERSVR},
		{Headers{HK_ACCEPT_VERSION, "1.3"}, Headers{HK_VERSION, "1.1"}, EBADVERCLI},
		{Headers{HK_ACCEPT_VERSION, "1.0,1.1,1.2"}, Headers{}, nil},
	}
)

const (
// None at present.
)

// conndisc_test END

// data_test BEGIN

type (
// None at present.
)

var (
	suptests = []supdata{
		{SPL_10, true},
		{SPL_11, true},
		{SPL_12, true},
		{"1.3", false},
		{"2.0", false},
		{"2.1", false},
	}
)

const (
// None at present.
)

// data_test END

// hb_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
	//
	hbs = 45
)

// hb_test END

// headers_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// headers_test END

// logger_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// logger_test END

// misc_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// misc_test END

// nack_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// nack_test END

// send_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// send_test END
// sendbytes_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// sendbytes_test END
// sub_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// sub_test END
// suppress_test BEGIN

type (
// None at present.
)

var (
	tsclData = []struct {
		ba     []uint8
		wanted string
	}{
		{
			[]uint8{0x61, 0x62, 0x63, 0x64, 0x65, 0x66},
			"abcdef",
		},
		{
			[]uint8{0x61, 0x62, 0x63, 0x00, 0x64, 0x65, 0x66},
			"abc",
		},
		{
			[]uint8{0x64, 0x65, 0x66, 0x00},
			"def",
		},
		{
			[]uint8{0x00, 0x64, 0x65, 0x66, 0x00},
			"",
		},
	}
	tsctData = []struct {
		body       string
		doSuppress bool
		wanted     bool
	}{
		{
			"some data",
			true,
			false,
		},
		{
			"other data",
			false,
			true,
		},
	}
)

const (
// None at present.
)

// suppress_test END

// testdata_test BEGIN

type (
	//
	sendRecvCodecData struct {
		sid string
		sk  []string          // send keys
		sv  []string          // send values
		rv  map[string]string // expected receive value
	}
	supdata struct {
		v string // version
		s bool   // is supported
	}
)

var (
	srcdList10 = []sendRecvCodecData{
		{sid: "sub10a",
			sk: []string{"keya"},
			sv: []string{"valuea"},
			rv: map[string]string{"keya": "valuea"}},
		{sid: "sub10b",
			sk: []string{"key:one"},
			sv: []string{"value:a"},
			rv: map[string]string{"key": "one:value:a"}},
		{sid: "sub10c",
			sk: []string{"key"},
			sv: []string{"valuec"},
			rv: map[string]string{"key": "valuec"}},
	}

	srcdList1p = []sendRecvCodecData{
		{sid: "sub1xa",
			sk: []string{"keya"},
			sv: []string{"valuea"},
			rv: map[string]string{"keya": "valuea"}},
		{sid: "sub1xb",
			sk: []string{"key:one", "key/ntwo", "key:three/naaa\\bbb"},
			sv: []string{"value\\one", "value:two\\back:slash", "value\\three:aaa/nbbb"},
			rv: map[string]string{"key:one": "value\\one",
				"key/ntwo":            "value:two\\back:slash",
				"key:three/naaa\\bbb": "value\\three:aaa/nbbb"}},
	}

	srcdmap = map[string][]sendRecvCodecData{SPL_10: srcdList10,
		SPL_11: srcdList1p,
		SPL_12: srcdList1p}
)

const (
//
)

// testdata_test END

// trans_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// trans_test END
// unsub_test BEGIN

type (
	unsubData struct {
		p string // protocol
		e error  // error
	}
)

var (
	unsubListNoHdr = []unsubData{
		{SPL_10, EREQDIUNS},
		{SPL_11, EREQDIUNS},
		{SPL_12, EREQDIUNS},
	}

	unsubBadId = []unsubData{
		{SPL_11, EBADSID},
		{SPL_12, EBADSID},
	}

	unsubNoId = []unsubData{
		{SPL_11, EUNOSID},
		{SPL_12, EUNOSID},
	}
)

const (
// None at present.
)

// unsub_test END
// utils_test BEGIN

type (
	multi_send_data struct {
		conn  *Connection // this connection
		dest  string      // queue/topic name
		mpref string      // message prefix
		count int         // number of messages
	}
)

var ()

const (
// None at present.
)

// utils_test END
// wrsubrduns_test BEGIN

type (
// None at present.
)

var (
// None at present.
)

const (
// None at present.
)

// wrsubrduns_test END

// For use by all
var (
	TEST_HEADERS     = Headers{HK_LOGIN, "guest", HK_PASSCODE, "guest"}
	TEST_TDESTPREF   = "/queue/test.pref."
	TEST_TRANID      = "TransactionA"
	md               MessageData
	hv               string
	ok               bool
	empty_headers    = Headers{}
	testuser         = "guest" // "guest" is required by some brokers
	testpw           = "guest"
	login_headers    = Headers{HK_LOGIN, testuser, HK_PASSCODE, testpw}
	rid              = "receipt-12345"
	oneOnePlusProtos = []string{SPL_11, SPL_12}
	e                error
	n                net.Conn
	conn             *Connection
	sc               <-chan MessageData
)
