package hystrix

import (
	"awesomeProject/gin-contrib"
	"net/http"
	"time"
)

func Wrapper(size,reqThreshold int,filedThreshold float64,brokentime time.Duration)gin.HandlerFunc{
	r := NewRollingWindow(size,reqThreshold,filedThreshold,brokentime)
	r.Start()
	r.Monitor()
	r.ShowRequestsli()
	return func(c *gin.Context){
		if r.Broken() {
			c.String(http.StatusInternalServerError,"reject by hystri")
			c.Abort()
			return
		}
		c.Next()
		if c.Writer.Status() != http.StatusOK {
			r.RecordReqResult(false)
		}else {
			r.RecordReqResult(true)
		}
	}
}