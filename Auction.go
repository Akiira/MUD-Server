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
	msg.addMessage2("\nAttention Players! There is an auction going on. The current status of the auction is the following:\n")
	msg.addMessage2(fmt.Sprintf("\tItem: %-15s ", a.itemUp.getName()))

	if a.highestBid != nil {
		msg.addMessage2(fmt.Sprintf("Current Bid: %d     Bidder: %-15s", a.highestBid.amount, a.highestBid.bidder.getCharactersName()))
	} else {
		msg.addMessage2("\tCurrent Bid: None")
	}
	msg.addMessage2("\tTime left: " + a.endTime.Sub(time.Now()).String() + "\n")
	msg.addMessage2("If you would like to bid on this item then type 'bid [amount]' where amount is the amount of gold you want to bid.\n")

	return newServerMessageFS(msg.fmtedStrings)
}

func (a *Auction) bidOnItem(amount int, bidder *ClientConnection, timeOfBid time.Time) []FormattedString {

	estimatedTime := timeOfBid.Add(-1 * bidder.GetAverageRoundTripTime())
	distance := a.endTime.Sub(estimatedTime)

	if distance > 0 {
		bid := new(Bid)
		bid.bidder = bidder
		bid.amount = amount
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
