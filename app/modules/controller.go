package modules

import (
	"github.com/gin-gonic/gin"
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

	router.GET("/", func(ctx *gin.Context) {
		r := struct {
			Hostname string `json:"hostname"`
		}{c.hostname}
		ctx.JSON(http.StatusOK, r)
	})

	router.GET(counterPath, func(ctx *gin.Context) {
		ids, err := c.counter.ListAllCounterId()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorFormatter(err))
		}
		r := struct {
			Ids []string `json:"ids"`
		}{ids}
		ctx.JSON(http.StatusOK, r)
	})

	router.POST(counterPath, func(ctx *gin.Context) {
		// TODO(kenji-kondo)
		// * error handling and messages
		// * handle if query is empty
		// * error handlingがうまく行っていない。両方とも通っているので、それはぜせいする
		to := ctx.Query(toQueryKey)
		toInt64, err := strconv.ParseInt(to, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorFormatter(err))
		}
		id, errGenerateCounter := c.counter.GenerateCounter(toInt64)
		if errGenerateCounter != nil {
			ctx.JSON(http.StatusInternalServerError, errorFormatter(errGenerateCounter))
		}
		r := struct {
			Id string `json:"id"`
		}{id}
		ctx.JSON(http.StatusOK, r)

	})

	router.GET(counterPath + "/:id", func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		r, err := c.counter.GetCounter(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorFormatter(err))
		}
		ctx.JSON(http.StatusOK, r)
	})

	router.POST(counterPath + "/:id" + stopPath, func(ctx *gin.Context) {
		id := ctx.Params.ByName("id")
		err := c.counter.DeleteCounter(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorFormatter(err))
		}
		ctx.JSON(http.StatusNoContent, "")
	})

	c.router = router
}

func (c *Controller) Run() {
	c.router.Run(":" + c.listenPort)
}

func errorFormatter(err error) interface{} {
	return struct {
		Error string `json:"error"`
	}{err.Error()}
}