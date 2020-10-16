package release

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v32/github"
)

const user = "via-justa"
const repo = "admiral"

// CheckForUpdates report if new release available to download
func CheckForUpdates(currentVersion string) {
	ctx := context.Background()
	client := github.NewClient(nil)

	releases, _, err := client.Repositories.ListReleases(ctx, user, repo, nil)
	if err != nil {
		fmt.Println(err)
	}

	if currentVersion != releases[0].GetName() {
		log.Printf("\nNew version available! Running: %v Available: %v\n\n", currentVersion, releases[0].GetName())
	}
}
