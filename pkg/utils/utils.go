package utils

// formatMessageCount returns the formatted number of
// messages containing a thousands comma
func FormatMessageCount(messageCount string) string {
	if len(messageCount) == 9 {
		return messageCount[:3] + "," + messageCount[3:6] + "," + messageCount[6:]
	} else if len(messageCount) == 8 {
		return messageCount[:2] + "," + messageCount[2:5] + "," + messageCount[5:]
	} else if len(messageCount) == 7 {
		return messageCount[:1] + "," + messageCount[1:4] + "," + messageCount[4:]
	} else if len(messageCount) == 6 {
		return messageCount[:3] + "," + messageCount[3:]
	} else if len(messageCount) == 5 {
		return messageCount[:2] + "," + messageCount[2:]
	} else if len(messageCount) == 4 {
		return messageCount[:1] + "," + messageCount[1:]
	} else {
		return messageCount
	}
}

