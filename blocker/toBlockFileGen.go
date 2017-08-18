package blocker

import (
	"fmt"
	"os"
)

func imeisToBlockFileGen() {
	fo, err := os.Create("imeisToBlock.csv")
	checkErr(err, "Error opening output file ")
	defer fo.Close()
	for _, m := range imeiResults {
		l := fmt.Sprintf(
			"%s|%s|%s|%s|%s \n",
			m.detected.imei,
			m.detected.imsi,
			m.detected.msisdn,
			m.detected.treatment,
			m.status)
		_, err := fo.WriteString(l)
		checkErr(err, "Error writing line: "+l)
	}
}

func imeiResultsToFile(e []imeiResult, f string) {
	fo, err := os.Create(f)
	checkErr(err, "Error opening output file ")
	defer fo.Close()
	for _, m := range e {
		l := fmt.Sprintf(
			"%s|%s|%s|%s|%s \n",
			m.detected.imei,
			m.detected.imsi,
			m.detected.msisdn,
			m.detected.treatment,
			m.status)
		_, err := fo.WriteString(l)
		checkErr(err, "Error writing line: "+l)
	}
}
