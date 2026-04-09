package client

import (
	"context"
	"errors"
	"io"
	"net/url"
)

func (h *HttpClient) Send(ctx context.Context, payload Payload) (Response, error) {
	// Create request
	request, err := h.createRequest(h.url, payload)
	if err != nil {
		return Response{}, err
	}

	// Create context with max timeout
	ctx, cancel := context.WithTimeout(ctx, h.timeoutAndRetryOption.maxWaitPerTry)
	defer cancel()
	request = request.WithContext(ctx) // limit each HTTP request timeout option for per try

	// Execute request
	resp, err := h.pool.Do(request)
	if err != nil {
		// Handle error
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return Response{}, ErrOverflowMaxWait
		}

		// evaluate if err is caused by connection timeout
		var urlErr *url.Error
		if !errors.As(err, &urlErr) || !urlErr.Timeout() {
			if errors.Is(err, context.Canceled) {
				return Response{}, ErrOperationContextCanceled
			}
			return Response{}, err
		}

		// check if we need to retry on timeout or not
		if !h.timeoutAndRetryOption.retryOnTimeout {
			return Response{}, ErrTimeout
		}

		return Response{}, ErrTimeout
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return Response{}, ErrOverflowMaxWait
		}

		// evaluate if err is caused by connection timeout
		if errors.Is(err, context.Canceled) {
			return Response{}, ErrOperationContextCanceled
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return Response{}, ErrTimeout
		}

		return Response{}, err
	}

	return Response{
		Status: resp.StatusCode,
		Body:   respBody,
		Header: resp.Header,
	}, nil
}
