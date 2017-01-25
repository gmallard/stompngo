package stompngo

import (
	"log"
	"net"
	"os"
)

//=============================================================================
//= ack_test type =============================================================
//=============================================================================
type (
	terrData struct {
		proto   string
		headers Headers
		errval  Error
	}
)

//=============================================================================
//= ack_test var ==============================================================
//=============================================================================
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

//=============================================================================
//= ack_test const ============================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= codec_test type ===========================================================
//=============================================================================
type (
	//
	testdata struct {
		encoded string
		decoded string
	}
)

//=============================================================================
//= codec_test var ============================================================
//=============================================================================
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

//=============================================================================
//= codec_test const ==========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= connbv_test type ==========================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= connbv_test var ===========================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= connbv_test const =========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= conndisc_test type ========================================================
//=============================================================================
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

//=============================================================================
//= conndisc_test var =========================================================
//=============================================================================
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

//=============================================================================
//= conndisc_test const =======================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= data_test type ============================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= data_test var =============================================================
//=============================================================================
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

//=============================================================================
//= data_test const ===========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= hb_test type ==============================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= hb_test var ===============================================================
//=============================================================================
var (
	testhbl = false // Run long heartbeat tests
)

//=============================================================================
//= hb_test const =============================================================
//=============================================================================
const (
	//
	hbs = 45 // Wait time (secs)
)

//=============================================================================
//= headers_test type =========================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= headers_test var ==========================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= headers_test const ========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= logger_test type ==========================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= logger_test var ===========================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= logger_test const =========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= misc_test type ============================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= misc_test var =============================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= misc_test const ===========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= nack_test type ============================================================
//=============================================================================
type (
	nackData struct {
		proto   string
		headers Headers
		errval  Error
	}
)

//=============================================================================
//= nack_test var =============================================================
//=============================================================================
var (
	nackList = []nackData{
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/a"},
			EBADVERNAK},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/a"},
			EREQSUBNAK},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/a", HK_SUBSCRIPTION, "sub11a"},
			EREQMIDNAK},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/a"},
			EREQIDNAK},
	}
)

//=============================================================================
//= nack_test const ===========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= send_test type ============================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= send_test var =============================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= send_test const ===========================================================
//=============================================================================
const (
// None at present.
)

// send_test END
// sendbytes_test BEGIN

//=============================================================================
//= sendbytes_test type =======================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= sendbytes_test var ========================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= sendbytes_test const ======================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= sub_test type =============================================================
//=============================================================================
type (
	subNoHeaderData struct {
		proto string
		exe   error
	}

	subNoIDData struct {
		proto string
		subh  Headers
		exe   error
	}

	subPlainData struct {
		proto  string
		subh   Headers
		unsubh Headers
		exe1   error
		exe2   error
	}

	subTwiceData struct {
		proto string
		subh  Headers
		exe1  error
		exe2  error
	}

	subAckData struct {
		proto string
		subh  Headers
		exe   error
	}
)

//=============================================================================
//= sub_test var ==============================================================
//=============================================================================
var (
	subNoHeaderDataList = []subNoHeaderData{
		{SPL_10, EREQDSTSUB},
		{SPL_11, EREQDSTSUB},
		{SPL_12, EREQDSTSUB},
	}
	subNoIDDataList = []subNoIDData{
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subNoIDTest.10"},
			nil},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subNoIDTest.11"},
			nil},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subNoIDTest.12"},
			nil},
	}
	subPlainDataList = []subPlainData{
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subPlainTest.10",
				HK_ID, "subPlainTest.10"},
			Headers{HK_DESTINATION, "/queue/subPlainTest.10",
				HK_ID, "subPlainTest.10"},
			nil, nil},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subPlainTest.11",
				HK_ID, "subPlainTest.11"},
			Headers{HK_DESTINATION, "/queue/subPlainTest.11",
				HK_ID, "subPlainTest.11"},
			nil, nil},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subPlainTest.12",
				HK_ID, "subPlainTest.11"},
			Headers{HK_DESTINATION, "/queue/subPlainTest.12",
				HK_ID, "subPlainTest.11"},
			nil, nil},
	}

	subTwiceDataList = []subTwiceData{
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subTwiceTest.10",
				HK_ID, "subTwiceTest.10"},
			nil, EDUPSID},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subTwiceTest.11",
				HK_ID, "subTwiceTest.11"},
			nil, EDUPSID},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subTwiceTest.12",
				HK_ID, "subTwiceTest.11"},
			nil, EDUPSID},
	}

	subAckDataList = []subAckData{
		// 1.0
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subAckTest.10.1"},
			nil},
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subAckTest.10.2",
				HK_ACK, AckModeAuto},
			nil},
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subAckTest.10.3",
				HK_ACK, AckModeClient},
			nil},
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subAckTest.10.3",
				HK_ACK, AckModeClientIndividual},
			ESBADAM},
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/subAckTest.10.4",
				HK_ACK, badam},
			ESBADAM},
		// 1.1
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subAckTest.11.1"},
			nil},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subAckTest.11.2",
				HK_ACK, AckModeAuto},
			nil},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subAckTest.11.3",
				HK_ACK, AckModeClient},
			nil},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subAckTest.11.4",
				HK_ACK, AckModeClientIndividual},
			nil},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/subAckTest.11.5",
				HK_ACK, badam},
			ESBADAM},
		// 1.2
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subAckTest.12.1"},
			nil},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subAckTest.12.2",
				HK_ACK, AckModeAuto},
			nil},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subAckTest.12.3",
				HK_ACK, AckModeClient},
			nil},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subAckTest.12.4",
				HK_ACK, AckModeClientIndividual},
			nil},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/subAckTest.12.5",
				HK_ACK, badam},
			ESBADAM},
	}
)

//=============================================================================
//= sub_test const ============================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= suppress_test type ========================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= suppress_test var =========================================================
//=============================================================================
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

//=============================================================================
//= suppress_test const =======================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= testdata_test type ========================================================
//=============================================================================
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

//=============================================================================
//= testdata_test var =========================================================
//=============================================================================
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

//=============================================================================
//= testdata_test const =======================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= trans_test type ===========================================================
//=============================================================================
type (
	transBasicData struct {
		action string
		th     Headers
		te     error
	}
	transSendCommitData struct {
		tid string
		exe error
	}
	transSendAbortData struct {
		tid string
		exe error
	}
	transMessageOrderData struct {
		sh Headers
		se error
	}
)

//=============================================================================
//= trans_test var ============================================================
//=============================================================================
var (
	transBasicList = []transBasicData{
		{BEGIN, Headers{}, EREQTIDBEG},
		{COMMIT, Headers{}, EREQTIDCOM},
		{ABORT, Headers{}, EREQTIDABT},
		{BEGIN, Headers{HK_TRANSACTION, ""}, ETIDBEGEMT},
		{COMMIT, Headers{HK_TRANSACTION, ""}, ETIDCOMEMT},
		{ABORT, Headers{HK_TRANSACTION, ""}, ETIDABTEMT},
	}
	transSendCommitList = []transSendCommitData{
		{"trans.send.commit", nil},
	}

	transSendAbortList = []transSendAbortData{
		{"trans.send.abort", nil},
	}
	transMessageOrderList = []transMessageOrderData{
		{Headers{HK_DESTINATION, "/queue/tsrbdata.q"}, nil},
	}
)

//=============================================================================
//= trans_test const ==========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= unsub_test type ===========================================================
//=============================================================================
type (
	unsubNoHeaderData struct {
		proto string
		exe   error
	}
	unsubNoIDData struct {
		proto  string
		unsubh Headers
		exe    error
	}
	unsubBoolData struct {
		proto    string
		subfirst bool
		subh     Headers
		unsubh   Headers
		exe1     error
		exe2     error
	}
)

//=============================================================================
//= unsub_test var ============================================================
//=============================================================================
var (
	unsubNoHeaderDataList = []subNoHeaderData{
		{SPL_10, EREQDIUNS},
		{SPL_11, EREQDIUNS},
		{SPL_12, EREQDIUNS},
	}
	unsubNoIDDataList = []subNoIDData{
		// 1.0
		{SPL_10,
			Headers{},
			EREQDIUNS},
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/unsubIDTest.10.1"},
			EUNOSID},
		{SPL_10,
			Headers{HK_DESTINATION, "/queue/unsubIDTest.10.2",
				HK_ID, "unsubIDTest.10.2"},
			EBADSID},
		// 1.1
		{SPL_11,
			Headers{},
			EREQDIUNS},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/unsubIDTest.11.1"},
			EUNOSID},
		{SPL_11,
			Headers{HK_DESTINATION, "/queue/unsubIDTest.11.2",
				HK_ID, "unsubIDTest.11.2"},
			EBADSID},
		// 1.2
		{SPL_12,
			Headers{},
			EREQDIUNS},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/unsubIDTest.12.1"},
			EUNOSID},
		{SPL_12,
			Headers{HK_DESTINATION, "/queue/unsubIDTest.12.2",
				HK_ID, "unsubIDTest.12.2"},
			EBADSID},
	}

	// REQIDUNS = Error("id required, UNSUBSCRIBE")
	// REQDIUNS  = Error("destination required, UNSUBSCRIBE")

	unsubBoolDataList = []unsubBoolData{
		// 1.0
		{SPL_10, false,
			Headers{},
			Headers{},
			EREQDSTUNS, EREQDSTUNS},
		{SPL_10, false,
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			EREQIDUNS, EREQIDUNS},
		{SPL_10, true,
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			nil, nil},
		// 1.1
		{SPL_11, false,
			Headers{},
			Headers{},
			EREQDSTUNS, EREQDSTUNS},
		{SPL_11, false,
			Headers{HK_DESTINATION, "/queue/PlainDataTest.11.1"},
			Headers{HK_DESTINATION, "/queue/PlainDataTest.11.1"},
			EREQIDUNS, EREQIDUNS},
		{SPL_11, true,
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			nil, EREQIDUNS},
		// 1.2
		{SPL_12, false,
			Headers{},
			Headers{},
			EREQDIUNS, EREQDSTUNS},
		{SPL_12, false,
			Headers{HK_DESTINATION, "/queue/PlainDataTest.12.1"},
			Headers{HK_DESTINATION, "/queue/PlainDataTest.12.1"},
			EREQIDUNS, EREQIDUNS},
		{SPL_12, true,
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			Headers{HK_DESTINATION, "/queue/PlainDataTest.10.1"},
			nil, EREQIDUNS},
	}
)

//=============================================================================
//= unsub_test const ==========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= utils_test type ===========================================================
//=============================================================================
type (
	multi_send_data struct {
		conn  *Connection // this connection
		dest  string      // queue/topic name
		mpref string      // message prefix
		count int         // number of messages
	}
)

//=============================================================================
//= utils_test var ============================================================
//=============================================================================
var (
// None at present.
)

//=============================================================================
//= utils_test const ==========================================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= shovel_dup_headers_test type ==============================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= shovel_dup_headers_test var ===============================================
//=============================================================================
var (
	tsdhHeaders = Headers{
		"dupkey1", "value0",
		"dupkey1", "value1",
		"dupkey1", "value2",
	}
	wantedDupeV1 = Headers{
		"dupkey1", "value1",
	}
	wantedDupeV2 = Headers{
		"dupkey1", "value2",
	}
	wantedDupeVAll = Headers{
		"dupkey1", "value1",
		"dupkey1", "value2",
	}
)

//=============================================================================
//= shovel_dup_headers_test const =============================================
//=============================================================================
const (
// None at present.
)

//=============================================================================
//= for use by all type =======================================================
//=============================================================================
type (
// None at present.
)

//=============================================================================
//= for use by all var ========================================================
//=============================================================================
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
	sp               string
	badam            = "AckModeInvalid"
	brokerid         int
	tm               = "A Test Message."
	tlg              = log.New(os.Stderr, "TLG|", log.Ldate|log.Lmicroseconds)
)

//=============================================================================
//= for use by all const ======================================================
//=============================================================================
const (
	TEST_ANYBROKER = iota
	TEST_AMQ       = iota
	TEST_RMQ       = iota
	TEST_ARTEMIS   = iota
	TEST_APOLLO    = iota
)
