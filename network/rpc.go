package network

import (
	"bufio"
	"context"
	"io"
	"os"
)

type RPC struct {
	ctx      context.Context
	mediator NetworkMediator
}

func NewRPC(ctx context.Context, mediator NetworkMediator) *RPC {
	return &RPC{
		ctx:      ctx,
		mediator: mediator,
	}
}

func (rpc *RPC) Start() {
	// TODO: Start RPC server.

	rpc.rpcMessageHandler()
}

func (rpc *RPC) rpcMessageHandler() {
	reader := bufio.NewReader(os.Stdin)

	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			continue
		}
		if err != nil {
			panic(err)
		}

		rpc.mediator.MessageHandler(message, SourceRPC)
	}
}
