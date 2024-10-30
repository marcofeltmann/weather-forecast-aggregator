# Usage

You simply `go run cmd/server --weather-api-key==<your key>` and in another 
terminal (session) just
`curl 'http://localhost:8080/weather?lat=42.6493934&lon=-8.8201753'`
if you want to see the forecast for my region.

# Metrics

Although this is a single sample server app running on your device instead of
"The Cloud" I still feel unlucky without metering at least some performance
indexes.

For the sake of lightweight solutions and small dependency trees I'll stick to
exporting them via the [expvar](https://pkg.go.dev/expvar) standard library
package.

For viewing the metrics I highly recommend the awesome [expvarmon](https://github.com/divan/expvarmon)
project.

# Dependencies

Although it is risky to depend on the work of others it helps alot to speed
things up.

## [github.com/ardanlabs/conf/v3](https://pkg.go.dev/github.com/ardanlabs/conf/v3) by ArdanLabs

This package enables configuration of your Go binaries by taking parameters,
env vars and config files into account, by following the best practices for
respecting the order of these configurations.

It also enables quick access to help and version information on CLI apps.

### Risk Analysis

As ArdanLabs makes a living as a software and kubernetes consulting company
it is doubtful the code will disappear, but there is a change.

ArdanLabs might decide to bound the access to this package to one of
their subscriptions. That way the whole package might disappear.

To me both cases seem really unlikely to happen.

### Risk Minimization Options

Forking the repository
can reduce the impact of this risk by adding technical dept with keeping up to
package updates

Using `go mod vendor` copys the required code of that package into this 
repository so we have life-time access to the used version with the little
overhead of a slightly larger app repository

## Application [k6](https://k6.io/) by Grafana

Grafana k6 is an end-to-end testing solution with focus on kubernetes and web services.  
It has a broad feature set with performance tests, API and UI tests, stress and
spike tests.

It is able to communicate with API, GraphQL, WebSockets and gRPC services and is
able to manage tresholds or handle cookies.


### Risk Analysis

As Grafana makes a living with testing and monitoring system-as-a-service
offers they know that less-featured FOSS options are great to get customers
hooked to the paid services.

It feels very unlikely to me that the OpenSource version of k6 will disappear.  
But there never is a guarantee.

### Risk Minimization Options

Checking in all the node modules to keep it running after `npm install` would
reduce the risk alot. 
As this is such a small project and QA isn't really involved I'll take the risk
to rewrite all the blackbox tests in Golang whenever the worst case happens.
Not sure if I'm still interested in this project at that time, though.

# License

Normally I choose *public domain* or *creative commons 0* for samples like these.

If I felt really funny I also used to use the *what the fuck you want license*
or the *beer* respective *coffee license* as there is not a high level of 
ingenuity in this sample.

But there is **AI** and at least in Germany one precedent decided that the `AGPLv3`
is a binding license according to law.  
So with this approach there might be some fun in the future.

For further information on the whole topic feel free to listen to Dylan Beattie
for an hour on his talk [Open Source, Open Mind: The Cost of Free Software](https://www.youtube.com/watch?v=vzYqxo13I1U&t=1174s)
at NDC Oslo 2024.


# Assorted Thoughts

uber MAXPROCS, cloud deployment
