package server

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	abci "github.com/Fincor-Blockchain/pou/abci/types"
	"github.com/Fincor-Blockchain/pou/libs/log"
	tmtypes "github.com/Fincor-Blockchain/pou/types"
	dbm "github.com/Fincor-Blockchain/tm-db"

	sdk "github.com/Fincor-Blockchain/cusp-sdk/types"
)

type (
	// AppCreator is a function that allows us to lazily initialize an
	// application using various configurations.
	AppCreator func(log.Logger, dbm.DB, io.Writer) abci.Application

	// AppExporter is a function that dumps all app state to
	// JSON-serializable structure and returns the current validator set.
	AppExporter func(log.Logger, dbm.DB, io.Writer, int64, bool, []string) (json.RawMessage, []tmtypes.GenesisValidator, error)
)

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := sdk.NewLevelDB("application", dataDir)
	return db, err
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile != "" {
		w, err = os.OpenFile(
			traceWriterFile,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0666,
		)
		return
	}
	return
}
