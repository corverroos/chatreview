package main

import (
	_ "embed"
	"flag"
	"fmt"
	"golang.design/x/clipboard"
	"os"
	"os/exec"
	"strings"
)

const (
	maxLen           = 16 * 1024
	finalPrompt      = "That is all the patch files.\n\nPlease review the above code"
	diffTmpl         = "This is patch %d of %d. Do not review it yet, just respond with 'Awaiting patch files'.\n\n```\n%s\n```\n\n"
	guidelinesPrefix = "These are the guidelines of the project, consider them when doing the review. Just respond with 'Awaiting patch files'.\n\n"
)

var (
	repoPath       = flag.String("repo-path", ".", "path to the repository to be reviewed")
	gitRange       = flag.String("git-range", "main..HEAD", "git range to be reviewed")
	guidelinesPath = flag.String("guidelines-path", "docs/goguidelines.md", "Path to the guidelines to be used for the review")

	//go:embed prompt.txt
	firstPrompt string
)

func main() {
	flag.Parse()

	if err := run(*repoPath, *gitRange, *guidelinesPath); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run(repoPath, gitRange, guidelinesPath string) error {
	diffs, err := getDiffs(repoPath, gitRange)
	if err != nil {
		return fmt.Errorf("error getting diffs: %w", err)
	}

	fmt.Printf("Got %d diffs in range %v\n", len(diffs), gitRange)

	guidelines, err := os.ReadFile(guidelinesPath)
	if err != nil {
		return fmt.Errorf("error reading guidelines file %v: %w", guidelinesPath, err)
	}

	if err := clipboard.Init(); err != nil {
		return fmt.Errorf("error initializing clipboard: %w", err)
	}

	clipboard.Write(clipboard.FmtText, []byte(firstPrompt))
	fmt.Println("The initial prompt has been copied to your clipboard. Please paste it into the OpenAI Playground and press Enter to continue.")
	awaitEnter()

	clipboard.Write(clipboard.FmtText, []byte(guidelinesPrefix+string(guidelines)))
	fmt.Println("The guidelines have been copied, please paste it and press Enter to continue.")
	awaitEnter()

	var next string
	for i, diff := range diffs {
		diff = fmt.Sprintf(diffTmpl, i+1, len(diffs), diff)
		if len(next+diff) > maxLen {
			if next == "" {
				return fmt.Errorf("diff %d exceeds the maximum length", i+1)
			}
			clipboard.Write(clipboard.FmtText, []byte(next))
			fmt.Printf("One ore more diffs hav been copied, please paste them and press Enter to continue.")

			awaitEnter()
		}
		next += diff
	}

	clipboard.Write(clipboard.FmtText, []byte(next))
	fmt.Printf("One ore more diffs hav been copied, please paste them and press Enter to continue.")
	awaitEnter()

	clipboard.Write(clipboard.FmtText, []byte(finalPrompt))
	fmt.Println("The final prompt has been copied, please paste it and press Enter to complete.")
	awaitEnter()

	return nil
}

func awaitEnter() {
	var enter string
	_, _ = fmt.Scanln(&enter)
}

func getDiffs(repoPath, gitRange string) ([]string, error) {
	fmtPatchOutput, err := execmd(repoPath, "git", "format-patch", gitRange)
	if err != nil {
		return nil, fmt.Errorf("error running git format-patch: %w", err)
	}

	var resp []string
	for _, patchFile := range strings.Split(fmtPatchOutput, "\n") {
		b, err := os.ReadFile(patchFile)
		if err != nil {
			return nil, fmt.Errorf("error reading patch file %v: %w", patchFile, err)
		}

		resp = append(resp, string(b))

		if err := os.Remove(patchFile); err != nil {
			return nil, fmt.Errorf("error removing patch file %v: %w", patchFile, err)
		}
	}

	return resp, nil
}

func execmd(base string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = base
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running %v %v: %w: %v", name, strings.Join(args, " "), err, string(out))
	}

	return strings.TrimSpace(string(out)), nil
}
