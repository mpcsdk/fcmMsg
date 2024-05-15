module fcmMsg

go 1.15

require (
	firebase.google.com/go/v4 v4.14.0
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.7.0 // indirect
	github.com/gogf/gf/contrib/drivers/pgsql/v2 v2.7.0
	github.com/gogf/gf/contrib/nosql/redis/v2 v2.7.0
	github.com/gogf/gf/v2 v2.7.0
	github.com/mpcsdk/mpcCommon v0.0.0
	github.com/nats-io/nats.go v1.34.1
	google.golang.org/api v0.170.0
)

replace github.com/mpcsdk/mpcCommon v0.0.0 => ./mpcCommon
