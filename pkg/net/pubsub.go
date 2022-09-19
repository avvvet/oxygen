package net

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type PubsubInst struct {
	Host         *host.Host
	PubSub       *pubsub.PubSub
	Topic        *pubsub.Topic
	TopicName    string
	Subscription *pubsub.Subscription
}

type BroadcastData struct {
	Type string
	Data []byte
}

func InitPubSub(ctx context.Context, h host.Host, topicName string) (*PubsubInst, error) {
	pubsub, err1 := pubsub.NewGossipSub(ctx, h)
	if err1 != nil {
		return nil, err1
	}

	topic, err2 := pubsub.Join(topicName)
	if err2 != nil {
		return nil, err2
	}

	subscription, err3 := topic.Subscribe()
	if err2 != nil {
		return nil, err3
	}
	return &PubsubInst{&h, pubsub, topic, topicName, subscription}, nil
}
