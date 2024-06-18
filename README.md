# golang-websocket-protocol

Implementing my own websocket using the net package in GO

### What it does
A server can receive a request from the client to establish a tcp connection and to connect a user to a user-defined namespace, a namespace contains other clients that can send and receive messages pushed by other clients to the server, the server only pushes messages from to the clients that are connected to the same namespace, clients outside the namespace will be ignored.


### Why I built this
I was studying Computer Networking A top-down Approach and was really facinated in controlling how these sockets work for two network processes that are communicating with each other, and how both these entities can conform to any specific structure (application layer protocols) to establish a well defined communication link between these two network process.
