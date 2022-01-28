# Contribution Process

The following is a set of guidelines for contributing to Stencil. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request. Here are some important resources:

* [Stencil quick start guide](https://odpf.gitbook.io/stencil/quick_start) to get an idea about all components of stencil
* Github [issues](https://github.com/odpf/stencil/issues) track the ongoing and reported issues.

## How can I contribute?

We use RFCS and GitHub issues to communicate ideas.

* You can report a bug or suggest a feature enhancement or can just ask questions. Reach out on Github discussions for this purpose.
* You are also welcome to add new features, improve monitoring,logging and code quality.
* You can help with documenting new features or improve existing documentation.
* You can also review and accept other contributions if you are a maintainer.

Please submit a PR to the master branch of the Stencil repository once you are ready to submit your contribution. Code submission to Stencil, including submission from project maintainers, require review and approval from maintainers or code owners. PRs that are submitted by the general public need to pass the build. Once build is passed community members will help to review the pull request.

## Becoming a maintainer

We are always interested in adding new maintainers. What we look for is a series of contributions, good taste, and an ongoing interest in the project.

* maintainers will have write access to the Stencil repositories.
* There is no strict protocol for becoming a maintainer or PMC member. Candidates for new maintainers are typically people that are active contributors and community members.
* Candidates for new maintainers can also be suggested by current maintainers or PMC members.
* If you would like to become a maintainer, you should start contributing to Stencil in any of the ways mentioned. You might also want to talk to other maintainers and ask for their advice and guidance.

## Guidelines

Please follow these practices for you change to get merged fast and smoothly:

* Contributions can only be accepted if they contain appropriate testing \(Unit and Integration Tests.\)
* If you are introducing a completely new feature or making any major changes in an existing one, we recommend to start with an RFC and get consensus on the basic design first.
* Make sure your local build is running with all the tests and [golint](https://github.com/golang/lint) checks passing.
* Docs live in the code repo under [`docs`](https://github.com/odpf/raccoon/docs/README.md) so that changes to that can be done in the same PR as changes to the code.
* We follow [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages.
