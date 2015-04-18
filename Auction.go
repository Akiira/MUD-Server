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

	return a
}

func (a *Auction) determineWinner() *Bid {
	return a.highestBid
}

func (a *Auction) timeTillOver() time.Duration {
	return a.endTime.Sub(time.Now())
}
func (a *Auction) isOver() bool {
	return time.Now().Sub(a.endTime).Seconds() > 0
}

func (a *Auction) awardItemToWinner(winner *Bid) {
	winner.bidder.Write(newServerMessageS("You won the auction."))
	winner.bidder.giveItem(a.itemUp)
}

func (a *Auction) getAuctionInfo() ServerMessage {
	msg := newFormattedStringCollection()
	msg.addMessage2("\n\tItem:" + a.itemUp.getName())
	if a.highestBid != nil {
		msg.addMessage2(fmt.Sprint("\tCurrent Bid: ", a.highestBid.amount))
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

		if a.highestBid == nil || a.highestBid.amount <= amount {
			a.highestBid = bid
			return newFormattedStringSplice2(ct.Green, "Your bid was recorded for time: "+estimatedTime.String()+"\n")
		} else {
			return newFormattedStringSplice2(ct.Red, "Your bid was too low.")
		}
	} else {
		return newFormattedStringSplice2(ct.Red, "The auction ended before you could place your bid.")
	}
}
