# Developer's Guide to Contributing to the Muxer Package

Thank you for your interest in contributing to the Muxer package! This
guide will help you get started with programming against the Muxer package
and provide you with the necessary information to contribute to the
project.

## Getting Started

To get started, make sure you have Go installed on your system. You can
download and install Go from the official Go website: https://golang.org/

Once Go is installed, you can set up your development environment and
start programming against the Muxer package.

## Development Environment Setup

1. Clone the Muxer repository:
```
git clone https://github.com/shellfu/muxer.git
```

2. Change into the cloned directory:
```
cd muxer
```

3. Install project dependencies:
```
go mod download
```

4. Verify that all tests pass:
```
go test -cover ./...
```

The output should indicate that the coverage for both the `muxer` and `middleware` packages.

Now you're ready to start working with the Muxer package!

## Contributing Guidelines

To contribute to the Muxer package, please follow these guidelines:

1. Fork the Muxer repository on GitHub.
2. Create a new branch for your changes:
```
git checkout -b my-feature
```
3. Make your changes to the codebase. Ensure that your changes follow the existing code style and conventions.
4. Write tests for your changes to maintain code coverage. Run the tests to ensure they pass:
```
go test -cover ./...
```
5. Commit your changes and push them to your forked repository:
```
git commit -m "Add my feature"
git push origin my-feature
```
6. Open a pull request (PR) against the main Muxer repository on GitHub.
7. Provide a clear and descriptive explanation of your changes in the PR description.
8. Wait for feedback and iterate on your changes if necessary.

## Code Coverage

The Muxer package aims to maintain a high code coverage. When making changes
or adding new features, please ensure that your changes are covered by
tests. Run the tests with coverage to verify the coverage percentage:
```
go test -cover ./...
```
Ensure that the coverage is adequate for your new feature / bug fix.

## Documentation

Documentation is an essential part of any project. If you're adding new
features or modifying existing ones, please update the relevant
documentation to reflect the changes. This includes updating code
comments, README files, and any other relevant documentation files.

## Style Guide

The Muxer package follows the standard Go code style guidelines. Please
ensure that your code follows these guidelines to maintain consistency.
You can refer to the following for more details:

- https://github.com/golang/go/wiki/CodeReviewComments
- https://go.dev/doc/effective_go
- http://www.catb.org/~esr/writings/taoup/html/ch01s06.html


## Contact and Support

If you have any questions or need support while working on the Muxer
package, feel free to reach out to the project maintainers through the
GitHub repository. They will be happy to assist you.

Thank you for your contributions to the Muxer package! Together, we can
make it even better.

Happy coding!
