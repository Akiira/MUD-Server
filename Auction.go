package main

import (
	"github.com/daviddengcn/go-colortext"
	"time"
)

//TODO just messing around with interfaces, prolly wont keep it
type Auctioner interface {
	newAuction(item *Item) *Auction
	beginAuction()
	getAuctionInfo() *ServerMessage
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

func newAuction(item *Item) *Auction {
	a := new(Auction)
	a.endTime = time.Now().Add(time.Minute * 2) //This could also come in from user
	a.highestBid = nil
	a.itemUp = item
	a.recentBids = make([]*Bid, 5)

	return a
}

func (a *Auction) beginAuction() {

}

func (a *Auction) getAuctionInfo() *ServerMessage {

}

func (a *Auction) bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString {
	bid := new(Bid)
	bid.bidder = bidder

	estimatedTime := timeOfBid.Add(-1 * bidder.getAverageRoundTripTime())

	bid.durationFromEnd = a.endTime.Sub(estimatedTime)

	msg := "Your bid was recorded for time: " + estimatedTime.String() + "\n"
	return newFormattedStringSplice2(ct.Green, msg)
}
