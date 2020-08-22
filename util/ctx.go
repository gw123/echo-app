package echoapp_util

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type ctxLogger struct {
	logger *logrus.Entry
	fields logrus.Fields
}

type Record struct {
	TraceId   int64  `json:"trace_id"`
	CreatedAt int64  `json:"uint_64"`
	Point     string `json:"point"`
	Place     string `json:"place"`
}

const ctxRequestIdKey = "request_id"

var (
	defaultLogger *logrus.Logger
	isDebug       = false
)

const (
	//ctxkeys
	ctxUserKey       = "&userKey{}"
	ctxUserIdKey     = "&userIdKey{}"
	ctxAddrIdKey     = "&addrIdKey{}"
	ctxComKey        = "&comKey{}"
	ctxUserRolesKey  = "&userRolesKey{}"
	ctxLoggerKey     = "&loggerKey{}"
	ctxJwtPayloadKey = "&jwtPayloadKey{}"
)

func GetCtxComId(c echo.Context) uint {
	comId, _ := strconv.Atoi(c.Param("com_id"))
	if comId == 0 {
		comId, _ = strconv.Atoi(c.QueryParam("com_id"))
	}
	return uint(comId)
}

func GetCtxClientUUID(c echo.Context) string {
	clientUuid := c.Request().Header.Get("Client_UUID")
	return clientUuid
}

// 分页时候使用 lastId 最后一个id ，limit分页大小
func GetCtxListParams(c echo.Context) (lastId uint, limit int) {
	lastID, _ := strconv.Atoi(c.QueryParam("last_id"))
	lastId = uint(lastID)
	limit, _ = strconv.Atoi(c.QueryParam("limit"))
	if limit < 2 || limit > 100 {
		limit = 10
	}
	return lastId, limit
}

func SetCtxUserId(ctx echo.Context, userId int64) {
	AddField(ctx, "user_id", strconv.FormatInt(userId, 10))
	ctx.Set(ctxUserIdKey, userId)
}

func GetCtxtUserId(ctx echo.Context) (int64, error) {
	userId, ok := ctx.Get(ctxUserIdKey).(int64)
	if !ok {
		return 0, errors.New("get ctxUserId flied")
	}
	return userId, nil
}
func SetCtxAddrId(ctx echo.Context, addrId int64) {
	AddField(ctx, "addr_id", strconv.FormatInt(addrId, 10))
	ctx.Set(ctxAddrIdKey, addrId)
}

func GetCtxtAddrId(ctx echo.Context) (int64, error) {
	addrId, ok := ctx.Get(ctxAddrIdKey).(int64)
	if !ok {
		return 0, errors.New("get ctxAddrId flied")
	}
	return addrId, nil
}

func SetCtxUser(ctx echo.Context, user *echoapp.User) {
	ctx.Set(ctxUserKey, user)
}

func GetCtxtUser(ctx echo.Context) (*echoapp.User, error) {
	user, ok := ctx.Get(ctxUserKey).(*echoapp.User)
	if !ok {
		return nil, errors.New("get ctxUser flied")
	}
	return user, nil
}

func SetCtxCompany(ctx echo.Context, company *echoapp.Company) {
	ctx.Set(ctxComKey, company)
	AddField(ctx, "com_id", strconv.Itoa(int(company.Id)))
}

func GetCtxCompany(ctx echo.Context) (*echoapp.Company, error) {
	company, ok := ctx.Get(ctxComKey).(*echoapp.Company)
	if !ok {
		return nil, errors.New("get ctxCompany flied")
	}
	return company, nil
}

func SetCtxUserRoles(ctx echo.Context, company []echoapp.Role) {
	ctx.Set(ctxComKey, company)
}

func GetCtxtUserRoles(ctx echo.Context) ([]echoapp.Role, error) {
	roles, ok := ctx.Get(ctxComKey).([]echoapp.Role)
	if !ok {
		return nil, errors.New("get userRoles flied")
	}
	return roles, nil
}

func SetCtxJwsPayload(ctx echo.Context, payload string) {
	ctx.Set(ctxJwtPayloadKey, payload)
}

func GetCtxtJwsPayload(ctx echo.Context) (string, error) {
	payload, ok := ctx.Get(ctxJwtPayloadKey).(string)
	if !ok {
		return "", errors.New("get ctxPayload flied")
	}
	return payload, nil
}

func SetDebug(flag bool) {
	isDebug = flag
}

// 为了方便创建一个默认的Logger
func DefaultLogger() *logrus.Logger {
	if defaultLogger == nil {
		defaultLogger = logrus.New()
	}
	defaultLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		FieldMap:         nil,
		CallerPrettyfier: nil,
	})
	return defaultLogger
}

// 为了方便创建一个默认的Logger
func DefaultJsonLogger() *logrus.Logger {
	if defaultLogger == nil {
		defaultLogger = logrus.New()
	}
	defaultLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableTimestamp: false,
		DataKey:          "",
		FieldMap:         nil,
		CallerPrettyfier: nil,
		PrettyPrint:      isDebug,
	})

	return defaultLogger
}

func NewDefaultEntry() *logrus.Entry {
	//return logrus.NewEntry(DefaultLogger())
	return logrus.NewEntry(DefaultJsonLogger())
}

// 添加logrus.Entry到context, 这个操作添加的logrus.Entry在后面AddFields和Extract都会使用到
func ToContext(ctx echo.Context, entry *logrus.Entry) {
	l := &ctxLogger{
		logger: entry,
		fields: logrus.Fields{},
	}
	ctx.Set(ctxLoggerKey, l)
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddFields(ctx echo.Context, fields logrus.Fields) {
	l, ok := ctx.Get(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}

//添加日志字段到日志中间件(ctx_logrus)，添加的字段会在后面调用 info，debug，error 时候输出
func AddField(ctx echo.Context, key, val string) {
	l, ok := ctx.Get(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	l.fields[key] = val
}

// 添加一个追踪规矩id 用来聚合同一次请求, 注意要用返回的contxt 替换传入的ctx
func AddRequestId(ctx echo.Context, requestId string) {
	ctx.Set(ctxRequestIdKey, requestId)
}

//导出requestId
func ExtractRequestId(ctx echo.Context) string {
	l, ok := ctx.Get(ctxRequestIdKey).(string)
	if !ok {
		return ""
	}
	return l
}

//导出ctx_logrus日志库
func ExtractEntry(ctx echo.Context) *logrus.Entry {
	l, ok := ctx.Get(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return logrus.NewEntry(logrus.New())
	}

	fields := logrus.Fields{}
	for k, v := range l.fields {
		fields[k] = v
	}

	requestId := ExtractRequestId(ctx)
	if requestId != "" {
		fields[ctxRequestIdKey] = requestId
	}
	return l.logger.WithFields(fields)
}

// 选择一个最合适的静态资源前缀 如果是内网访问直接走内网地址 ,如果是外网访问直接走线上的环境
func GetOptimalPublicHost(ctx echo.Context, asset echoapp.Asset) string {
	ip := ctx.RealIP()
	for _, ipPrefix := range asset.InnerIpPrefix {
		if strings.HasPrefix(ip, ipPrefix) {
			return asset.PublicHostInner
		}
	}
	return asset.PublicHost
}
