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

To get started, send me a 5-digit bus stop code like "02151". For a tour of what else you can do with Bus Eta Bot, use the /tour command. 

I hope that you will find Bus Eta Bot useful! For questions or feedback, feel free to reach out to @jiayu.`
	stringTourStart = `Hey there! This tour will bring you through Bus Eta Bot's major features.

1. ETA queries
2. Filtering ETA queries
3. Refresh and resend
4. Arriving bus details
5. Favourites
6. Inline queries

Use the button below to get started:`
	stringTourTitleETAQueries = "ETA queries"
	stringTourETAQueries      = `<strong>ETA queries</strong>

To get ETAs for a bus stop, just send me the 5-digit bus stop code directly. There's no need to use a command like /eta. 

Give it a try now by sending me the text "02151".

When you're done, come back to this message and hit the next button.`
	stringTourTitleFilteringETAQueries = "Filtering ETA queries"
	stringTourFilteringETAQueries      = `<strong>Filtering ETA queries</strong>

If you're only interested in ETAs for some of the buses at a bus stop, you can specify them as well.

Try sending me the text "02151 36 106". You will only see ETAs for those two services.

Use the next button to move on.`
	stringTourTitleRefreshResend = "Refresh and resend"
	stringTourRefreshResend      = `<strong>Refresh and resend</strong>

The Refresh button on an ETA message updates the ETAs in the message. The time an ETA message was sent displayed by Telegram does not change, so refer to the last updated time in the message instead.

The Resend button also updates ETAs, but sends a new message to the chat instead of editing the old message. Use it when you're worried about ETA messages getting lost in your message history.

Try refreshing and resending the ETAs you previously requested.`
	stringTourTitleArrivingBusDetails = "Arriving bus details"
	stringTourArrivingBusDetails      = `<strong>Arriving bus details</strong>

Switch to the arriving bus details view using the "Show arriving bus details" button. This shows additional information about each incoming bus: whether it is single or double decked, how full it is and whether it is wheelchair accessible.

You can switch back to the summary view from the arriving bus details view using the "Show arriving bus summary" button.`
	stringTourTitleFavourites = "Favourites"
	stringTourFavourites      = `<strong>Favourites</strong>

If you frequently make certain ETA queries, you can add them as favourites for quicker access. Use the /favourites command to manage your favourites.`
	stringTourTitleInlineQueries = "Inline queries"
	stringTourInlineQueries      = `<strong>Inline queries</strong>

You can send ETA messages to any chat, including this one, using Bus Eta Bot's inline mode. Give it a try now by starting a message with @BusEtaBot followed by a space---you should see a list of bus stops appear. Selecting any of them will send an ETA message to the current chat.`
)
