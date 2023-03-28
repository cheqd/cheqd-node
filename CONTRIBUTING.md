# Contributing to cheqd

We would love for you to contribute to cheqd and help make it even better than it is today!
As a contributor, here are the guidelines we would like you to follow.

## 🧑‍⚖️ Code of Conduct

Help us keep cheqd open and inclusive.
Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md)

## ❓ Got a question or problem?

You can get help for any questions or problems that you have on through the following channels:

- Check if your question/problem is already covered under our documentation sites:
  - [Node / cheqd network](https://docs.cheqd.io/node)
  - [Identity features / SDKs](https://docs.cheqd.io/identity)
  - [Governance Framework](https://gov.cheqd.io)
  - [Product Suite & Updates](https://product.cheqd.io)
  - [Learn about cheqd](https://learn.cheqd.io) (basics for general audience)
- Raise a bug report or feature request using the "Issue" tab on Github
- Ask the question on our [Community Slack](http://cheqd.link/join-cheqd-slack) or [Discord](http://cheqd.link/discord-github)

## 🪲 Found a bug?

If you find a bug in the source code, you can help us by [submitting an issue](#submit-issue) to our [GitHub Repository][github].
Even better, you can [submit a Pull Request](#submit-pr) with a fix.

## ❣️ Missing a feature?

You can *request* a new feature by [submitting an issue](#submit-issue) to our GitHub Repository.
If you would like to *implement* a new feature, please consider the size of the change in order to determine the right steps to proceed:

* For a **Major Feature**, first open an issue and outline your proposal so that it can be discussed. This process allows us to better coordinate our efforts, prevent duplication of work, and help you to craft the change so that it is successfully accepted into the project.
* **Small Features** can be crafted and directly [submitted as a Pull Request](#submit-pr).

## Submission Guidelines

### Submitting an Issue

Before you submit an issue, please search the issue tracker. An issue for your problem might already exist and the discussion might inform you of workarounds readily available.

We want to fix all the issues as soon as possible, but before fixing a bug, we need to reproduce and confirm it.
In order to reproduce bugs, we require that you provide a minimal reproduction.
Having a minimal reproducible scenario gives us a wealth of important information without going back and forth to you with additional questions.

A minimal reproduction allows us to quickly confirm a bug (or point out a coding problem) as well as confirm that we are fixing the right problem.

We require a minimal reproduction to save maintainers' time and ultimately be able to fix more bugs.
Often, developers find coding problems themselves while preparing a minimal reproduction.
We understand that sometimes it might be hard to extract essential bits of code from a larger codebase but we really need to isolate the problem before we can fix it.

Unfortunately, we are not able to investigate / fix bugs without a minimal reproduction, so if we don't hear back from you, we are going to close an issue that doesn't have enough info to be reproduced.

You can file new issues by selecting from our [new issue templates](https://github.com/angular/angular/issues/new/choose) and filling out the issue template.

### Submitting a Pull Request (PR)

Before you submit your Pull Request (PR) consider the following guidelines:

1. Search [GitHub](https://github.com/angular/angular/pulls) for an open or closed PR that relates to your submission. You don't want to duplicate existing efforts.
2. Be sure that an issue describes the problem you're fixing, or documents the design for the feature you'd like to add. Discussing the design upfront helps to ensure that we're ready to accept your work.
3. [Fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo) the [cheqd-node](https://github.com/cheqd/cheqd-node) repository.
4. In your forked repository, make your changes in a new git branch:

     ```bash
     git checkout -b my-fix-branch main
     ```

5. Create your patch, **including appropriate test cases**.
6. Check that all workflow actions for linting / build / test pass.
7. Commit your changes using a descriptive commit message that follows [Conventional Commits convention](https://www.conventionalcommits.org/en/v1.0.0/)
8. Push your branch to GitHub:

    ```shell
    git push origin my-fix-branch
    ```

9. In GitHub, send a pull request to `cheqd-node:main`.
