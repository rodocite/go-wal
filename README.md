# Go Write-Ahead Log (WAL)

A simple Write-Ahead Log (WAL) implementation in Go. The WAL is used to ensure data durability and consistency in case of a crash or system failure. It writes all the changes to a log before they are applied to the main data store. In case of a failure, the log can be used to replay the changes and recover the data store to a consistent state.

## Features

- Simple key-value storage
- Write operations are logged before being applied
- Log replay functionality for recovery
