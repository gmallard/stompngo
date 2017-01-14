package stompngo

// ack_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// ack_test END

// codec_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// codec_test END

// connbv_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// connbv_test END

// conndisc_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// conndisc_test END

// data_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// data_test END

// hb_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// hb_test END

// headers_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// headers_test END

// logger_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// logger_test END

// misc_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// misc_test END

// nack_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// nack_test END

// send_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// send_test END
// sendbytes_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// sendbytes_test END
// sub_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// sub_test END
// suppress_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// suppress_test END
// trans_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// trans_test END
// unsub_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// unsub_test END
// utils_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// utils_test END
// wrsubrduns_test BEGIN

type (
//
)

var (
//
)

const (
//
)

// wrsubrduns_test END

type sendRecvCodecData struct {
	sid string
	sk  []string          // send keys
	sv  []string          // send values
	rv  map[string]string // expected receive value
}

type testdata struct {
	encoded string
	decoded string
}

type multi_send_data struct {
	conn  *Connection // this connection
	dest  string      // queue/topic name
	mpref string      // message prefix
	count int         // number of messages
}

type frameData struct {
	data string
	resp error
}

type supdata struct {
	v string // version
	s bool   // is supported
}

type verData struct {
	ch Headers // Client headers
	sh Headers // Server headers
	e  error   // Expected error
}

type unsubData struct {
	p string // protocol
	e error  // error
}

var (
	TEST_HEADERS   = Headers{HK_LOGIN, "guest", HK_PASSCODE, "guest"}
	TEST_TDESTPREF = "/queue/test.pref."
	TEST_TRANID    = "TransactionA"
	//

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

	md MessageData
	hv string
	ok bool

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

	empty_headers = Headers{}

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

	suptests = []supdata{
		{SPL_10, true},
		{SPL_11, true},
		{SPL_12, true},
		{"1.3", false},
		{"2.0", false},
		{"2.1", false},
	}

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
	hbs = 45
)
