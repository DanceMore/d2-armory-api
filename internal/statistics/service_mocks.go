// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package statistics

import (
	"context"
	"github.com/nokka/d2-armory-api/internal/domain"
	"sync"
)

// Ensure, that statisticsRepositoryMock does implement statisticsRepository.
// If this is not the case, regenerate this file with moq.
var _ statisticsRepository = &statisticsRepositoryMock{}

// statisticsRepositoryMock is a mock implementation of statisticsRepository.
//
// 	func TestSomethingThatUsesstatisticsRepository(t *testing.T) {
//
// 		// make and configure a mocked statisticsRepository
// 		mockedstatisticsRepository := &statisticsRepositoryMock{
// 			GetByCharacterFunc: func(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
// 				panic("mock out the GetByCharacter method")
// 			},
// 			UpsertFunc: func(ctx context.Context, stat domain.StatisticsRequest) error {
// 				panic("mock out the Upsert method")
// 			},
// 		}
//
// 		// use mockedstatisticsRepository in code that requires statisticsRepository
// 		// and then make assertions.
//
// 	}
type statisticsRepositoryMock struct {
	// GetByCharacterFunc mocks the GetByCharacter method.
	GetByCharacterFunc func(ctx context.Context, character string) (*domain.CharacterStatistics, error)

	// UpsertFunc mocks the Upsert method.
	UpsertFunc func(ctx context.Context, stat domain.StatisticsRequest) error

	// calls tracks calls to the methods.
	calls struct {
		// GetByCharacter holds details about calls to the GetByCharacter method.
		GetByCharacter []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Character is the character argument value.
			Character string
		}
		// Upsert holds details about calls to the Upsert method.
		Upsert []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Stat is the stat argument value.
			Stat domain.StatisticsRequest
		}
	}
	lockGetByCharacter sync.RWMutex
	lockUpsert         sync.RWMutex
}

// GetByCharacter calls GetByCharacterFunc.
func (mock *statisticsRepositoryMock) GetByCharacter(ctx context.Context, character string) (*domain.CharacterStatistics, error) {
	if mock.GetByCharacterFunc == nil {
		panic("statisticsRepositoryMock.GetByCharacterFunc: method is nil but statisticsRepository.GetByCharacter was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Character string
	}{
		Ctx:       ctx,
		Character: character,
	}
	mock.lockGetByCharacter.Lock()
	mock.calls.GetByCharacter = append(mock.calls.GetByCharacter, callInfo)
	mock.lockGetByCharacter.Unlock()
	return mock.GetByCharacterFunc(ctx, character)
}

// GetByCharacterCalls gets all the calls that were made to GetByCharacter.
// Check the length with:
//     len(mockedstatisticsRepository.GetByCharacterCalls())
func (mock *statisticsRepositoryMock) GetByCharacterCalls() []struct {
	Ctx       context.Context
	Character string
} {
	var calls []struct {
		Ctx       context.Context
		Character string
	}
	mock.lockGetByCharacter.RLock()
	calls = mock.calls.GetByCharacter
	mock.lockGetByCharacter.RUnlock()
	return calls
}

// Upsert calls UpsertFunc.
func (mock *statisticsRepositoryMock) Upsert(ctx context.Context, stat domain.StatisticsRequest) error {
	if mock.UpsertFunc == nil {
		panic("statisticsRepositoryMock.UpsertFunc: method is nil but statisticsRepository.Upsert was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Stat domain.StatisticsRequest
	}{
		Ctx:  ctx,
		Stat: stat,
	}
	mock.lockUpsert.Lock()
	mock.calls.Upsert = append(mock.calls.Upsert, callInfo)
	mock.lockUpsert.Unlock()
	return mock.UpsertFunc(ctx, stat)
}

// UpsertCalls gets all the calls that were made to Upsert.
// Check the length with:
//     len(mockedstatisticsRepository.UpsertCalls())
func (mock *statisticsRepositoryMock) UpsertCalls() []struct {
	Ctx  context.Context
	Stat domain.StatisticsRequest
} {
	var calls []struct {
		Ctx  context.Context
		Stat domain.StatisticsRequest
	}
	mock.lockUpsert.RLock()
	calls = mock.calls.Upsert
	mock.lockUpsert.RUnlock()
	return calls
}
