package api

import (
	"os"
	"fmt"
	"time"

	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/api/applications"
	"github.com/cloudfoundry/cli/cf/api/authentication"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/i18n/detection"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
	"github.com/cloudfoundry/cli/cf/trace"
)

type CloudController struct {
	config  core_config.Repository
	gateway net.Gateway
	appRepo applications.CloudControllerApplicationRepository
	appSummaryRepo api.CloudControllerAppSummaryRepository
}

func NewCloudController() (cc CloudController) {

	errorHandler := func(err error) {
		if err != nil {
			fmt.Sprintf("Config error: %s", err)
		}
	}
	cc.config = core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), errorHandler)
	cc.gateway = net.NewCloudControllerGateway(cc.config, time.Now, nil)
	cc.gateway.SetTokenRefresher(authentication.NewUAAAuthenticationRepository(cc.gateway, cc.config))
	cc.appRepo = applications.NewCloudControllerApplicationRepository(cc.config, cc.gateway)
	cc.appSummaryRepo = api.NewCloudControllerAppSummaryRepository(cc.config, cc.gateway)

	// I18N usage in the library will cause the app to crash unless this is initialized
	i18n.T = i18n.Init(cc.config, &detection.JibberJabberDetector{})

	if os.Getenv("CF_TRACE") != "" {
		trace.Logger = trace.NewLogger(os.Getenv("CF_TRACE"))
	} else {
		trace.Logger = trace.NewLogger(cc.config.Trace())
	}

	return
}

func (cc *CloudController) GetApplication(appName string) (app models.Application, err error) {
	app, err = cc.appRepo.Read(appName)
	if nil != err {
		return
	}
	app, err = cc.appSummaryRepo.GetSummary(app.Guid)
	return
}

func (cc *CloudController) UpdateApplication(app *models.Application, params models.AppParams) (err error) {
	*app, err = cc.appRepo.Update(app.Guid, params)
	return
}

func (cc *CloudController) StartApplication(app *models.Application) (err error) {
	state := "STARTED"
	*app, err = cc.appRepo.Update(app.Guid, models.AppParams{State: &state})
	*app, err = cc.appSummaryRepo.GetSummary(app.Guid)
	for nil == err && app.InstanceCount > app.RunningInstances {
		time.Sleep(time.Duration(1 * time.Second))
		*app, err = cc.appSummaryRepo.GetSummary(app.Guid)
	}
	return
}

func (cc *CloudController) StopApplication(app *models.Application) (err error) {
	state := "STOPPED"
	*app, err = cc.appRepo.Update(app.Guid, models.AppParams{State: &state})
	return
}
