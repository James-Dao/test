package store_test

import (
	. "github.com/cloudfoundry/loggregatorlib/store"
	"github.com/cloudfoundry/storeadapter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path"

	"github.com/cloudfoundry/loggregatorlib/appservice"
	"github.com/cloudfoundry/loggregatorlib/store/cache"
)

var _ = Describe("AppServiceStore", func() {
	var store *AppServiceStore
	var adapter storeadapter.StoreAdapter
	var incomingChan chan appservice.AppServices

	var app1Service1 appservice.AppService
	var app1Service2 appservice.AppService
	var app2Service1 appservice.AppService

	assertInStore := func(appServices ...appservice.AppService) {
		for _, appService := range appServices {
			Eventually(func() error {
				_, err := adapter.Get(path.Join("/loggregator/services/", appService.AppId, appService.Id()))
				return err
			}).ShouldNot(HaveOccurred())
		}
	}

	assertNotInStore := func(appServices ...appservice.AppService) {
		for _, appService := range appServices {
			Eventually(func() error {
				_, err := adapter.Get(path.Join("/loggregator/services/", appService.AppId, appService.Id()))
				return err
			}).Should(Equal(storeadapter.ErrorKeyNotFound))
		}
	}

	assertAppNotInStore := func(appIds ...string) {
		for _, appId := range appIds {
			Eventually(func() error {
				_, err := adapter.Get(path.Join("/loggregator/services/", appId))
				return err
			}).Should(Equal(storeadapter.ErrorKeyNotFound))
		}
	}

	BeforeEach(func() {
		adapter = etcdRunner.Adapter()

		incomingChan = make(chan appservice.AppServices)
		c := cache.NewAppServiceCache()

		store = NewAppServiceStore(adapter, c)
	})

	AfterEach(func() {
		err := adapter.Disconnect()
		Expect(err).NotTo(HaveOccurred())

		if incomingChan != nil {
			close(incomingChan)
		}
	})

	Context("when the incoming chan is closed", func() {
		BeforeEach(func() {
			close(incomingChan)
		})

		AfterEach(func() {
			incomingChan = nil //want the clean-up in AfterEach not to panic
		})

		It("should return", func(done Done) {
			store.Run(incomingChan)
			close(done)
		})
	})

	Context("when the store has data", func() {
		BeforeEach(func() {
			go store.Run(incomingChan)

			app1Service1 = appservice.AppService{AppId: "app-1", Url: "syslog://example.com:12345"}
			app1Service2 = appservice.AppService{AppId: "app-1", Url: "syslog://example.com:12346"}
			app2Service1 = appservice.AppService{AppId: "app-2", Url: "syslog://example.com:12345"}

			incomingChan <- appservice.AppServices{
				AppId: app1Service1.AppId,
				Urls:  []string{app1Service1.Url, app1Service2.Url},
			}
			incomingChan <- appservice.AppServices{
				AppId: app2Service1.AppId,
				Urls:  []string{app2Service1.Url},
			}
			Eventually(incomingChan).Should(BeEmpty())
			assertInStore(app1Service1, app1Service2, app2Service1)
		})

		It("does not modify the store, if the incoming data is already there", func(done Done) {
			events, stop, _ := adapter.Watch("/loggregator/services")

			incomingChan <- appservice.AppServices{
				AppId: app1Service1.AppId,
				Urls:  []string{app1Service1.Url, app1Service2.Url},
			}

			Consistently(events).Should(BeEmpty())

			stop <- true

			close(done)
		}, 2)

		Context("when there is new data for the store", func() {
			Context("when an existing app has a new service", func() {
				It("adds that service to the store", func(done Done) {
					app2Service2 := appservice.AppService{AppId: app2Service1.AppId, Url: "syslog://new.example.com:12345"}

					incomingChan <- appservice.AppServices{
						AppId: app2Service1.AppId,
						Urls:  []string{app2Service1.Url, app2Service2.Url},
					}

					assertInStore(app2Service1, app2Service2)

					close(done)
				})
			})

			Context("when a new app appears", func() {
				It("adds that app and its services to the store", func(done Done) {
					app3Service1 := appservice.AppService{AppId: "app-3", Url: "syslog://app3.example.com:12345"}
					app3Service2 := appservice.AppService{AppId: "app-3", Url: "syslog://app3.example.com:12346"}

					incomingChan <- appservice.AppServices{
						AppId: app3Service1.AppId,
						Urls:  []string{app3Service1.Url, app3Service2.Url},
					}

					assertInStore(app3Service1, app3Service2)

					close(done)
				})
			})
		})

		Context("when a service or app should be removed", func() {
			Context("when an existing app loses one of its services", func() {
				It("removes that service from the store", func(done Done) {
					incomingChan <- appservice.AppServices{
						AppId: app1Service1.AppId,
						Urls:  []string{app1Service1.Url},
					}

					assertInStore(app1Service1)
					assertNotInStore(app1Service2)

					close(done)
				})
			})

			Context("when an existing app loses all of its services", func() {
				It("removes the app entirely", func(done Done) {
					incomingChan <- appservice.AppServices{
						AppId: app1Service1.AppId,
						Urls:  []string{},
					}

					assertNotInStore(app1Service1, app1Service2)
					assertAppNotInStore(app1Service1.AppId)

					incomingChan <- appservice.AppServices{
						AppId: app1Service1.AppId,
						Urls:  []string{app1Service1.Url, app1Service2.Url},
					}

					assertInStore(app1Service1, app1Service2)

					close(done)
				})
			})
		})

		Describe("with multiple updates to the same app-id", func() {
			It("should perform the updates correctly in the store", func(done Done) {
				// Remove service 2
				incomingChan <- appservice.AppServices{
					AppId: app1Service1.AppId,
					Urls:  []string{app1Service1.Url},
				}

				assertInStore(app1Service1)
				assertNotInStore(app1Service2)

				// Add service 2 back
				incomingChan <- appservice.AppServices{
					AppId: app1Service1.AppId,
					Urls:  []string{app1Service1.Url, app1Service2.Url},
				}

				assertInStore(app1Service1, app1Service2)

				// Remove service 1
				incomingChan <- appservice.AppServices{
					AppId: app1Service1.AppId,
					Urls:  []string{app1Service2.Url},
				}

				assertInStore(app1Service2)
				assertNotInStore(app1Service1)

				close(done)
			})
		})
	})
})
