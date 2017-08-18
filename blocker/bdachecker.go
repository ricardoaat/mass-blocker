package blocker

import (
	"regexp"
	"strconv"

	"sync"

	"fmt"

	"github.com/ricardoaat/mass-blocker/SOAPservices"
	"github.com/ricardoaat/mass-blocker/db"
	log "github.com/sirupsen/logrus"
)

var countFound int
var countRegistered int
var countBlacklisted int
var countWorkerDone int
var countInvalid int

var detectionlist []detected
var imeiResults []imeiResult

type imeiResult struct {
	detected detected
	status   string
}

type detected struct {
	imei, imsi, msisdn, treatment string
}

const maxServiceRequest int = 100

func retriveNotBlockedhandsets() {

	s :=
		`
		select dl.detectedimei, dl.detectedimsi, dl.detectedmsisdn, dl.treatment
		from detection_list dl
		left join blacklist bl on dl.detectedimei = bl.imei
		left join whitelist wl on dl.detectedimei = wl.imei
		where dl.providerid = 4
		and dl.plannedblockingdate < '20170401'
		and dl.approvaldate is null
		and dl.registrationdate is null
		and bl.providerid is null
		and wl.providerid is null
		`

	rows, err := db.DbHandset.Query(s)
	checkErr(err, "Error fetching detection list DB ")

	defer rows.Close()

	if rows.Next() {
		countFound++
		var imei, imsi, msisdn, treatment string
		err := rows.Scan(&imei, &imsi, &msisdn, &treatment)
		checkErr(err, "Error scaning column IMEI")
		d := detected{imei, imsi, msisdn, treatment}
		detectionlist = append(detectionlist, d)

		for rows.Next() {
			countFound++
			err := rows.Scan(&imei, &imsi, &msisdn, &treatment)
			checkErr(err, "Error scaning column IMEI")
			d := detected{imei, imsi, msisdn, treatment}
			detectionlist = append(detectionlist, d)
		}
	} else {
		log.Fatal("No Handset on detection list query")
	}

	log.WithFields(log.Fields{
		"Found Query": countFound,
		"ImeiSlice":   len(detectionlist),
	}).Info("Fetched from DB")
}

func serviceConsumerDispatcher() {
	var wg sync.WaitGroup
	cr := make(chan imeiResult, maxServiceRequest)
	cet := make(chan error, maxServiceRequest)
	crs := make(chan []imeiResult)
	cwl := 0
	cdw := 0
	log.Info(fmt.Sprintf("Bout to dispatch %d workers", len(detectionlist)))
	go workerWatcher(cr, cet, crs)

	for i, d := range detectionlist {
		if cwl == maxServiceRequest {
			log.Info("Bout to wait " + strconv.Itoa(i))
			wg.Wait()
			cwl = 0
		}
		wg.Add(1)
		go checkImeiOnBda(d, cr, cet, &wg)
		cwl++
		cdw = i + 1
	}

	log.Info(fmt.Sprintf("Getting Worker Watcher Result, dispatched: %d", cdw))
	imeiResults = <-crs
	log.Info("Finished consuming Status")
	log.Info(fmt.Sprintf("To BLOCK:%d Already blocked:%d", len(imeiResults), len(detectionlist)-len(imeiResults)))
	log.Info(fmt.Sprintf("Registered to other:%d Blacklisted:%d Unregistered but invalid:%d",
		countRegistered,
		countBlacklisted,
		countInvalid,
	))
}

func checkImeiOnBda(d detected, cr chan imeiResult, cet chan error, wg *sync.WaitGroup) {
	r, err := SOAPservices.CallService(d.imei)

	if err != nil {
		cet <- err
		panic(err)
	}
	re, err := regexp.Compile("<status>(.*)</status>")
	wg.Done()
	if err != nil {
		cet <- err
		panic(err)
	} else {
		status := re.FindStringSubmatch(r)
		cr <- imeiResult{d, status[1]}
	}

}

func workerWatcher(cr chan imeiResult, cet chan error, crs chan []imeiResult) {
	var imeir []imeiResult
	for i := 1; i <= len(detectionlist); i++ {
		select {
		case ires := <-cr:
			log.Info(fmt.Sprintf("%d Imei: %s status: %s", i, ires.detected, ires.status))
			if ires.status == "UNREGISTERED" {
				if ires.detected.treatment == "INVALID_IMEI" {
					countInvalid++
				}
				imeir = append(imeir, ires)
			} else if ires.status == "REGISTERED_TO_OTHER" {
				countRegistered++
			} else if ires.status == "BLACKLISTED" {
				countBlacklisted++
			} else {
				log.Warning(ires)
			}
		case errorOnThread := <-cet:
			log.Error(errorOnThread)
		}
		countWorkerDone++
	}
	crs <- imeir
}
