package blocker

//Start Runs the all mighty massblocker
func Start() {
	retriveNotBlockedhandsets()
	serviceConsumerDispatcher()
	imeisToBlockFileGen()
}
