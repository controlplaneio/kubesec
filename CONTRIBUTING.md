# Contributing to Kubesec

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

Kubesec is Apache 2.0 licensed and accepts contributions via GitHub pull requests.

The following is a set of guidelines for contributing to kubesec and it's related projects. We generally have stricter rules
as it's a security tool but don't let that discourage you from creating your PR, it can be incrementally fixed to fit the
rules. Also feel free to propose changes to this document in a pull request.

## Table Of Contents

- [Code of Conduct](#code-of-conduct)
- [I Don't Want To Read This Whole Thing I Just Have a Question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)
- [What Should I Know Before I Get Started?](#what-should-i-know-before-i-get-started)
  - [Kubesec and Related Projects](#kubesec-and-related-projects)
- [How Can I Contribute?](#how-can-i-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Pull Requests](#pull-requests)
- [Style Guides](#style-guides)
  - [Git Commit Messages](#git-commit-messages)
  - [GoLang Style Guide](#golang-style-guide)
  - [bash/bats Style Guide](#bashbats-style-guide)
  - [Documentation Style Guide](#documentation-style-guide)

---

## Code of Conduct

This project and everyone participating are governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you
are expected to uphold this code. Please report unacceptable behaviour to [andy@control-plane.io](mailto:andy@control-plane.io).

## I Don't Want To Read This Whole Thing I Just Have a Question!!!

We have an official message board with a detailed FAQ and where the community chimes in with helpful advice if you have questions.

We also have an issue template for questions [here](https://github.com/controlplaneio/kubesec/issues/new).

## What Should I Know Before I Get Started?

### Kubesec and Related Projects

- [controlplaneio/kubesec](https://github.com/controlplaneio/kubesec)
  - The main Kubesec repository! The main command-line tool for local scanning or running as a HTTP service. You should
    also use this repository for feedback related to the API and for large, overarching design proposals
- [controlplaneio/kubectl-kubesec](https://github.com/controlplaneio/kubectl-kubesec)
  - A `kubectl` plugin that can feed your deployments, pods, etc into Kubesec

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for Kubesec. Following these guidelines helps maintainers and the
community understand your report, reproduce the behaviour, and find related reports.

Before creating bug reports, please check [this list](#before-submitting-a-bug-report) as you might find out that you
don't need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report).
Fill out the issue template for bugs, the information it asks for helps us resolve issues faster.

> **Note:** If you find a **Closed** issue that seems like it is the same thing that you're experiencing, open a new issue
> and include a link to the original issue in the body of your new one.

#### Before Submitting a Bug Report

- **Determine [which repository the problem should be reported in](#kubesec-and-related-projects)**
- **Perform a [cursory search](https://github.com/search?q=+is:issue+user:controlplaneio)** to see if the problem has already
  been reported. If it has **and the issue is still open**, add a comment to the existing issue instead of opening a new
  one

#### How Do I Submit a (Good) Bug Report?

Bugs are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](#kubesec-and-related-projects)
your bug is related to, create an issue on that repository and provide the following information by filling in the issue
template [here](https://github.com/controlplaneio/kubesec/issues/new).

Explain the problem and include additional details to help maintainers reproduce the problem:

- **Use a clear and descriptive title** for the issue to identify the problem
- **Describe the exact steps which reproduce the problem** in as many details as possible. For example, start by explaining
- how you started `kubectl`, e.g. which command you used in the terminal, or how you started Kubesec otherwise
- **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable
  snippets, which you use in those examples. If you're providing snippets in the issue, use [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines)
- **Describe the behaviour you observed after following the steps** and point out what exactly is the problem with that behaviour
- **Explain which behaviour you expected to see instead and why.**

Provide more context by answering these questions:

- **Did the problem start happening recently** (e.g. after updating to a new version of Kubesec) or was this always a problem?
- If the problem started happening recently, **can you reproduce the problem in an older version of Kubesec?** What's the
  most recent version in which the problem doesn't happen? You can download older versions of Kubesec from
  [the releases page](https://github.com/controlplaneio/kubesec/releases)
- **Can you reliably reproduce the issue?** If not, provide details about how often the problem happens and under which conditions
  it normally happens
- If the problem is related to scanning files, **does the problem happen for all files and projects or only some?** Is there
  anything else special about the files you are using? Please include them in your report, censor any sensitive information
  but ensure the issue still exists with the censored file

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Kubesec, including completely new features and minor
improvements to existing functionality. Following these guidelines helps maintainers and the community understand your suggestion
and find related suggestions.

Before creating enhancement suggestions, please check [this list](#before-submitting-an-enhancement-suggestion) as you might
find out that you don't need to create one. When you are creating an enhancement suggestion, please
[include as many details as possible](#how-do-i-submit-a-good-enhancement-suggestion). Fill in the template feature request
template, including the steps that you imagine you would take if the feature you're requesting existed.

#### Before Submitting an Enhancement Suggestion

- **Check if there's already [project covering that enhancement](#kubesec-and-related-projects) if it's a larger enhancement**
- **Determine [which repository the enhancement should be suggested in](#kubesec-and-related-projects)**
- **Perform a [cursory search](https://github.com/search?q=+is:issue+user:controlplaneio)** to see if the enhancement has
  already been suggested. If it has, add a comment to the existing issue instead of opening a new one

#### How Do I Submit A (Good) Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined
[which repository](#kubesec-and-related-projects) your enhancement suggestion is related to, create an issue on that repository
and provide the following information:

- **Use a clear and descriptive title** for the issue to identify the suggestion
- **Provide a step-by-step description of the suggested enhancement** in as many details as possible
- **Provide specific examples to demonstrate the steps**. Include copy/pasteable snippets which you use in those examples,
  as [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines)
- **Describe the current behaviour** and **explain which behaviour you expected to see instead** and why
- **Explain why this enhancement would be useful** to most Kubesec users and isn't something that can or should be implemented
  as a separate community project
- **List some other tools where this enhancement exists.**
- **Specify which version of Kubesec you're using.** You can get the exact version by running `kubesec version` in your terminal
- **Specify the name and version of the OS you're using.**

### Your First Code Contribution

Unsure where to begin contributing to Kubesec? You can start by looking through these `Good First Issue` and `Help Wanted`
issues:

- [Good First Issue issues][good_first_issue] - issues which should only require a few lines of code, and a test or two
- [Help wanted issues][help_wanted] - issues which should be a bit more involved than `Good First Issue` issues

Both issue lists are sorted by total number of comments. While not perfect, number of comments is a reasonable proxy for
impact a given change will have.

#### Development

To build the project you can use `make build`. The resulting binary will be in `./dist`.

To test the project you can run `make test` for unit and command-line acceptance testing. For http testing also run `make test-remote`.

### Pull Requests

The process described here has several goals:

- Maintain Kubesec's quality
- Fix problems that are important to users
- Engage the community in working toward the best possible Kubesec
- Enable a sustainable system for Kubesec's maintainers to review contributions

Please follow these steps to have your contribution considered by the maintainers:

<!-- markdownlint-disable no-inline-html -->

1. Follow all instructions in the template
2. Follow the [style guides](#style-guides)
3. After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/)
   are passing
   <details>
    <summary>What if the status checks are failing?</summary>
    If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on
    the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for
    you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our
    status check suite.
   </details>

<!-- markdownlint-enable no-inline-html -->

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to
complete additional tests, or other changes before your pull request can be ultimately accepted.

## Style Guides

### Git Commit Messages

- It's strongly preferred you [GPG Verify][commit_signing] your commits if you can
- Follow [Conventional Commits](https://www.conventionalcommits.org)
- Use the present tense ("add feature" not "added feature")
- Use the imperative mood ("move cursor to..." not "moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### General Style Guide

Look at installing an `.editorconfig` plugin or configure your editor to match the `.editorconfig` file in the root of the
repository.

### GoLang Style Guide

All Go code is linted with [golangci-lint](https://golangci-lint.run/).

For formatting rely on `gofmt` to handle styling.

### Bash/Bats Style Guide

We follow the [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html).
All bash/bats code is linted with [shellcheck](https://www.shellcheck.net/).
In the future it will also be formatted with [shfmt](https://github.com/mvdan/sh).

### Documentation Style Guide

All markdown code is linted with [markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli).

[good_first_issue]:https://github.com/controlplaneio/kubesec/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22+sort%3Acomments-desc
[help_wanted]: https://github.com/controlplaneio/kubesec/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22help+wanted%22

[commit_signing]: https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/managing-commit-signature-verification
