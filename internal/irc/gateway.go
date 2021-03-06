package main

import (
	"emersyx.net/common/pkg/api"
	"github.com/BurntSushi/toml"
	goirc "github.com/fluffle/goirc/client"
)

// gateway is the type which defines an irc.Gateway implementation, namely an IRC resource and receptor for the emersyx
// platform.
type gateway struct {
	api.PeripheralBase
	api      *goirc.Conn
	config   *goirc.Config
	messages chan api.Event
}

// NewPeripheral creates a new api.IRCGateway instance and applies to configuration specified in the arguments.
func NewPeripheral(opts api.PeripheralOptions) (api.Peripheral, error) {
	var err error

	// create a new gateway and initialize the base
	gw := new(gateway)
	gw.InitializeBase(opts)

	// create the messages channel
	gw.messages = make(chan api.Event)

	// create a Config object for the underlying library
	gw.config = goirc.NewConfig("placeholder")

	// override several default values from the underlying library
	gw.config.Me.Ident = "emersyx"
	gw.config.Me.Name = "emersyx"
	gw.config.Version = "emersyx"
	gw.config.SSL = false
	gw.config.QuitMessage = "bye"

	// standard function for generating new nicks
	gw.config.NewNick = func(n string) string { return n + "^" }

	// apply the extended options from the config file
	c := new(config)
	if _, err = toml.DecodeFile(opts.ConfigPath, c); err != nil {
		return nil, err
	}
	if err = c.validate(); err != nil {
		return nil, err
	}
	c.apply(gw)

	// create the underlying Conn object
	gw.api = goirc.Client(gw.config)

	// initialize callbacks
	gw.initCallbacks()

	// connect to the server
	gw.connect()

	return gw, nil
}

// initCallbacks sets the callback functions for the internally used goirc library.
func (gw *gateway) initCallbacks() {
	gw.api.HandleFunc(goirc.PRIVMSG, channelCallback(gw))
	gw.api.HandleFunc(goirc.JOIN, channelCallback(gw))
	gw.api.HandleFunc(goirc.QUIT, channelCallback(gw))
	gw.api.HandleFunc(goirc.PART, channelCallback(gw))
}
