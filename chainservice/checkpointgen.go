package chainservice

import (
	"html/template"
	"os"

	"github.com/breez/breez/config"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/neutrino/headerfs"
)

// GenerateCheckpoints generates a code that populates an array of checkpoints.
// it does so by iterating the neutrino db and using a template file for code
// generation. A sample for the result can be observed in checkpoint.go
func GenerateCheckpoints(workingDir string, tplFilePath string, outputFilePath string) error {

	config, err := config.GetConfig(workingDir)
	if err != nil {
		return err
	}

	params, err := chainParams(config.Network)
	if err != nil {
		return err
	}

	neutrinoDataDir, db, err := getNeutrinoDB(workingDir)
	if err != nil {
		return err
	}

	blockHeaderStore, err := headerfs.NewBlockHeaderStore(neutrinoDataDir, db, params)
	if err != nil {
		return err
	}

	filterHeaderStore, err := headerfs.NewFilterHeaderStore(neutrinoDataDir, db, headerfs.RegularFilter, params, nil)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(tplFilePath)
	if err != nil {
		return err
	}
	writer, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, height, err := blockHeaderStore.ChainTip()
	if err != nil {
		return err
	}
	for i := 0; i < int(height/wire.CFCheckptInterval); i++ {
		height := uint32(i * wire.CFCheckptInterval)
		wireHeader, err := blockHeaderStore.FetchHeaderByHeight(height)
		filterHeader, err := filterHeaderStore.FetchHeaderByHeight(height)
		if err != nil {
			return err
		}

		err = tmpl.Execute(writer, Checkpoint{
			Height:       height,
			BlockHeader:  wireHeader,
			FilterHeader: filterHeader,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
