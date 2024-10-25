package handlers

import (
	"fmt"
	"github.com/Archetarcher/metrics.git/internal/server/domain"
	"github.com/go-resty/resty/v2"
)

func ExampleMetricsHandler_UpdateMetrics() {
	c := int64(2896127014)

	req := resty.New().R()
	req.Method = "POST"
	req.URL = conf.server.URL + "/update/{type}/{name}/{value}"
	params := map[string]string{"type": "counter", "name": "counter_value", "value": fmt.Sprintf("%d", c)}
	req.SetPathParams(params)
	resp, _ := req.Send()

	fmt.Println(resp)

	params2 := map[string]string{"type": "gauge", "name": "name", "value": "value"}
	req.SetPathParams(params2)
	resp2, _ := req.Send()

	fmt.Println(resp2)

}

func ExampleMetricsHandler_GetMetrics() {
	req := resty.New().R()
	req.Method = "GET"
	req.URL = conf.server.URL + "/value/{type}/{name}"
	params := map[string]string{"type": "counter", "name": "counter_value"}
	req.SetPathParams(params)
	resp, _ := req.Send()

	fmt.Println(resp)

}

func ExampleMetricsHandler_UpdateMetricsJSON() {
	c := int64(2896127014)

	req := resty.New().R()
	req.Method = "POST"
	req.URL = conf.server.URL + "/update/"
	body := domain.Metrics{
		ID:    "counter_value_2",
		MType: "counter",
		Delta: &c,
		Value: nil,
	}
	req.SetBody(body)
	resp, _ := req.Send()

	fmt.Println(resp)

}

func ExampleMetricsHandler_GetMetricsJSON() {
	req := resty.New().R()
	req.Method = "POST"
	req.URL = conf.server.URL + "/value/"
	body := map[string]string{"type": "counter", "id": "counter_value_2"}
	req.SetBody(body)
	resp, _ := req.Send()

	fmt.Println(resp)

}

func ExampleMetricsHandler_UpdatesMetrics() {
	c := int64(2896127014)
	g := 0.31167763133187076

	req := resty.New().R()
	req.Method = "POST"
	req.URL = conf.server.URL + "/updates/"
	body := []domain.Metrics{
		{
			ID:    "gauge_value_2",
			MType: "gauge",
			Delta: nil,
			Value: &g,
		},
		{
			ID:    "counter_value_3",
			MType: "counter",
			Delta: &c,
			Value: nil,
		},
	}

	req.SetBody(body)
	resp, _ := req.Send()

	fmt.Println(resp)

}
