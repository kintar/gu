package ddbutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type QueryResult[T any] struct {
	Items []T
	Error error
}

var (
	ErrQueryContextCanceled = errors.New("context canceled")
	ErrQueryDdbClientError  = errors.New("dynamodb client returned error")
	ErrQueryCannotUnmarshal = errors.New("failed to unmarshal item")
)

// QueryAsync is a generic worker that creates a result channel and a goroutine which will run the provided
// dynamodb.QueryInput in a loop until LastEvaluatedKey is empty or the context is canceled. Each page of results is
// collected into a slice and returned on the channel in a QueryResult. When the worker is done, the channel is closed.
// Note that T must be a struct that can be unmarshalled by the attributevalue package. Passing in a pointer or interface
// type will cause the first unmarshal call to abort the worker and return an error.
// The channel has a buffer of 2, so if somehow your processing is slower than the I/O with DynamoDB, the worker will
// stall until you read from the channel.
// All error values that can be placed on the QueryResult will be one of the ErrQuery* types defined in this package.
func QueryAsync[T any](ctx context.Context, client *dynamodb.Client, input dynamodb.QueryInput) <-chan QueryResult[T] {
	out := make(chan QueryResult[T], 2)

	go func() {
		defer close(out)

		p := 0

		for {
			select {
			case <-ctx.Done():
				out <- QueryResult[T]{Error: ErrQueryContextCanceled}
				return
			default:
				r, e := client.Query(ctx, &input)
				if e != nil {
					out <- QueryResult[T]{Error: fmt.Errorf("%w: %v", ErrQueryDdbClientError, e)}
					return
				}
				result := QueryResult[T]{Items: make([]T, 0, len(r.Items))}
				for _, item := range r.Items {
					var t T
					if err := attributevalue.UnmarshalMap(item, &t); err != nil {
						result.Error = fmt.Errorf("%w: %v", ErrQueryCannotUnmarshal, err)
						out <- result
						return
					}
					result.Items = append(result.Items, t)
				}
				out <- result
				input.ExclusiveStartKey = r.LastEvaluatedKey
				// if there is no data in the start key, there's no more data to read
				if len(input.ExclusiveStartKey) == 0 {
					return
				}
				p++
			}
		}
	}()

	return out
}
