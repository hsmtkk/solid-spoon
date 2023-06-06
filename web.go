package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/iterator"
)

const COLLECTION = "items"

func main() {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT env var is not defined")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("failed to parse %s as int: %v", portStr, err)
	}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("PROJECT_ID env var is not defined")
	}

	handler := newHandler(projectID)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler.List)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type handler struct {
	projectID string
}

func newHandler(projectID string) *handler {
	return &handler{projectID}
}

type item struct {
	ID   int    `json:"id" firestore:"id"`
	Text string `json:"text" firestore:"text"`
}

func (h *handler) List(ectx echo.Context) error {
	ctx := ectx.Request().Context()
	client, err := firestore.NewClient(ctx, h.projectID)
	if err != nil {
		return fmt.Errorf("failed to init firestore client: %w", err)
	}
	defer client.Close()
	items := []item{}
	iter := client.Collection(COLLECTION).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return fmt.Errorf("iteration failed: %w", err)
		}
		item := item{}
		if err := doc.DataTo(&item); err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}
		items = append(items, item)
	}
	return ectx.JSON(http.StatusOK, items)
}
