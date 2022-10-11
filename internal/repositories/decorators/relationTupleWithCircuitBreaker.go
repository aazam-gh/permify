package decorators

import (
	"github.com/afex/hystrix-go/hystrix"
	"golang.org/x/net/context"

	"github.com/Permify/permify/internal/repositories"
	"github.com/Permify/permify/pkg/errors"
	base "github.com/Permify/permify/pkg/pb/base/v1"
	"github.com/Permify/permify/pkg/tuple"
)

// RelationTupleWithCircuitBreaker -
type RelationTupleWithCircuitBreaker struct {
	repository repositories.IRelationTupleRepository
}

// NewRelationTupleWithCircuitBreaker -.
func NewRelationTupleWithCircuitBreaker(relationTupleRepository repositories.IRelationTupleRepository) *RelationTupleWithCircuitBreaker {
	return &RelationTupleWithCircuitBreaker{repository: relationTupleRepository}
}

// Migrate -
func (r *RelationTupleWithCircuitBreaker) Migrate() (err errors.Error) {
	return nil
}

// ReverseQueryTuples -
func (r *RelationTupleWithCircuitBreaker) ReverseQueryTuples(ctx context.Context, entity string, relation string, subjectEntity string, subjectIDs []string, subjectRelation string) (Iterator tuple.ITupleIterator, err errors.Error) {
	type circuitBreakerResponse struct {
		Iterator tuple.ITupleIterator
		Error    errors.Error
	}

	output := make(chan circuitBreakerResponse, 1)

	hystrix.ConfigureCommand("relationTupleRepository.reverseQueryTuples", hystrix.CommandConfig{Timeout: 1000})
	bErrors := hystrix.Go("entityConfigRepository.reverseQueryTuples", func() error {
		tup, cErr := r.repository.ReverseQueryTuples(ctx, entity, relation, subjectEntity, subjectIDs, subjectRelation)
		output <- circuitBreakerResponse{Iterator: tup, Error: cErr}
		return nil
	}, func(err error) error {
		return nil
	})

	select {
	case out := <-output:
		return out.Iterator, out.Error
	case <-bErrors:
		return Iterator, errors.CircuitBreakerError
	}
}

// QueryTuples -
func (r *RelationTupleWithCircuitBreaker) QueryTuples(ctx context.Context, entity string, objectID string, relation string) (Iterator tuple.ITupleIterator, err errors.Error) {
	type circuitBreakerResponse struct {
		Iterator tuple.ITupleIterator
		Error    errors.Error
	}

	output := make(chan circuitBreakerResponse, 1)
	hystrix.ConfigureCommand("relationTupleRepository.queryTuples", hystrix.CommandConfig{Timeout: 1000})
	bErrors := hystrix.Go("entityConfigRepository.queryTuples", func() error {
		tup, cErr := r.repository.QueryTuples(ctx, entity, objectID, relation)
		output <- circuitBreakerResponse{Iterator: tup, Error: cErr}
		return nil
	}, func(err error) error {
		return nil
	})

	select {
	case out := <-output:
		return out.Iterator, out.Error
	case <-bErrors:
		return Iterator, errors.CircuitBreakerError
	}
}

// Read -
func (r *RelationTupleWithCircuitBreaker) Read(ctx context.Context, filter *base.TupleFilter) (collection tuple.ITupleCollection, err errors.Error) {
	type circuitBreakerResponse struct {
		Collection tuple.ITupleCollection
		Error      errors.Error
	}

	output := make(chan circuitBreakerResponse, 1)
	hystrix.ConfigureCommand("relationTupleRepository.read", hystrix.CommandConfig{Timeout: 1000})
	bErrors := hystrix.Go("entityConfigRepository.read", func() error {
		tup, cErr := r.repository.Read(ctx, filter)
		output <- circuitBreakerResponse{Collection: tup, Error: cErr}
		return nil
	}, func(err error) error {
		return nil
	})

	select {
	case out := <-output:
		return out.Collection, out.Error
	case <-bErrors:
		return collection, errors.CircuitBreakerError
	}
}

// Write -
func (r *RelationTupleWithCircuitBreaker) Write(ctx context.Context, iterator tuple.ITupleIterator) (err errors.Error) {
	outputErr := make(chan errors.Error, 1)
	hystrix.ConfigureCommand("relationTupleRepository.write", hystrix.CommandConfig{Timeout: 1000})
	bErrors := hystrix.Go("entityConfigRepository.write", func() error {
		err = r.repository.Write(ctx, iterator)
		outputErr <- err
		return nil
	}, func(err error) error {
		return nil
	})

	select {
	case err = <-outputErr:
		return err
	case <-bErrors:
		return errors.CircuitBreakerError
	}
}

// Delete -
func (r *RelationTupleWithCircuitBreaker) Delete(ctx context.Context, iterator tuple.ITupleIterator) (err errors.Error) {
	outputErr := make(chan errors.Error, 1)
	hystrix.ConfigureCommand("relationTupleRepository.delete", hystrix.CommandConfig{Timeout: 1000})
	bErrors := hystrix.Go("entityConfigRepository.delete", func() error {
		err = r.repository.Delete(ctx, iterator)
		outputErr <- err
		return nil
	}, func(err error) error {
		return nil
	})

	select {
	case err = <-outputErr:
		return err
	case <-bErrors:
		return errors.CircuitBreakerError
	}
}
