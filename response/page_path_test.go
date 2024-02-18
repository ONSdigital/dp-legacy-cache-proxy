package response

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestCase struct {
	uri      string
	expected string
}

func TestGetPagePath(t *testing.T) {
	Convey("Given a list of URIs", t, func() {
		ctx := context.Background()

		visualisationsEndpoints := []TestCase{
			{
				uri:      "/visualisations/dvc1945/seasonalflu/index.html",
				expected: "/visualisations/dvc1945",
			},
			{
				uri:      "/visualisations/dvc1945/duetoinvolving/fallback.png",
				expected: "/visualisations/dvc1945",
			},
			{
				uri:      "/visualisations/dvc2775/fig5/css/tabstyles.css",
				expected: "/visualisations/dvc2775",
			},
			{
				uri:      "/visualisations/dvc2775/lib2/chosen.order.jquery.js",
				expected: "/visualisations/dvc2775",
			},
			{
				uri:      "/visualisations/dvc2775/lib2/globalStyle.css",
				expected: "/visualisations/dvc2775",
			},
			{
				uri:      "/visualisations/segment-1/segment-2/segment-3/",
				expected: "/visualisations/segment-1",
			},
		}
		resourceEndpoints := []TestCase{
			{
				uri:      "/generator?uri=/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024/426e63a0&format=csv",
				expected: "/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024",
			},
			{
				uri:      "/chartimage?uri=/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024/426e63a0",
				expected: "/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024",
			},
			{
				uri:      "/file?uri=/aboutus/whatwedo/programmesandprojects/sustainabledevelopmentgoals/c3ac8759.png",
				expected: "/aboutus/whatwedo/programmesandprojects/sustainabledevelopmentgoals",
			},
			{
				uri:      "/generator?format=csv&uri=/economy/grossvalueaddedgva/timeseries/abml/pn2/previous/v29",
				expected: "/economy/grossvalueaddedgva/timeseries/abml/pn2",
			},
			{
				uri:      "/chartconfig?uri=/test",
				expected: "/test",
			},
			{
				uri:      "/embed?uri=/test",
				expected: "/test",
			},
			{
				uri:      "/chart?uri=/test",
				expected: "/test",
			},
			{
				uri:      "/resource?uri=/test",
				expected: "/test",
			},
			{
				uri:      "/export?uri=/test",
				expected: "/test",
			},
		}
		resourceEndpointsWithEncodedURI := []TestCase{
			{
				uri:      "/generator?uri=%2Feconomy%2Feconomicoutputandproductivity%2Foutput%2Fbulletins%2Feconomicactivityandsocialchangeintheukrealtimeindicators%2F15february2024%2F426e63a0&format=csv",
				expected: "/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024",
			},
			{
				uri:      "/chartimage?uri=%2Feconomy%2Feconomicoutputandproductivity%2Foutput%2Fbulletins%2Feconomicactivityandsocialchangeintheukrealtimeindicators%2F15february2024%2F426e63a0",
				expected: "/economy/economicoutputandproductivity/output/bulletins/economicactivityandsocialchangeintheukrealtimeindicators/15february2024",
			},
			{
				uri:      "/file?uri=%2Faboutus%2Fwhatwedo%2Fprogrammesandprojects%2Fsustainabledevelopmentgoals%2Fc3ac8759.png",
				expected: "/aboutus/whatwedo/programmesandprojects/sustainabledevelopmentgoals",
			},
			{
				uri:      "/generator?format=csv&uri=%2Feconomy%2Fgrossvalueaddedgva%2Ftimeseries%2Fabml%2Fpn2%2Fprevious%2Fv29",
				expected: "/economy/grossvalueaddedgva/timeseries/abml/pn2",
			},
			{
				uri:      "/embed?&uri=%2Furi%2Fwith%2Ftrailing%2Fslash%2F",
				expected: "/uri/with/trailing/slash",
			},
		}
		bulletinsPaths := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/bulletins/producerpriceinflation/latest",
				expected: "/economy/inflationandpriceindices/bulletins/producerpriceinflation/latest",
			},
			{
				uri:      "/economy/inflationandpriceindices/bulletins/producerpriceinflations/latest/data",
				expected: "/economy/inflationandpriceindices/bulletins/producerpriceinflations/latest",
			},
			{
				uri:      "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022",
				expected: "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022",
			},
			{
				uri:      "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/",
				expected: "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022",
			},
			{
				uri:      "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/data",
				expected: "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022",
			},
			{
				uri:      "/generator?uri=/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022/30d7d6c2&format=csv",
				expected: "/economy/inflationandpriceindices/bulletins/producerpriceinflation/october2022",
			},
			{
				uri:      "/economy/inflationandpriceindices/bulletins/consumerpriceinflation/january2024",
				expected: "/economy/inflationandpriceindices/bulletins/consumerpriceinflation/january2024",
			},
			{
				uri:      "/segment-0/bulletins/segment-1/segment-2/segment-3/",
				expected: "/segment-0/bulletins/segment-1/segment-2",
			},
		}
		articlesPaths := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/articles/researchanddevelopmentsinthetransformationofukconsumerpricestatistics/december2023",
				expected: "/economy/inflationandpriceindices/articles/researchanddevelopmentsinthetransformationofukconsumerpricestatistics/december2023",
			},
			{
				uri:      "/economy/inflationandpriceindices/articles/researchanddevelopmentsinthetransformationofukconsumerpricestatistics/latest",
				expected: "/economy/inflationandpriceindices/articles/researchanddevelopmentsinthetransformationofukconsumerpricestatistics/latest",
			},
			{
				uri:      "/economy/inflationandpriceindices/articles/costofliving/latestinsights",
				expected: "/economy/inflationandpriceindices/articles/costofliving/latestinsights",
			},
			{
				uri:      "/economy/inflationandpriceindices/articles/trackingthelowestcostgroceryitemsukexperimentalanalysis/april2021toseptember2022/pdf",
				expected: "/economy/inflationandpriceindices/articles/trackingthelowestcostgroceryitemsukexperimentalanalysis/april2021toseptember2022",
			},
			{
				uri:      "/segment-0/articles/segment-1/segment-2/segment-3/",
				expected: "/segment-0/articles/segment-1/segment-2",
			},
		}
		methodologiesPaths := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/methodologies/consumerpriceinflationincludesall3indicescpihcpiandrpiqmi",
				expected: "/economy/inflationandpriceindices/methodologies/consumerpriceinflationincludesall3indicescpihcpiandrpiqmi",
			},
			{
				uri:      "/economy/inflationandpriceindices/methodologies/consumerpriceinflationincludesall3indicescpihcpiandrpiqmi/data",
				expected: "/economy/inflationandpriceindices/methodologies/consumerpriceinflationincludesall3indicescpihcpiandrpiqmi",
			},
			{
				uri:      "/file?uri=economy/inflationandpriceindices/methodologies/consumerpriceinflationincludesall3indicescpihcpiandrpiqmi/17a6d562.xls",
				expected: "economy/inflationandpriceindices/methodologies/consumerpriceinflationincludesall3indicescpihcpiandrpiqmi",
			},
			{
				uri:      "/segment-0/methodologies/segment-1/segment-2/segment-3/",
				expected: "/segment-0/methodologies/segment-1",
			},
		}
		qmisPaths := []TestCase{
			{
				uri:      "/economy/grossdomesticproductgdp/qmis/annualacquisitionsanddisposalsofcapitalassetssurveyqmi",
				expected: "/economy/grossdomesticproductgdp/qmis/annualacquisitionsanddisposalsofcapitalassetssurveyqmi",
			},
			{
				uri:      "/businessindustryandtrade/changestobusiness/mergersandacquisitions/qmis/mergersandacquisitionsmaqmi",
				expected: "/businessindustryandtrade/changestobusiness/mergersandacquisitions/qmis/mergersandacquisitionsmaqmi",
			},
			{
				uri:      "/employmentandlabourmarket/peopleinwork/publicsectorpersonnel/qmis/publicsectoremploymentqmi",
				expected: "/employmentandlabourmarket/peopleinwork/publicsectorpersonnel/qmis/publicsectoremploymentqmi",
			},
			{
				uri:      "/segment-0/qmis/segment-1/segment-2/segment-3/",
				expected: "/segment-0/qmis/segment-1",
			},
		}
		adHocsPaths := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/adhocs/009581cpiinflationbetween2010and2018",
				expected: "/economy/inflationandpriceindices/adhocs/009581cpiinflationbetween2010and2018",
			},
			{
				uri:      "/peoplepopulationandcommunity/armedforcescommunity/adhocs/56testingadhoc",
				expected: "/peoplepopulationandcommunity/armedforcescommunity/adhocs/56testingadhoc",
			},
			{
				uri:      "/file?uri=/economy/inflationandpriceindices/adhocs/009581cpiinflationbetween2010and2018/cpiinflationbetween2010and2018.xls",
				expected: "/economy/inflationandpriceindices/adhocs/009581cpiinflationbetween2010and2018",
			},
			{
				uri:      "/segment-0/adhocs/segment-1/segment-2/segment-3/",
				expected: "/segment-0/adhocs/segment-1",
			},
		}
		endsWithDataPaths := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/pricequotesseptember2023/data",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/pricequotesseptember2023",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/data",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindices/data",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindices",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindices/data/",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindices",
			},
			{
				uri:      "/test/data",
				expected: "/test",
			},
		}
		endsWithFileExtPaths := []TestCase{
			{
				uri:      "/file?uri=/economy/inflationandpriceindices/datasets/consumerpriceindices/current/mm23.csv",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindices/current",
			},
			{
				uri:      "/file?uri=/employmentandlabourmarket/peopleinwork/earningsandworkinghours/datasets/thisisatest/feb24/3a88ffc6.xlsx",
				expected: "/employmentandlabourmarket/peopleinwork/earningsandworkinghours/datasets/thisisatest/feb24",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/pricequotesseptember2023/upload-pricequotes202309.csv",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/pricequotesseptember2023",
			},
			{
				uri:      "/test/image.png",
				expected: "/test",
			},
		}
		timeSeriesPaths := []TestCase{
			{
				uri:      "/economy/grossdomesticproductgdp/timeseries/abmi/pn2/linechartconfig",
				expected: "/economy/grossdomesticproductgdp/timeseries/abmi/pn2",
			},
			{
				uri:      "/economy/grossdomesticproductgdp/timeseries/abmi/pn2/data/linechartconfig",
				expected: "/economy/grossdomesticproductgdp/timeseries/abmi/pn2",
			},
			{
				uri:      "/economy/grossdomesticproductgdp/timeseries/abmi/pn2/data/linechartconfig/",
				expected: "/economy/grossdomesticproductgdp/timeseries/abmi/pn2",
			},
			{
				uri:      "/economy/grossdomesticproductgdp/timeseries/ihyq",
				expected: "/economy/grossdomesticproductgdp/timeseries/ihyq",
			},
			{
				uri:      "/economy/grossdomesticproductgdp/timeseries/abmi/pn2",
				expected: "/economy/grossdomesticproductgdp/timeseries/abmi/pn2",
			},
			{
				uri:      "/economy/grossdomesticproductgdp/timeseries/abmi/pn2/data/linechartimage",
				expected: "/economy/grossdomesticproductgdp/timeseries/abmi/pn2",
			},
			{
				uri:      "/generator?format=csv&uri=/economy/inflationandpriceindices/timeseries/l55o/mm23",
				expected: "/economy/inflationandpriceindices/timeseries/l55o/mm23",
			},
		}
		datasetsPaths := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/pricequotesseptember2023",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindicescpiandretailpricesindexrpiitemindicesandpricequotes/pricequotesseptember2023",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindices",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindices",
			},
			{
				uri:      "/economy/inflationandpriceindices/datasets/consumerpriceindices/current",
				expected: "/economy/inflationandpriceindices/datasets/consumerpriceindices/current",
			},
			{
				uri:      "/employmentandlabourmarket/peopleinwork/earningsandworkinghours/datasets/analysisofwageandpriceincreases",
				expected: "/employmentandlabourmarket/peopleinwork/earningsandworkinghours/datasets/analysisofwageandpriceincreases",
			},
			{
				uri:      "/businessindustryandtrade/changestobusiness/mergersandacquisitions/datasets/timeseries/15march2024",
				expected: "/businessindustryandtrade/changestobusiness/mergersandacquisitions/datasets/timeseries/15march2024",
			},
			{
				uri:      "/segment-0/datasets/segment-1/segment-2/segment-3/",
				expected: "/segment-0/datasets/segment-1/segment-2",
			},
		}
		// The default behaviour for any URI that does not fit any of the previous groups is to remain unchanged
		defaultBehaviour := []TestCase{
			{
				uri:      "/economy/inflationandpriceindices/publications?filter=bulletin",
				expected: "/economy/inflationandpriceindices/publications?filter=bulletin",
			},
			{
				uri:      "/economy/regionalaccounts/grossdisposablehouseholdincome/datalist",
				expected: "/economy/regionalaccounts/grossdisposablehouseholdincome/datalist",
			},
			{
				uri:      "/economy/regionalaccounts/grossdisposablehouseholdincome",
				expected: "/economy/regionalaccounts/grossdisposablehouseholdincome",
			},
		}

		groupedTestCases := [][]TestCase{
			visualisationsEndpoints,
			resourceEndpoints,
			resourceEndpointsWithEncodedURI,
			bulletinsPaths,
			articlesPaths,
			methodologiesPaths,
			qmisPaths,
			adHocsPaths,
			endsWithDataPaths,
			endsWithFileExtPaths,
			timeSeriesPaths,
			datasetsPaths,
			defaultBehaviour,
		}

		Convey("When the 'getPagePath' function is called", func() {
			for _, testCases := range groupedTestCases {
				for _, tc := range testCases {
					pagePath, err := getPagePath(ctx, tc.uri)

					Convey("Then it should resolve the page path for "+tc.uri, func() {
						So(pagePath, ShouldEqual, tc.expected)
						So(err, ShouldBeNil)
					})
				}
			}
		})
	})
}
