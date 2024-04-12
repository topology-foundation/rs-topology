package rpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/topology-gg/gram/config"
	ex "github.com/topology-gg/gram/execution"
	p2p "github.com/topology-gg/gram/network/p2p"
)

type RPC struct {
	ctx      context.Context
	errCh    chan error
	executor ex.Execution
	server   *http.Server
	p2p      *p2p.P2P
}

func NewRPC(ctx context.Context, errCh chan error, executor ex.Execution, config *config.RpcConfig, p2p *p2p.P2P) (*RPC, error) {
	r := mux.NewRouter()

	rpc := &RPC{
		ctx:      ctx,
		errCh:    errCh,
		executor: executor,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: r,
		},
		p2p: p2p,
	}

	r.HandleFunc("/rpc", rpc.rpcMessageHandler).Methods("POST")

	return rpc, nil
}

func (rpc *RPC) Start() {
	fmt.Println("Starting RPC server on", rpc.server.Addr)
	if err := rpc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		rpc.errCh <- err
		return
	}
}

func (rpc *RPC) Shutdown() error {
	if err := rpc.server.Shutdown(rpc.ctx); err != nil {
		return err
	}

	fmt.Println("RPC server successfully shutted down")
	return nil
}
