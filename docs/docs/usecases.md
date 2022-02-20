# Usecases

This page describes popular Stencil use cases and provides related resources that you can use to create Stencil workflows.

## Event-driven architecture

Event-driven architecture is a software paradigm that promotes using events as a means of communication between decoupled services. Events are the records of a change in state, e.g. a customer booking an order, driver, a driver confirming the booking, etc. Events are immutable and are usually ordered in the sequence of their creation.

Event-driven architecture usually has three key components: producers, message brokers, and consumers. Consumer and producer services are loosely where event producers don't know which events consumers are listening to. This decoupling allows producers and consumers to evolve, scale, and deploy independently.

But this also opens a challenge for managing data schema across consumers and producers.
