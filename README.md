# CQRS/ES

This is my attempt at a simple CQRS/ES framework. It doesn't offer much
at the moment and certainly isn't production ready.

You can find an example of how you might use the framework in the
`examples` directory.

This framework only ships with a very simple in-memory event store
(`stores.NewMemoryStore()`). You could try implementing your own store
which satisfies the `cqrs_es.Store` interface if you want something
with persistence.