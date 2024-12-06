package stat

import "link-shortner-api/pkg/event"

type StatServiceDependencies struct {
	EventBus *event.EventBus
	*StatRepository
}

type StatService struct {
	EventBus *event.EventBus
	*StatRepository
}

func NewStatService(dependencies *StatServiceDependencies) *StatService {
	return &StatService{
		EventBus:       dependencies.EventBus,
		StatRepository: dependencies.StatRepository,
	}
}

// метод, бесконечно слушает события и вызывает методы репозитория
func (s *StatService) ListenLinkEvents() {
	for {
		message := <-s.EventBus.Subscribe()
		switch message.Type {
		case event.EventLinKVisited:
			s.StatRepository.AddClick(message.Payload.(uint)) // продьюсили идентификатор ссылки, но лучше сделать маппинг типов типов сообщений со значениями payload
		}
	}
}
