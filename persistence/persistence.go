package persistence

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dandeliondeathray/nona/game"

	"github.com/coreos/etcd/clientv3"
)

// RecoveryHandler handles recovery from persisted state, after a restart.
type RecoveryHandler interface {
	OnRoundRecovered(seed int64)
}

type Persistence struct {
	team string
}

func (p *Persistence) StoreNewRound(seed int64) {
	err := p.setCurrentRound(seed)
	if err != nil {
		log.Printf("Failed to set new round: %v", err)
	}
}

func (p *Persistence) ResolvePlayerState(player game.Player, resolution game.PlayerStateResolution) {
	f := func() {
		resolution.PlayerStateResolved(game.PlayerState{0})
	}
	go f()
}

func (p *Persistence) Recover(handler RecoveryHandler) error {
	log.Printf("Recovering round information for team %s", p.team)
	seed, err := p.getCurrentRound()
	if err != nil {
		log.Printf("Failed to recover round: %v", err)
		return err
	}
	handler.OnRoundRecovered(seed)
	return nil
}

func (p *Persistence) getCurrentRound() (int64, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
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
	seed, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 64)
	if err != nil {
		return 0, err
	}
	return seed, nil
}

func (p *Persistence) setCurrentRound(seed int64) error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
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

func NewPersistence(team string) *Persistence {
	return &Persistence{team}
}
