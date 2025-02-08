package wa_service

var WAMessageTemplates = map[string]interface{}{
	"defaultMessage": `*Hey, please send any of the following options*
    *1. Register yourself in the app*
    *2. Summary of the transactions*
    *3. Contact us*
    
After registering you can start sending your transactions by following ways:
    1. Share the transaction from payment applications
    2. Share debit statement SMS from the bank`,

	"defaultErrorMessage": `Oops, something went sideways! Blame it on the glitch gremlins. We'll sort it out soon! 🚨`,

	"register": map[string]string{
		"success":            `You're registered ✅ 🎊, Now you can start sending transactions our way 🤝`,
		"already_registered": `Welcome back 👋, You're already registered! Start sending transactions our way 🤝`,
		"error":              `Some error occurred 😬, Please try again after some time or contact us at 🤷‍♂️`,
	},

	"transactions": map[string]string{
		"invalid_input":        `Invalid input`,
		"input_received":       `Transaction received ✅`,
		"media_download_error": `Error in processing this transaction, Please try again`,
		"one_day_maxed_out":    `Whoa there, speed racer! You've hit today’s limit. But don't worry, the message fairy will refill your stash tomorrow. 📨✨`,
	},

	"testSeries": map[string]string{
		"input_received": `Added ✅`,
	},
}
