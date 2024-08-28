package matchmaker

import (
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

type Player struct {
	Name      string
	Skill     float64
	Latency   float64
	QueueTime time.Time
}

type Matcher struct {
	storage   Storage
	groupSize int
	mu        sync.Mutex
}

func NewMatcher(storage Storage, groupSize int) *Matcher {
	return &Matcher{
		storage:   storage,
		groupSize: groupSize,
	}
}

func (m *Matcher) AddPlayer(player Player) {
	player.QueueTime = time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()

	m.storage.AddPlayer(player)

	if m.storage.PlayerCount() >= m.groupSize {
		m.formGroup()
	}
}

func (m *Matcher) formGroup() {
	players := m.storage.GetPlayers()

	sort.Slice(players, func(i, j int) bool {
		return players[i].QueueTime.Before(players[j].QueueTime)
	})

	group := players[:m.groupSize]
	m.storage.RemovePlayers(group)

	m.printGroupStats(group)
}

func (m *Matcher) printGroupStats(players []Player) {
	var (
		minSkill, maxSkill, sumSkill         = math.MaxFloat64, -math.MaxFloat64, 0.0
		minLatency, maxLatency, sumLatency   = math.MaxFloat64, -math.MaxFloat64, 0.0
		minQueueTime, maxQueueTime, sumQueue = math.MaxFloat64, -math.MaxFloat64, 0.0
	)

	now := time.Now()

	for _, p := range players {
		if p.Skill < minSkill {
			minSkill = p.Skill
		}
		if p.Skill > maxSkill {
			maxSkill = p.Skill
		}
		sumSkill += p.Skill

		if p.Latency < minLatency {
			minLatency = p.Latency
		}
		if p.Latency > maxLatency {
			maxLatency = p.Latency
		}
		sumLatency += p.Latency

		queueTime := now.Sub(p.QueueTime).Seconds()
		if queueTime < minQueueTime {
			minQueueTime = queueTime
		}
		if queueTime > maxQueueTime {
			maxQueueTime = queueTime
		}
		sumQueue += queueTime
	}

	log.Printf("Group formed:\n")
	log.Printf("Players: %v\n", getPlayerNames(players))
	log.Printf("Skill - min: %.2f, max: %.2f, avg: %.2f\n", minSkill, maxSkill, sumSkill/float64(len(players)))
	log.Printf("Latency - min: %.2f, max: %.2f, avg: %.2f\n", minLatency, maxLatency, sumLatency/float64(len(players)))
	log.Printf("Queue Time - min: %.2f, max: %.2f, avg: %.2f\n", minQueueTime, maxQueueTime, sumQueue/float64(len(players)))
}

func getPlayerNames(players []Player) []string {
	names := make([]string, len(players))
	for i, p := range players {
		names[i] = p.Name
	}
	return names
}
