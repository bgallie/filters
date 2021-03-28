package filters

import "log"

func checkFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkFatalMsg(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}
