package modules

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Controller struct {
	counter Counter
	router *gin.Engine
	listenPort string
	hostname string
}

const (
	counterPath string = "/counter"
	stopPath string = "/stop"
	toQueryKey string = "to"
)

// Initialize Controller instance. You would do this method first.
func NewController(counter Counter, listenPort string, hostname string) *Controller {
	c := &Controller{
		counter:    counter,
		listenPort: listenPort,
		hostname:   hostname,
	}
	c.setupRouter()
	return c
}

func (c *Controller) setupRouter() {
	router := gin.Default()

	// Return hostname against "GET /"
	router.GET("/", func(ctx *gin.Context) {
		r := struct {
			Hostname string `json:"hostname"`
		}{c.hostname}
		ctx.JSON(http.StatusOK, r)
	})

	// Return all registered counter IDs against "GET /counter"
	router.GET(counterPath, func(ctx *gin.Context) {
		ids, err := c.counter.ListAllCounterId()

		// Return 500 if it got some errors when IDs from DB
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, errorFormatter(http.StatusText(http.StatusInternalServerError)))
			return
		}

		r := struct {
			Ids []string `json:"ids"`
		}{ids}
		ctx.JSON(http.StatusOK, r)
	})


	// Generate a new counter and return its counter ID against "POST /counter?to=[int]"
	router.POST(counterPath, func(ctx *gin.Context) {
		to := ctx.Query(toQueryKey)

		// Return 400 if "to" param is empty
		if to == "" {
			ctx.JSON(http.StatusBadRequest, errorFormatter("param to is required"))
			return
		}

		toInt64, err := strconv.ParseInt(to, 10, 64)
		// Return 400 if the value of the param "to" is invalid.
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorFormatter(fmt.Sprintf("the value %s is invalid", to)))
			return
		}

		id, errGenerateCounter := c.counter.GenerateCounter(toInt64)
		// Return 500 if it failed to generate counter by some internal reasons.
		if errGenerateCounter != nil {
			logrus.Error(errGenerateCounter)
			ctx.JSON(http.StatusInternalServerError, errorFormatter(http.StatusText(http.StatusInternalServerError)))
			return
		}

		r := struct {
			Id string `json:"id"`
		}{id}
		ctx.JSON(http.StatusCreated, r)
	})

	// Return counter corresponding to the specified ID against "GET /counter/:id"
	router.GET(counterPath + "/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		r, err := c.counter.GetCounter(id)

		// Return 500 if internal error occurs
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, errorFormatter(http.StatusText(http.StatusInternalServerError)))
			return
		}

		// Return 404 if such counter doesn't exist.
		if !r.counterExistence  {
			ctx.JSON(http.StatusNotFound, errorFormatter(fmt.Sprintf("no such counter with %s", id)))
			return
		}

		ctx.JSON(http.StatusOK, r)
	})

	// Delete the counter with the given ID and return no content against "POST /counter/:id/stop"
	router.POST(counterPath + "/:id" + stopPath, func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		err := c.counter.DeleteCounter(id)
		// Return 500 if it failed to delete a counter.
		if err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusInternalServerError, errorFormatter(http.StatusText(http.StatusInternalServerError)))
		}
		ctx.JSON(http.StatusNoContent, nil)
	})

	// Return 404 Not Found against no route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, errorFormatter(http.StatusText(http.StatusNotFound)))
	})

	c.router = router
}

// Run API server
func (c *Controller) Run() error {
	err := c.router.Run(":" + c.listenPort)
	if err != nil {
		return err
	}
	return nil
}

func errorFormatter(s string) interface{} {
	return struct {
		Error string `json:"error"`
	}{s}
}