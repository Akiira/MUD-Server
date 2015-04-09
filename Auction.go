package main

import (
	"github.com/daviddengcn/go-colortext"
	"time"
)

//TODO just messing around with interfaces, prolly wont keep it
type Auctioner interface {
	newAuction(item *Item, endTime time.Time)
	beginAuction()
	sendAuctionInfoToWorld()
	bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString
	determineWinner()
	awardItemToWinner()
}

type Auction struct {
	highestBid *Bid //TODO think about best type for this
	itemUp     *Item
	endTime    time.Time
	recentBids []*Bid
}

type Bid struct {
	bidder          *ClientConnection
	durationFromEnd time.Duration
	amount          int
}

func newAuction(item *Item, endTime time.Time) {
	a := new(Auction)
	a.endTime = endTime
	a.highestBid = nil
	a.itemUp = item
	a.recentBids = make([]*Bid, 5)
}

func (a *Auction) bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString {
	bid := new(Bid)
	bid.bidder = bidder

	estimatedTime := timeOfBid.Add(-1 * bidder.getAverageRoundTripTime())

	bid.durationFromEnd = a.endTime.Sub(estimatedTime)

	msg := "Your bid was recorded for time: " + estimatedTime.String() + "\n"
	return newFormattedStringSplice2(ct.Green, msg)
}
