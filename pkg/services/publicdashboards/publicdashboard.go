package publicdashboards

import (
	"context"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/publicdashboards/api"
	"github.com/grafana/grafana/pkg/services/publicdashboards/service"
	"github.com/grafana/grafana/pkg/setting"
)

func ProvideService(features featuremgmt.FeatureToggles, cfg *setting.Cfg, store PublicDashboardStore, routeRegister routing.RouteRegister,
) (*PublicDashboardService, error) {
	var errDeclareRoles error

	s := service.ProvidePublicDashboardService(cfg, store)

	if !s.IsDisabled() {
		api := api.PublicDashboardAPI{
			RouteRegister:          routeRegister,
			PublicDashboardService: s,
		}
		api.RegisterAPIEndpoints()
	}

	return s, errDeclareRoles
}

//go:generate mockery --name DashboardProvisioningService --structname FakeDashboardProvisioning --inpackage --filename dashboard_provisioning_mock.go
// DashboardProvisioningService is a service for operating on provisioned dashboards.
type PublicDashboardService interface {
	GetPublicDashboard(ctx context.Context, accessToken string) (*models.Dashboard, error)
	GetPublicDashboardConfig(ctx context.Context, orgId int64, dashboardUid string) (*models.PublicDashboard, error)
	SavePublicDashboardConfig(ctx context.Context, dto *SavePublicDashboardConfigDTO) (*models.PublicDashboard, error)
	BuildPublicDashboardMetricRequest(ctx context.Context, dashboard *models.Dashboard, publicDashboard *models.PublicDashboard, panelId int64) (dtos.MetricRequest, error)
}

//go:generate mockery --name Store --structname FakeDashboardStore --inpackage --filename store_mock.go
// Store is a dashboard store.
type PublicDashboardStore interface {
	GetPublicDashboard(ctx context.Context, accessToken string) (*models.PublicDashboard, *models.Dashboard, error)
	GetPublicDashboardConfig(ctx context.Context, orgId int64, dashboardUid string) (*models.PublicDashboard, error)
	GenerateNewPublicDashboardUid(ctx context.Context) (string, error)
	SavePublicDashboardConfig(ctx context.Context, cmd models.SavePublicDashboardConfigCommand) (*models.PublicDashboard, error)
	UpdatePublicDashboardConfig(ctx context.Context, cmd models.SavePublicDashboardConfigCommand) error
}
