package main

import (
    "log"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "estimate-backend/internal/interface/controller"
    "estimate-backend/internal/usecase"
    // TODO: Add repository implementations
)

func main() {
    // Initialize Echo
    e := echo.New()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // TODO: Initialize repositories
    // For now, we'll use mock repositories

    // Initialize use cases
    processUseCase := usecase.NewProcessUseCase(nil) // TODO: Add process repository
    estimateUseCase := usecase.NewEstimateUseCase(nil, nil, nil, nil, nil) // TODO: Add repositories
    cocomoUseCase := usecase.NewCOCOMOUseCase(nil) // TODO: Add COCOMO repository

    // Initialize controllers
    processController := controller.NewProcessController(processUseCase)
    estimateController := controller.NewEstimateController(estimateUseCase)
    cocomoController := controller.NewCOCOMOController(cocomoUseCase)

    // Register routes
    processController.RegisterRoutes(e)
    estimateController.RegisterRoutes(e)
    cocomoController.RegisterRoutes(e)

    // Start server
    log.Fatal(e.Start(":8080"))
}