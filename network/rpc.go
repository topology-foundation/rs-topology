package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type RPC struct {
	ctx      context.Context
	mediator NetworkMediator
	server   *http.Server
}

func NewRPC(ctx context.Context, mediator NetworkMediator) *RPC {
	return &RPC{
		ctx:      ctx,
		mediator: mediator,
		server: &http.Server{
			Addr: ":8080", // TODO: get port from config
		},
	}
}

func (rpc *RPC) Start() {
	http.HandleFunc("/rpc", rpc.rpcMessageHandler) // Set the RPC endpoint

	fmt.Println("Starting RPC server on", rpc.server.Addr)
	if err := rpc.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
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

func (rpc *RPC) Shutdown() {
	if err := rpc.server.Shutdown(rpc.ctx); err != nil {
		fmt.Printf("RPC server shutdown error: %v\n", err)
	} else {
		fmt.Println("RPC server successfully shut down")
	}
}
