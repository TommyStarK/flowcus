# Flowcus

Flowcus is an easy to use data flow testing framework to implement Black Box, White Box or Gray Box Testing

## Context

As a devOps trainee in a company providing analytics services, my first assignment was to write a software to perform integration tests on the data collection process. The data goes through different stages of treatment, my goal was to test these steps as well as the entire process flow.
The real challenge was to correctly setup the test and being able to perform it, not really to write the test itself. So I have
decided to rewrite the 1st version of my software to make it generic and reusable.

## Contribution

Each Contribution is welcomed and encouraged. I do not claim to cover each use cases nor completely master the Golang. If you encounter a non sense or any trouble, you can open an issue and I will be happy to discuss about it :)

## Install

```bash
    go get -u github.com/TommyStarK/flowcus
```

## Test

```bash
    go test -v --cover ./...
```

## Usage

Implementation examples:

* [Black Box Testing - equivalency](https://github.com/TommyStarK/flowcus/blob/v2/docs/bbox-equivalency.md)
* [Gray Box Testing - exploratory](https://github.com/TommyStarK/flowcus/blob/v2/docs/gbox-exploratory.md)

