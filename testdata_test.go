package stompngo

type sendRecvCodecData struct {
	sid string
	sk  []string          // send keys
	sv  []string          // send values
	rv  map[string]string // expected receive value
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
)
