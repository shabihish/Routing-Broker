# Threaded-Queue-Broker
A GO queue broker implemented for efficient message passing from client to server and backwards. The broker uses multiple threads for direct communication with both the client and server, and passes the queue messages to each whenever their own queues become free. This implementation uses a 2-way queue for bidirectional message passing.
