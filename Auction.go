package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"time"
)

type Auction struct {
	highestBid *Bid
	itemUp     Item_I
	endTime    time.Time
	recentBids []*Bid
}

type Bid struct {
	bidder          *ClientConnection
	durationFromEnd time.Duration
	amount          int
}

func newAuction(item Item_I) *Auction {
	a := new(Auction)
	a.endTime = time.Now().Add(time.Minute * 1) //This could also come in from user
	a.highestBid = nil
	a.itemUp = item
	a.recentBids = make([]*Bid, 5)

	return a
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

func (a *Auction) timeTillOver() time.Duration {
	return a.endTime.Sub(time.Now())
}
func (a *Auction) isOver() bool {
	return time.Now().Sub(a.endTime).Seconds() > 0
}

func (a *Auction) awardItemToWinner(winner *Bid) {
	winner.bidder.sendMsgToClient(newServerMessageS("You won the auction."))
	winner.bidder.giveItem(a.itemUp)
}

func (a *Auction) getAuctionInfo() ServerMessage {
	msg := newFormattedStringCollection()
	msg.addMessage2("\n\tItem:" + a.itemUp.getName())
	if a.highestBid != nil {
		msg.addMessage2(fmt.Sprint("\tCurrent Bid: %d", a.highestBid.amount))
	}
	msg.addMessage2("\tTime left: " + a.endTime.Sub(time.Now()).String() + "\n")

	return newServerMessageFS(msg.fmtedStrings)
}

func (a *Auction) bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString {

	estimatedTime := timeOfBid.Add(-1 * bidder.getAverageRoundTripTime())
	distance := a.endTime.Sub(estimatedTime)

	if distance > 0 {
		bid := new(Bid)
		bid.bidder = bidder
		bid.durationFromEnd = distance
		a.recentBids = append(a.recentBids, bid)

		if a.highestBid == nil || a.highestBid.amount <= amount {
			a.highestBid = bid
			return newFormattedStringSplice2(ct.Green, "Your bid was recorded for time: "+estimatedTime.String()+"\n")
		} else {
			return newFormattedStringSplice2(ct.Green, "Your bid was too low.")
		}
	} else {
		return newFormattedStringSplice2(ct.Green, "The auction ended before you could place your bid.")
	}
}
