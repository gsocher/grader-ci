package route

// func TestGithubWebhook(t *testing.T) {
// 	router := mux.NewRouter()
// 	amqpClient := service.NewMockClient()
// 	statusRepo := repo.NewInMemoryStatusRepo()
// 	builder := service.NewBuilder(amqpClient, statusRepo)

// 	// add the route to the router
// 	RegisterGithubWebhookRoutes(router, builder)

// 	ts := httptest.NewServer(router)
// 	defer ts.Close()

// 	// create a test webhook payload
// 	var payload githubWebhookRequest
// 	payload.Repository.FullName = "docker/docker"

// 	b, _ := json.Marshal(payload)
// 	buf := bytes.NewBuffer(b)

// 	// execute post
// 	r, err := http.Post(ts.URL+pathURLGithubWebhookAPI, "application/json", buf)
// 	if err != nil {
// 		t.Fatalf("unexpected http client err: %v", err)
// 	}

// 	// check error code
// 	if r.StatusCode != http.StatusOK {
// 		t.Fatalf("expected %v got %v", http.StatusOK, r.StatusCode)
// 	}

// }