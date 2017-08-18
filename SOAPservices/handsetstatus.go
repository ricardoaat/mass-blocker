package SOAPservices

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ricardoaat/mass-blocker/config"
)

//CallService Calls example status service
func CallService(imei string) (string, error) {

	ad := config.Conf.Services["leservice"].Serviceaddrs
	po := config.Conf.Services["leservice"].Serverport
	se := config.Conf.Services["leservice"].ServiceEndpoint
	url := "http://" + ad + ":" + po + se

	s :=
		`
        <soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:col="http://bla.bla.com/le/service">
        <soap:Header/>
        <soap:Body>
            <col:blablaService>
                <request>
                    <document>
    	                <id>123456789</id>
	                    <type>1</type>
                    </document>            
                    <imei>%s</imei>
                </request>
            </col:blablaService>
        </soap:Body>
        </soap:Envelope>
    	`
	m := []byte(fmt.Sprintf(s, imei))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	req.Header.Add("Accept", "application/soap+xml, application/dime, multipart/related, text/*")
	if err != nil {
		return "", err
	}

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	r := string(b)

	return r, nil
}
