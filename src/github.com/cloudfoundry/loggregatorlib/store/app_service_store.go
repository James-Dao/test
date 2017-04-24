package store

import (
	"github.com/cloudfoundry/loggregatorlib/appservice"
	"github.com/cloudfoundry/loggregatorlib/cfcomponent"
	"github.com/cloudfoundry/loggregatorlib/store/cache"
	"github.com/cloudfoundry/storeadapter"

	"path"
)

type AppServiceStore struct {
	adapter storeadapter.StoreAdapter
	cache   cache.AppServiceCache
}

func NewAppServiceStore(adapter storeadapter.StoreAdapter, cache cache.AppServiceCache) *AppServiceStore {
	return &AppServiceStore{
		adapter: adapter,
		cache:   cache,
	}
}

func (s *AppServiceStore) Run(incomingChan <-chan appservice.AppServices) {

	for appServices := range incomingChan {
		cfcomponent.Logger.Debugf("AppStore: New services for %s", appServices.AppId)
		if len(appServices.Urls) == 0 {
			s.removeAppFromStore(appServices.AppId)
			continue
		}

		cachedAppServices := s.cache.Get(appServices.AppId)

		appServiceToAdd := []appservice.AppService{}
		appServiceToRemove := []appservice.AppService{}
		serviceUrls := make(map[string]bool)

		for _, serviceUrl := range appServices.Urls {
			serviceUrls[serviceUrl] = true

			appService := appservice.AppService{AppId: appServices.AppId, Url: serviceUrl}
			if !s.cache.Exists(appService) {
				appServiceToAdd = append(appServiceToAdd, appService)
			}
		}

		for _, appService := range cachedAppServices {
			if !serviceUrls[appService.Url] {
				appServiceToRemove = append(appServiceToRemove, appService)
			}
		}

		s.addToStore(appServiceToAdd)
		s.removeFromStore(appServiceToRemove)
		cfcomponent.Logger.Debugf("AppStore: Successfully updated app service %s", appServices.AppId)
	}
}

func (s *AppServiceStore) addToStore(appServices []appservice.AppService) {
	nodes := make([]storeadapter.StoreNode, len(appServices))
	for i, appService := range appServices {
		s.cache.Add(appService)
		nodes[i] = storeadapter.StoreNode{
			Key:   path.Join("/loggregator/services", appService.AppId, appService.Id()),
			Value: []byte(appService.Url),
		}
	}

	s.adapter.SetMulti(nodes)
}

func (s *AppServiceStore) removeFromStore(appServices []appservice.AppService) {
	keys := make([]string, len(appServices))
	for i, appService := range appServices {
		s.cache.Remove(appService)
		keys[i] = path.Join("/loggregator/services", appService.AppId, appService.Id())
	}
	if len(keys) > 0 {
		s.adapter.Delete(keys...)
	}
}

func (s *AppServiceStore) removeAppFromStore(appId string) {
	cfcomponent.Logger.Debugf("AppStore: removing app %s", appId)
	removedApps := s.cache.RemoveApp(appId)
	if len(removedApps) > 0 {
		s.adapter.Delete(path.Join("/loggregator/services", appId))
	}
}
