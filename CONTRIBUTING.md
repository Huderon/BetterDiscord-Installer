# Contributing to BetterDiscord's Installer

Thanks for taking the time to contribute!

The following is a set of guidelines for contributing to BetterDiscord's Installer. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request. These guidelines have been adapted from [Atom](https://github.com/atom/atom/blob/master/CONTRIBUTING.md).

#### Table Of Contents

[Code of Conduct](#code-of-conduct)

[What should I know before I get started?](#what-should-i-know-before-i-get-started)
  * [Development Setup](#development-setup)

[How Can I Contribute?](#how-can-i-contribute)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Pull Requests](#pull-requests)

[AI-assisted contributions](#ai-assisted-contributions)

[Styleguides](#styleguides)
  * [Git Commit Messages](#git-commit-messages)
  * [Frontend (Svelte + TypeScript) Styleguide](#frontend-svelte--typescript-styleguide)
  * [Go Styleguide](#go-styleguide)

[Additional Notes](#additional-notes)
  * [Issue Labels](#issue-labels)

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct from the Contributor Covenant](https://www.contributor-covenant.org/version/1/4/code-of-conduct.html). By participating, you are expected to uphold this code. Please report unacceptable behavior.

## What should I know before I get started?

This installer is built with [Wails](https://wails.io/) (Go backend + Svelte frontend). The Go runtime embeds the compiled frontend assets from `frontend/build` (see `//go:embed` in `main.go`).

The repository is organized into:

```
.
├── api                     // Wails bindings and backend runtime helpers.
├── build                   // Build assets and platform-specific packaging files.
├── frontend                // Svelte UI bundled by Vite (Bun-powered).
├── types                   // Shared Go types used by bindings.
├── utils                   // Backend utilities.
├── main.go                 // Wails application entry point.
├── wails.json              // Wails configuration + frontend hooks.
```

### Development Setup

Prerequisites:

* [Go](https://go.dev/) (matches `go.mod`)
* [Bun](https://bun.sh/)
* [Wails CLI](https://wails.io/docs/gettingstarted/installation)

Linux prerequisites:

* Wails requires GTK/WebKit and build tooling. Follow the official Linux dependencies guide:
  https://wails.io/docs/gettingstarted/installation#linux
* Fedora note: you may need to set `PKG_CONFIG_PATH` to include system pkgconfig directories:
  ```sh
  export PKG_CONFIG_PATH="/usr/lib64/pkgconfig:/usr/share/pkgconfig"
  ```

Common commands:

* `wails dev` - run the app locally with backend bindings.
* `wails build` - build distributable binaries.
* `bun run lint` (from `frontend/`) - lint frontend changes.
* `bun run check` (from `frontend/`) - typecheck and static checks.
* `go test ./...` - run backend tests.
* `gofmt -w .` - format Go files you touch.

Build and development details:

* Frontend assets are built by Vite into `frontend/build` and embedded by Go via `//go:embed` in `main.go`.
* Wails uses the hooks in `wails.json` to install and build frontend assets during `wails dev` and `wails build`.

## How Can I Contribute?

### Reporting Bugs

Please search for existing issues first. If you find a match, add any extra context you have (logs, OS details, screenshots).

#### Before Submitting A Bug Report

* Reproduce on the latest `main` branch or a recent release.
* Collect logs from the installer UI (copy any error output shown).
* Note your OS and architecture (Windows/macOS/Linux, x64/arm64).

#### How Do I Submit A (Good) Bug Report?

Include clear steps to reproduce, expected vs. actual behavior, and any logs shown in the UI.

### Suggesting Enhancements

Open an issue describing the problem you are trying to solve, the proposed solution, and any alternatives considered.

#### Before Submitting An Enhancement Suggestion

Confirm the change aligns with the existing installer flow and doesn't require a new runtime permission without a clear UX plan.

#### How Do I Submit A (Good) Enhancement Suggestion?

Provide context, screenshots/mockups (if UI changes), and a minimal spec of expected behavior.

### Your First Code Contribution

Unsure where to begin contributing? You can start by looking through `help-wanted` issues or any issues labelled `can't reproduce`.

### Pull Requests

Please follow these steps to have your contribution considered by the maintainers:

1. Use a pull request template, if one exists.
2. Follow the [styleguides](#styleguides)
3. After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/) are passing <details><summary>What if the status checks are failing?</summary>If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our status check suite.</details>

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to complete additional design work, tests, or other changes before your pull request can be ultimately accepted.

## AI-assisted contributions

AI tools are fine to use. Submitting output with little or no personal review is not.

- **Review the diff yourself.** If you can't explain why each changed line is correct, the PR isn't ready.
- **Run the required checks locally.** Don't rely on CI to catch problems you could catch before pushing.
- **Disclose AI involvement in your PR description.** If an AI wrote a meaningful portion of the code, say so and briefly describe what you reviewed. This isn't meant as a penalty, it's useful context for the reviewer.
- **If you used an autonomous agent with minimal personal review of the output, say that explicitly.** PRs that appear to be unreviewed agent output and don't disclose this will be closed without detailed feedback.

## Styleguides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* When only changing documentation, include `[ci skip]` in the commit title

### Frontend (Svelte + TypeScript) Styleguide

All frontend code must adhere to the [ESLint rules](https://github.com/BetterDiscord/Installer/blob/main/frontend/eslint.config.js) of the repo.

Some other style related points not covered by ESLint:

* Use verbose variable names
* Inline `export`s with expressions whenever possible
  ```js
  // Use this:
  export default class ClassName {
  
  }

  // Instead of:
  class ClassName {

  }
  export default ClassName
  ```
* Place class properties in the following order:
    * Class methods and properties (methods starting with `static`)
    * Instance methods and properties
* Place imports in the following order:
    * Built in modules (such as `path`)
    * Third-party dependencies
    * Local modules (relative paths or `$lib` aliases)
* Prefer to import whole modules instead of singular functions when working in Node/Go tooling
```js
const fs = require("fs"); // Use this
const {readFile, writeFile} = require("fs"); // Avoid this

import Utilities from "./utilities"; // Use this
import {deepclone, isEmpty} from "./utilties"; // Avoid this
```

### Go Styleguide

* Run `gofmt` on any Go files you touch.
* Keep standard library imports first, then a blank line, then external deps, then local packages.
* Prefer early returns for error handling and emit Wails events (`log`, `success`, `failure`, `reset`) instead of silently swallowing failures.

## Additional Notes

### Releases

Releases are tag-driven. Create a tag like `vX.Y.Z` and push it to trigger the release workflow. The workflow uses the tag to:

* Set the runtime version via `-ldflags="-X main.version=vX.Y.Z"`.
* Override `wails.json` `info.productVersion` during CI builds so Windows metadata matches the tag.

If you build locally for a release, pass the same ldflags and update `wails.json` to keep metadata aligned:

```sh
wails build -ldflags="-X main.version=vX.Y.Z"
```

### Issue Labels

Use labels to clarify impact (`bug`, `enhancement`, `maintenance`) and scope (`frontend`, `backend`, `release`).

#### Type of Issue and Issue State

If you can, include an expected fix timeline or milestone to help triage.
