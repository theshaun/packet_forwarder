// Copyright © 2017 The Things Network. Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package util

import (
	"net/http"
	"os"
	"time"

	cliHandler "github.com/TheThingsNetwork/go-utils/handlers/cli"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"
	esHandler "github.com/TheThingsNetwork/ttn/utils/elasticsearch/handler"
	"github.com/apex/log"
	levelHandler "github.com/apex/log/handlers/level"
	multiHandler "github.com/apex/log/handlers/multi"
	"github.com/spf13/viper"
	elastic "github.com/tj/go-elastic"
)

func GetLogger() ttnlog.Interface {
	logLevel := log.InfoLevel
	if viper.GetBool("verbose") {
		logLevel = log.DebugLevel
	}

	handlers := []log.Handler{levelHandler.New(cliHandler.New(os.Stdout), logLevel)}

	if viper.GetBool("elasticsearch.enable") {
		esClient := elastic.New(viper.GetString("elasticsearch.address"))
		esClient.HTTPClient = &http.Client{Timeout: 5 * time.Second}
		handlers = append(handlers, levelHandler.New(esHandler.New(&esHandler.Config{
			Client:     esClient,
			Prefix:     viper.GetString("elasticsearch.prefix"),
			BufferSize: 10,
		}), logLevel))
	}

	ctx := apex.Wrap(&log.Logger{Handler: multiHandler.New(handlers...)})
	return ctx
}
