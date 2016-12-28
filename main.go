/*
Copyright 2016 Cameron Martin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type EnvoyS struct {
	Production []struct {
		Type             string  `json:"type"`
		ActiveCount      int     `json:"activeCount"`
		ReadingTime      int     `json:"readingTime"`
		WNow             float64 `json:"wNow"`
		WhLifetime       float64 `json:"whLifetime"`
		MeasurementType  string  `json:"measurementType,omitempty"`
		VarhLeadLifetime float64 `json:"varhLeadLifetime,omitempty"`
		VarhLagLifetime  float64 `json:"varhLagLifetime,omitempty"`
		VahLifetime      float64 `json:"vahLifetime,omitempty"`
		RmsCurrent       float64 `json:"rmsCurrent,omitempty"`
		RmsVoltage       float64 `json:"rmsVoltage,omitempty"`
		ReactPwr         float64 `json:"reactPwr,omitempty"`
		ApprntPwr        float64 `json:"apprntPwr,omitempty"`
		PwrFactor        float64 `json:"pwrFactor,omitempty"`
		WhToday          float64 `json:"whToday,omitempty"`
		WhLastSevenDays  float64 `json:"whLastSevenDays,omitempty"`
		VahToday         float64 `json:"vahToday,omitempty"`
		VarhLeadToday    float64 `json:"varhLeadToday,omitempty"`
		VarhLagToday     float64 `json:"varhLagToday,omitempty"`
	} `json:"production"`
	Consumption []struct {
		Type             string  `json:"type"`
		ActiveCount      int     `json:"activeCount"`
		MeasurementType  string  `json:"measurementType"`
		ReadingTime      int     `json:"readingTime"`
		WNow             float64 `json:"wNow"`
		WhLifetime       float64 `json:"whLifetime"`
		VarhLeadLifetime float64 `json:"varhLeadLifetime"`
		VarhLagLifetime  float64 `json:"varhLagLifetime"`
		VahLifetime      float64 `json:"vahLifetime"`
		RmsCurrent       float64 `json:"rmsCurrent"`
		RmsVoltage       float64 `json:"rmsVoltage"`
		ReactPwr         float64 `json:"reactPwr"`
		ApprntPwr        float64 `json:"apprntPwr"`
		PwrFactor        float64 `json:"pwrFactor"`
		WhToday          float64 `json:"whToday"`
		WhLastSevenDays  float64 `json:"whLastSevenDays"`
		VahToday         float64 `json:"vahToday"`
		VarhLeadToday    float64 `json:"varhLeadToday"`
		VarhLagToday     float64 `json:"varhLagToday"`
	} `json:"consumption"`
	Storage []struct {
		Type        string `json:"type"`
		ActiveCount int    `json:"activeCount"`
		ReadingTime int    `json:"readingTime"`
		WNow        int    `json:"wNow"`
		WhNow       int    `json:"whNow"`
		State       string `json:"state"`
		PercentFull int    `json:"percentFull"`
	} `json:"storage"`
}

var envoyIP string
var envoyPort int
var pvoutputApiKey string
var pvoutputSystemId int
var pollIntervalSeconds int
var timezone string

func init() {
	flag.StringVar(&envoyIP, "ENVOYIP", "", "IP Address of Envoy S to retrieve data from")
	flag.IntVar(&envoyPort, "ENVOYPORT", 80, "Port of the Envoy S to retrieve data from")
	flag.StringVar(&pvoutputApiKey, "PVOUTPUTAPIKEY", "", "PVOutput.org API Key to use to post data")
	flag.IntVar(&pvoutputSystemId, "PVOUTPUTSYSTEMID", 0, "PVOutput.org system ID for the Envoy S")
	flag.IntVar(&pollIntervalSeconds, "POLLINTERVALSECONDS", 300, "Polling interval in seconds")
	flag.StringVar(&timezone, "TIMEZONE", "", "Timezone of the Envoy S. If unset, same as current local timezone")
	flag.Parse()

	if envoyIP == "" || pvoutputApiKey == "" || pvoutputSystemId == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {

	for {
		time.Sleep(time.Second * time.Duration(pollIntervalSeconds))

		url := "http://" + envoyIP + ":" + strconv.Itoa(envoyPort) + "/production.json"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("NewRequest: ", err)
			continue
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Do: ", err)
			continue
		}

		var reading EnvoyS
		if err := json.NewDecoder(resp.Body).Decode(&reading); err != nil {
			log.Println(err)
			continue
		}

		resp.Body.Close()

		//if a timezone has been set use it, otherwise local
		ti := time.Now()
		if timezone != "" {
			l, err := time.LoadLocation(timezone)
			if err != nil {
				log.Println("LoadLocation: ", err)
				continue
			}
			ti = ti.In(l)
		}

		d := ti.Format("20060102")
		t := ti.Format("15:04")
		v1 := reading.Production[1].WhLifetime  //production to date in watt hours
		v3 := reading.Consumption[0].WhLifetime //consumption to date in watt hours

		submitUrl := "http://pvoutput.org/service/r2/addstatus.jsp"

		req2, err := http.NewRequest("GET", submitUrl, nil)

		req2.Header.Add("X-PVOutput-APIKey", pvoutputApiKey)
		req2.Header.Add("X-PVOutput-SystemID", strconv.Itoa(pvoutputSystemId))

		q := req2.URL.Query()
		q.Add("d", d)
		q.Add("t", t)
		q.Add("v1", strconv.FormatFloat(v1, 'f', -1, 64))
		q.Add("v3", strconv.FormatFloat(v3, 'f', -1, 64))
		q.Add("c1", "1")
		req2.URL.RawQuery = q.Encode()

		resp2, err2 := client.Do(req2)
		if err2 != nil {
			log.Println("Do: ", err)
		} else {
			fmt.Print(".")
		}

		resp2.Body.Close()

	}
}
