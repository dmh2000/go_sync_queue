package queue

import (
	"testing"
)

// test variables
const rqsize int = 100


// ===================
// Race Detector Tests
// ===================

func TestChannelRace(t *testing.T) {
	async1(t,NewChannelQueue(rqsize))
	async3(t,NewChannelQueue(rqsize))
}

func TestListRace(t *testing.T) {
	async1(t,NewSyncList(rqsize))
	async3(t,NewSyncList(rqsize))
}
func TestCircularRace(t *testing.T) {
	async1(t,NewSyncCircular(rqsize))
	async3(t,NewSyncCircular(rqsize))
}

func TestRingRace(t *testing.T) {
	async1(t,NewSyncRing(rqsize))
	async3(t,NewSyncRing(rqsize))
}

func TestSliceRace(t *testing.T) {
	async1(t,NewSyncSlice(rqsize))
	async3(t,NewSyncSlice(rqsize))
}

func TestComboRace(t *testing.T) {
	async1(t,NewSyncCircular(rqsize))
	async3(t,NewSyncList(rqsize))
}

func TestNativeRace(t *testing.T) {
	async2(t,NewNativeQueue(rqsize))
	async4(t,NewNativeQueue(rqsize))
}