package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.iman.uz/imandev/bnpl_contracts/genproto/bnpl_payment_service"
	"gitlab.iman.uz/imandev/common_package/pkg/server"
	"google.golang.org/grpc"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "grpc-server",
	Run: func(cmd *cobra.Command, args []string) {
		//ctx := context.Background()
		//application, terminateFn := service.NewApplication(ctx)
		//defer func() {
		//	err := terminateFn()
		//	if err != nil {
		//		panic(err)
		//	}
		//}()

		_ = bnpl_payment_service.AddCardResponse{}

		server.RunGRPCServer(func(server *grpc.Server) {

			//svc := ports.NewGrpcServer(application)
			//telegram_service.RegisterTelegramServiceServer(server, svc)
		})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
