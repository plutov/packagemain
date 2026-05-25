import "dotenv/config";
import { Agent, createTool, type AgentToolContext } from "@cline/sdk";
import { getRecentCommits } from "./git.js";
import { buildSystemPrompt, buildUserPrompt } from "./prompt.js";

const TOOL_DESCRIPTION =
  "Read recent git commits from the current repository, including commit subjects, bodies, authors, dates, and changed file paths.";

function parseSince(argv: string[]) {
  const value = Number(argv[argv.indexOf("--since") + 1]);

  if (!Number.isInteger(value) || value <= 0) {
    throw new Error("Usage: npm run draft-release -- --since 20");
  }

  return value;
}

function createRecentCommitsTool() {
  return createTool<{ since: number }, Awaited<ReturnType<typeof getRecentCommits>>>({
    name: "get_recent_commits",
    description: TOOL_DESCRIPTION,
    inputSchema: {
      type: "object",
      properties: {
        since: { type: "number", description: "How many recent commits to inspect." },
      },
      required: ["since"],
      additionalProperties: false,
    },
    execute(input: { since: number }, context: AgentToolContext) {
      context.emitUpdate?.(`Reading last ${input.since} commits from git`);
      return getRecentCommits(process.cwd(), input.since);
    },
  });
}

async function main() {
  const apiKey = process.env.OPENAI_API_KEY;

  if (!apiKey) throw new Error("Missing OPENAI_API_KEY environment variable");

  // Create an agent with the OpenAI provider, the system prompt, and the tool for reading recent commits
  const agent = new Agent({
    providerId: "openai-native",
    modelId: process.env.OPENAI_MODEL ?? "gpt-4.1-mini",
    apiKey,
    systemPrompt: buildSystemPrompt(),
    tools: [createRecentCommitsTool()],
    maxIterations: 6,
  });

  const result = await agent.run(buildUserPrompt(parseSince(process.argv.slice(2))));

  if (result.status !== "completed") {
    throw result.error ?? new Error(`Agent ended with status: ${result.status}`);
  }

  process.stdout.write(`${result.outputText.trim()}\n`);
}

main().catch((error) => {
  process.stderr.write(`${error instanceof Error ? error.message : String(error)}\n`);
  process.exit(1);
});
