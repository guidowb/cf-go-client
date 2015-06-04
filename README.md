# Cloud Foundry client library for Golang

This is a gentle wrapper around the Golang libraries used in the Cloud Foundry CLI implementation.
It makes it much easier to consume the API directly from CLI plugins, rather than having to invoke
it through CliConnection.CliCommand[WithoutTerminalOutput] and parsing the output of that.

The benefits of using the wrapper over the CLI libraries directly are:
- Everything is properly initialized for you to start invoking APIs. No trial and error to figure out what is and is not required for any particular API
- No need to create all the different repository and helper classes to invoke different parts of the API. Everything is available from the one and only client package

