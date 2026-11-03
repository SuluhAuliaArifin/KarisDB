package main

import (
	"context"
	"io"
	"os"
	"time"

	"karisdb/internal/app"
	"karisdb/internal/protocol"
)

func main() {
	input, err := io.ReadAll(io.LimitReader(os.Stdin, protocol.MaxRequestBytes+1))
	if err != nil {
		_, _ = os.Stderr.WriteString("karisdb-core: could not read request\n")
		os.Exit(2)
	}

	router, err := app.NewRouter()
	var response protocol.Response
	if err != nil {
		response = protocol.Response{
			ProtocolVersion: protocol.Version,
			Success:         false,
			Error:           &protocol.ErrorDetail{Code: "CORE_INIT_FAILED", Message: "core could not initialize"},
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		response = protocol.Handle(ctx, input, router)
	}

	if err := protocol.WriteResponse(os.Stdout, response); err != nil {
		_, _ = os.Stderr.WriteString("karisdb-core: could not write response\n")
		os.Exit(2)
	}
	if !response.Success {
		os.Exit(1)
	}
}
