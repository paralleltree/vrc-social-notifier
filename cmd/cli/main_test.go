package main

import (
	"reflect"
	"sync"
	"testing"

	"github.com/paralleltree/vrc-social-notifier/testlib"
)

func TestMakeChGenerator_BroadcastsReveivedElementsFromSourceChannel(t *testing.T) {
	t.Parallel()

	// arrange
	wantValues := []int{10, 20}
	sourceCh := make(chan int)
	chGenerator := makeChGenerator(sourceCh)

	ch1 := chGenerator()
	ch2 := chGenerator()

	wg := &sync.WaitGroup{}
	assertFunc := func(t *testing.T, ch chan int, wantValues []int) {
		t.Helper()
		defer wg.Done()
		got := testlib.ConsumeChannel(ch)
		if !reflect.DeepEqual(wantValues, got) {
			t.Errorf("unexpected value: want %v, but got %v", wantValues, got)
		}
	}

	// act
	wg.Add(2)
	// assert received values asynchronusly
	go assertFunc(t, ch1, wantValues)
	go assertFunc(t, ch2, wantValues)
	sourceCh <- wantValues[0]
	sourceCh <- wantValues[1]

	close(sourceCh)
	wg.Wait()
}
