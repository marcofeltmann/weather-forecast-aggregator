# The package `internal`

I love the idea of ArdanLab's Bill Kennedy to structure Go directories in a way
that reflects the Domain-Design of the problem domain.

His workshops use this [ardanlabs/service](https://github.com/ardanlabs/service)
sample layout.

Unfortunately it's pretty hard to understand what this means if you're used to
the standard Go application project layouts out there in the world.

As the `internal` package is guarded from import in other products it looks like
the perfect place to store all the app's **business** logic.

If another system needs to import this app's business logic something changed,
so we should have a conversation about porting that function into a **foundation**
level `pkg`.
