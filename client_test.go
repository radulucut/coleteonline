package coleteonline

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var currentTime time.Time = time.Now()

func Test_Client(t *testing.T) {
	server, err := startServer()
	if err != nil {
		t.Error(err)
	}
	url := "http://" + server.Addr
	t.Run("GetAuthBearer", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		token, err := client.GetAuthBearer()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*token, "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix())); diff != "" {
			t.Errorf("Token mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetAuthBearer_Timeout", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       100 * time.Millisecond,
		})
		client.authURL = url + "/auth/timeout/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		_, err := client.GetAuthBearer()
		if diff := cmp.Diff(err.Error(), "Post \"http://"+server.Addr+"/auth/timeout/token\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)"); diff != "" {
			t.Errorf("Error mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CreateOrder", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		order := newOrder()
		res, err := client.CreateOrder(&order)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newOrderResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("CreateOrder_Timeout", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       100 * time.Millisecond,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1/timeout"
		client.timeNow = func() time.Time {
			return currentTime
		}
		order := newOrder()
		_, err := client.CreateOrder(&order)
		if diff := cmp.Diff(err.Error(), "Post \"http://"+server.Addr+"/v1/timeout/order\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)"); diff != "" {
			t.Errorf("Error mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("CreateOrder_UseAddressIds", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1/order-with-address-id"
		client.timeNow = func() time.Time {
			return currentTime
		}
		order := newOrderWithAddressId()
		res, err := client.CreateOrder(&order)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newOrderResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrderOrice", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		order := newOrder()
		res, err := client.OrderPrice(&order)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newOrderPriceResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("OrderPrice_UseAddressIds", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1/order-with-address-id"
		client.timeNow = func() time.Time {
			return currentTime
		}
		order := newOrderWithAddressId()
		res, err := client.OrderPrice(&order)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newOrderPriceResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrderStatus", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		id := "id_1234"
		res, err := client.OrderStatus(&id)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newOrderStatusResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("AddressList", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		res, err := client.AddressList(1)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newAddressListResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ServiceList", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		res, err := client.ServiceList()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(res, newServiceListResponse()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("UserBalance", func(t *testing.T) {
		t.Parallel()
		client := NewClient(Config{
			ClientId:      "client_id",
			ClientSecret:  "client_secret",
			UseProduction: true,
			Timeout:       10 * time.Second,
		})
		client.authURL = url + "/auth/token"
		client.apiURL = url + "/v1"
		client.timeNow = func() time.Time {
			return currentTime
		}
		res, err := client.UserBalance()
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(*res, newUserBalance()); diff != "" {
			t.Errorf("Order response mismatch (-want +got):\n%s", diff)
		}
	})
}

func startServer() (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/auth/token", authTokenHandler)
	mux.HandleFunc("/auth/timeout/token", timeoutHandler)
	mux.HandleFunc("/v1/order", createOrderHandler)
	mux.HandleFunc("/v1/timeout/order", timeoutHandler)
	mux.HandleFunc("/v1/order-with-address-id/order", createOrderWithAddressIdHandler)
	mux.HandleFunc("/v1/order/price", orderPriceHandler)
	mux.HandleFunc("/v1/order-with-address-id/order/price", orderPriceWithAddressIdHandler)
	statusPath := "/v1/order/status/"
	mux.Handle(
		statusPath,
		http.StripPrefix(statusPath, http.HandlerFunc(orderStatusHandler)),
	)
	mux.HandleFunc("/v1/address", addressListHandler)
	mux.HandleFunc("/v1/service", serviceListHandler)
	mux.HandleFunc("/v1/user/balance", userBalanceHandler)
	server := &http.Server{
		Addr:    ":9876",
		Handler: mux,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Error starting server: %v", err))
		}
	}()
	tries := 0
	for {
		time.Sleep(100 * time.Millisecond)
		_, err := http.Get("http://" + server.Addr)
		if err == nil {
			break
		}
		tries++
		if tries > 5 {
			panic(fmt.Errorf("Expected server to start, got %v", err))
		}
	}
	return server, nil
}

func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
}

func authTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: content must be application/x-www-form-urlencoded",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Basic "+base64.StdEncoding.EncodeToString([]byte("client_id:client_secret")) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_client",
			Description: "Invalid client credentials",
		})
		w.Write(b)
		return
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) != "grant_type=client_credentials" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Missing parameter: `grant_type`",
		})
		w.Write(b)
		return
	}
	b, _ = json.Marshal(AuthToken{
		AccessToken: getTestJWT(currentTime.Add(2 * time.Hour).Unix()),
	})
	w.Write(b)
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be POST",
		})
		w.Write(b)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: content must be application/json",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expectedOrder, _ := json.Marshal(newOrder())
	if string(b) != string(expectedOrder) {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: invalid order",
		})
		w.Write(b)
		return
	}
	b, _ = json.Marshal(newOrderResponse())
	w.Write(b)
}

func createOrderWithAddressIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be POST",
		})
		w.Write(b)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: content must be application/json",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expectedOrder, _ := json.Marshal(newOrderWithAddressId())
	if string(b) != string(expectedOrder) {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: invalid order",
		})
		w.Write(b)
		return
	}
	b, _ = json.Marshal(newOrderResponse())
	w.Write(b)
}

func orderPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be POST",
		})
		w.Write(b)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: content must be application/json",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expectedOrder, _ := json.Marshal(newOrder())
	if string(b) != string(expectedOrder) {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: invalid order",
		})
		w.Write(b)
		return
	}
	b, _ = json.Marshal(newOrderPriceResponse())
	w.Write(b)
}

func orderPriceWithAddressIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be POST",
		})
		w.Write(b)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: content must be application/json",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expectedOrder, _ := json.Marshal(newOrderWithAddressId())
	if string(b) != string(expectedOrder) {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: invalid order",
		})
		w.Write(b)
		return
	}
	b, _ = json.Marshal(newOrderPriceResponse())
	w.Write(b)
}

func orderStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be GET",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	uniqueId := r.URL.Path
	if uniqueId != "id_1234" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: invalid uniqueId",
		})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(newOrderStatusResponse())
	w.Write(b)
}

func addressListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be GET",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	page := r.URL.Query().Get("page")
	if page != "1" {
		w.WriteHeader(http.StatusBadRequest)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: invalid page",
		})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(newAddressListResponse())
	w.Write(b)
}

func serviceListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be GET",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(newServiceListResponse())
	w.Write(b)
}

func userBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_request",
			Description: "Invalid request: method must be GET",
		})
		w.Write(b)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Bearer "+getTestJWT(currentTime.Add(2*time.Hour).Unix()) {
		w.WriteHeader(http.StatusUnauthorized)
		b, _ := json.Marshal(AuthResponseError{
			Name:        "invalid_token",
			Description: "Invalid token",
		})
		w.Write(b)
		return
	}
	b, _ := json.Marshal(newUserBalance())
	w.Write(b)
}

func newOrder() Order {
	return Order{
		Sender: Sender{
			Contact: &Contact{
				Name:    "Sender Name",
				Phone:   "0123456789",
				Phone2:  "0123456789",
				Company: "Company",
				Email:   "sender@test.com",
			},
			Address: &Address{
				CountryCode:    "RO",
				PostalCode:     "123456",
				City:           "City",
				County:         "County",
				CountyCode:     "B",
				Street:         "Street",
				Number:         "1",
				Building:       "Building",
				Entrance:       "Entrance",
				Floor:          "Floor",
				Apartment:      "Apartment",
				AdditionalInfo: "Additional info",
			},
			ValidationStrategy: ValidationStrategyTypeMinimal,
		},
		Recipient: Recipient{
			Contact: &Contact{
				Name:    "Recipient Name",
				Phone:   "0123456789",
				Phone2:  "0123456789",
				Company: "Company",
				Email:   "recipient@test.com",
			},
			Address: &Address{
				CountryCode:    "RO",
				PostalCode:     "123456",
				City:           "City",
				County:         "County",
				CountyCode:     "B",
				Street:         "Street",
				Number:         "1",
				Building:       "Building",
				Entrance:       "Entrance",
				Floor:          "Floor",
				Apartment:      "Apartment",
				AdditionalInfo: "Additional info",
			},
			ValidationStrategy: ValidationStrategyTypeMinimal,
		},
		Packages: []Package{
			{
				Weight: 1.0,
				Width:  1.0,
				Height: 1.0,
				Length: 1.0,
			},
		},
		ServiceSelectStrategy: ServiceSelectStrategy{
			SelectionType: ServiceTypeDirectId,
			ServiceIds:    []int64{1},
		},
		ExtraOptions: []map[string]interface{}{
			{
				"id": ExtraOptionIdOpenAtDelivery,
			},
		},
	}
}

func newOrderWithAddressId() Order {
	return Order{
		Sender: Sender{
			AddressId: 1,
		},
		Recipient: Recipient{
			AddressId: 2,
		},
		Packages: []Package{
			{
				Weight: 1.0,
				Width:  1.0,
				Height: 1.0,
				Length: 1.0,
			},
		},
		ServiceSelectStrategy: ServiceSelectStrategy{
			SelectionType: ServiceTypeDirectId,
			ServiceIds:    []int64{1},
		},
		ExtraOptions: []map[string]interface{}{
			{
				"id": ExtraOptionIdOpenAtDelivery,
			},
		},
	}
}

func newOrderResponse() OrderResponse {
	return OrderResponse{
		Service: OrderService{
			Price: ServicePrice{
				Total: 10.0,
				NoVat: 8.0,
			},
			Service: ServiceDetails{
				Id:          1,
				CourierName: "Courier Name",
				Name:        "Service Name",
			},
		},
		UniqueId:            "id_1234",
		AWB:                 "awb_1234",
		EstimatedPickupDate: "2023-01-01",
	}
}

func newOrderPriceResponse() OrderPriceResponse {
	return OrderPriceResponse{
		Selected: OrderService{
			Price: ServicePrice{
				Total: 10.0,
				NoVat: 8.0,
			},
			Service: ServiceDetails{
				Id:          1,
				CourierName: "Courier Name",
				Name:        "Service Name",
			},
		},
		List: []OrderService{
			{
				Price: ServicePrice{
					Total: 10.0,
					NoVat: 8.0,
				},
				Service: ServiceDetails{
					Id:          1,
					CourierName: "Courier Name",
					Name:        "Service Name",
				},
			},
		},
	}
}

func newOrderStatusResponse() OrderStatusResponse {
	return OrderStatusResponse{
		Summary: StatusSummary{
			UniqueId: "id_1234",
			AWB:      "awb_1234",
		},
		History: []StatusHistory{
			{
				DateTime:     time.Unix(1672574400, 0),
				UnixDateTime: 1672574400,
				StatusTextParts: StatusTextParts{
					Ro: StatusTextPart{
						Name:   "Status Name",
						Reason: "Status Reason",
					},
				},
				StatusComment: StatusComment{
					Ro: "Status Comment",
				},
				Code: 1,
			},
		},
	}
}

func newAddressListResponse() AddressListResponse {
	return AddressListResponse{
		Data: []OrderAddress{
			{
				Contact: Contact{
					Name:    "Sender Name",
					Phone:   "0123456789",
					Phone2:  "0123456789",
					Company: "Company",
					Email:   "sender@test.com",
				},
				Address: Address{
					CountryCode:    "RO",
					PostalCode:     "123456",
					City:           "City",
					County:         "County",
					CountyCode:     "B",
					Street:         "Street",
					Number:         "1",
					Building:       "Building",
					Entrance:       "Entrance",
					Floor:          "Floor",
					Apartment:      "Apartment",
					AdditionalInfo: "Additional info",
				},
				ValidationStrategy: ValidationStrategyTypeMinimal,
			},
		},
		Pagination: Pagination{
			TotalItems:  1,
			CurrentPage: 1,
			TotalPages:  1,
		},
	}
}

func newServiceListResponse() []ServiceResponse {
	return []ServiceResponse{
		{
			Id:          1,
			CourierName: "Courier Name",
			Name:        "Service Name",
			ExtraOptions: []ServiceExtraOption{
				{
					Id:             1,
					Name:           "Extra Option Name",
					RequiredFields: []string{"field1", "field2"},
					OptionalFields: []string{"field3", "field4"},
				},
			},
		},
	}
}

func newUserBalance() UserBalance {
	return UserBalance{
		Amount: 10.0,
		Bonus:  10.0,
	}
}

func getTestJWT(exp int64) string {
	return "header." +
		base64.RawURLEncoding.EncodeToString(
			[]byte(`{"exp":`+fmt.Sprintf("%d", exp)+`}`),
		) +
		".signature"
}
