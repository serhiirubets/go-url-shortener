package stat

import (
	"api/pkg/event"
	"log"
)

type StatServiceDeps struct {
	EventBus *event.EventBus
	StatRepo *StatRepository
}

type StatService struct {
	EventBus *event.EventBus
	StatRepo *StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus: deps.EventBus,
		StatRepo: deps.StatRepo,
	}
}

func (s *StatService) AddClick() {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type == event.EventLinkVisited {
			id, ok := msg.Data.(uint)
			if !ok {
				log.Fatalln("Bad EventLinkVisited Data: ", msg.Data)
			}
			s.StatRepo.AddClick(id)
		}
	}
}
