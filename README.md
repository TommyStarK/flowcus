# Flowcus

Flowcus is an easy to use data flow testing framework to implement Black Box Testing.

## Context

As a devOps trainee in a company providing analytics services, my first assignment was to write a software to perform integration tests on the data collection process. The data goes through different stages of treatment, my goal was to test these steps as well as the entire process flow.
The real challenge was to correctly setup the test and being able to perform it, not really to write the test itself. So I have
decided to rewrite the 1st version of my software to make it generic and reusable.

## Contribution

Each Contribution is welcomed and encouraged. I do not claim to cover each use cases nor completely master the Golang. If you encounter a non sense or any trouble, you can open an issue and I will be happy to discuss about it :)

## Install

```bash
    go get github.com/TommyStarK/Flowcus
```

## Test

```bash
    go test -v --cover ./...
```

## Usage

Implementation examples:

* [Black Box Testing - linear](https://github.com/TommyStarK/Flowcus/blob/master/docs/linear.md)
* [Black Box Testing - exploratory](https://github.com/TommyStarK/Flowcus/blob/master/docs/exploratory.md)
* [Black Box Testing - nonlinear](https://github.com/TommyStarK/Flowcus/blob/master/docs/nonlinear.md)
