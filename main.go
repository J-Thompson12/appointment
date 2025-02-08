package main

import (
	"github.com/justinthompson/appointment/pkg/appointment"
	"github.com/justinthompson/appointment/pkg/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	appManager, err := appointment.NewAppointmentManager()
	if err != nil {
		e.Logger.Fatal(e.Start(":8000"))
	}

	handlers.BuildRouter(e, appManager)
	e.Logger.Fatal(e.Start(":8000"))
}
