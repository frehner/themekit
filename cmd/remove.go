package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"

	"github.com/Shopify/themekit/kit"
)

var removeCmd = &cobra.Command{
	Use:   "remove <filenames>",
	Short: "Remove theme file(s) from shopify",
	Long: `Remove will delete all specified files from shopify servers.

For more documentation please see http://shopify.github.io/themekit/commands/#remove
	`,
	RunE: forEachClient(remove),
}

func remove(client kit.ThemeClient, filenames []string, wg *sync.WaitGroup) {
	defer wg.Done()

	if client.Config.ReadOnly {
		kit.LogErrorf("[%s]environment is reaonly", kit.GreenText(client.Config.Environment))
		return
	}

	for _, filename := range filenames {
		wg.Add(1)
		go performRemove(client, kit.Asset{Key: filename}, wg)
	}
}

func performRemove(client kit.ThemeClient, asset kit.Asset, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := client.DeleteAsset(asset)
	if err != nil {
		kit.LogErrorf("[%s]%s", kit.GreenText(client.Config.Environment), err)
	} else {
		kit.Printf(
			"[%s] Successfully removed file %s from %s",
			kit.GreenText(client.Config.Environment),
			kit.BlueText(asset.Key),
			kit.YellowText(resp.Host),
		)
		os.Remove(filepath.Join(client.Config.Directory, asset.Key))
	}
}
