# sreconf-temporal-demo

This project is a simple demo utilizing [Temporal](https://temporal.io) to run a workflow that tags AWS Subnets.

## Setup

In order to run this, besides needing Go, you will also need Temporal.

* Install Temporal (on a Mac) with `brew install temporal`

## Running

* Start Temporal with: `make start-temporal`
* Start the worker with: `make worker`

## Running a Workflow

The workflow provided by this code is called `TagAllSubnets`. It will use _whatever_ AWS credentials the SDK has access to,
and will iterate through all VPCs, find the subnets, then tag them with the provided `key` and `value`.

An example of running the workflow from the CLI would look like:

    temporal workflow execute --task-queue default --type "TagAllSubnets" --input '{"key": "owner", "value": "bob"}
