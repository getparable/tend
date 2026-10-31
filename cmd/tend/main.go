// Copyright (c) 2026 Parable, LLC
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/getparable/tend/internal/ipc"
	"github.com/getparable/tend/internal/protocol"
	"github.com/getparable/tend/internal/service"
)

func main() { os.Exit(run()) }
func run() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: tend serve|request [--socket PATH]")
		return 2
	}
	f := flag.NewFlagSet(os.Args[1], flag.ContinueOnError)
	path := f.String("socket", ipc.DefaultPath(), "private Unix socket")
	if f.Parse(os.Args[2:]) != nil || f.NArg() != 0 {
		return 2
	}
	switch os.Args[1] {
	case "serve":
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer stop()
		if err := service.Serve(ctx, *path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	case "request":
		var req protocol.Request
		var resp protocol.Response
		if err := protocol.Read(os.Stdin, &req); err != nil {
			resp = protocol.Failure(req.ID, "invalid_request", err.Error())
		} else if e := protocol.Validate(req); e != nil {
			resp = protocol.Failure(req.ID, e.Code, e.Message)
		} else {
			resp = ipc.Call(*path, req)
		}
		if err := protocol.Write(os.Stdout, resp); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if resp.Error != nil {
			return 1
		}
		return 0
	default:
		fmt.Fprintln(os.Stderr, "unknown command")
		return 2
	}
}
