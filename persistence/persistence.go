package persistence

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dandeliondeathray/nona/game"

	"github.com/coreos/etcd/clientv3"
)

//go:generate mockgen -destination=../mock/mock_persistence.go -package=mock github.com/dandeliondeathray/nona/persistence RecoveryHandler

// RecoveryHandler handles recovery from persisted state, after a restart.
type RecoveryHandler interface {
	OnRoundRecovered(seed int64)
}

type Persistence struct {
	team      string
	seed      int64
	endpoints []string
}

func (p *Persistence) StoreNewRound(seed int64) {
	err := p.setCurrentRound(seed)
	if err != nil {
		log.Printf("Failed to set new round: %v", err)
		return
	}
	p.seed = seed
}

func (p *Persistence) ResolvePlayerState(player game.Player, resolution game.PlayerStateResolution) {
	f := func() {
		playerState, skipped, err := p.getPlayerState(player)
		if err == nil {
			resolution.PlayerStateResolved(game.PlayerState{PuzzleIndex: playerState, Skipped: skipped})
		} else {
			log.Println("Error when getting player state: ", err)
		}
	}
	go f()
}

func (p *Persistence) ResolveAllPlayerStates(resolution game.AllPlayerStatesResolution) {
	f := func() {
		states, err := p.getAllPlayerStates()
		if err == nil {
			resolution.AllPlayerStatesResolved(states)
		} else {
			log.Println("Error when getting all player states:", err)
		}
	}
	go f()
}

func (p *Persistence) PlayerSolvedPuzzle(player game.Player, newPuzzleIndex int) {
	f := func() {
		err := p.setPlayerState(player, newPuzzleIndex)
		if err != nil {
			log.Println("Error when getting player state: ", err)
		}
	}
	go f()
}

func (p *Persistence) PlayerSkippedPuzzle(player game.Player, newPuzzleIndex int, skipped int) {
	f := func() {
		err := p.setPlayerStateWithSkipped(player, newPuzzleIndex, skipped)
		if err != nil {
			log.Println("Error when getting player state: ", err)
		}
	}
	go f()
}

func (p *Persistence) Recover(handler RecoveryHandler, done chan<- bool) {
	f := func() {
		log.Printf("Recovering round information for team %s", p.team)
		seed, err := p.getCurrentRound()
		if err != nil {
			log.Printf("Failed to recover round: %v", err)
			done <- false
		} else {
			p.seed = seed
			handler.OnRoundRecovered(seed)
			done <- true
		}
	}
	go f()
}

func (p *Persistence) getClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   p.endpoints,
		DialTimeout: 5 * time.Second,
	})
}

func (p *Persistence) getCurrentRound() (int64, error) {
	cli, err := p.getClient()
	if err != nil {
		return 0, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	kvc := clientv3.NewKV(cli)
	resp, err := kvc.Get(ctx, fmt.Sprintf("%s/current_round", p.team))
	cancel()
	if err != nil {
		return 0, err
	}
	// use the response
	log.Printf("Get Response is: %v", resp)
	if len(resp.Kvs) == 0 {
		return 0, fmt.Errorf("No round set")
	}
	seed, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
	if err != nil {
		return 0, err
	}
	return seed, nil
}

func (p *Persistence) setCurrentRound(seed int64) error {
	cli, err := p.getClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	kvc := clientv3.NewKV(cli)
	resp, err := kvc.Put(ctx, fmt.Sprintf("%s/current_round", p.team), fmt.Sprintf("%d", seed))
	cancel()
	if err != nil {
		return err
	}
	// use the response
	log.Printf("Put Response is: %v", resp)
	return nil
}

func (p *Persistence) getPlayerState(player game.Player) (int, int, error) {
	cli, err := p.getClient()
	if err != nil {
		return 0, 0, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	kvc := clientv3.NewKV(cli)
	resp, err := kvc.Get(ctx, fmt.Sprintf("%s/%d/index/%s", p.team, p.seed, player))
	respSkipped, err := kvc.Get(ctx, fmt.Sprintf("%s/%d/skipped/%s", p.team, p.seed, player))
	cancel()
	if err != nil {
		return 0, 0, err
	}
	// use the response
	log.Printf("Player state Response is: %v", resp)
	log.Printf("Player state skipped Response is: %v", respSkipped)
	if len(resp.Kvs) == 0 {
		return 0, 0, nil
	}
	index, err := strconv.Atoi(string(resp.Kvs[0].Value))
	if err != nil {
		return 0, 0, err
	}

	skipped := 0
	if len(respSkipped.Kvs) != 0 {
		skipped, err = strconv.Atoi(string(respSkipped.Kvs[0].Value))
		if err != nil {
			return 0, 0, err
		}
	}

	return index, skipped, nil
}

func (p *Persistence) setPlayerState(player game.Player, index int) error {
	cli, err := p.getClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	kvc := clientv3.NewKV(cli)
	resp, err := kvc.Put(ctx, fmt.Sprintf("%s/%d/index/%s", p.team, p.seed, player), fmt.Sprintf("%d", index))
	cancel()
	if err != nil {
		return err
	}
	// use the response
	log.Printf("Put Response is: %v", resp)
	return nil
}

func (p *Persistence) setPlayerStateWithSkipped(player game.Player, index int, skipped int) error {
	cli, err := p.getClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	kvc := clientv3.NewKV(cli)
	resp, err := kvc.Put(ctx, fmt.Sprintf("%s/%d/index/%s", p.team, p.seed, player), fmt.Sprintf("%d", index))
	respSkipped, err := kvc.Put(ctx, fmt.Sprintf("%s/%d/skipped/%s", p.team, p.seed, player), fmt.Sprintf("%d", skipped))
	cancel()
	if err != nil {
		return err
	}
	// use the response
	log.Printf("Put Response is: %v", resp)
	log.Printf("Put Skipped Response is: %v", respSkipped)
	return nil
}

func (p *Persistence) getAllPlayerStates() (map[game.Player]game.PlayerState, error) {
	cli, err := p.getClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	kvc := clientv3.NewKV(cli)
	resp, err := kvc.Get(ctx, fmt.Sprintf("%s/%d/", p.team, p.seed), clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	// use the response
	log.Printf("All player states response is: %v", resp)
	if len(resp.Kvs) == 0 {
		return map[game.Player]game.PlayerState{}, nil
	}

	result := make(map[game.Player]game.PlayerState)
	for _, kv := range resp.Kvs {
		isIndex, err := isIndexKey(string(kv.Key))
		if err == nil {
			player, err := playerFromIndexKey(string(kv.Key))
			_, ok := result[player]
			if !ok {
				result[player] = game.PlayerState{PuzzleIndex: 0, Skipped: 0}
			}
			if err == nil {
				if isIndex {
					index, err := strconv.Atoi(string(kv.Value))
					if err == nil {
						currentState := result[player]
						currentState.PuzzleIndex = index
						result[player] = currentState
					} else {
						log.Printf("getAllPlayerStates: Player %s: Could not parse puzzle index: %s", player, string(kv.Value))
					}
				} else {
					skipped, err := strconv.Atoi(string(kv.Value))
					if err == nil {
						currentState := result[player]
						currentState.Skipped = skipped
						result[player] = currentState
					} else {
						log.Printf("getAllPlayerStates: Player %s: Could not parse puzzles skipped: %s", player, string(kv.Value))
					}
				}
			} else {
				log.Printf("getAllPlayerStates: %v", err)
			}
		}
	}
	return result, nil
}

func NewPersistence(team string, endpoints []string) *Persistence {
	return &Persistence{team, 0, endpoints}
}

func playerFromIndexKey(key string) (game.Player, error) {
	tokens := strings.Split(key, "/")
	if len(tokens) != 4 {
		return game.Player(""), fmt.Errorf("Expected key on form team/seed/index/player, but got wrong number of tokens %v", tokens)
	}
	return game.Player(tokens[3]), nil
}

func isIndexKey(key string) (bool, error) {
	tokens := strings.Split(key, "/")
	if len(tokens) != 4 {
		return true, fmt.Errorf("Expected key on form team/seed/index/player, but got wrong number of tokens %v", tokens)
	}
	return tokens[2] == "index", nil
}
