package starter

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"startkit/library/gins"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/acme/autocert"
)

var (
	errInternetConnection = errors.New("Not connected to the network")
)

type Server struct {
	Engine
	Mode                   string
	TLSCert                string
	IsNoCert               bool
	IsPerfamceCheck        bool
	Host                   string
	Port                   int
	Domain                 string
	RequestTimeout         time.Duration
	TimeFormat             string
	TimeZone               string
	StaticPath             string
	StaticHTMLPath         string
	ServerExternalIP       string
	CookieKey              string
	SessionsKey            string
	SessionExpiryInMin     int
	SessionServeRootPath   string
	SessionServeSecureMode bool
	JWTIssuer              string
	JWTSignedString        string
	CrtFilePath            string
	KeyFilePath            string
}

type Engine struct {
	*gin.Engine
	HandlersFuncs []gin.HandlerFunc
}

// TODO: Map to Domain, later regester
func (m *Server) Builder(c *Content) error {
	ip, err := getLocalExternalIP()
	if err != nil {
		return err
	} else {
		m.ServerExternalIP = ip
	}
	for {
		if !checkPortAvailable(strconv.Itoa(m.Port)) {
			m.Port++
		} else {
			break
		}
	}
	m.RequestTimeout = m.RequestTimeout * time.Second
	if local, err := time.LoadLocation(m.TimeZone); err != nil {
		time.Local = time.UTC
	} else {
		time.Local = local
	}
	m.newGinEngine(c)
	return nil
}

func (m *Server) ListenAndServe() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", m.Port),
		Handler:      m.Engine,
		ReadTimeout:  m.RequestTimeout,
		WriteTimeout: m.RequestTimeout,
	}
	go server.ListenAndServe()
}

func (m *Server) SecurityListenAndServe() {

}

func (m *Server) newGinEngine(c *Content) {
	m.Engine.Engine = gin.New()
	m.useMiddlewares()
	m.static(c)
	m.utilsServerSetting(c)
}

func (m *Server) static(c *Content) {
	if m.StaticHTMLPath != "" {
		m.Engine.Engine.LoadHTMLGlob(m.StaticHTMLPath)
	}
	if m.StaticPath != "" {
		m.Engine.Engine.Use(static.Serve("/", static.LocalFile("./"+m.StaticPath, true)))
	}
	return
}

func (m *Server) utilsServerSetting(c *Content) {
	m.perfomanceCheck()
	gin.DefaultWriter = c.Logger.HTTPMessagesFile
	gin.DefaultErrorWriter = c.Logger.HTTPMessagesFile
	gin.SetMode(c.Server.Mode)
	return
}

func (m *Server) perfomanceCheck() {
	if m.IsPerfamceCheck {
		pprof.Register(m.Engine.Engine, "debug/pprof")
	}
	return
}

func (m *Server) useMiddlewares() {
	m.Engine.Use(
		gins.CORS(),
		gin.Logger(),
		gin.Recovery(),
		gins.GinErrors(),
		m.SessionMiddleware(),
	)
	return
}

func (m *Server) SessionMiddleware() gin.HandlerFunc {
	store := cookie.NewStore([]byte(m.CookieKey))
	store.Options(sessions.Options{
		Path:   m.SessionServeRootPath,
		MaxAge: m.SessionExpiryInMin * 60,
		Secure: m.SessionServeSecureMode,
	})
	return sessions.Sessions(m.SessionsKey, store)
}

func (m *Server) IsSessionsKeysExisted(keys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			count   = 0
			session = sessions.Default(c)
			values  = make([]interface{}, len(keys))
		)
		if count = len(keys); count > 0 {
			for i := 0; i < count; i++ {
				values[i] = session.Get(keys[i])
			}
			for i := 0; i < count; i++ {
				if values[i] == nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": keys[i] + " Not Exist"})
					return
				}
			}
		} else {
			c.Next()
			return
		}
	}
}

func (m *Server) SessionVarification(key string, DB *gorm.DB, obj interface{}, checkers []func(obj interface{}) (error, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			session   = sessions.Default(c)
			value, ok = session.Get(key).(string)
		)
		type Resp struct {
			Message interface{} `json:"message,omitempty"`
			Data    interface{} `json:"data,omitempty"`
		}
		if value == "" || !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_message": Resp{
					Message: "The Session Value Associated To The Given Key Is Not Existed",
					Data:    "Key '" + key + "' Value In Session Is Not Exist In The Claim",
				},
			})
			return
		}
		if token, err := jwt.ParseWithClaims(value, &gins.Claim{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unofficial Sign Method for Token")
			}
			return m.JWTSignedString, nil
		}); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_message": Resp{
					Message: "Cannot Parse JWT With The Claims Key",
					Data:    "Error: " + err.Error(),
				},
			})
			return
		} else if !token.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_message": Resp{
					Message: "Token Is Not Validated",
					Data:    "",
				},
			})
			return
		} else if token == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_message": Resp{
					Message: "Token Is Not Existed",
					Data:    "",
				},
			})
			return
		} else if claim, ok := token.Claims.(*gins.Claim); ok && token.Valid {
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error_message": Resp{
						Message: "Token Claim Is Not Able To Get",
						Data:    "Error: " + err.Error(),
					},
				})
				return
			}
			if claim == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error_message": Resp{
						Message: "Token Claim Is Not Existed",
						Data:    "",
					},
				})
				return
			}
			if claim.CheckByObject(DB, obj); err != nil && err != gorm.ErrRecordNotFound {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error_message": Resp{
						Message: "Error Occured While Validating Check",
						Data:    "Error: " + err.Error(),
					},
				})
				return
			}
			if count := len(checkers); count > 0 {
				for i := 0; i < count; i++ {
					if err, ok := checkers[i](obj); !ok {
						c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
							"error_message": Resp{
								Message: "Fail Validating Check",
								Data:    "",
							},
						})
						return
					} else if err != nil {
						c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
							"error_message": Resp{
								Message: "Fail Validating Check",
								Data:    "Error: " + err.Error(),
							},
						})
						return
					}
				}
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"is_validated": true,
			})
		}
	}
}

func (m *Server) JWTVarification(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ()
	}
}

func checkPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func getLocalExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errInternetConnection
}

func (m *Server) Starter(c *Content) error {
	if m.IsNoCert {
		m.StartNoCert()
	}

	return nil
}

func (m *Server) Start() error {
	if m.IsNoCert {
		m.StartNoCert()
	} else {
		m.StartTLS()
	}
	return nil
}

func (m *Server) Router(r Router) {
	m.Router(r)
}

func (m *Server) StartNoCert() {
	fmt.Fprintf(os.Stderr, "--- Started At Port [:%d] ---\n", m.Port)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", m.Port),
		Handler:      m.Engine,
		ReadTimeout:  m.RequestTimeout,
		WriteTimeout: m.RequestTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error = %s\n", err.Error())
	}
}

func (m *Server) StartTLS() error {
	var err error
	config := tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(
		m.CrtFilePath,
		m.KeyFilePath,
	)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:         strconv.Itoa(m.Port),
		Handler:      m.Engine,
		TLSConfig:    &config,
		ReadTimeout:  m.RequestTimeout,
		WriteTimeout: m.RequestTimeout,
		IdleTimeout:  m.RequestTimeout,
	}
	return server.ListenAndServeTLS("", "")
}

// Run: one-line LetsEncrypt HTTPS servers
func (m *Server) Run(trustedDomains []string) error {
	return http.Serve(autocert.NewListener(trustedDomains...), m.Engine)
}

// RunWithManager: custom autocert manager
func (m *Server) RunWithManager(manager *autocert.Manager) error {
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: manager.GetCertificate},
		Handler:   m.Engine,
	}
	go http.ListenAndServe(":http", manager.HTTPHandler(nil))
	return s.ListenAndServeTLS("", "")
}

func (m *Server) AutoCrtManager(trustedDomains []string, cacheDir string) *autocert.Manager {
	if cacheDir == "" {
		cacheDir = "/var/www/.cache"
	}
	manager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(trustedDomains...),
		Cache:      autocert.DirCache(cacheDir),
	}
	return &manager
}
