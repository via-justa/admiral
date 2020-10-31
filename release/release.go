package release

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
)

const user = "via-justa"
const repo = "admiral"
const releaseURL = "https://github.com/via-justa/admiral/releases/tag/"

// CheckForUpdates report if new release available to download
func CheckForUpdates(currentVersion string) string {
	ctx := context.Background()
	client := github.NewClient(nil)

	tag, _, err := client.Repositories.ListTags(ctx, user, repo, nil)
	if err != nil {
		fmt.Println(err)
	}

	releases, _, err := client.Repositories.ListReleases(ctx, user, repo, nil)
	if err != nil {
		fmt.Println(err)
	}

	if currentVersion != releases[0].GetName() {
		return fmt.Sprintf("Admiral version %v. New version %v Available\nYou "+
			"can download the latest version here: %v%v\n",
			currentVersion, releases[0].GetName(), releaseURL, tag[0].GetName())
	}

	return ""
}
