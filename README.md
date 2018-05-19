# Flowcus

Flowcus is a data processing testing framework implementing various testing patterns such as Black Box Testing, or
Gray Box Testing.

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

* [Black Box Testing - equivalency](https://github.com/TommyStarK/flowcus/blob/v2/bbox-equivalency.md)
* [Gray Box Testing](https://github.com/TommyStarK/flowcus/blob/v2/gbox.md)

