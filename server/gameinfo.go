package server

import (
	"fmt"

	"github.com/rumblefrog/go-a2s"
)

type Player struct {
	Name  string  `json:"name"`
	Score uint32  `json:"score"`
	Time  float32 `json:"time"`
}

type GameInfoStruct struct {
	State       bool     `json:"state"`
	Error       string   `json:"error"`
	PlayerCount uint8    `json:"player_count"`
	MaxPlayers  uint8    `json:"max_players"`
	CurrentMap  string   `json:"current_map"`
	Players     []Player `json:"players"`
}

func (s *Server) GetGameInfo() GameInfoStruct {
	ip, port := s.cfg.Allocations.DefaultMapping.Ip, s.cfg.Allocations.DefaultMapping.Port
	client, err := a2s.NewClient(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return returnError(false, err)
	}
	defer client.Close()

	info, err := client.QueryInfo() // QueryInfo, QueryPlayer, QueryRules
	if err != nil {
		return returnError(false, err)
	}
	players, err := client.QueryPlayer()
	if err != nil {
		return returnError(false, err)
	}

	GameInfo := GameInfoStruct{
		State:       true,
		Error:       "",
		PlayerCount: info.Players,
		MaxPlayers:  info.MaxPlayers,
		CurrentMap:  info.Map,
	}

	for _, player := range players.Players {
		GameInfo.Players = append(GameInfo.Players, Player{
			Name:  player.Name,
			Score: player.Score,
			Time:  player.Duration,
		})
	}

	return GameInfo
}

func returnError(state bool, err error) GameInfoStruct {
	return GameInfoStruct{
		State: state,
		Error: err.Error(),
	}
}
