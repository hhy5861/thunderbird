package thunderbird

type RoomChannel struct {
	tb *Thunderbird
}

func (rc *RoomChannel) Received(event Event) {
	switch event.Type {
	case "message":
		rc.tb.Broadcast(event)
	}
}
