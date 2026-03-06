package main

import (
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

func StartPeriodicFetching() {
	var err error

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	// Fetch immediately
	fmt.Println("Fetching calendar...")
	cal, err = ics.ParseCalendarFromUrl(generateURL())
	addCustomEvent()
	if err != nil {
		fmt.Println("Error fetching calendar:", err)
		return
	}

	// Then fetch every hour
	for range ticker.C {
		fmt.Println("Fetching calendar...")
		cal, err = ics.ParseCalendarFromUrl(generateURL())
		addCustomEvent()
		if err != nil {
			fmt.Println("Error fetching calendar:", err)
			continue
		}
	}
}

func generateURL() string {
	return "https://ade-web-consult.univ-amu.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?projectId=8&resources=1819&calType=ical&firstDate=2025-09-01&lastDate=2026-09-19"
}

func addCustomEvent() {
	event := cal.AddEvent(fmt.Sprintf("thankyou-%d@vamary.fr", time.Now().Unix()))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())

	event.SetStartAt(time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC))
	event.SetAllDayStartAt(time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC))
	event.SetEndAt(time.Date(2026, 4, 10, 23, 59, 59, 0, time.UTC))
	event.SetSummary("Merci d'avoir utilisé mon script !")
	event.SetDescription("Merci d'avoir utilisé mon script !")
	event.SetURL("https://vamary.fr/")
	event.SetOrganizer("valentin.mary@proton.me", ics.WithCN("Valentin Mary"))
}
