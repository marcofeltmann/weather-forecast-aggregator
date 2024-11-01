# Architectural Decision Records (ADR)

Inspired by Michael Nygard's idea of [Architectural Decision Records](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html)
I'll collect mine in this single file.

Since there is no team around I'll simply `accept` every decision in the first
place.

----

# ADR-00: Frank Sinatra

No, I shouldn't name ADRs like that. This ADR kinda shows why. ;)

## Context

There are a lot of apps out there that try to do something similar.  
i.e. The brilliant [wego terminal client](https://github.com/schachmat/wego)
provides an `json` output.

## Decision

I'll do it my way.

## Status

Accepted.

## Consequences

### Positive

- escalate with all the fancy tech stuff
- abandon it as soon as it becomes too much work

### Negative

- some life time flies by that could have been used in a different way

----

# ADR-01: WeatherAPI API key is a server configuration

## Context

To access the API provided by <https://www.weatherapi.com/> an API key is required.
This key can be handed over in different ways, at least:

- configuration via app start
- request parameter
- part of the request header

## Decision

The WeatherAPI API key is provided as a server configuration during app startup.

## Status

Accepted.

## Consequences

### Positive

- request URIs stay short
- easy to extend when new backends like [MeteoGalicia](https://www.meteogalicia.gal/web/home) should be added that require a key as well
- if the server itself will require an API key it can be provided inside the request without much confusion

### Negative

- API keys are credentials and I won't upload mine to give the users a quick-start, so everyone needs to get one on their own

----

# ADR-02: Metrics via expvars

## Context

This app will be a server. So it requires some metrics. It makes me feel better
this way.

There are a lot of options out there, just to name the two biggest players:

1. Prometheus Exporter
2. OpenTelemetry

There is a builtin option called *expvars*, short for **exp**orted **var**iable**s**.
Yet exporting variables might be considered a risk for leaking data.

## Decision

Metrics are provided via expvars on the associated endpoint `/debug/vars`.

## Status

Accepted.

## Consequences

### Positive

- having the best practice of metrics even on the local running small server
- avoiding OpenTelemetry dependencies reduces code complexity
- avoiding Prometheus exporter reduces server binary size

### Negative

- doesn't scale in the cloud or on the cluster
- I didn't find any Grafana importer for these metrics
- might lead to wrong metrics during loops

## Related

- ADR-03: Meter according to RED

----

# ADR-03: Meter according to RED

## Context

When it comes to meter performance indexes there are a lot of options to follow.  
USE for Utilization, Saturation, Errors  
RED for Rate, Errors, Duration  
The whole Google Book on Site Reliability Engineering (SRE)

## Status

Accepted.

## Consequences

### Positive

- The focus of RED is on the service performance itself and this is a service
- Reducing the amount of metrics enables cleaner code that's easier to understand

### Negative

- As soon as the service is going cloud-native the metrics won't suffice anymore, tech debt incoming

## Related

- ADR-02: Metrics via expvars

## Further Information

A good summary with further references is [The USE and RED method](https://pagertree.com/learn/devops/what-is-observability/use-and-red-method)

# ADR-04: Whitebox testing with standard library's `testing` package

## Context

Whitebox testing is testing code parts directly with access to that code.  
You can look at all the parts, call functions and methods directly and test in
insolation. The most famous whitebox test is a Unit Test.

There are at least to different approaches on whitebox testing in Go:

- The internal [`testing`](https://pkg.go.dev/testing) package
- Third-party packages implementing `assert`-like functions, like Mat Ryer's
[`testify`](https://pkg.go.dev/github.com/stretchr/testify)

## Decision

Use standard library's `testing` package for whitebox testing.

## Status

Accepted.

## Consequences

### Positive

- less dependencies
- no need to learn new internal API that's not used in any productive code

### Neutral

- might be unfamiliar for people coming from other languages like Java or Swift


# ADR-05: Blackbox tests with Grafana k6

## Context

Blackbox testing is done without accessing the code. They require the software
under testing to run, accept inputs and generate outputs.

For blackbox testing I only know two options:

- Use the `testing` framework with benchmarks
- If you develop a web service there is [k6 by Grafana](https://k6.io/)

## Decision

Use Grafana k6 for blackbox testing.

## Consequences

### Positive

- k6 has a broad adoption rate with a lot of samples to plug'n'play with
- the open source version runs local in npm, so it can be integrated into CI/CD
- the SaaS version runs on cloud-scale, so the test cases are future-proof

### Neutral

- test cases are written in `json`, not in Golang
- this implies more complexety for the Golang devs
- but enables outsourcing the task to a QA department

### Negative

- additional dependency
- requires node and npm whose dependency tree is considered highly intransparent

----

# ADR-06: Skip automated API documentation

## Context

It is a pretty good idea to document APIs with OpenAPI, Swagger or similar tools.
This way consumers of the API can get a comprehensive overview and even test the
API.

Devs can also use that documentation to generate their own clients and servers.

This introduces some overhead and dependencies, though. The sample and it's API
is considered rather simple.

## Decision

I'll skip that overhead for my one or two endpoints.

## Status

Accepted

## Consequences

### Positive

- less dependencies
- more time to hack

## Neutral

- can be generated later from the rather simple code
- can't imagine there's somebody sad 'cuz I not expect anyone to re-invent this 
server

----

# ADR-07: Use structured logging with `log/slog`

## Context

When it comes to app logging there are a few options and two mutually exclusive
opinions: 

- Logging is meant for people, so write prose.
- Logs are read by machines, so structure them in machine-readable formats

## Decision

I use structured logging with the Go `log/slog` standard library package.

## Status

Accepted.

## Consequences

### Positive

- Even if it's structured the logs are human-readable
- Because it is structured it's easy to extend information to it

### Neutral

- By changing a few lines the `log/slog` can be configured to output JSON
- Not a real drop-in replacement for `log`, but I don't have legacy code yet

