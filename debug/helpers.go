package debug

func DebugErr(err error) {
	if err == nil {
		return
	}
	// ENABLE when debugging
	// log.Println(err)
}

func DebugLog(v interface{}) {
	// ENABLE when debugging
	// log.Println(v)
}
