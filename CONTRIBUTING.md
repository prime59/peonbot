# Development

This bot was developed in an Ubuntu 18 virtual machine with the following
tools:
> go version go1.11.10 linux/amd64

This development guide is for **linux only**. I am not knowledgeable enough
to advise how to develop Go programs on Windows. I'm sure there are plenty
of guides on the internet, if you don't have access to a linux environment.

## Dev Environment
The easiest way to create a dev environment is use the
[Makefile](../Makefile) in the project root to drop into a container shell.
This will mount all source and config files into a single location enabling
you to run the bot from source:
```
$ make shell
```

# Testing

The `peonbot` package defined in this project is the primary module of the
bot, and is close to 80% unit tested:
```
ok      peonBot/peonbot 0.007s  coverage: 78.0% of statements
```

In the small off-chance that someone wants to contribute back to this
project, please keep the test functions updated as source code changes
are made.

## Running Tests
* Run unit tests
> ```
> $ make unit-test
> ```

* Generate a command-line coverage report
> ```
> $ make coverage
> ```

* Generate a html coverage report
> ```
> $ make coverage-html # Spits out a coverage.html file in a coverage folder.
> ```

# Building

I've provided both a 64-bit linux, and a 64-bit Windows binary for easy
access to the bot. However, if you need to rebuild the binaries:
* Linux
> ```
> $ make bot
> ```

* Windows
> ```
> $ make bot-win
> ```

## Other OS's and Architectures
If for whatever reason the two provided binaries are unable to run on your
system, you can build the bot for different operating systems by modifying
the value of the environment variable `GOOS`, and for different system
architectures by modifying the value of the environment variable `GOARCH`.

* For other possible `GOOS` and `GOARCH` values, please see
https://golang.org/doc/install/source.