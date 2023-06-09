package actions

func StartActions() {
	go LoginMonitor()
	go CustomActions()
	go OperatorMonitor()
}
