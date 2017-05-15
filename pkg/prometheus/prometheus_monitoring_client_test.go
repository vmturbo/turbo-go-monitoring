package prometheus

import (
	"flag"
	"github.com/chlam4/monitoring/pkg/client"
	"github.com/chlam4/monitoring/pkg/model/entity"
	"github.com/chlam4/monitoring/pkg/model/property"
	"github.com/chlam4/monitoring/pkg/model/resource"
	"github.com/chlam4/monitoring/pkg/repository"
	"github.com/chlam4/monitoring/pkg/repository/simpleRepo"
	"github.com/chlam4/monitoring/pkg/template"
	"testing"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "5")
}

func TestPrometheusMonitor(t *testing.T) {
	//
	// What metrics do you want Prometheus to collect?
	//
	metricDefs := template.MonitoringTemplate{
		template.MakeMetricMetaWithDefaultSetter(entity.NODE, resource.CPU, property.USED),
		template.MakeMetricMetaWithDefaultSetter(entity.NODE, resource.MEM, property.USED),
		template.MakeMetricMetaWithDefaultSetter(entity.NODE, resource.MEM, property.CAP),
		template.MakeMetricMetaWithDefaultSetter(entity.NODE, resource.MEM, property.AVERAGE),
		template.MakeMetricMetaWithDefaultSetter(entity.NODE, resource.MEM, property.PEAK),
		template.MakeMetricMetaWithDefaultSetter(entity.NODE, resource.NETWORK, property.USED),
		template.MakeMetricMetaWithDefaultSetter(entity.POD, resource.MEM, property.USED),
	}
	//
	// What entities do you want Prometheus to monitor?
	//
	repoEntities := []repository.RepositoryEntity{
		simpleRepo.NewSimpleMetricRepoEntity(entity.NODE, "192.168.99.100"),
		simpleRepo.NewSimpleMetricRepoEntity(entity.NODE, "10.10.172.235"),
		simpleRepo.NewSimpleMetricRepoEntity(entity.POD, "10.0.2.15"),
	}
	repo := simpleRepo.NewSimpleMetricRepo()
	repo.SetEntityInstances(repoEntities)
	//
	// Construct the monitor target
	//
	monitorTarget := client.MakeMonitorTarget(repo, metricDefs)
	//
	// Call Prometheus to collect metrics
	//
	promServerUrl := "http://192.168.99.100:30900"
	promMonitor, err := NewPrometheusMonitor(promServerUrl)
	if err != nil {
		t.Errorf("Error instantiating a Prometheus Monitor instance: %s", err)
	}
	promMonitor.Monitor(&monitorTarget)
	//
	// Process the collected metrics
	//
	for _, repoEntity := range repo.GetAllEntityInstances() {
		t.Logf("Metrics collected for (%v, %v) are as follows:\n %s",
			repoEntity.GetType(), repoEntity.GetId(), repoEntity.GetAllMetrics())
	}
}
