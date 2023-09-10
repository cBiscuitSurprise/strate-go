package storage

import (
	"context"
	"time"

	"github.com/cBiscuitSurprise/strate-go/internal/game"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type StrategoRedisClient struct {
	name    string
	client  *redis.Client
	context context.Context
}

func NewStrategoRedisClient(name string) *StrategoRedisClient {
	return &StrategoRedisClient{
		name: name,
	}
}

func (r *StrategoRedisClient) IsConnected() bool {
	return r.client != nil
}

func (r *StrategoRedisClient) Connect() {
	r.context = context.Background()
	r.client = redis.NewClient(&redis.Options{
		Addr:       "localhost:6379",
		Password:   "", // no password set
		DB:         0,  // use default DB
		ClientName: r.name,
	})
}

func (r *StrategoRedisClient) PublishPieceMoveEvent(gameId string, move game.Move) {
	r.client.XAdd(r.context, &redis.XAddArgs{
		Stream: GameMoveStreamKey(gameId),
		Values: map[string]string{"move": game.SerializeMove(move)},
	})
}

// ListenForPieceMoveEvent is intended to run as a goroutine
// all moves received from the stream are sent to the `moves` channel
// if `true` is sent to the `quit` channel, we will exit the loop
func (r *StrategoRedisClient) ListenForPieceMoveEvent(gameId string, sinceId string, moves chan game.Move, quit chan bool) {
	currentId := sinceId
	for {
		for {
			select {
			case <-quit:
				return
			default:
				if response, err := r.client.XRead(r.context, &redis.XReadArgs{
					Streams: []string{GameMoveStreamKey(gameId), currentId},
					Block:   time.Duration(time.Duration.Minutes(1)),
				}).Result(); err == nil {
					for _, result := range response {
						for _, message := range result.Messages {
							currentId = message.ID

							log.Trace().
								Str("method", "StrategoRedisClient.ListenForPieceMoveEvent").
								Str("id", message.ID).
								Str("gameId", gameId).
								Msg("received move from stream")

							if messageStr, ok := message.Values["move"].(string); ok {
								if move, err := game.DeserializeMove(messageStr); err == nil {
									moves <- move
								} else {
									log.Error().
										Err(err).
										Str("method", "StrategoRedisClient.ListenForPieceMoveEvent").
										Str("message", messageStr).
										Str("gameId", gameId).
										Msg("failed to deserialize move")
									return
								}
							} else {
								log.Warn().
									Str("method", "StrategoRedisClient.ListenForPieceMoveEvent").
									Str("gameId", gameId).
									Type("messageType", message.Values["move"]).
									Msg("recieved non-string `move` value from move stream... skipping")
							}
						}
					}
				} else {
					return
				}
			}
		}
	}
}
