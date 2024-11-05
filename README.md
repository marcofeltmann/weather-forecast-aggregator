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

## [golang.org/x/sync](https://pkg.go.dev/golang.org/x/sync) by Google

The sync package provides ErrorGroups that act like WaitGroups, but they are
able to collect errors that happened along the path.

As a bunch of things might fail during HTTP and Network I prefer ErrorGroups for
orchestrating GoRoutines.

### Risk Analysis

The `golang.org/x` package is the "area to mature before it makes it into the
standard library".

It is kind of the staging area for the standard library, so I don't see any
other risks beside it's disappearance after being transferred into the standard
library.

### Risk Minimization Options

Using `go mod vendor` enables the usage of the `x/sync` package without code
changes even after it was transferred to the standard library.

## [github.com/google/go-cmp](https://pkg.go.dev/github.com/google/go-cmp@v0.6.0) by The Go Authors

This packages provides an easier and more reliable way to compare structs and
other complex data types for semantical equality than a self-coded solution with
the standard library's `reflect` package does.

### Risk Analysis

It's not directly bound to the official Google repositories, only created by the
Go Authors.  
The BSD license does not exclude a potencial disappearance of this package.

### Risk Minimization Options

Using `go mod vendor` enables the usage of the `co-cmp` package even when it was
discontinued.

Limiting usage on testing reduces the amount of rewrite work of the functionality.

## Primer: Manage Application Dependencies

As I use Linux for a pretty long time and I believe managing all the tools for a
specific project seems to be a lot of manual amount I switched to the Nix package
manager a couple of years ago.

Here you describe your preferred system in the Nix language like you'd describe
your Kubernetes resources in a YAML file.
Nix Package Manager does all the things needed to obtain the tools you described.

If you combine Nix with the great `direnv` application, simply changing into the
project directory via terminal will begin downloading all the required tools.

Without `direnv` you could simply `nix-shell ./default.nix` to get a shell with
all the required application dependencies.

Finally, if you need to know which software you need but don't want to fire up 
Nix, simply peek into the `default.nix` file.

## Application [make](https://www.gnu.org/software/make/) by GNU

The GNU Make Build System is used to reduce the mental load on the build process.
Instead of remembering all the tools and how to call them for a desired output,
all you need to remember is the usage of Make and how to read the `Makefile`.

### Risk Analysis

There is nearly no risk at all. As GNU Make is licensed under the GPL there is
nearly 0 chance of inaccessibility in the near or far future.  
Nearly everything of GNU and Linux is built using the GNU build toolchain, with
GNU Make being an essential part of.

Even if all the FOSS developers decide to no longer work on this thing, the
current state of software will stay available.

### Risk Minimization Options

Once you installed GNU Make, never uninstall it.  
That should be safe enough for the next 20 years.

## Application [golangci-lint](https://golangci-lint.run/) by golangci

Golang continuous integration linter collection to get automated suggestions on
well-known common mistakes and security issues.

### Risk Analysis

The application is licensed under GPL, so it is very unlikely to disappear.  
On the other hand it contains a lot of shared code from ~300 other projects, so 
it includes a total of 13 other licenses, as analyzed by [FOSSA](https://app.fossa.io/projects/git%2Bgithub.com%2Fgolangci%2Fgolangci-lint).

There is a high risk of subdependencies doing evil things, as I don't think the
maintainers can keep track of each and every dependency in the stack.

### Risk Minimization Options

The golangci-lint only runs locally before `git push` on the developer's machine.  
It's not going to run in kubernetes or anywhere on a production server.

This is a short-running tool scanning over source files, so any malicious code
has just a few seconds to execute.

The codebase itself is hosted in a public repository on GitHub, so the risk of
source code leeking kinda doesn't exist.

## Application [k6](https://k6.io/) by Grafana

Grafana k6 is an end-to-end testing solution for kubernetes and web services.  
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
