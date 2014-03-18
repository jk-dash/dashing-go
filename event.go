package dashing

import (
    _ "log"
)

type Event struct {
    Id string
    Body map[string]interface{}
    Target string
}

type Broker struct {
    // Create a map of clients, the keys of the map are the channels
    // over which we can push messages to attached clients. (The values
    // are just booleans and are meaningless)
    clients map[chan *Event]bool

    // Channel into which new clients can be pushed
    newClients chan chan *Event

    // Channel into which disconnected clients should be pushed
    defunctClients chan chan *Event

    // Channel into which events are pushed to be broadcast out
    // to attached clients
    events chan *Event
}

func (b *Broker) Start() {
    go func() {
        // Loop endlessly
        for {
            // Block until we receive from one of the
            // three following channels.
            select {
            case s := <-b.newClients:
                // There is a new client attached and we
                // want to start sending them events.
                b.clients[s] = true
                // log.Println("Added new client")
            case s := <-b.defunctClients:
                // A client has detached and we want to
                // stop sending them events.
                delete(b.clients, s)
                // log.Println("Removed client")
            case event := <-b.events:
                // There is a new event to send. For each
                // attached client, push the new event
                // into the client's channel.
                for s, _ := range b.clients {
                    s <- event
                }
                // log.Printf("Broadcast event to %d clients", len(b.clients))
            }
        }
    }()
}

func NewBroker() *Broker {
    return &Broker{
        make(map[chan *Event]bool),
        make(chan (chan *Event)),
        make(chan (chan *Event)),
        make(chan *Event),
    }
}
