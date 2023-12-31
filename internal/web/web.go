package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/adriancrafter/todoapp/internal/am"
	"github.com/adriancrafter/todoapp/internal/auth"
)

const (
	authPath = "/{tenantID}/auth/"
	todoPath = "/{tenantID}/"
)

type (
	Server struct {
		am.Core
		opts []am.Option
		http.Server
		authService auth.Service
		//todoService    todo.Service
		router         *am.Router
		authController auth.WebController
		//todoController todo.WebController
		fs   am.Filesystems
		i18n *am.I18N
	}
)

var (
	cfgKey = am.Key
)

func NewServer(router *am.Router, fs am.Filesystems, i18n *am.I18N, opts ...am.Option) (server *Server) {
	name := "web-server"

	return &Server{
		Core:   am.NewCore(name, opts...),
		router: router,
		fs:     fs,
		i18n:   i18n,
	}
}

func (srv *Server) Setup(ctx context.Context) (err error) {
	staticFS, _ := srv.fs.Get(am.WebFSKey)

	err = srv.authController.Setup(ctx)
	if err != nil {
		return errors.Wrap(err, "auth controller setup error")
	}

	//err = srv.todoController.Setup(ctx)
	//if err != nil {
	//	return errors.Wrap(err, "todo controller setup error")
	//}

	reqLog := am.NewReqLoggerMiddleware(srv.Log())

	srv.router.Use(am.TenantID)
	srv.router.Use(am.RequestID)
	srv.router.Use(am.RealIP)
	srv.router.Use(reqLog)
	srv.router.Use(am.Recoverer)
	srv.router.Use(srv.i18n.Middleware)
	//srv.router.Use(am.HTMLFormat)

	srv.AddFileServer(staticFS)

	return nil
}

func (srv *Server) Start(ctx context.Context) error {
	srv.Server = http.Server{
		Addr:    srv.Address(),
		Handler: srv.Router(),
	}

	group, errGrpCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		srv.Log().Infof("%s started listening at %s", srv.Name(), srv.Address())
		defer srv.Log().Errorf("%s shutdown", srv.Name())

		err := srv.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	group.Go(func() error {
		<-errGrpCtx.Done()
		srv.Log().Errorf("%s shutdown", srv.Name())

		ctx, cancel := context.WithTimeout(context.Background(), srv.ShutdownTimeout())
		defer cancel()

		if err := srv.Server.Shutdown(ctx); err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}

func (srv *Server) SetRouter(r *am.Router) {
	srv.router = r
}

func (srv *Server) Router() (router *am.Router) {
	return srv.router
}

func (srv *Server) SetAuthController(c auth.WebController) {
	srv.authController = c
	srv.addController(c)
}

//func (srv *Server) SetVenueController(c todo.WebController) {
//	srv.todoWebController = c
//}

func (srv *Server) addController(c am.Controller) {
	srv.router.PathPrefix(c.Path()).Handler(c)
}

func (srv *Server) Address() string {
	host := srv.Cfg().GetString("http.web.server.host")
	port := srv.Cfg().GetInt("http.web.server.port")
	return fmt.Sprintf("%s:%d", host, port)
}

func TenantIDFromContext(ctx context.Context) (tenantID string, ok bool) {
	if tenantID, ok := ctx.Value(am.TenantIDKey).(string); ok {
		return tenantID, true
	}
	return "", false
}

func (srv *Server) ShutdownTimeout() time.Duration {
	secs := time.Duration(srv.Cfg().GetInt(cfgKey.WebServerTimeout))
	return secs * time.Second
}

func (srv *Server) sampleHandler(text string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, text)
	})
}
