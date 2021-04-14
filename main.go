package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	regularLicenseBaseNextAppURL = "https://telegov.njportal.com/njmvc/CustomerCreateAppointments/GetNextAvailableDate?appointmentTypeId=11&locationId="
	appointmentLinkURL           = "https://telegov.njportal.com/njmvc/AppointmentWizard/11/"
	dateRegex                    = `\d{2}/\d{2}/\d{4}`
	slackWebhook                 = `` //`https://hooks.slack.com/services/...`
	numberOfDays                 = 15
	slackUserID                  = `` //`<@slackuseID>`
	sleepTimeInSeconds           = 60
)

func main() {
	listOfZipCodes := []string{"07002", "07047", "07644", "07505", "07470", "07114", "07065", "08817"}
	re := regexp.MustCompile(dateRegex)
	lList := GetLocationList()
	for {
		fmt.Println(time.Now())
		for _, v := range lList {
			if !contains(listOfZipCodes, v.Zip) {
				continue
			}
			aptDetails := appointmentDetails{}
			err := genericCall(getRegularApptURL(v), "GET", nil, &aptDetails)
			if err != nil {
				fmt.Println("Failed to query: " + err.Error())
				continue
			}
			if !re.MatchString(aptDetails.Next) {
				// no Appointments found
				continue
			}

			submatchall := re.FindAllString(aptDetails.Next, 1)
			aptTime, _ := time.Parse("01/02/2006", submatchall[0])
			if aptTime.Sub(time.Now()).Hours() > 24*numberOfDays {
				continue
			}

			sData := getSlackFormatedMessage(aptDetails.Next, v)
			slackData, _ := json.Marshal(sData)
			err = genericCall(slackWebhook, "POST", slackData, nil)
			if err != nil {
				fmt.Println("Failed to post to slack" + err.Error())
			}

			fmt.Println(" -- ")
			fmt.Printf("Appointment: %v\n", aptDetails.Next)
			fmt.Printf("Name: %v\n", v.Name)
			fmt.Printf("Street: %v, %v \n", v.Street1, v.Street2)
			fmt.Printf("City: %v\n", v.City)
			fmt.Printf("Zip: %v\n", v.Zip)
			fmt.Printf("Link to book: %v\n", appointmentLinkURL+strconv.Itoa(v.LocAppointments[0].LocationID))
			fmt.Println(" -- ")

		}
		fmt.Printf("-------------- \n\n\n")
		time.Sleep(sleepTimeInSeconds * time.Second)
	}
}

func genericCall(url, operation string, payload []byte, returnInterface interface{}) error {

	req, err := http.NewRequest(operation, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	requestBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodyString := string(requestBody[:])
	if strings.EqualFold(bodyString, "ok") {
		return nil
	}
	err = json.Unmarshal([]byte(bodyString), &returnInterface)
	if err != nil {
		return fmt.Errorf("Serialization failed : %v", err.Error())
	}
	if resp.StatusCode != 200 || err != nil {
		return fmt.Errorf("%v: request to %v failed with code %v", operation, url, resp.Status)
	}
	return nil
}

func getRegularApptURL(l LocationData) string {
	return regularLicenseBaseNextAppURL + strconv.Itoa(l.LocAppointments[0].LocationID)
}

func contains(arr []string, val string) bool {
	if len(arr) <= 0 {
		return true
	}

	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}
