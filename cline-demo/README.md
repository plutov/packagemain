# Building an Agent with the Cline SDK

Cline is an open-source AI coding agent focused on real software work. Most developers first encounter Cline as an assistant in the editor or terminal, but Cline is broader than a single interface.

It has both a **CLI** and an **SDK**:

- the **CLI** is for running agent workflows directly from the terminal
- the **SDK** is for embedding the same agent runtime inside your own scripts, products, CI jobs, and internal tools

That distinction matters. In the [Cline SDK launch post](https://cline.bot/blog/introducing-cline-sdk-the-upgraded-agent-runtime), the team explains that the original product started inside the VS Code extension, and over time the runtime became harder to separate from the IDE around it. The SDK is the answer to that problem: the agent runtime is treated as a shared service rather than an implementation detail hidden inside one app.

That is also why the SDK is more interesting than a thin API wrapper. According to the launch post and SDK docs, the same runtime now powers Cline across the CLI and IDE surfaces, while staying open for other teams to embed in their own products. The low-level agent loop stays reusable, and the stateful runtime around it becomes more durable, portable, and product-agnostic.

The SDK documentation describes Cline SDK as an open-source **TypeScript** framework for building agentic applications. The launch post adds a few more important ideas:

- Cline 2.0 is built as a layered TypeScript stack
- teams can start with a small surface area and add more runtime pieces later
- the runtime is extensible through tools, plugins, MCP servers, skills, and hooks
- provider choice is not locked to one model vendor

That makes TypeScript the natural choice for a first project.

In this article, we will build a small release notes generator that uses the Cline SDK to inspect recent git history and turn it into readable markdown.

## Why the Cline SDK is a good fit here

This release-notes project is small, but it matches the SDK well because it uses the part of Cline that matters most: the runtime for tool-using agents.

We are not trying to rebuild the whole Cline product. We are only borrowing the runtime shape:

- define a focused tool
- give that tool to an agent
- let the agent inspect real project state
- turn the result into a useful artifact

That pattern lines up closely with how the Cline team positions the SDK in the launch post: something you can embed in scripts, internal tools, CI workflows, and other product surfaces, not just in an IDE.

## What we are building

The project lives in `cline-demo/` and exposes one command:

```bash
npm run draft-release -- --since 20
```

It does one job well:

1. start a Cline agent
2. give the agent one custom tool, `get_recent_commits`
3. let the tool read recent git history from the current repository
4. have the agent turn that data into release notes
5. print the result to stdout

The interesting part is not the CLI itself. The interesting part is the architecture: our application provides a narrow, useful capability, and the Cline runtime decides how to use it.

That is exactly the kind of problem the SDK is meant for. As the launch post puts it, the runtime is no longer supposed to live only inside one UI surface. It is something you can pull into your own stack.

## Project structure

```text
cline-demo/
├── .env.example
├── .gitignore
├── package.json
├── README.md
├── tsconfig.json
└── src/
    ├── git.ts
    ├── index.ts
    └── prompt.ts
```

## Setup

The Cline SDK requires **Node.js 22+**.

Install dependencies:

```bash
cd cline-demo
npm install
```

Create your environment file:

```bash
cp .env.example .env
```

Then fill in your OpenAI key:

```env
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-4.1-mini
```

This project uses Cline's `openai-native` provider.

Run it from inside any git repository:

```bash
npm run draft-release -- --since 20
```

To save the output:

```bash
npm run draft-release -- --since 20 > release-notes.md
```

## How the code works

### `src/index.ts`

This is the entry point. It does three things:

- parses the `--since` argument
- creates the custom `get_recent_commits` tool
- runs a Cline `Agent` and prints the final result

The core shape is intentionally small:

```ts
const agent = new Agent({
  providerId: "openai-native",
  modelId: process.env.OPENAI_MODEL ?? "gpt-4.1-mini",
  apiKey,
  systemPrompt: buildSystemPrompt(),
  tools: [createRecentCommitsTool()],
  maxIterations: 6,
})

const result = await agent.run(buildUserPrompt(parseSince(process.argv.slice(2))))
process.stdout.write(`${result.outputText.trim()}\n`)
```

This is the mental model to remember: **your app defines tools, and Cline supplies the agent runtime**.

In other words, we are not reimplementing agent orchestration ourselves. We are reusing the same idea Cline uses internally: a runtime that can reason, call tools, and produce a final artifact.

### `src/git.ts`

This file keeps the repository access logic out of the main program.

It uses:

- `git log` to collect recent commits
- `git show --name-only` to collect changed file paths per commit

Each commit is returned as structured data:

- sha
- shortSha
- author
- date
- subject
- body
- files

That structure matters. It gives the model enough context to infer whether a change is a feature, a fix, a maintenance task, or something that may require an upgrade note.

### `src/prompt.ts`

This file contains the prompt contract.

The system prompt tells the agent to:

- call `get_recent_commits` before answering
- use only tool evidence
- return markdown only
- organize the answer into:
  - Overview
  - Features
  - Fixes
  - Docs / Chore
  - Upgrade Notes

Keeping the prompt separate makes the project easier to explain and modify. The runtime code stays small, while the output rules live in one place.

## The custom tool

The single most important part of the project is the tool definition:

```ts
function createRecentCommitsTool() {
  return createTool({
    name: "get_recent_commits",
    description:
      "Read recent git commits from the current repository, including commit subjects, bodies, authors, dates, and changed file paths.",
    inputSchema: {
      type: "object",
      properties: {
        since: { type: "number", description: "How many recent commits to inspect." },
      },
      required: ["since"],
      additionalProperties: false,
    },
    execute(input, context) {
      context.emitUpdate?.(`Reading last ${input.since} commits from git`)
      return getRecentCommits(process.cwd(), input.since)
    },
  })
}
```

Without this tool, the agent would only be rephrasing whatever text we pasted into the prompt. With the tool, it can actively inspect the repository through a controlled interface.

That is where the SDK becomes interesting: it is not just a text wrapper around a model. It is a runtime for tool-using agents.

And if this project needed to grow later, the SDK already has room for that. The Cline team highlights plugins, custom tools, MCP integration, skills, and multi-agent capabilities as extension points. We are deliberately not using all of that here, but it is useful to know that the simple version and the more advanced version can live on the same foundation.

## Why release notes are a good SDK use case

Release notes sit in a sweet spot for agent automation:

- the input is messy but structured enough to inspect
- the output has a clear shape
- the task is useful in real projects
- the problem is narrow enough to understand quickly

In other words, this is not a toy chatbot, but it is also not an overbuilt autonomous system. It is a believable piece of SDLC automation.

## Example output

Here is the kind of markdown the tool produces:

```md
# Release Notes

## Overview
This release focused on improving authentication flows, tightening API validation, and cleaning up project documentation.

## Features
- Added a token refresh path for expired sessions.
- Introduced a reusable API client helper for authenticated requests.

## Fixes
- Fixed inconsistent validation errors in the user settings endpoint.
- Resolved a bug where logout did not fully clear local session state.

## Docs / Chore
- Updated onboarding docs for local development.
- Refactored auth-related file organization for easier maintenance.

## Upgrade Notes
- If you rely on the old auth client helper, update imports to the new shared client module.
```

## Final takeaway

The CLI version of Cline is about **using** an agent from the terminal. The SDK version is about **embedding** that agent into your own software.

That is the main idea behind the Cline SDK launch as well: pull the runtime out of a single product surface, make it reusable, and let other developers build on top of it.

This project shows that idea in a compact form:

- define one useful tool
- hand it to a Cline agent
- let the agent inspect real project data
- turn the result into a polished artifact

Once this pattern clicks, the same structure can be reused for PR summaries, changelog drafting, test-plan generation, issue triage, and other software delivery workflows.
