// Package server implements the API server for Bytebase.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	// embed will embeds the acl policy.
	_ "embed"

	"github.com/blang/semver/v4"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	scas "github.com/qiangmzsx/string-adapter/v2"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"

	"github.com/bytebase/bytebase/backend/api/auth"
	v1 "github.com/bytebase/bytebase/backend/api/v1"
	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/component/activity"
	"github.com/bytebase/bytebase/backend/component/config"
	"github.com/bytebase/bytebase/backend/component/dbfactory"
	"github.com/bytebase/bytebase/backend/component/state"
	enterpriseAPI "github.com/bytebase/bytebase/backend/enterprise/api"
	enterpriseService "github.com/bytebase/bytebase/backend/enterprise/service"
	api "github.com/bytebase/bytebase/backend/legacyapi"
	"github.com/bytebase/bytebase/backend/metric"
	metricCollector "github.com/bytebase/bytebase/backend/metric/collector"
	"github.com/bytebase/bytebase/backend/migrator"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	advisorDb "github.com/bytebase/bytebase/backend/plugin/advisor/db"
	"github.com/bytebase/bytebase/backend/plugin/app/feishu"
	"github.com/bytebase/bytebase/backend/plugin/db"
	metricPlugin "github.com/bytebase/bytebase/backend/plugin/metric"
	bbs3 "github.com/bytebase/bytebase/backend/plugin/storage/s3"
	"github.com/bytebase/bytebase/backend/resources/mongoutil"
	"github.com/bytebase/bytebase/backend/resources/mysqlutil"
	"github.com/bytebase/bytebase/backend/resources/postgres"
	"github.com/bytebase/bytebase/backend/runner/anomaly"
	"github.com/bytebase/bytebase/backend/runner/approval"
	"github.com/bytebase/bytebase/backend/runner/apprun"
	"github.com/bytebase/bytebase/backend/runner/backuprun"
	"github.com/bytebase/bytebase/backend/runner/mail"
	"github.com/bytebase/bytebase/backend/runner/metricreport"
	"github.com/bytebase/bytebase/backend/runner/rollbackrun"
	"github.com/bytebase/bytebase/backend/runner/schemasync"
	"github.com/bytebase/bytebase/backend/runner/slowquerysync"
	"github.com/bytebase/bytebase/backend/runner/taskcheck"
	"github.com/bytebase/bytebase/backend/runner/taskrun"
	"github.com/bytebase/bytebase/backend/store"
	_ "github.com/bytebase/bytebase/docs/openapi" // initial the swagger doc
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"

	// Register clickhouse driver.
	"github.com/bytebase/bytebase/backend/plugin/db/clickhouse"
	"github.com/bytebase/bytebase/backend/plugin/db/util"

	// Register mysql driver.
	"github.com/bytebase/bytebase/backend/plugin/db/mysql"
	// Register postgres driver.
	"github.com/bytebase/bytebase/backend/plugin/db/pg"
	// Register snowflake driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/snowflake"
	// Register sqlite driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/sqlite"
	// Register mongodb driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/mongodb"
	// Register spanner driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/spanner"
	// Register redis driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/redis"
	// Register oracle driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/oracle"
	// Register mssql driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/mssql"
	// Register redshift driver.
	_ "github.com/bytebase/bytebase/backend/plugin/db/redshift"
	// Register pingcap parser driver.
	_ "github.com/pingcap/tidb/types/parser_driver"
	// Register fake advisor.
	_ "github.com/bytebase/bytebase/backend/plugin/advisor/fake"
	// Register mysql advisor.
	_ "github.com/bytebase/bytebase/backend/plugin/advisor/mysql"
	// Register postgresql advisor.
	_ "github.com/bytebase/bytebase/backend/plugin/advisor/pg"

	// Register mysql differ driver.
	_ "github.com/bytebase/bytebase/backend/plugin/parser/sql/differ/mysql"
	// Register postgres differ driver.
	_ "github.com/bytebase/bytebase/backend/plugin/parser/sql/differ/pg"
	// Register mysql edit driver.
	_ "github.com/bytebase/bytebase/backend/plugin/parser/sql/edit/mysql"
	// Register postgres edit driver.
	_ "github.com/bytebase/bytebase/backend/plugin/parser/sql/edit/pg"
	// Register postgres parser driver.
	_ "github.com/bytebase/bytebase/backend/plugin/parser/sql/engine/pg"
	// Register mysql transform driver.
	_ "github.com/bytebase/bytebase/backend/plugin/parser/sql/transform/mysql"
)

const (
	// internalAPIPrefix is the API prefix for Bytebase internal, used by the UX.
	internalAPIPrefix = "/api"
	// webhookAPIPrefix is the API prefix for Bytebase webhook.
	webhookAPIPrefix = "/hook"
	// openAPIPrefix is the API prefix for Bytebase OpenAPI.
	openAPIPrefix = "/v1"
)

// Server is the Bytebase server.
type Server struct {
	// Asynchronous runners.
	TaskScheduler      *taskrun.Scheduler
	TaskCheckScheduler *taskcheck.Scheduler
	MetricReporter     *metricreport.Reporter
	SchemaSyncer       *schemasync.Syncer
	SlowQuerySyncer    *slowquerysync.Syncer
	MailSender         *mail.SlowQueryWeeklyMailSender
	BackupRunner       *backuprun.Runner
	AnomalyScanner     *anomaly.Scanner
	ApplicationRunner  *apprun.Runner
	RollbackRunner     *rollbackrun.Runner
	ApprovalRunner     *approval.Runner
	runnerWG           sync.WaitGroup

	ActivityManager *activity.Manager

	licenseService enterpriseAPI.LicenseService

	// SchemaVersion is the bytebase's schema version
	SchemaVersion *semver.Version

	profile         config.Profile
	e               *echo.Echo
	grpcServer      *grpc.Server
	metaDB          *store.MetadataDB
	store           *store.Store
	dbFactory       *dbfactory.DBFactory
	startedTs       int64
	secret          string
	errorRecordRing api.ErrorRecordRing

	// MySQL utility binaries
	mysqlBinDir string
	// MongoDB utility binaries
	mongoBinDir string
	// Postgres utility binaries
	pgBinDir string

	s3Client       *bbs3.Client
	feishuProvider *feishu.Provider

	// stateCfg is the shared in-momory state within the server.
	stateCfg *state.State

	// boot specifies that whether the server boot correctly
	cancel context.CancelFunc
}

//go:embed acl_casbin_model.conf
var casbinModel string

//go:embed acl_casbin_policy_owner.csv
var casbinOwnerPolicy string

//go:embed acl_casbin_policy_dba.csv
var casbinDBAPolicy string

//go:embed acl_casbin_policy_developer.csv
var casbinDeveloperPolicy string

// Use following cmd to generate swagger doc
// swag init -g ./backend/server.go -d ./backend/server --output docs/openapi --parseDependency

// @title Bytebase OpenAPI
// @version 1.0
// @description The OpenAPI for bytebase.
// @termsOfService https://www.bytebase.com/terms

// @contact.name API Support
// @contact.url https://github.com/bytebase/bytebase/
// @contact.email support@bytebase.com

// @license.name MIT
// @license.url https://github.com/bytebase/bytebase/blob/main/LICENSE

// @host localhost:8080
// @BasePath /v1/
// @schemes http

// NewServer creates a server.
func NewServer(ctx context.Context, profile config.Profile) (*Server, error) {
	s := &Server{
		profile:         profile,
		startedTs:       time.Now().Unix(),
		errorRecordRing: api.NewErrorRecordRing(),
	}

	// Display config
	log.Info("-----Config BEGIN-----")
	log.Info(fmt.Sprintf("mode=%s", profile.Mode))
	log.Info(fmt.Sprintf("dataDir=%s", profile.DataDir))
	log.Info(fmt.Sprintf("resourceDir=%s", profile.ResourceDir))
	log.Info(fmt.Sprintf("readonly=%t", profile.Readonly))
	log.Info(fmt.Sprintf("debug=%t", profile.Debug))
	log.Info(fmt.Sprintf("demoName=%s", profile.DemoName))
	log.Info(fmt.Sprintf("backupStorageBackend=%s", profile.BackupStorageBackend))
	log.Info(fmt.Sprintf("backupBucket=%s", profile.BackupBucket))
	log.Info(fmt.Sprintf("backupRegion=%s", profile.BackupRegion))
	log.Info(fmt.Sprintf("backupCredentialFile=%s", profile.BackupCredentialFile))
	log.Info("-----Config END-------")

	serverStarted := false
	defer func() {
		if !serverStarted {
			_ = s.Shutdown(ctx)
		}
	}()

	var err error
	// Install mysqlutil
	s.mysqlBinDir, err = mysqlutil.Install(profile.ResourceDir)
	if err != nil {
		return nil, errors.Wrap(err, "cannot install mysql utility binaries")
	}

	s.mongoBinDir, err = mongoutil.Install(profile.ResourceDir)
	if err != nil {
		return nil, errors.Wrap(err, "cannot install mongo utility binaries")
	}

	// Installs the Postgres and utility binaries and creates the 'activeProfile.pgUser' user/database
	// to store Bytebase's own metadata.
	s.pgBinDir, err = postgres.Install(profile.ResourceDir)
	if err != nil {
		return nil, err
	}

	// Start a Postgres sample server. This is used for onboarding users without requiring them to
	// configure an external instance.
	log.Info("-----Sample Postgres Instance BEGIN-----")
	sampleDataDir := common.GetPostgresSampleDataDir(profile.DataDir)
	log.Info(fmt.Sprintf("sampleDatabasePort=%d", profile.SampleDatabasePort))
	log.Info(fmt.Sprintf("sampleDataDir=%s", sampleDataDir))
	if err := postgres.StartSampleInstance(ctx, s.pgBinDir, sampleDataDir, profile.SampleDatabasePort, profile.Mode); err != nil {
		return nil, err
	}
	log.Info("-----Sample Postgres Instance END-----")

	// New MetadataDB instance.
	if profile.UseEmbedDB() {
		pgDataDir := common.GetPostgresDataDir(profile.DataDir, profile.DemoName)
		log.Info("-----Embedded Postgres Config BEGIN-----")
		log.Info(fmt.Sprintf("datastorePort=%d", profile.DatastorePort))
		log.Info(fmt.Sprintf("pgDataDir=%s", pgDataDir))
		log.Info("-----Embedded Postgres Config END-----")
		if err := postgres.InitDB(s.pgBinDir, pgDataDir, profile.PgUser); err != nil {
			return nil, err
		}
		s.metaDB = store.NewMetadataDBWithEmbedPg(profile.PgUser, pgDataDir, s.pgBinDir, profile.DemoName, profile.Mode)
	} else {
		s.metaDB = store.NewMetadataDBWithExternalPg(profile.PgURL, s.pgBinDir, profile.DemoName, profile.Mode)
	}

	// Connect to the instance that stores bytebase's own metadata.
	storeDB, err := s.metaDB.Connect(profile.DatastorePort, profile.Readonly, profile.Version)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect metadb")
	}

	if err := storeDB.Open(ctx); err != nil {
		// return s so that caller can call s.Close() to shut down the postgres server if embedded.
		return nil, errors.Wrap(err, "cannot open metadb")
	}
	storeInstance := store.New(storeDB)
	if profile.Readonly {
		log.Info("Database is opened in readonly mode. Skip migration and demo data setup.")
	} else {
		metadataVersion, err := migrator.MigrateSchema(ctx, storeDB, !profile.UseEmbedDB(), s.pgBinDir, profile.DemoName, profile.Version, profile.Mode)
		if err != nil {
			return nil, err
		}
		s.SchemaVersion = metadataVersion
	}

	s.stateCfg = &state.State{
		InstanceDatabaseSyncChan:       make(chan *api.Instance, 100),
		InstanceSlowQuerySyncChan:      make(chan *api.Instance, 100),
		InstanceOutstandingConnections: make(map[int]int),
	}
	s.store = storeInstance
	s.licenseService, err = enterpriseService.NewLicenseService(profile.Mode, storeInstance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create license service")
	}
	// Cache the license.
	s.licenseService.LoadSubscription(ctx)

	config, err := s.getInitSetting(ctx, storeInstance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}
	s.secret = config.secret

	s.ActivityManager = activity.NewManager(storeInstance)
	s.dbFactory = dbfactory.New(s.mysqlBinDir, s.mongoBinDir, s.pgBinDir, profile.DataDir, s.secret)
	e := echo.New()
	e.Debug = profile.Debug
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: func(c echo.Context) bool {
			// Skip grpc and webhook calls.
			return strings.HasPrefix(c.Request().URL.Path, "/bytebase.v1.") ||
				strings.HasPrefix(c.Request().URL.Path, webhookAPIPrefix) ||
				strings.HasPrefix(c.Request().URL.Path, "/api/sheet/")
		},
		Timeout: 30 * time.Second,
	}))
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 30, Burst: 60, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}))

	// Disallow to be embedded in an iFrame.
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XFrameOptions: "DENY",
	}))

	// MetricReporter middleware.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				// only update for authorized request
				id, ok := c.Get(getPrincipalIDContextKey()).(int)
				if !ok || id <= 0 {
					return
				}
				s.profile.LastActiveTs = time.Now().Unix()
				ctx := c.Request().Context()

				s.MetricReporter.Report(ctx, &metricPlugin.Metric{
					Name:  metric.APIRequestMetricName,
					Value: 1,
					Labels: map[string]any{
						"path":   c.Request().URL.Path,
						"method": c.Request().Method,
					},
				})
			}()
			return next(c)
		}
	})

	embedFrontend(e)
	s.e = e

	if profile.BackupBucket != "" {
		credentials, err := bbs3.GetCredentialsFromFile(ctx, profile.BackupCredentialFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get credentials from file")
		}
		s3Client, err := bbs3.NewClient(ctx, profile.BackupRegion, profile.BackupBucket, credentials)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create AWS S3 client")
		}
		s.s3Client = s3Client
	}

	s.MetricReporter = metricreport.NewReporter(s.store, s.licenseService, &s.profile, false)
	if !profile.Readonly {
		s.SchemaSyncer = schemasync.NewSyncer(storeInstance, s.dbFactory, s.stateCfg, profile)
		s.SlowQuerySyncer = slowquerysync.NewSyncer(storeInstance, s.dbFactory, s.stateCfg, profile)
		// TODO(p0ny): enable Feishu provider only when it is needed.
		s.feishuProvider = feishu.NewProvider(profile.FeishuAPIURL)
		s.ApplicationRunner = apprun.NewRunner(storeInstance, s.ActivityManager, s.feishuProvider, profile)
		s.BackupRunner = backuprun.NewRunner(storeInstance, s.dbFactory, s.s3Client, s.stateCfg, &profile)
		s.RollbackRunner = rollbackrun.NewRunner(storeInstance, s.dbFactory, s.stateCfg)
		s.ApprovalRunner = approval.NewRunner(storeInstance, s.dbFactory, s.stateCfg, s.ActivityManager, s.licenseService)

		s.MailSender = mail.NewSender(s.store, s.stateCfg)

		s.TaskScheduler = taskrun.NewScheduler(storeInstance, s.ApplicationRunner, s.SchemaSyncer, s.ActivityManager, s.licenseService, s.stateCfg, profile, s.MetricReporter)
		s.TaskScheduler.Register(api.TaskGeneral, taskrun.NewDefaultExecutor())
		s.TaskScheduler.Register(api.TaskDatabaseCreate, taskrun.NewDatabaseCreateExecutor(storeInstance, s.dbFactory, s.SchemaSyncer, profile))
		s.TaskScheduler.Register(api.TaskDatabaseSchemaBaseline, taskrun.NewSchemaBaselineExecutor(storeInstance, s.dbFactory, s.ActivityManager, s.licenseService, s.stateCfg, s.SchemaSyncer, profile))
		s.TaskScheduler.Register(api.TaskDatabaseSchemaUpdate, taskrun.NewSchemaUpdateExecutor(storeInstance, s.dbFactory, s.ActivityManager, s.licenseService, s.stateCfg, s.SchemaSyncer, profile))
		s.TaskScheduler.Register(api.TaskDatabaseSchemaUpdateSDL, taskrun.NewSchemaUpdateSDLExecutor(storeInstance, s.dbFactory, s.ActivityManager, s.licenseService, s.stateCfg, s.SchemaSyncer, profile))
		s.TaskScheduler.Register(api.TaskDatabaseDataUpdate, taskrun.NewDataUpdateExecutor(storeInstance, s.dbFactory, s.ActivityManager, s.licenseService, s.stateCfg, profile))
		s.TaskScheduler.Register(api.TaskDatabaseBackup, taskrun.NewDatabaseBackupExecutor(storeInstance, s.dbFactory, s.s3Client, profile))
		s.TaskScheduler.Register(api.TaskDatabaseSchemaUpdateGhostSync, taskrun.NewSchemaUpdateGhostSyncExecutor(storeInstance, s.stateCfg, s.secret))
		s.TaskScheduler.Register(api.TaskDatabaseSchemaUpdateGhostCutover, taskrun.NewSchemaUpdateGhostCutoverExecutor(storeInstance, s.dbFactory, s.ActivityManager, s.licenseService, s.stateCfg, s.SchemaSyncer, profile))
		s.TaskScheduler.Register(api.TaskDatabaseRestorePITRRestore, taskrun.NewPITRRestoreExecutor(storeInstance, s.dbFactory, s.s3Client, s.SchemaSyncer, s.stateCfg, profile))
		s.TaskScheduler.Register(api.TaskDatabaseRestorePITRCutover, taskrun.NewPITRCutoverExecutor(storeInstance, s.dbFactory, s.SchemaSyncer, s.BackupRunner, s.ActivityManager, profile))

		s.TaskCheckScheduler = taskcheck.NewScheduler(storeInstance, s.licenseService, s.stateCfg)
		statementSimpleExecutor := taskcheck.NewStatementAdvisorSimpleExecutor(storeInstance)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseStatementFakeAdvise, statementSimpleExecutor)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseStatementSyntax, statementSimpleExecutor)
		statementCompositeExecutor := taskcheck.NewStatementAdvisorCompositeExecutor(storeInstance, s.dbFactory)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseStatementAdvise, statementCompositeExecutor)
		statementTypeExecutor := taskcheck.NewStatementTypeExecutor(storeInstance, s.dbFactory)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseStatementType, statementTypeExecutor)
		databaseConnectExecutor := taskcheck.NewDatabaseConnectExecutor(storeInstance, s.dbFactory)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseConnect, databaseConnectExecutor)
		ghostSyncExecutor := taskcheck.NewGhostSyncExecutor(storeInstance, s.secret)
		s.TaskCheckScheduler.Register(api.TaskCheckGhostSync, ghostSyncExecutor)
		pitrMySQLExecutor := taskcheck.NewPITRMySQLExecutor(storeInstance, s.dbFactory)
		s.TaskCheckScheduler.Register(api.TaskCheckPITRMySQL, pitrMySQLExecutor)
		statementTypeReportExecutor := taskcheck.NewStatementTypeReportExecutor(storeInstance)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseStatementTypeReport, statementTypeReportExecutor)
		statementAffectedRowsExecutor := taskcheck.NewStatementAffectedRowsReportExecutor(storeInstance, s.dbFactory)
		s.TaskCheckScheduler.Register(api.TaskCheckDatabaseStatementAffectedRowsReport, statementAffectedRowsExecutor)

		// Anomaly scanner
		s.AnomalyScanner = anomaly.NewScanner(storeInstance, s.dbFactory, s.licenseService)

		// Metric reporter
		s.initMetricReporter()
	}

	// Middleware
	//
	// Error recorder middleware.
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		errorRecorderMiddleware(err, s, c, e)
	}
	// API logger middleware.
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			if s.profile.Mode == common.ReleaseModeProd && !s.profile.Debug {
				return true
			}
			return !common.HasPrefixes(c.Path(), internalAPIPrefix, openAPIPrefix, webhookAPIPrefix)
		},
		Format: `{"time":"${time_rfc3339}",` +
			`"method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}"}` + "\n",
	}))
	// Panic recovery middleware.
	e.Use(recoverMiddleware)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	webhookGroup := e.Group(webhookAPIPrefix)
	s.registerWebhookRoutes(webhookGroup)

	apiGroup := e.Group(internalAPIPrefix)
	// API JWT authentication middleware.
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return JWTMiddleware(internalAPIPrefix, s.store, next, profile.Mode, config.secret)
	})

	m, err := model.NewModelFromString(casbinModel)
	if err != nil {
		return nil, err
	}
	sa := scas.NewAdapter(strings.Join([]string{casbinOwnerPolicy, casbinDBAPolicy, casbinDeveloperPolicy}, "\n"))
	ce, err := casbin.NewEnforcer(m, sa)
	if err != nil {
		return nil, err
	}
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return aclMiddleware(s, internalAPIPrefix, ce, next, profile.Readonly)
	})
	s.registerDebugRoutes(apiGroup)
	s.registerSettingRoutes(apiGroup)
	s.registerOAuthRoutes(apiGroup)
	s.registerPrincipalRoutes(apiGroup)
	s.registerMemberRoutes(apiGroup)
	s.registerPolicyRoutes(apiGroup)
	s.registerProjectRoutes(apiGroup)
	s.registerProjectWebhookRoutes(apiGroup)
	s.registerEnvironmentRoutes(apiGroup)
	s.registerInstanceRoutes(apiGroup)
	s.registerDatabaseRoutes(apiGroup)
	s.registerIssueRoutes(apiGroup)
	s.registerIssueSubscriberRoutes(apiGroup)
	s.registerTaskRoutes(apiGroup)
	s.registerStageRoutes(apiGroup)
	s.registerActivityRoutes(apiGroup)
	s.registerInboxRoutes(apiGroup)
	s.registerBookmarkRoutes(apiGroup)
	s.registerSQLRoutes(apiGroup)
	s.registerVCSRoutes(apiGroup)
	s.registerPlanRoutes(apiGroup)
	s.registerSheetRoutes(apiGroup)
	s.registerSheetOrganizerRoutes(apiGroup)
	s.registerAnomalyRoutes(apiGroup)

	// Register healthz endpoint.
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK!\n")
	})

	// Setup the gRPC and grpc-gateway.
	authProvider := auth.New(s.store, s.secret, s.licenseService, profile.Mode)
	aclProvider := v1.NewACLInterceptor(s.store, s.secret, s.licenseService, profile.Mode)
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(authProvider.AuthenticationInterceptor, aclProvider.ACLInterceptor),
	)
	v1pb.RegisterAuthServiceServer(s.grpcServer, v1.NewAuthService(s.store, s.secret, s.licenseService, s.MetricReporter, &profile,
		func(ctx context.Context, user *store.UserMessage, firstEndUser bool) error {
			if s.profile.TestOnlySkipOnboardingData {
				return nil
			}
			// Only generate onboarding data after the first enduser signup.
			if firstEndUser {
				if err := s.generateOnboardingData(ctx, user.ID); err != nil {
					return status.Errorf(codes.Internal, "failed to prepare onboarding data, error: %v", err)
				}
			}
			return nil
		}))
	v1pb.RegisterActuatorServiceServer(s.grpcServer, v1.NewActuatorService(s.store, &s.profile))
	v1pb.RegisterSubscriptionServiceServer(s.grpcServer, v1.NewSubscriptionService(
		s.store,
		&s.profile,
		s.MetricReporter,
		s.licenseService))
	v1pb.RegisterEnvironmentServiceServer(s.grpcServer, v1.NewEnvironmentService(s.store, s.licenseService))
	v1pb.RegisterInstanceServiceServer(s.grpcServer, v1.NewInstanceService(
		s.store,
		s.licenseService,
		s.MetricReporter,
		s.secret,
		s.stateCfg,
		s.dbFactory))
	v1pb.RegisterProjectServiceServer(s.grpcServer, v1.NewProjectService(s.store))
	v1pb.RegisterDatabaseServiceServer(s.grpcServer, v1.NewDatabaseService(s.store, s.BackupRunner, s.licenseService))
	v1pb.RegisterInstanceRoleServiceServer(s.grpcServer, v1.NewInstanceRoleService(s.store, s.dbFactory))
	v1pb.RegisterOrgPolicyServiceServer(s.grpcServer, v1.NewOrgPolicyService(s.store, s.licenseService))
	v1pb.RegisterIdentityProviderServiceServer(s.grpcServer, v1.NewIdentityProviderService(s.store, s.licenseService))
	v1pb.RegisterSettingServiceServer(s.grpcServer, v1.NewSettingService(s.store, &s.profile, s.licenseService, s.stateCfg))
	v1pb.RegisterAnomalyServiceServer(s.grpcServer, v1.NewAnomalyService(s.store))
	v1pb.RegisterSQLServiceServer(s.grpcServer, v1.NewSQLService())
	v1pb.RegisterExternalVersionControlServiceServer(s.grpcServer, v1.NewExternalVersionControlService(s.store))
	v1pb.RegisterRiskServiceServer(s.grpcServer, v1.NewRiskService(s.store, s.licenseService))
	v1pb.RegisterReviewServiceServer(s.grpcServer, v1.NewReviewService(s.store, s.ActivityManager, s.stateCfg))
	v1pb.RegisterRoleServiceServer(s.grpcServer, v1.NewRoleService(s.store, s.licenseService))
	reflection.Register(s.grpcServer)

	// REST gateway proxy.
	grpcEndpoint := fmt.Sprintf(":%d", profile.GrpcPort)
	grpcConn, err := grpc.Dial(grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	mux := runtime.NewServeMux(runtime.WithForwardResponseOption(auth.GatewayResponseModifier))
	if err := v1pb.RegisterAuthServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterActuatorServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterSubscriptionServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterEnvironmentServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterInstanceServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterProjectServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterDatabaseServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterInstanceRoleServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterOrgPolicyServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterIdentityProviderServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterSettingServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterAnomalyServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterSQLServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterExternalVersionControlServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	if err := v1pb.RegisterRoleServiceHandler(ctx, mux, grpcConn); err != nil {
		return nil, err
	}
	e.Any("/v1/*", echo.WrapHandler(mux))
	// GRPC web proxy.
	options := []grpcweb.Option{
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	}
	wrappedGrpc := grpcweb.WrapServer(s.grpcServer, options...)
	e.Any("/bytebase.v1.*", echo.WrapHandler(wrappedGrpc))

	// Register open API routes
	s.registerOpenAPIRoutes(e, ce, profile)

	// Register pprof endpoints.
	pprof.Register(e)
	// Register prometheus metrics endpoint.
	p := prometheus.NewPrometheus("api", nil)
	p.Use(e)

	// TODO: remove this on May 2023
	go s.backfillInstanceChangeHistory(ctx)

	serverStarted = true
	return s, nil
}

func (s *Server) registerOpenAPIRoutes(e *echo.Echo, ce *casbin.Enforcer, prof config.Profile) {
	jwtMiddlewareFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
		return JWTMiddleware(openAPIPrefix, s.store, next, prof.Mode, s.secret)
	}
	aclMiddlewareFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
		return aclMiddleware(s, openAPIPrefix, ce, next, prof.Readonly)
	}
	metricMiddlewareFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
		return openAPIMetricMiddleware(s, next)
	}
	e.POST("/v1/sql/advise", s.sqlCheckController)
	e.POST("/v1/sql/schema/diff", schemaDiff)
	e.PATCH("/v1/instances/:instanceName/databases/:database", s.updateInstanceDatabase, jwtMiddlewareFunc, aclMiddlewareFunc, metricMiddlewareFunc)
	e.POST("/v1/issues", s.createIssueByOpenAPI, jwtMiddlewareFunc, aclMiddlewareFunc, metricMiddlewareFunc)
}

// initMetricReporter will initial the metric scheduler.
func (s *Server) initMetricReporter() {
	metricReporter := metricreport.NewReporter(s.store, s.licenseService, &s.profile, s.profile.EnableMetric)
	metricReporter.Register(metric.InstanceCountMetricName, metricCollector.NewInstanceCountCollector(s.store))
	metricReporter.Register(metric.IssueCountMetricName, metricCollector.NewIssueCountCollector(s.store))
	metricReporter.Register(metric.ProjectCountMetricName, metricCollector.NewProjectCountCollector(s.store))
	metricReporter.Register(metric.PolicyCountMetricName, metricCollector.NewPolicyCountCollector(s.store))
	metricReporter.Register(metric.TaskCountMetricName, metricCollector.NewTaskCountCollector(s.store))
	metricReporter.Register(metric.DatabaseCountMetricName, metricCollector.NewDatabaseCountCollector(s.store))
	metricReporter.Register(metric.SheetCountMetricName, metricCollector.NewSheetCountCollector(s.store))
	metricReporter.Register(metric.MemberCountMetricName, metricCollector.NewMemberCountCollector(s.store))
	s.MetricReporter = metricReporter
}

// retrieved via the SettingService upon startup.
type workspaceConfig struct {
	// secret used to sign the JWT auth token
	secret string
	// workspaceID used to initial the identify for a new workspace.
	workspaceID string
}

func (s *Server) getInitSetting(ctx context.Context, datastore *store.Store) (*workspaceConfig, error) {
	// secretLength is the length for the secret used to sign the JWT auto token.
	const secretLength = 32

	// initial branding
	if _, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingBrandingLogo,
		Value:       "",
		Description: "The branding logo image in base64 string format.",
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	conf := &workspaceConfig{}

	// initial JWT token
	value, err := common.RandomString(secretLength)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random JWT secret")
	}
	authSetting, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingAuthSecret,
		Value:       value,
		Description: "Random string used to sign the JWT auth token.",
	}, api.SystemBotID)
	if err != nil {
		return nil, err
	}
	conf.secret = authSetting.Value

	// initial workspace
	workspaceSetting, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingWorkspaceID,
		Value:       uuid.New().String(),
		Description: "The workspace identifier",
	}, api.SystemBotID)
	if err != nil {
		return nil, err
	}
	conf.workspaceID = workspaceSetting.Value

	// initial license
	if _, _, err = datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingEnterpriseLicense,
		Value:       "",
		Description: "Enterprise license",
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	// initial feishu app
	if _, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingAppIM,
		Value:       "",
		Description: "",
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	// initial watermark setting
	if _, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingWatermark,
		Value:       "0",
		Description: "Display watermark",
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	// initial OpenAI key setting
	if _, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingPluginOpenAIKey,
		Value:       "",
		Description: "API key to request OpenAI (ChatGPT)",
	}, api.SystemBotID); err != nil {
		return nil, err
	}
	if _, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name:        api.SettingPluginOpenAIEndpoint,
		Value:       "",
		Description: "API Endpoint for OpenAI",
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	// initial workspace approval setting
	approvalSettingValue, err := protojson.Marshal(&storepb.WorkspaceApprovalSetting{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal initial workspace approval setting")
	}
	if _, _, err := datastore.CreateSettingIfNotExistV2(ctx, &store.SettingMessage{
		Name: api.SettingWorkspaceApproval,
		// Value is ""
		Value:       string(approvalSettingValue),
		Description: "The workspace approval setting",
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	// initial workspace profile setting
	settingName := api.SettingWorkspaceProfile
	workspaceProfileSetting, err := s.store.GetSettingV2(ctx, &store.FindSettingMessage{
		Name:    &settingName,
		Enforce: true,
	})
	if err != nil {
		return nil, err
	}

	workspaceProfilePayload := &storepb.WorkspaceProfileSetting{
		ExternalUrl: s.profile.ExternalURL,
	}
	if workspaceProfileSetting != nil {
		workspaceProfilePayload = new(storepb.WorkspaceProfileSetting)
		if err := protojson.Unmarshal([]byte(workspaceProfileSetting.Value), workspaceProfilePayload); err != nil {
			return nil, err
		}
		if s.profile.ExternalURL != "" {
			workspaceProfilePayload.ExternalUrl = s.profile.ExternalURL
		}
	}

	bytes, err := protojson.Marshal(workspaceProfilePayload)
	if err != nil {
		return nil, err
	}

	if _, err := datastore.UpsertSettingV2(ctx, &store.SetSettingMessage{
		Name:  api.SettingWorkspaceProfile,
		Value: string(bytes),
	}, api.SystemBotID); err != nil {
		return nil, err
	}

	return conf, nil
}

// Run will run the server.
func (s *Server) Run(ctx context.Context, port int) error {
	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	if !s.profile.Readonly {
		if err := s.TaskScheduler.ClearRunningTasks(ctx); err != nil {
			return errors.Wrap(err, "failed to clear existing RUNNING tasks before starting the task scheduler")
		}
		// runnerWG waits for all goroutines to complete.
		s.runnerWG.Add(1)
		go s.TaskScheduler.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.TaskCheckScheduler.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.SchemaSyncer.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.SlowQuerySyncer.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.MailSender.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.BackupRunner.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.AnomalyScanner.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.ApplicationRunner.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.RollbackRunner.Run(ctx, &s.runnerWG)
		s.runnerWG.Add(1)
		go s.ApprovalRunner.Run(ctx, &s.runnerWG)

		s.runnerWG.Add(1)
		go s.MetricReporter.Run(ctx, &s.runnerWG)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port+1))
	if err != nil {
		return err
	}
	go func() {
		if err := s.grpcServer.Serve(listen); err != nil {
			log.Error("grpc server listen error", zap.Error(err))
		}
	}()
	return s.e.Start(fmt.Sprintf(":%d", port))
}

// Shutdown will shut down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	log.Info("Stopping Bytebase...")
	log.Info("Stopping web server...")

	// Close the metric reporter
	if s.MetricReporter != nil {
		s.MetricReporter.Close()
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Cancel the worker
	if s.cancel != nil {
		s.cancel()
	}

	// Shutdown echo
	if s.e != nil {
		if err := s.e.Shutdown(ctx); err != nil {
			s.e.Logger.Fatal(err)
		}
	}
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	// Wait for all runners to exit.
	s.runnerWG.Wait()

	// Close db connection
	if s.store != nil {
		if err := s.store.Close(ctx); err != nil {
			return err
		}
	}

	// Shutdown postgres sample instance.
	if err := postgres.Stop(s.pgBinDir, common.GetPostgresSampleDataDir(s.profile.DataDir)); err != nil {
		log.Error("Failed to stop postgres sample instance", zap.Error(err))
	}

	// Shutdown postgres server if embed.
	if s.metaDB != nil {
		s.metaDB.Close()
	}
	log.Info("Bytebase stopped properly")

	return nil
}

// GetEcho returns the echo server.
func (s *Server) GetEcho() *echo.Echo {
	return s.e
}

// getSampleSQLReviewPolicy returns a sample SQL review policy for preparing onboardign data.
func getSampleSQLReviewPolicy() *advisor.SQLReviewPolicy {
	policy := &advisor.SQLReviewPolicy{
		Name: "SQL Review Sample Policy",
	}

	ruleList := []*advisor.SQLReviewRule{}

	// Add DropEmptyDatabase rule for MySQL, TiDB, MariaDB.
	for _, e := range []advisorDb.Type{advisorDb.MySQL, advisorDb.TiDB, advisorDb.MariaDB} {
		ruleList = append(ruleList, &advisor.SQLReviewRule{
			Type:    advisor.SchemaRuleDropEmptyDatabase,
			Level:   advisor.SchemaRuleLevelError,
			Engine:  e,
			Payload: "{}",
		})
	}

	// Add ColumnNotNull rule for MySQL, TiDB, MariaDB, Postgres.
	for _, e := range []advisorDb.Type{advisorDb.MySQL, advisorDb.TiDB, advisorDb.MariaDB, advisorDb.Postgres} {
		ruleList = append(ruleList, &advisor.SQLReviewRule{
			Type:    advisor.SchemaRuleColumnNotNull,
			Level:   advisor.SchemaRuleLevelWarning,
			Engine:  e,
			Payload: "{}",
		})
	}

	// Add TableDropNamingConvention rule for MySQL, TiDB, MariaDB Postgres.
	for _, e := range []advisorDb.Type{advisorDb.MySQL, advisorDb.TiDB, advisorDb.MariaDB, advisorDb.Postgres} {
		ruleList = append(ruleList, &advisor.SQLReviewRule{
			Type:    advisor.SchemaRuleTableDropNamingConvention,
			Level:   advisor.SchemaRuleLevelError,
			Engine:  e,
			Payload: "{\"format\":\"_del$\"}",
		})
	}

	policy.RuleList = ruleList
	return policy
}

// generateOnboardingData generates onboarding data after the first signup.
func (s *Server) generateOnboardingData(ctx context.Context, userID int) error {
	project, err := s.store.CreateProjectV2(ctx, &store.ProjectMessage{
		ResourceID: "project-sample",
		Title:      "Sample Project",
		Key:        "SAM",
		Workflow:   api.UIWorkflow,
		Visibility: api.Public,
		TenantMode: api.TenantModeDisabled,
	}, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to create onboarding project")
	}

	instance, err := s.store.CreateInstanceV2(ctx, api.DefaultProdEnvironmentID, &store.InstanceMessage{
		ResourceID:   postgres.SampleInstanceResourceID,
		Title:        "Postgres Sample Instance",
		Engine:       db.Postgres,
		ExternalLink: "",
		DataSources: []*store.DataSourceMessage{
			{
				Title:              api.AdminDataSourceName,
				Type:               api.Admin,
				Username:           postgres.SampleUser,
				ObfuscatedPassword: common.Obfuscate("", s.secret),
				Host:               common.GetPostgresSocketDir(),
				Port:               strconv.Itoa(s.profile.SampleDatabasePort),
				Database:           postgres.SampleDatabase,
			},
		},
	}, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to create onboarding instance")
	}

	// Sync the instance schema so we can transfer the sample database later.
	if _, err := s.SchemaSyncer.SyncInstance(ctx, instance); err != nil {
		return errors.Wrapf(err, "failed to sync onboarding instance")
	}

	// Transfer sample database to the just created project.
	transferDatabaseMessage := &store.UpdateDatabaseMessage{
		EnvironmentID: api.DefaultProdEnvironmentID,
		InstanceID:    instance.ResourceID,
		DatabaseName:  postgres.SampleDatabase,
		ProjectID:     &project.ResourceID,
	}
	_, err = s.store.UpdateDatabase(ctx, transferDatabaseMessage, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to transfer sample database")
	}

	dbName := postgres.SampleDatabase
	envID := api.DefaultProdEnvironmentID
	database, err := s.store.GetDatabaseV2(ctx, &store.FindDatabaseMessage{
		EnvironmentID: &envID,
		InstanceID:    &instance.ResourceID,
		DatabaseName:  &dbName,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to find onboarding instance")
	}

	// Need to sync database schema so we can configure sensitive data policy and create the schema
	// update issue later.
	if err := s.SchemaSyncer.SyncDatabaseSchema(ctx, database, true /* force */); err != nil {
		return errors.Wrapf(err, "failed to sync sample database schema")
	}

	// Add a sample SQL Review policy to the prod environment. This pairs with the following schema
	// change issue to demonstrate the SQL Review feature.
	policyPayload, err := json.Marshal(*getSampleSQLReviewPolicy())
	if err != nil {
		return errors.Wrapf(err, "failed to marshal onboarding SQL Review policy")
	}

	_, err = s.store.CreatePolicyV2(ctx, &store.PolicyMessage{
		ResourceUID:       api.DefaultProdEnvironmentUID,
		ResourceType:      api.PolicyResourceTypeEnvironment,
		Payload:           string(policyPayload),
		Type:              api.PolicyTypeSQLReview,
		InheritFromParent: true,
		// Enforce cannot be false while creating a policy.
		Enforce: true,
	}, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to create onboarding SQL Review policy")
	}

	// Create a schema update issue.
	createContext, err := json.Marshal(
		&api.MigrationContext{
			DetailList: []*api.MigrationDetail{
				{
					MigrationType: db.Migrate,
					DatabaseID:    database.UID,
					// This will violate the NOT NULL SQL Review policy configured above and emit a
					// warning. Thus to demonstrate the SQL Review capability.
					Statement: "ALTER TABLE employee ADD COLUMN IF NOT EXISTS email TEXT DEFAULT '';",
				},
			},
		})
	if err != nil {
		return errors.Wrapf(err, "failed to marshal sample schema update issue context")
	}

	issueCreate := &api.IssueCreate{
		ProjectID: project.UID,
		Name:      "👉👉👉 [START HERE] Add email column to Employee table",
		Type:      api.IssueDatabaseSchemaUpdate,
		Description: `A sample issue to showcase how to review database schema change.
		
Click "Approve" button to apply the schema update.`,
		AssigneeID:            userID,
		AssigneeNeedAttention: true,
		CreateContext:         string(createContext),
	}

	// Use system bot as the creator so that the issue only appears in the user's assignee list
	issue, err := s.createIssue(ctx, issueCreate, api.SystemBotID)
	if err != nil {
		return errors.Wrapf(err, "failed to create sample issue")
	}

	// Bookmark the issue.
	if _, err := s.store.CreateBookmarkV2(ctx, &store.BookmarkMessage{
		Name: "Sample Issue",
		Link: fmt.Sprintf("/issue/%s-%d", slug.Make(issue.Name), issue.ID),
	}, userID); err != nil {
		return errors.Wrapf(err, "failed to bookmark sample issue")
	}

	// Add a sensitive data policy to pair it with the sample query below. So that user can
	// experience the sensitive data masking feature from SQL Editor.
	sensitiveDataPolicy := api.SensitiveDataPolicy{
		SensitiveDataList: []api.SensitiveData{
			{
				Schema: "public",
				Table:  "salary",
				Column: "amount",
				Type:   api.SensitiveDataMaskTypeDefault,
			},
		},
	}
	policyPayload, err = json.Marshal(sensitiveDataPolicy)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal onboarding sensitive data policy")
	}

	_, err = s.store.CreatePolicyV2(ctx, &store.PolicyMessage{
		ResourceUID:       database.UID,
		ResourceType:      api.PolicyResourceTypeDatabase,
		Payload:           string(policyPayload),
		Type:              api.PolicyTypeSensitiveData,
		InheritFromParent: true,
		// Enforce cannot be false while creating a policy.
		Enforce: true,
	}, userID)
	if err != nil {
		return errors.Wrapf(err, "failed to create onboarding sensitive data policy")
	}

	// Create a SQL sheet with sample queries.
	sheetCreate := &api.SheetCreate{
		CreatorID:  userID,
		ProjectID:  project.UID,
		DatabaseID: &database.UID,
		Name:       "Sample Sheet",
		Statement:  "SELECT * FROM salary;",
		Visibility: api.ProjectSheet,
		Source:     api.SheetFromBytebase,
		Type:       api.SheetForSQL,
	}
	_, err = s.store.CreateSheet(ctx, sheetCreate)
	if err != nil {
		return errors.Wrapf(err, "failed to create sample sheet")
	}

	return nil
}

func (s *Server) backfillInstanceChangeHistory(ctx context.Context) {
	err := func() error {
		migratedInstanceList, err := s.store.ListInstanceHavingInstanceChangeHistory(ctx)
		if err != nil {
			return err
		}
		instanceMigrated := make(map[int]bool)
		for _, id := range migratedInstanceList {
			instanceMigrated[id] = true
		}

		instanceList, err := s.store.ListInstancesV2(ctx, &store.FindInstanceMessage{})
		if err != nil {
			return err
		}

		issueList, err := s.store.ListIssueV2(ctx, &store.FindIssueMessage{})
		if err != nil {
			return err
		}
		hasIssueID := make(map[int]bool)
		for _, issue := range issueList {
			hasIssueID[issue.UID] = true
		}

		principalList, err := s.store.GetPrincipalList(ctx)
		if err != nil {
			return err
		}
		nameToPrincipal := make(map[string]*api.Principal)
		for _, principal := range principalList {
			nameToPrincipal[principal.Name] = principal
		}

		var errList error
		for _, instance := range instanceList {
			err := func(instance *store.InstanceMessage) error {
				limit := 10
				offset := 0

				if !(instance.Engine == db.MySQL || instance.Engine == db.Postgres || instance.Engine == db.ClickHouse) {
					return nil
				}
				if instanceMigrated[instance.UID] {
					return nil
				}

				databaseList, err := s.store.ListDatabases(ctx, &store.FindDatabaseMessage{InstanceID: &instance.ResourceID, ShowDeleted: true})
				if err != nil {
					return err
				}
				nameToDatabase := make(map[string]*store.DatabaseMessage)
				for i, db := range databaseList {
					nameToDatabase[db.DatabaseName] = databaseList[i]
				}

				driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, instance, "bytebase")
				if err != nil {
					return err
				}
				defer driver.Close(ctx)

				for {
					var history []*db.MigrationHistory
					if instance.Engine == db.MySQL {
						myDriver := driver.(*mysql.Driver)
						history, err = myDriver.FindMigrationHistoryList(ctx, &db.MigrationHistoryFind{
							InstanceID: &instance.UID,
							Limit:      &limit,
							Offset:     &offset,
						})
						if err != nil {
							return err
						}
					} else if instance.Engine == db.Postgres {
						pgDriver := driver.(*pg.Driver)
						history, err = pgDriver.FindMigrationHistoryList(ctx, &db.MigrationHistoryFind{
							InstanceID: &instance.UID,
							Limit:      &limit,
							Offset:     &offset,
						})
						if err != nil {
							return err
						}
					} else if instance.Engine == db.ClickHouse {
						cDriver := driver.(*clickhouse.Driver)
						history, err = cDriver.FindMigrationHistoryList(ctx, &db.MigrationHistoryFind{
							InstanceID: &instance.UID,
							Limit:      &limit,
							Offset:     &offset,
						})
						if err != nil {
							return err
						}
					}
					if len(history) == 0 {
						break
					}
					offset += limit

					var creates []*store.InstanceChangeHistoryMessage
					for _, h := range history {
						var databaseID *int
						if database, ok := nameToDatabase[h.Namespace]; ok {
							databaseID = &database.UID
						}

						var issueID *int
						if id, err := strconv.Atoi(h.IssueID); err != nil {
							errList = multierr.Append(errList, err)
						} else if hasIssueID[id] {
							// Has FK constraint on issue_id.
							// Set to id if issue exists.
							issueID = &id
						}

						creatorID := api.SystemBotID
						updaterID := api.SystemBotID
						if principal, ok := nameToPrincipal[h.Creator]; ok {
							creatorID = principal.ID
						}
						if principal, ok := nameToPrincipal[h.Updater]; ok {
							updaterID = principal.ID
						}

						storedVersion, err := util.ToStoredVersion(h.UseSemanticVersion, h.Version, h.SemanticVersionSuffix)
						if err != nil {
							return err
						}

						changeHistory := store.InstanceChangeHistoryMessage{
							CreatorID:           creatorID,
							CreatedTs:           h.CreatedTs,
							UpdaterID:           updaterID,
							UpdatedTs:           h.UpdatedTs,
							InstanceID:          &instance.UID,
							DatabaseID:          databaseID,
							IssueID:             issueID,
							ReleaseVersion:      h.ReleaseVersion,
							Sequence:            int64(h.Sequence),
							Source:              h.Source,
							Type:                h.Type,
							Status:              h.Status,
							Version:             storedVersion,
							Description:         h.Description,
							Statement:           h.Statement,
							Schema:              h.Schema,
							SchemaPrev:          h.SchemaPrev,
							ExecutionDurationNs: h.ExecutionDurationNs,
							Payload:             h.Payload,
						}

						creates = append(creates, &changeHistory)
					}

					if _, err := s.store.CreateInstanceChangeHistory(ctx, creates...); err != nil {
						return err
					}
				}
				return nil
			}(instance)
			if err != nil {
				// New instances may not have the "bytebase" database.
				if strings.Contains(err.Error(), "database \"bytebase\" does not exist") {
					return nil
				}
				errList = multierr.Append(errList, err)
			}
		}
		return errList
	}()

	if err != nil {
		log.Warn("failed to backfill migration history", zap.Error(err))
	}
}
