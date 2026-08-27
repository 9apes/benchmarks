// Package main pulls in a large module graph for cache-bench CI compile time.
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	root := &cobra.Command{
		Use:   "cache-bench-go",
		Short: "Synthetic binary for Go CI cache benchmarking",
		Run: func(cmd *cobra.Command, args []string) {
			id := uuid.NewString()
			logger.Info("bench", zap.String("id", id))

			_, _ = clientcmd.BuildConfigFromFlags("", "")
			_, _ = kubernetes.NewForConfig(nil)

			_, _ = clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})

			r := gin.Default()
			r.GET("/metrics", gin.WrapH(promhttp.Handler()))
			r.GET("/", func(c *gin.Context) {
				c.String(http.StatusOK, fmt.Sprintf("cache-bench-go %s", id))
			})
			_ = r.Run(":8080")
		},
	}
	_ = root.Execute()
}
