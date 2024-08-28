package matchmaker

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

type Storage interface {
	AddPlayer(player Player)
	GetPlayers() []Player
	RemovePlayers(players []Player)
	PlayerCount() int
}

// In-memory implementation
type InMemoryStorage struct {
	players []Player
}

type PostgresStorage struct {
	db *sql.DB
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		players: []Player{},
	}
}

func (s *InMemoryStorage) AddPlayer(player Player) {
	s.players = append(s.players, player)
}

func (s *InMemoryStorage) GetPlayers() []Player {
	return s.players
}

func (s *InMemoryStorage) RemovePlayers(players []Player) {
	var newPlayers []Player
	for _, p := range s.players {
		shouldRemove := false
		for _, rp := range players {
			if p.Name == rp.Name {
				shouldRemove = true
				break
			}
		}
		if !shouldRemove {
			newPlayers = append(newPlayers, p)
		}
	}
	s.players = newPlayers
}

func (s *InMemoryStorage) PlayerCount() int {
	return len(s.players)
}

func NewPostgresStorage(cfg DBConfig) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) AddPlayer(player Player) {
	_, err := s.db.Exec(
		"INSERT INTO players (name, skill, latency, queue_time) VALUES ($1, $2, $3, $4)",
		player.Name, player.Skill, player.Latency, time.Now(),
	)
	if err != nil {
		log.Fatalf("Failed to add player to database: %v", err)
	}
}

// Реализация метода GetPlayers
func (s *PostgresStorage) GetPlayers() []Player {
	rows, err := s.db.Query("SELECT name, skill, latency, queue_time FROM players")
	if err != nil {
		log.Fatalf("Failed to retrieve players from database: %v", err)
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var player Player
		var queueTime time.Time
		if err := rows.Scan(&player.Name, &player.Skill, &player.Latency, &queueTime); err != nil {
			log.Fatalf("Failed to scan player data: %v", err)
		}
		player.QueueTime = queueTime
		players = append(players, player)
	}
	return players
}

// Реализация метода RemovePlayers
func (s *PostgresStorage) RemovePlayers(players []Player) {
	for _, player := range players {
		_, err := s.db.Exec("DELETE FROM players WHERE name = $1", player.Name)
		if err != nil {
			log.Fatalf("Failed to remove player from database: %v", err)
		}
	}
}

// Реализация метода PlayerCount
func (s *PostgresStorage) PlayerCount() int {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM players").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to count players in database: %v", err)
	}
	return count
}
