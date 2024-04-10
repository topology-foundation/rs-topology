package network

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/log"
)

type RPC struct {
	ctx      context.Context
	errCh    chan error
	mediator NetworkMediator
	server   *http.Server
	mux      *http.ServeMux
}

func NewRPC(ctx context.Context, errCh chan error, mediator NetworkMediator, config *config.RpcConfig) (*RPC, error) {
	mux := http.NewServeMux()
	rpc := &RPC{
		ctx:      ctx,
		errCh:    errCh,
		mediator: mediator,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Port),
			Handler: mux,
		},
		mux: mux,
	}

	mux.HandleFunc("/rpc", rpc.rpcMessageHandler)

	return rpc, nil
}

func (rpc *RPC) Start() {
	log.Info("(RPC Server)", "address", rpc.server.Addr)
	if err := rpc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		rpc.errCh <- err
		return
	}
}

func (rpc *RPC) rpcMessageHandler(w http.ResponseWriter, req *http.Request) {
	// Only accept POST or GET requests
	if req.Method == http.MethodPost {
		// Read the message from the body of the request
		body, err := io.ReadAll(req.Body) //ioutil.ReadAll has deprecated
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()

		message := string(body)

		// Process the message using the mediator's MessageHandler method
		rpc.mediator.MessageHandler(message, SourceRPC)

		// Optionally, write a response back to the client
		fmt.Fprintf(w, "Message processed: %s", message)
	} else if req.Method == http.MethodGet {
		message := req.URL.Query().Get("message")
		rpc.mediator.MessageHandler(message, SourceRPC)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func (rpc *RPC) Shutdown() error {
	if err := rpc.server.Shutdown(rpc.ctx); err != nil {
		return err
	}

	log.Info("RPC server successfully shutted down")
	return nil
}
