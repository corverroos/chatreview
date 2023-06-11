# chatreview

This is a go command line tool to generate prompts for an AI to code review a range of git commits (PR/branch).

This doesn't use the OpenAI API, since it doesn't support GPT4 at time of writing.

Instead, it copies the prompts to your clipboard allowing you to paste them into the OpenAI web UI which 
is the only way to use GPT4 at time of writing.

## Install
Install with go directly:
```shell
go install github.com/corverroos/chatreview
```

Or build from source:
```shell
git checkout https://github.com/corverroos/chatreview.git 
cd chatreview
go install .
```
Ensure `chatreview` is installed and in your path:
```shell
which chatreview
```

## Usage
Navigate to your git repo and checkout the branch to review.
```shell
cd charon
git checkout feature/charon-123
```

Then run `chatreview`:
```shell
chatreview
```
You will be instructed to paste multiple prompts into the web UI, pressing `ENTER` in the command line after each to continue.

## Flags
The following flags are supported:
- `--repo-path`: Path to the repo to review. Defaults to current working directory.
- `--git-range`: The git range to review. Defaults to `main..HEAD`.
- `--guidelines-path`: Path to the repos coding guidelines markdown file. Defaults to `docs/goguidelines.md`.


