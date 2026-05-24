import { execFile } from "node:child_process";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);
const MAX_BUFFER = 10 * 1024 * 1024;
const RECORD_SEPARATOR = "\u001e";
const FIELD_SEPARATOR = "\u001f";

export type CommitSummary = {
  sha: string;
  shortSha: string;
  author: string;
  date: string;
  subject: string;
  body: string;
  files: string[];
};

export type RecentCommitsSnapshot = {
  repoRoot: string;
  branch: string;
  totalCommits: number;
  commits: CommitSummary[];
};

async function git(cwd: string, args: string[]) {
  const { stdout } = await execFileAsync("git", args, { cwd, maxBuffer: MAX_BUFFER });
  return stdout.trim();
}

function parseCommitMetadata(raw: string): Omit<CommitSummary, "files">[] {
  return raw
    .split(RECORD_SEPARATOR)
    .map((record) => record.trim())
    .filter(Boolean)
    .map((record) => {
      const [sha, author, date, subject, body = ""] = record.split(FIELD_SEPARATOR);
      return { sha, shortSha: sha.slice(0, 7), author, date, subject: subject.trim(), body: body.trim() };
    });
}

async function getChangedFiles(repoRoot: string, sha: string) {
  const raw = await git(repoRoot, ["show", "--pretty=format:", "--name-only", "--no-renames", sha]);
  return [...new Set(raw.split("\n").map((line) => line.trim()).filter(Boolean))];
}

export async function getRecentCommits(cwd: string, since: number): Promise<RecentCommitsSnapshot> {
  if (!Number.isInteger(since) || since <= 0) throw new Error("--since must be positive integer");

  const repoRoot = await git(cwd, ["rev-parse", "--show-toplevel"]);
  const branch = (await git(repoRoot, ["branch", "--show-current"])) || "detached-head";
  const metadata = await git(repoRoot, [
    "log",
    `-${since}`,
    "--date=short",
    `--pretty=format:%H%x1f%an%x1f%ad%x1f%s%x1f%b%x1e`,
  ]);

  const commits = await Promise.all(
    parseCommitMetadata(metadata).map(async (commit) => ({ ...commit, files: await getChangedFiles(repoRoot, commit.sha) })),
  );

  return { repoRoot, branch, totalCommits: commits.length, commits };
}
