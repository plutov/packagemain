# LightStep Community Contributing Guide

First, 🎉 **thanks for contributing!** 🎉

## Issues

You're encouraged to log issues for any questions or problems you might have. When in doubt, log an issue. The exception to this rule is [security disclosures](#reporting-security-issues).

Generally speaking, the more context you can provide, the better. Please add information such as what **version** you're using, **stack traces** and/or **logs** (to the extent that you're able to share them), and whatever else you think may be relevant. Project maintainers may ask for additional clarification, logs, and other pertinent metadata before we can address your issue.

For bug submission, we especially appreciate **details on how to reproduce the bug** to the extent you're able to provide them, e.g., an isolated repo or [gist](https://gist.github.com).

### Reporting Security Issues

If you find a security issue, please **do not** file a public issue for it. Instead, send your report to us privately at [security@lightstep.com](mailto:security@lightstep.com).

## Contributions

All contributions big and small are welcome, from typo corrections to bug fixes to suggested improvements!

Any changes to project resources in this repository must be made through a pull request. This includes, but is not limited to, changes affecting:

- Documentation
- Source code
- Binaries
- Sample projects or other examples

No pull request can be merged without at least one review from a maintainer.

By default, contributions are accepted once no committers object to the PR. Specific contributors may be suggested or required to review a pull request based on repository settings.

In the event of objections or disagreement, everyone involved should seek to arrive at a consensus around the expressed objections. These can take the form of addressing concerns through changes, compromising around the change, or withdrawing it entirely.

## Development

### Testing

To run the tests:

```
make test
```

### Style Guide

Please make sure that your code is formatted using `gofmt -s`.

## Submitting a Pull Request

1. [Fork the repository.](https://help.github.com/en/github/getting-started-with-github/fork-a-repo)
1. Create a new branch.
1. Add tests for your change.
1. [Run the tests](#testing) to make sure that they don't already pass. If they do (and you're not backfilling test coverage), please modify them.
1. Implement the change such that your new tests pass.
1. Make sure that your code conforms to the [style guide](#style-guide).
1. [Commit and push your changes.](https://guides.github.com/introduction/flow/)
1. [Submit your pull request.](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests)
1. Adjust your pull request based on feedback.
1. Get it merged! 🎉

We're happy to help with any questions you may have on the git or GitHub side, e.g., how to push a branch to your fork. Just create an issue and we'll try to help answer them :)
