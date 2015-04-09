package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"time"
)

type Auction struct {
	highestBid *Bid
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

func (a *Auction) isOver() bool {
	return time.Now().Sub(a.endTime).Seconds() > 0
}

func (a *Auction) awardItemToWinner(winner *Bid) {
	winner.bidder.sendMsgToClient(newServerMessageS("You won the auction."))
	winner.bidder.giveItem(a.itemUp)
}

func (a *Auction) getAuctionInfo() ServerMessage {
	msg := newFormattedStringCollection()
	msg.addMessage2("\tItem:" + a.itemUp.name)
	msg.addMessage2(fmt.Sprint("\tCurrent Bid: %d", a.highestBid.amount))
	msg.addMessage2("\tTime left: " + a.endTime.Sub(time.Now()).String())

	return newServerMessageFS(msg.fmtedStrings)
}

func (a *Auction) bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString {

	if !a.isOver() && a.highestBid.amount <= amount {
		bid := new(Bid)
		bid.bidder = bidder

		estimatedTime := timeOfBid.Add(-1 * bidder.getAverageRoundTripTime())
		bid.durationFromEnd = a.endTime.Sub(estimatedTime)

		a.highestBid = bid
		a.recentBids = append(a.recentBids, bid)

		msg := "Your bid was recorded for time: " + estimatedTime.String() + "\n"
		return newFormattedStringSplice2(ct.Green, msg)
	} else if a.highestBid.amount > amount {
		return newFormattedStringSplice2(ct.Green, "Your bid was too low.")
	} else {
		return newFormattedStringSplice2(ct.Green, "The auction ended before you could place your bid.")
	}
}
