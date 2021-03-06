// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/dpolansky/grader-ci/pkg/model"
)

type FakeBuildMessageService struct {
	SendBuildStub        func(*model.BuildStatus) error
	sendBuildMutex       sync.RWMutex
	sendBuildArgsForCall []struct {
		arg1 *model.BuildStatus
	}
	sendBuildReturns struct {
		result1 error
	}
	sendBuildReturnsOnCall map[int]struct {
		result1 error
	}
	ListenForBuildMessagesStub        func(die chan struct{}) <-chan int
	listenForBuildMessagesMutex       sync.RWMutex
	listenForBuildMessagesArgsForCall []struct {
		die chan struct{}
	}
	listenForBuildMessagesReturns struct {
		result1 <-chan int
	}
	listenForBuildMessagesReturnsOnCall map[int]struct {
		result1 <-chan int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBuildMessageService) SendBuild(arg1 *model.BuildStatus) error {
	fake.sendBuildMutex.Lock()
	ret, specificReturn := fake.sendBuildReturnsOnCall[len(fake.sendBuildArgsForCall)]
	fake.sendBuildArgsForCall = append(fake.sendBuildArgsForCall, struct {
		arg1 *model.BuildStatus
	}{arg1})
	fake.recordInvocation("SendBuild", []interface{}{arg1})
	fake.sendBuildMutex.Unlock()
	if fake.SendBuildStub != nil {
		return fake.SendBuildStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.sendBuildReturns.result1
}

func (fake *FakeBuildMessageService) SendBuildCallCount() int {
	fake.sendBuildMutex.RLock()
	defer fake.sendBuildMutex.RUnlock()
	return len(fake.sendBuildArgsForCall)
}

func (fake *FakeBuildMessageService) SendBuildArgsForCall(i int) *model.BuildStatus {
	fake.sendBuildMutex.RLock()
	defer fake.sendBuildMutex.RUnlock()
	return fake.sendBuildArgsForCall[i].arg1
}

func (fake *FakeBuildMessageService) SendBuildReturns(result1 error) {
	fake.SendBuildStub = nil
	fake.sendBuildReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildMessageService) SendBuildReturnsOnCall(i int, result1 error) {
	fake.SendBuildStub = nil
	if fake.sendBuildReturnsOnCall == nil {
		fake.sendBuildReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendBuildReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBuildMessageService) ListenForBuildMessages(die chan struct{}) <-chan int {
	fake.listenForBuildMessagesMutex.Lock()
	ret, specificReturn := fake.listenForBuildMessagesReturnsOnCall[len(fake.listenForBuildMessagesArgsForCall)]
	fake.listenForBuildMessagesArgsForCall = append(fake.listenForBuildMessagesArgsForCall, struct {
		die chan struct{}
	}{die})
	fake.recordInvocation("ListenForBuildMessages", []interface{}{die})
	fake.listenForBuildMessagesMutex.Unlock()
	if fake.ListenForBuildMessagesStub != nil {
		return fake.ListenForBuildMessagesStub(die)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.listenForBuildMessagesReturns.result1
}

func (fake *FakeBuildMessageService) ListenForBuildMessagesCallCount() int {
	fake.listenForBuildMessagesMutex.RLock()
	defer fake.listenForBuildMessagesMutex.RUnlock()
	return len(fake.listenForBuildMessagesArgsForCall)
}

func (fake *FakeBuildMessageService) ListenForBuildMessagesArgsForCall(i int) chan struct{} {
	fake.listenForBuildMessagesMutex.RLock()
	defer fake.listenForBuildMessagesMutex.RUnlock()
	return fake.listenForBuildMessagesArgsForCall[i].die
}

func (fake *FakeBuildMessageService) ListenForBuildMessagesReturns(result1 <-chan int) {
	fake.ListenForBuildMessagesStub = nil
	fake.listenForBuildMessagesReturns = struct {
		result1 <-chan int
	}{result1}
}

func (fake *FakeBuildMessageService) ListenForBuildMessagesReturnsOnCall(i int, result1 <-chan int) {
	fake.ListenForBuildMessagesStub = nil
	if fake.listenForBuildMessagesReturnsOnCall == nil {
		fake.listenForBuildMessagesReturnsOnCall = make(map[int]struct {
			result1 <-chan int
		})
	}
	fake.listenForBuildMessagesReturnsOnCall[i] = struct {
		result1 <-chan int
	}{result1}
}

func (fake *FakeBuildMessageService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendBuildMutex.RLock()
	defer fake.sendBuildMutex.RUnlock()
	fake.listenForBuildMessagesMutex.RLock()
	defer fake.listenForBuildMessagesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBuildMessageService) recordInvocation(key string, args []interface{}) {
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
