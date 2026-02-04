package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ics "github.com/arran4/golang-ical"
)

var cal *ics.Calendar = ics.NewCalendar()

func main() {

	// Start the event fetcher in background
	go StartPeriodicFetching()

	// Define the front route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/index.html")
	})

	// Define the ics route
	http.HandleFunc("/api/v1/M1/calendar.ics", func(w http.ResponseWriter, r *http.Request) {

		// Get URL parameters
		query := r.URL.Query()

		// Get mentions
		var mentions []Mention
		json.Unmarshal([]byte(query.Get("mentions")), &mentions)

		// Get groups
		var groups []Group
		json.Unmarshal([]byte(query.Get("groups")), &groups)

		// Get options
		var options []Option
		json.Unmarshal([]byte(query.Get("options")), &options)

		// Get options-group
		var optionGroups []OptionGroup
		json.Unmarshal([]byte(query.Get("optionGroups")), &optionGroups)

		// Get the start/end dates selected by the user
		var startDate, errStartDate = time.Parse("2006-01-02", query.Get("startDate"))
		if errStartDate != nil {
			startDate = time.Date(2025, time.September, 1, 0, 0, 0, 0, time.Local)
		}
		var endDate, errEndDate = time.Parse("2006-01-02", query.Get("endDate"))
		if errEndDate != nil {
			endDate = time.Date(2026, time.August, 1, 0, 0, 0, 0, time.Local)
		}

		// Get the filtered calendar
		filteredCal := filterCalendar(
			mentions,
			groups,
			options,
			optionGroups,
			startDate,
			endDate)

		// Return the ics file to the user
		w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
		fmt.Fprint(w, filteredCal.Serialize())
	})

	// Start the http server
	fmt.Println("Starting HTTP server on :8080")
	http.ListenAndServe(":8080", nil)
}
