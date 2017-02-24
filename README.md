# ci

### Package Synopsis
Package | Description
--- | ---
`cmd` | Programs that build to executables, such as the web server or worker process 
`worker` | Worker process that listens for AMQP messages for builds, and then runs them with Docker
`server` | Web server which exposes a REST API and receives to Git webhooks
