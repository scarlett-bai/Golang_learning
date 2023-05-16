package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/auth/dao"
	"coolcar/auth/auth/token"
	"coolcar/auth/auth/wechat"
	"coolcar/shared/server"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/namsral/flag"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8081", "address to listen")
var mongoURI = flag.String("mongo_uri", "mongodb://localhost:27017", "mongo")
var privateKeyFile = flag.String("private_key_file", "auth/private.key", "private key file")
var wechatAppID = flag.String("wechat_app_id", "<APPID>", "wechat app id")
var wechatAppSecret = flag.String("wechat_app_secret", "<APPSECRET>", "wechat ap secret")

func main() {
	flag.Parse()
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger:%v\n", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	pkFile, err := os.Open(*privateKeyFile)
	if err != nil {
		logger.Fatal("cannot open private key")
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "auth",
		Addr:   *addr,
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					// AppID:     "wx7eef9c73c4911bbf",
					AppID: *wechatAppID,
					// AppSecret: "6d5e9a9dff1f924cb47de3efe48267a8",
					AppSecret: *wechatAppSecret,
				},
				Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
				Logger:         logger,
				TokenExpire:    2 * time.Hour,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privateKey),
			})
		},
	}))

}
