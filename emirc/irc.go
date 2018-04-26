package main

import (
	"emersyx.net/emersyx/api"
	"errors"
	"time"
)

// GetIdentifier returns the identifier of this receptor.
func (gw *ircGateway) GetIdentifier() string {
	return gw.identifier
}

// GetEventsInChannel returns the api.CoreEvent channel through which core events are received by the gateway
// instance.
func (gw *ircGateway) GetEventsInChannel() chan<- api.Event {
	// TODO implement a system where core events are received and appropriate actions are taken (i.e. disconnect from
	// the server when receiving the shutdown event).
	return nil
}

// GetEventsOutChannel returns the api.Event channel through which emersyx events are pushed by this gateway.
func (gw *ircGateway) GetEventsOutChannel() <-chan api.Event {
	return (<-chan api.Event)(gw.messages)
}

// Connect start the connection process of the IRC gateway to the server. This is a blocking call. If the gateway
// connects to the server without errors, then nil is returned. Otherwise an error with the appropriate message is
// returned.
func (gw *ircGateway) Connect() error {
	gw.log.Debugln("connecting to the server")
	err := gw.api.Connect()
	for gw.IsConnected() != true {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
	gw.log.Debugln("connected to the server")
	return err
}

// IsConnected returned true if the gateway is connected to the IRC server, otherwise it returns false.
func (gw *ircGateway) IsConnected() bool {
	return gw.api.Connected()
}

// Quit disconnects the IRC gateway from the server. If the gateway disconnects from the server without errors, then nil
// is returned. Otherwise an error with the appropriate message is returned.
func (gw *ircGateway) Quit() error {
	gw.log.Debugln("quitting to the server")
	err := gw.api.Close()
	return err
}

// Join sends the command for the IRC gateway to join a channel. The channel is specified in the method argument. If the
// gateway joins the channel without errors, then nil is returned. Otherwise an error with the appropriate message is
// returned.
func (gw *ircGateway) Join(ch string) error {
	if gw.IsConnected() {
		gw.log.Debugf("joining the \"%s\" channel\n", ch)
		gw.api.Join(ch)
		return nil
	}
	return errors.New("the ircGateway instance is not connected to any server")
}

// Privmsg sends either a message to an IRC channel or a private message to another user, depending on the method
// argument.
func (gw *ircGateway) Privmsg(to, msg string) error {
	if gw.IsConnected() {
		gw.log.Debugf("sending a PRIVMSG to \"%s\"\n", to)
		gw.api.Privmsg(to, msg)
		return nil
	}
	return errors.New("the ircGateway instance is not connected to any server")
}