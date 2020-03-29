package bete

const (
	stringETACommandPrompt                = "Send me the bus stop code you wish to see arrivals for, optionally followed by services to filter by. For example: \"96049 5 24\""
	stringErrorSorry                      = "Oh no, Something went wrong! Sorry about that, we're looking into it."
	stringAddFavouritePromptForName       = "Adding the query %q to your favourites with a custom name. Send me the name for this query."
	stringAddFavouritePromptForQuery      = "Send me the ETA query you wish to save as a favourite."
	stringAddFavouriteSuggestName         = "Adding the query %q to your favourites. What would you like to name it?"
	stringAddFavouriteAdded               = "Added the query %q to your favourites as %q!"
	stringCallbackRefresh                 = "Refresh"
	stringCallbackResend                  = "Resend"
	stringFavouritesAddNew                = "Add a new favourite"
	stringFavouritesHide                  = "Hide favourites keyboard"
	stringFavouritesDelete                = "Delete an existing favourite"
	stringFavouritesSetCustomName         = "Set a custom name"
	stringFavouritesShow                  = "Show favourites keyboard"
	stringFavouritesChooseAction          = "What would you like to do?"
	stringFavouritesOnlyPrivateChat       = "Sorry, you can only manage your favourites in a private chat."
	stringDeleteFavouriteDeleted          = "Removed %s from your favourites!"
	stringDeleteFavouritesNoFavourites    = "Oops, you don't have any favourites to delete! How about creating one first?"
	stringDeleteFavouritesChoose          = "Select a favourite to delete it:"
	stringShowFavouritesNoFavourites      = "Oops, you haven't added any favourites yet! How about creating one first?"
	stringShowFavouritesShowing           = "Showing favourites keyboard"
	stringHideFavouritesHiding            = "Hiding favourites keyboard"
	stringQueryContainsInvalidCharacters  = "An ETA query should contain only letters and numbers."
	stringQueryShouldStartWithBusStopCode = "An ETA query should start with a 5-digit bus stop code."
	stringQueryTooLong                    = "An ETA query should be less than 20 characters long (sorry about that)."
	stringRefreshETAsUpdated              = "ETAs updated!"
	stringResendETAsSent                  = "ETAs sent!"
	stringSomethingWentWrong              = "Something went wrong!"
	stringFormatSwitchSummary             = "Show arriving bus summary"
	stringFormatSwitchDetails             = "Show arriving bus details"
	stringWelcomeMessage                  = `Hi %s,

Bus Eta Bot tells you when your bus will arrive!

To get started immediately, send me a 5-digit bus stop code like "02151".

If you've got a little more time, use the button below to take a guided tour of Bus Eta Bot's features. 

I hope that you will find Bus Eta Bot useful! For questions or feedback, feel free to reach out to @jiayu.`
)
