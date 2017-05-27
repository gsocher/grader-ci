# grader-ci
[![Build Status](https://travis-ci.org/dpolansky/grader-ci.svg?branch=master
)](https://travis-ci.org/dpolansky/grader-ci)

A prototype CI server for grading programming assignments.

## Overview
To submit programming assignments at many universities, instructors rely on command-line applications which transfer files from a student's machine to a grading machine, where grading scripts are manually run to either run unit tests or diff program output. Current solutions for handing in assignments could be dramatically improved from both a student's perspective in terms of receiving feedback on whether their assignment builds on the grading machine and passes tests, and from an instructor's perspective in terms of automatically grading assignments and sending results to students.

When developing production-ready software, companies use continuous integration servers to build, test, and deploy code. CI servers are simple and easy to interact with: you make changes to your code, and the server ensures that the project builds in the target environment and passes tests. This project aims to explore how a continuous integration server could built to suit the task of handing in and grading programming assignments.

### Grading
To set up an assignment for grading, students simply create a Git repository for their assignment with a `.ci.yml` configuration file to specify build criteria, such as project dependencies and build/test scripts. When the student makes changes to their code and pushes their commits to a repository hosting service (such as Github), grader-ci receives a webhook containing the information necessary to clone and build the repository.

When testing the assignment against the instructor's test cases, the student creates a "binding" through the web interface from their source repository to the instructor's grading repository. A grading repository contains test files and/or source code that will be merged into the source repository when a build is triggered. This means that when a student makes a commit to the source repo, grader-ci will test it in combination with the grader's repository.

For example, a student's repository might look like:
- `fizzbuzz.go`
- `.ci.yml`

and the grading repository might look like:
- `fizzbuzz_test.go`
- `.ci.yml`

When the source repository is being built, `fizzbuzz.go` will be tested with the grader's `fizzbuzz_test.go`.

## Installation
Either clone the repository or use the go tool to fetch it:
```
$ go get github.com/dpolansky/grader-ci
```
`gvt` is used to manage dependencies. To install the dependencies, run:
```
$ gvt restore
```
The easiest way to setup the project is to use `vagarant` to provision a virtual machine with all of the project's dependencies (docker, rabbitmq, sqlite, etc).
```
$ vagrant up && vagrant ssh
```
`vagrant` will use the `Vagrantfile` to launch the VM, which uses the `setup.sh` script to provision the VM.

## Usage
The project relies on a single instance of `cmd/backend`, and potentially several instances of `cmd/worker`.

#### Worker
The worker process waits for messages from a local RabbitMQ instance containing a build task. To build and run a worker process:
```
$ cd cmd/worker/
$ go build
$ ./worker
```

#### Backend
The backend web server serves the project's frontend and API routes and relies on SQLite. To setup the database and run the backend:
```
$ cd cmd/backend/
$ ./refresh-db
$ go build
$ ./backend
````

## Dependencies
Dependency | Use
--- | ---
`docker` | Used to run builds in isolated environments with necessary dependencies.
`SQLite` | Used to store data for builds, repositories, and bindings.
`RabbitMQ` | Used to send messages between the backend server and worker instances.