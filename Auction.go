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
	determineWinner() *Bid
	awardItemToWinner(winner *Bid)
	isOver() bool
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

	for {
		time.Sleep(time.Second * 2)

		if time.Now().Sub(a.endTime).Seconds() > 0 {
			break
		}
	}

	winner := a.determineWinner()

	if winner != nil {
		a.awardItemToWinner(winner)
	}
}

func (a *Auction) determineWinner() *Bid {
	var highestBid *Bid

	for _, bid := range a.recentBids {
		if bid.amount > highestBid.amount {
			highestBid = bid
		}

		if bid.amount == highestBid.amount && bid.durationFromEnd < highestBid.durationFromEnd {
			highestBid = bid
		}
	}

	return highestBid
}

func (a *Auction) awardItemToWinner(winner *Bid) {

}

func (a *Auction) getAuctionInfo() *ServerMessage {
	return nil //TODO
}

func (a *Auction) bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString {
	bid := new(Bid)
	bid.bidder = bidder

	estimatedTime := timeOfBid.Add(-1 * bidder.getAverageRoundTripTime())

	bid.durationFromEnd = a.endTime.Sub(estimatedTime)

	msg := "Your bid was recorded for time: " + estimatedTime.String() + "\n"
	return newFormattedStringSplice2(ct.Green, msg)
}
