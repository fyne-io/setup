# Fyne Setup

## Synopsis

The repository hosts a graphical tool `setup` which aims
to help developers get their environment set up.

You can run this graphical interface to see if your computer is ready to go!

You should see something like the following screenshot:

![Setup screenshot](img/screenshot.png)

If your setup is not ready then it will give hints about how
to complete the configuration.

## Binaries

You may obtain binaries from `setup` tool from this [download page](https://geoffrey-artefacts.fynelabs.com/github/andydotxyz/fyne-io/setup/latest/).

## Installation

If you prefer to build from source, follow the instructions below.

### Prerequisites

Please make sure you have go, gcc and graphics library header files installed on your system.
You may consult the section [Prerequisites](https://docs.fyne.io/started/) in the documentation.

### Building and running

Use the commands below to build and run the tool from sources:

```
go install fyne.io/setup@latest
$(go env GOPATH)/bin/setup
```
