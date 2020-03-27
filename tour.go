package bete

// Tour section slugs
const (
	tourSectionStart               = "start"
	tourSectionETAQueries          = "eta_queries"
	tourSectionFilteringETAQueries = "filtering_eta_queries"
	tourSectionRefreshResend       = "refresh_resend"
	tourSectionArrivingBusDetails  = "bus_details"
	tourSectionFavourites          = "favourites"
	tourSectionInlineQueries       = "inline_queries"
	tourSectionFinish              = "finish"
)

// TourSectionData contains what to display for a tour section.
type TourSectionData struct {
	Text       string
	Navigation []TourSectionNavigation
}

// TourSectionNavigation describes a link between tour sections.
type TourSectionNavigation struct {
	Text   string
	Target string
}

type Tour map[string]TourSectionData

var tour = Tour{
	tourSectionStart: {
		Text: stringTourStart,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Next: " + stringTourTitleETAQueries,
				Target: tourSectionETAQueries,
			},
		},
	},
	tourSectionETAQueries: {
		Text: stringTourETAQueries,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Next: " + stringTourTitleFilteringETAQueries,
				Target: tourSectionFilteringETAQueries,
			},
		},
	},
	tourSectionFilteringETAQueries: {
		Text: stringTourFilteringETAQueries,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Next: " + stringTourTitleRefreshResend,
				Target: tourSectionRefreshResend,
			},
		},
	},
	tourSectionRefreshResend: {
		Text: stringTourRefreshResend,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Next: " + stringTourTitleArrivingBusDetails,
				Target: tourSectionArrivingBusDetails,
			},
		},
	},
	tourSectionArrivingBusDetails: {
		Text: stringTourArrivingBusDetails,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Next: " + stringTourTitleFavourites,
				Target: tourSectionFavourites,
			},
		},
	},
	tourSectionFavourites: {
		Text: stringTourFavourites,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Next: " + stringTourTitleInlineQueries,
				Target: tourSectionInlineQueries,
			},
		},
	},
	tourSectionInlineQueries: {
		Text: stringTourInlineQueries,
		Navigation: []TourSectionNavigation{
			{
				Text:   "Finish tour",
				Target: tourSectionFinish,
			},
		},
	},
}
