# Cloud Foundry client library for Golang

This is a start of a gentle wrapper around the Golang libraries used in the Cloud Foundry CLI implementation.
It makes it much easier to consume the API directly from CLI plugins, rather than having to invoke
it through CliConnection.CliCommand[WithoutTerminalOutput] and parsing the output of that.

What the wrapper tries to do for you:
- Initialize everything that needs to be initialized to safely call the APIs
- Hide all of the different internal classes (repositories and helpers)
- Provide access to useful feature that are in the CLI to ease consumption of the API

Examples of the latter:
- Panic handling
- Token refreshing
- Waiting for state to change (making API behavior synchronous)

The API is only as complete as needed for its consumers. I have made no attempt to extend it beyond the
specific calls I needed for my plugins.
