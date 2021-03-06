// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/dpolansky/grader-ci/pkg/model"
)

type FakeGithubWebhookService struct {
	HandleRequestStub        func(*model.GithubWebhookRequest) error
	handleRequestMutex       sync.RWMutex
	handleRequestArgsForCall []struct {
		arg1 *model.GithubWebhookRequest
	}
	handleRequestReturns struct {
		result1 error
	}
	handleRequestReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGithubWebhookService) HandleRequest(arg1 *model.GithubWebhookRequest) error {
	fake.handleRequestMutex.Lock()
	ret, specificReturn := fake.handleRequestReturnsOnCall[len(fake.handleRequestArgsForCall)]
	fake.handleRequestArgsForCall = append(fake.handleRequestArgsForCall, struct {
		arg1 *model.GithubWebhookRequest
	}{arg1})
	fake.recordInvocation("HandleRequest", []interface{}{arg1})
	fake.handleRequestMutex.Unlock()
	if fake.HandleRequestStub != nil {
		return fake.HandleRequestStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.handleRequestReturns.result1
}

func (fake *FakeGithubWebhookService) HandleRequestCallCount() int {
	fake.handleRequestMutex.RLock()
	defer fake.handleRequestMutex.RUnlock()
	return len(fake.handleRequestArgsForCall)
}

func (fake *FakeGithubWebhookService) HandleRequestArgsForCall(i int) *model.GithubWebhookRequest {
	fake.handleRequestMutex.RLock()
	defer fake.handleRequestMutex.RUnlock()
	return fake.handleRequestArgsForCall[i].arg1
}

func (fake *FakeGithubWebhookService) HandleRequestReturns(result1 error) {
	fake.HandleRequestStub = nil
	fake.handleRequestReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithubWebhookService) HandleRequestReturnsOnCall(i int, result1 error) {
	fake.HandleRequestStub = nil
	if fake.handleRequestReturnsOnCall == nil {
		fake.handleRequestReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.handleRequestReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithubWebhookService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.handleRequestMutex.RLock()
	defer fake.handleRequestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGithubWebhookService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
