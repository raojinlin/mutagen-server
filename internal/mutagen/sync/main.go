package sync

import (
	"context"
	"fmt"
	"github.com/mutagen-io/mutagen/pkg/prompting"
	"github.com/mutagen-io/mutagen/pkg/selection"
	promptingsvc "github.com/mutagen-io/mutagen/pkg/service/prompting"
	syncsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"google.golang.org/grpc"
	"time"
)

type Handler func(prompterIdentifier string) (interface{}, error)

func WithPrompting(conn *grpc.ClientConn, prompter prompting.Prompter, handler Handler) (interface{}, error) {
	promptingCtx, promptingCancel := context.WithCancel(context.Background())
	defer func() {
		promptingCancel()
	}()

	prompterIdentifier, promptingErrors, err := promptingsvc.Host(
		promptingCtx,
		promptingsvc.NewPromptingClient(conn),
		prompter,
		true,
	)

	if err != nil {
		return "", err
	}

	ret, err := handler(prompterIdentifier)
	if err != nil {
		return nil, err
	}

	errs := <-promptingErrors
	fmt.Println(errs)
	return ret, nil
}

func Create(conn *grpc.ClientConn, prompter prompting.Prompter, creation *syncsvc.CreationSpecification) (*syncsvc.CreateResponse, error) {
	resp, err := WithPrompting(conn, prompter, func(prompterIdentifier string) (interface{}, error) {
		request := &syncsvc.CreateRequest{
			Prompter:      prompterIdentifier,
			Specification: creation,
		}
		client := syncsvc.NewSynchronizationClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 9000*time.Second)
		defer cancel()

		return client.Create(ctx, request)
	})

	if resp == nil {
		return nil, err
	}

	return resp.(*syncsvc.CreateResponse), err
}

func List(conn *grpc.ClientConn, sel *selection.Selection) (*syncsvc.ListResponse, error) {
	request := &syncsvc.ListRequest{
		Selection:          sel,
		PreviousStateIndex: 0,
	}
	client := syncsvc.NewSynchronizationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()

	return client.List(ctx, request)
}

func Resume(conn *grpc.ClientConn, prompter prompting.Prompter, sel *selection.Selection) (*syncsvc.ResumeResponse, error) {
	result, err := WithPrompting(conn, prompter, func(prompterIdentifier string) (interface{}, error) {
		client := syncsvc.NewSynchronizationClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 109*time.Second)
		defer cancel()
		request := &syncsvc.ResumeRequest{
			Prompter:  prompterIdentifier,
			Selection: sel,
		}
		return client.Resume(ctx, request)
	})

	if result == nil {
		return nil, err
	}

	return result.(*syncsvc.ResumeResponse), err
}

func Pause(conn *grpc.ClientConn, prompter prompting.Prompter, sel *selection.Selection) (*syncsvc.PauseResponse, error) {
	result, err := WithPrompting(conn, prompter, func(prompterIdentifier string) (interface{}, error) {
		client := syncsvc.NewSynchronizationClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()
		request := &syncsvc.PauseRequest{
			Prompter:  prompterIdentifier,
			Selection: sel,
		}
		return client.Pause(ctx, request)
	})

	if result == nil {
		return nil, err
	}

	return result.(*syncsvc.PauseResponse), err
}

func Terminate(conn *grpc.ClientConn, prompter prompting.Prompter, sel *selection.Selection) (*syncsvc.TerminateResponse, error) {
	result, err := WithPrompting(conn, prompter, func(prompterIdentifier string) (interface{}, error) {
		client := syncsvc.NewSynchronizationClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()
		request := &syncsvc.TerminateRequest{
			Prompter:  prompterIdentifier,
			Selection: sel,
		}
		return client.Terminate(ctx, request)
	})

	if result == nil {
		return nil, err
	}

	return result.(*syncsvc.TerminateResponse), err
}

func Flush(conn *grpc.ClientConn, prompter prompting.Prompter, sel *selection.Selection) (*syncsvc.FlushResponse, error) {
	result, err := WithPrompting(conn, prompter, func(prompterIdentifier string) (interface{}, error) {
		client := syncsvc.NewSynchronizationClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()
		request := &syncsvc.FlushRequest{
			Prompter:  prompterIdentifier,
			Selection: sel,
		}
		return client.Flush(ctx, request)
	})

	if result == nil {
		return nil, err
	}

	return result.(*syncsvc.FlushResponse), err
}

func Reset(conn *grpc.ClientConn, prompter prompting.Prompter, sel *selection.Selection) (*syncsvc.ResetResponse, error) {
	result, err := WithPrompting(conn, prompter, func(prompterIdentifier string) (interface{}, error) {
		client := syncsvc.NewSynchronizationClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
		defer cancel()
		request := &syncsvc.FlushRequest{
			Prompter:  prompterIdentifier,
			Selection: sel,
		}
		return client.Flush(ctx, request)
	})

	if result == nil {
		return nil, err
	}

	return result.(*syncsvc.ResetResponse), err
}
