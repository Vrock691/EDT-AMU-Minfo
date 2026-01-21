package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
)

/*
Filters the calendar events based on the selected mentions, groups, options, and option groups
*/
func filterCalendar(
	mentions []Mention,
	groups []Group,
	options []Option,
	optionGroups []OptionGroup,
	startDate time.Time,
	endDate time.Time) ics.Calendar {

	// Create a regex expression to remove undesired events
	var eventToRemoveRegex []string

	// Check if calendar has been fetched
	if cal == nil {
		return *ics.NewCalendar()
	}

	// Create a new calendar (deep copy by building fresh)
	filteredCal := *ics.NewCalendar()

	// Add to regex every events that are not in the selected mentions
	for _, mention := range mentions {
		if codes, exists := mentionToCodesMap[string(mention)]; exists {
			eventToRemoveRegex = append(eventToRemoveRegex, codes...)
		}
	}

	// Add to regex every events that are not in the selected groups
	groupsToBan := make([]string, 0, len(groupToRegexMap))
	for _, pattern := range groupToRegexMap {
		groupsToBan = append(groupsToBan, pattern)
	}

	for _, group := range groups {
		if pattern, exists := groupToRegexMap[group]; exists {
			groupsToBan = removeStringFromList(groupsToBan, pattern)
		}
	}
	eventToRemoveRegex = append(eventToRemoveRegex, groupsToBan...)

	// Add to regex every events that are not in the selected options
	// As options are already part of others mentions, they're already in the regex, so we just have to remove them from it if they're selected
	for _, option := range options {
		switch option {
		case "cpp":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_PROG_CPP)
		case "crypto":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_CRYPTOGRAPHIE)
		case "intro-science-donnees":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_SDD_AA)
		case "methode-numeriques":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_METHODE_NUM_INFORMATIQUE)
		case "prog-fonctionnelle":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_PROG_FONCTIONNELLE)
		case "proba":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_PROBA_POUR_INFORMATIQUE)
		case "securite-des-apps":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_SECURITE_DES_APPS)
		case "adon":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_ANALYSE_DES_DONNEES)
		case "asr":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_ALGORITHMES_ET_SYSTEMES_ROBUSTES)
		case "iq":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_INFORMATIQUE_QUANTIQUE)
		case "tdc":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_THEORIE_DE_LA_COMPLEXITE)
		case "fl":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_FIABILITE_LOGICIELLE)
		case "baa":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_BASES_DE_APPRENTISSAGE_ARTIFICIEL)
		case "ihm":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_INTERFACE_HUMAIN_MACHINE)
		case "gla":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_GENIE_LOGICIEL_AVANCE)
		case "pgra":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_PROGRAMMATION_GRAPHIQUE)
		case "apro":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_ANALYSE_DE_PROGRAMMES)
		case "mgm":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_MODELISATION_GEOMETRIQUE_ET_MAILLAGES)
		case "tal":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_TRAITEMENT_AUTOMATIQUE_DES_LANGUES)
		case "apg":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_ALGORITHMES_A_PERFORMANCES_GARANTIES)
		case "cloud":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_DES_CONTENEURS_AU_CLOUD)
		case "idd":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_INTEGRATION_DES_DONNEES)
		case "tsg":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_THEORIE_STRUCTURELLE_DES_GRAPHES)
		case "calc":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_CALCULABILITE_AVANCEE)
		case "mrd":
			eventToRemoveRegex = removeStringFromList(eventToRemoveRegex, CODE_MODELISATION_ET_RESOLUTION_POUR_LA_DECISION)
		}
	}

	// Add to regex every events that are not in the selected option groups
	// We take all the regex of option groups, and remove the selected ones
	// Then we add the remaining ones to the regex to remove
	allOptionGroups := make([]string, 0, len(optionGroupNumberToRegexMap))
	for _, pattern := range optionGroupNumberToRegexMap {
		allOptionGroups = append(allOptionGroups, pattern)
	}
	for _, optionGroup := range optionGroups {
		if pattern, exists := optionGroupNumberToRegexMap[string(optionGroup)]; exists {
			allOptionGroups = removeStringFromList(allOptionGroups, pattern)
		}
	}
	eventToRemoveRegex = append(eventToRemoveRegex, allOptionGroups...)

	// Compile the regex once for performance
	regexPattern := strings.Join(eventToRemoveRegex, "|")
	fmt.Println(regexPattern)
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return *ics.NewCalendar()
	}

	// Add only events that should NOT be removed
	for _, event := range cal.Events() {
		var t, err = event.GetStartAt()
		if err != nil || !(t.Equal(startDate) || t.After(startDate)) && (t.Equal(endDate) || t.Before(endDate)) {
			continue
		}

		if event.GetProperty(ics.ComponentProperty(ics.PropertySummary)) != nil {
			summary := event.GetProperty(ics.ComponentProperty(ics.PropertySummary)).Value

			// Handle special case for GLA
			GLAregex, _ := regexp.Compile("TP [12] SINBU33DL: Genie logiciel avance")
			if GLAregex.MatchString(summary) {
				summary = strings.Replace(summary, "TP", "GRP", 1)
			}

			if !regex.MatchString(summary) {
				// Keep this event (it doesn't match the removal pattern)
				filteredCal.AddVEvent(event)
			}
		}
	}

	// Return the filtered calendar
	return filteredCal
}

/*
Removes a string from a list of strings
*/
func removeStringFromList(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
