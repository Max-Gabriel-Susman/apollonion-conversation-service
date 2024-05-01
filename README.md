# Apollonion Conversation Service

DONE - set the stage for channels:
	store identifiers and broadcast messages to all clients
	the client that sends the message should not receive it back
	test coverage
	should open up the door for using channels and then we can
	bring back in a good implementation of fanout orchestration

NEED TO START - controlling concurrency:
	cap goroutines, but still accept more clients to connect
	maybe store a slice of net.Conn to sorta cache them so we can get
	around to them later while staying under the goroutine cap
	test coverage
failure case coverage for existing and new logic:


[05/01/2024 notes]

wrap goroutines in an anonymous function always

mutex uneeded here because:

we're synchronously boradcasting message

so we'll have to wait for that message to be sent to every connection

parallel/async?

We need to actually write it back to the client sort of like an event for the client to receive

net.Conn has a writer maybe? it has a reader and the writer, uses bytes()

instead of mutex, in broadcaster() use a sync.Waitgroup and then fanout the inner for loop(faster)

just parallize inner for sake of simplicity(outer should be irl tho)
    - waitgroup instead of mutex
    - cap concurrency to n goroutines(use a buffered channel as a semaphore)
    - 

handleConnection:
    use conn.Read() function instead of scanner.Scan

we're gonna need a goroutine per connection

can be optimized

potential race condition: you potentially send to a connection after it's been closed before it's been removed from connection slice

global lock defeats purpose of concurrency

don't lock all connections, lock the slice when appending and deleting 

you want to have individual locks on each connection 

in the context of concurrency most data races come from two process writing to the same data at the same time

might not need to here

rework locking/unlocking strategy for sake of performance

so we can easily test and mock this we should better organize this into packages and then create an interface for the client that defines read and write so we can moq it up and make this generally better 

shoot for production readiness as a chat server