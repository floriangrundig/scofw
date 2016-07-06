package wktreeobserver

import (
	"github.com/floriangrundig/scofw/config"
	"github.com/floriangrundig/scofw/util"
)

type WorkTreeObserver struct {
	config *config.Config
	util   *util.Util
}

func New(config *config.Config, util *util.Util) *WorkTreeObserver {
	return &WorkTreeObserver{
		config: config,
		util:   util,
	}
}

func (observer *WorkTreeObserver) Start() {

}
