package rpc

import (
	"fmt"
	"io"
	"net/http"
)

func (rpc *RPC) rpcMessageHandler(w http.ResponseWriter, req *http.Request) {
	// Read the message from the body of the request
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	message := string(body)

	// Process the message
	rpc.executor.Execute(message)

	// Publish to p2p network
	rpc.p2p.Publish(fmt.Sprintf("%s: %s", rpc.p2p.HostId(), message))

	// Optionally, write a response back to the client
	fmt.Fprintf(w, "Message processed: %s", message)
}
