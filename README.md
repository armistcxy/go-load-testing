# `goltest` - Load Testing CLI

## Introduction
`goltest` is a command-line tool written in Go for load testing with customizable payloads. It allows you to execute HTTP request load tests based on configurations and visualize the results in a graph.

It's builded on top of 2 main packages [`vegeta`](https://github.com/tsenart/vegeta)  and [`faker`](https://github.com/go-faker/faker/tree/main)

## Installation
```bash
go install github.com/armistcxy/go-load-testing/command/goltest@latest
```

## Core function
1. Create random payload for request with body (only support for random field type currently)
2. Display plot of test results

## Usage

### `attack` command
Starts the load testing based on the configuration file
```bash
goltest attack -p <config-path> -f <frequency> --per <interval> -d <duration>
```
Use `goltest help attack` to see information about flags. 

Example 
```bash
goltest attack -p fig.json -f 100 --per 1 -d 60
```

Structure of the config file (*json format*)
```json
{
  "URL": "http://localhost:8080/users",
  "Method": "POST",
  "Header": {},
  "Fields": {
    "Field1": "Type1",
    "Field2": "Type2",
    "Field3": "Type3"
  }
}
```
Support type can be found [here](https://github.com/armistcxy/go-load-testing/blob/main/internal/attacker/support_type.md)

### `plot` command
```bash
goltest plot [-p <result-path>]
```

The tool supports Linux, Windows, and macOS platforms for opening the result plot in the default web browser.
