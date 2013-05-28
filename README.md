gopubsub
========

Provisional name for an HTTP long-polling pubsub server written in Go.

This application is pre-alpha, and is my first foray in to programming in Go.

To run:

 mkdir /var/www/html/gopubsub
 cp longpoll.html /var/www/html/gopubsub
 go run gopubsub.go

Then browser to:

 http://127.0.0.1:4000/longpoll.html

You will get periodic ping messages appearing on your browser, simulating events occuring.

