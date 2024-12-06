package event

const (
	EventLinkCreated = "link.created"
	EventLinKVisited = "link.visited"
)

type Event struct {
	Type    string
	Payload any
}

type EventBus struct {
	bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		bus: make(chan Event),
	}
}

func (eventBus *EventBus) Publish(event Event) {
	eventBus.bus <- event
}

func (eventBus *EventBus) Subscribe() <-chan Event {
	return eventBus.bus
}
