package coverage

import (
	"path"
	"runtime"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const petStorePath = "fixtures/petstore.json"

var _ = Describe("Swagger analysis", func() {

	FContext("With swagger petstore", func() {

		_, p, _, ok := runtime.Caller(0)
		if !ok {
			panic("Not possible to get test file path")
		}
		docPath := path.Join(path.Dir(p), petStorePath)

		Context("Without URI filter", func() {

			It("Should build full REST API structure", func() {
				expectedRestAPI := map[string]map[string]*RequestStats{
					"/pets": map[string]*RequestStats{
						"POST": &RequestStats{
							Body:         map[string]int{},
							Query:        map[string]int{},
							ParamsHit:    0,
							ParamsNum:    6,
							MethodCalled: false,
							Path:         "/pets",
							Method:       "POST",
						},
						"GET": &RequestStats{
							Body: nil,
							Query: map[string]int{
								"tags":  0,
								"limit": 0,
							},
							ParamsHit:    0,
							ParamsNum:    2,
							MethodCalled: false,
							Path:         "/pets",
							Method:       "GET",
						},
					},
					"/pets/{name}": map[string]*RequestStats{
						"GET": &RequestStats{
							Body:         nil,
							Query:        map[string]int{},
							ParamsHit:    0,
							ParamsNum:    0,
							MethodCalled: false,
							Path:         "/pets/{name}",
							Method:       "GET",
						},
						"DELETE": &RequestStats{
							Body:         nil,
							Query:        map[string]int{},
							ParamsHit:    0,
							ParamsNum:    0,
							MethodCalled: false,
							Path:         "/pets/{name}",
							Method:       "DELETE",
						},
					},
				}

				restAPI, err := getRESTApi(docPath, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(restAPI).To(Equal(expectedRestAPI))
			})

		})

		Context("With URI filter", func() {

			It("Should build filtered REST API structure", func() {
				expectedRestAPI := map[string]map[string]*RequestStats{
					"/pets/{name}": map[string]*RequestStats{
						"GET": &RequestStats{
							Body:         nil,
							Query:        map[string]int{},
							ParamsHit:    0,
							ParamsNum:    0,
							MethodCalled: false,
							Path:         "/pets/{name}",
							Method:       "GET",
						},
						"DELETE": &RequestStats{
							Body:         nil,
							Query:        map[string]int{},
							ParamsHit:    0,
							ParamsNum:    0,
							MethodCalled: false,
							Path:         "/pets/{name}",
							Method:       "DELETE",
						},
					},
				}

				restAPI, err := getRESTApi(docPath, "/pets/{name}")
				Expect(err).NotTo(HaveOccurred())
				Expect(restAPI).To(Equal(expectedRestAPI))

			})
		})

		It("Should add swagger params", func() {
			document, err := loads.JSONSpec(docPath)
			Expect(err).NotTo(HaveOccurred())

			By("Resolving referenced definitions")
			reqStats := RequestStats{
				Body:         nil,
				Query:        map[string]int{},
				ParamsHit:    0,
				ParamsNum:    0,
				MethodCalled: false,
				Path:         "/pets",
				Method:       "POST",
			}
			expectedReqStats := reqStats
			expectedReqStats.Body = map[string]int{}
			expectedReqStats.ParamsNum = 6

			addSwaggerParams(&reqStats, document.Analyzer.ParamsFor("POST", "/pets"), document.Spec().Definitions)
			Expect(reqStats).To(Equal(expectedReqStats))

			By("Not resolving referenced definitions")
			reqStats = RequestStats{
				Body:         nil,
				Query:        map[string]int{},
				ParamsHit:    0,
				ParamsNum:    0,
				MethodCalled: false,
				Path:         "/pets",
				Method:       "GET",
			}
			expectedReqStats = reqStats
			expectedReqStats.Query = map[string]int{
				"tags":  0,
				"limit": 0,
			}
			expectedReqStats.ParamsNum = 2

			addSwaggerParams(&reqStats, document.Analyzer.ParamsFor("GET", "/pets"), document.Spec().Definitions)
			Expect(reqStats).To(Equal(expectedReqStats))
		})

		It("Should count params from referenced models", func() {
			document, err := loads.JSONSpec(docPath)
			Expect(err).NotTo(HaveOccurred())

			params := document.Analyzer.ParamsFor("POST", "/pets")
			pCnt := countRefParams(params["body#Pet"].Schema, document.Spec().Definitions)
			Expect(pCnt).To(Equal(6))

			pCnt = countRefParams(params["body#Pet"].Schema, spec.Definitions{})
			Expect(pCnt).To(Equal(0))
		})
	})
})
