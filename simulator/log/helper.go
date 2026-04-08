package log

func MessageAndDetails(message string, details ...Simulation) Group {
	return Group{
		Message(message),
		Indent{Group(details)},
	}
}

func (m Message) AndDetails(details ...Simulation) Group {
	return Group{
		m,
		Indent{Group(details)},
	}
}
