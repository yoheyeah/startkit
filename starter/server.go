package starter

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"startkit/library/gins"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/acme/autocert"
	validator "gopkg.in/go-playground/validator.v8"
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
	StaticHTMLDomain       string
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
	JWTExpireAfterInMin    int
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
	}
	m.ServerExternalIP = ip
	for {
		if !checkPortAvailable(strconv.Itoa(m.Port)) {
			m.Port++
		} else {
			break
		}
	}
	m.RequestTimeout = m.RequestTimeout * time.Second
	time.Local = time.UTC
	if local, err := time.LoadLocation(m.TimeZone); err == nil {
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
		// gins.GinErrors(),
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

func (m *Server) ParseJWT(c *gin.Context, key string) (bool, string, *gins.Claim) {
	var (
		err        error
		reqBody    []byte
		headerSubs []string
		session    = sessions.Default(c)
		value, ok  = session.Get(key).(string)
		header     = c.Request.Header.Get("Authorization")
		query      = c.Query(key)
		token      = &jwt.Token{}
		req        = struct {
			JWTToken *string `json:"jwt_token"`
		}{}
	)
	if value == "" || !ok {
		switch true {
		case query != "":
			value = query
		case header != "":
			headerSubs = strings.SplitN(header, " ", 2)
			if !(len(headerSubs) == 2 && headerSubs[0] == "Bearer") {
				return false, "", nil
			}
			value = headerSubs[1]
		default:
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
			if err = c.BindJSON(&req); err != nil {
				_, ok := err.(validator.ValidationErrors)
				if !ok {
					return false, "", nil
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
				return false, "", nil
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
			if req.JWTToken != nil {
				if *req.JWTToken != "" {
					value = *req.JWTToken
				}
			}
		}
	}
	if value == "" {
		return false, "", nil
	}
	if token, err = jwt.ParseWithClaims(value, &gins.Claim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unofficial Sign Method for Token")
		}
		return []byte(m.JWTSignedString), nil
	}); err != nil {
		return false, "", nil
	} else if !token.Valid {
		return false, "", nil
	} else if token == nil {
		return false, "", nil
	} else if claim, ok := token.Claims.(*gins.Claim); ok && token.Valid {
		if err != nil {
			return false, "", nil
		}
		if claim == nil {
			return false, "", nil
		}
	}
	return true, value, token.Claims.(*gins.Claim)
}

type AuthResult struct {
	IsSuccessful bool                   `json:"is_successful"`
	Role         string                 `json:"role"`
	User         map[string]interface{} `json:"user"`
}

type AuthReq struct {
	JWT    string `json:"jwt_token"`
	Action string `json:"action"`
	Token  string `json:"token"`
}

func encode(ddat []byte) []byte {
	return ddat
}

func decode(edat []byte) ([]byte, error) {
	return edat, nil
}

func post(url string, content interface{}, res *AuthResult) error {
	dat, err := json.Marshal(content)
	if err != nil {
		return err
	}
	dat = encode(dat)
	var cli = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 30 * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := cli.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(dat))
	if err == nil {
		rdat, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = json.Unmarshal(rdat, res)
	}
	return err
}

func validateWithURL(url, jwt string) *AuthResult {
	var (
		req = AuthReq{
			Action: "validation",
			Token:  jwt,
		}
		res = AuthResult{}
		err = post(url, req, &res)
	)
	if err != nil {
		return &AuthResult{}
	}
	return &res
}

func (m *Server) GetUserRecord(c *gin.Context, url, key, token string) (*AuthResult, error) {
	type Resp struct {
		Message interface{} `json:"message"`
		Data    interface{} `json:"data"`
	}
	result := validateWithURL(url, token)
	if !result.IsSuccessful {
		return nil, errors.New("JWT Token Not Valid")
	}
	return result, nil
}

func (m *Server) AuthServiceVarification(url, key string, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Resp struct {
			Message interface{} `json:"message"`
			Data    interface{} `json:"data"`
		}
		succeed, token, _ := m.ParseJWT(c, key)
		if !succeed {
			c.AbortWithStatusJSON(
				http.StatusBadRequest, gin.H{
					"error_message": Resp{
						Message: "Failed",
						Data:    "Authorization Parameter Is Not Valid",
					},
				})
			return
		}
		result := validateWithURL(url, token)
		if !result.IsSuccessful {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{
					"error_message": Resp{
						Message: "Failed",
						Data:    "Authorization Failed",
					},
				})
			return
		}
		for _, v := range roles {
			if result.Role == v {
				return
			}
		}
		c.AbortWithStatusJSON(
			http.StatusUnauthorized, gin.H{
				"error_message": Resp{
					Message: "Failed",
					Data:    "Authorization Failed, Role Not Correct, Unexpected Role: " + result.Role,
				},
			})
		return
	}
}

func (m *Server) SessionVarification(key string, Mysql *Mysql, obj interface{}, checkers []func(obj interface{}) (error, bool)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer Mysql.Connector()()
		type Resp struct {
			Message interface{} `json:"message"`
			Data    interface{} `json:"data"`
		}
		var (
			err        error
			reqBody    []byte
			headerSubs []string
			p          = reflect.ValueOf(obj).Elem()
			session    = sessions.Default(c)
			value, ok  = session.Get(key).(string)
			header     = c.Request.Header.Get("Authorization")
			query      = c.Query(key)
			req        = struct {
				JWTToken *string `json:"jwt_token"`
			}{}
		)
		p.Set(reflect.Zero(p.Type()))
		if value == "" || !ok {
			switch true {
			case query != "":
				value = query
			case header != "":
				headerSubs = strings.SplitN(header, " ", 2)
				if !(len(headerSubs) == 2 && headerSubs[0] == "Bearer") {
					c.AbortWithStatusJSON(
						http.StatusBadRequest, gin.H{
							"error_message": Resp{
								Message: "Failed",
								Data:    "Authorization Header Is Not Valid",
							},
						})
					return
				}
				value = headerSubs[1]
			default:
				reqBody, _ = ioutil.ReadAll(c.Request.Body)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
				if err = c.BindJSON(&req); err != nil {
					_, ok := err.(validator.ValidationErrors)
					if !ok {
						c.AbortWithStatusJSON(
							http.StatusBadRequest, gin.H{
								"error_message": Resp{
									Message: "Failed",
									Data:    "Internal Server Error Reading 'jwt_token' For Authorization Check",
								},
							})
					} else {
						c.AbortWithStatusJSON(
							http.StatusBadRequest, gin.H{
								"error_message": Resp{
									Message: "Failed",
									Data:    "Invalid Request Reading 'jwt_token' For Authorization Check",
								},
							})
					}
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
					return
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
				if req.JWTToken != nil {
					if *req.JWTToken != "" {
						value = *req.JWTToken
					}
				}
			}
		}
		if value == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_message": Resp{
					Message: "No JSON Web Token",
					Data:    "Error: No JWT In Any Of Header, JSON Request Body And Session",
				},
			})
			return
		}
		if token, err := jwt.ParseWithClaims(value, &gins.Claim{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unofficial Sign Method for Token")
			}
			return []byte(m.JWTSignedString), nil
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
			if claim.FindByObject(Mysql.DB, obj); err != nil {
				if err == gorm.ErrRecordNotFound {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error_message": Resp{
							Message: "Record Not Exist",
							Data:    "Error: " + err.Error(),
						},
					})
					return
				}
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
					if err, ok := checkers[i](obj); ok || err != nil {
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

func (m *Server) Start() (err error) {
	if m.IsNoCert {
		err = m.StartNoCert()
	} else {
		go func() { m.StartNoCert() }()
		go func() { m.StartTLS() }()
		select {}
	}
	return
}

func (m *Server) Router(r Router) {
	m.Router(r)
}

func (m *Server) StartNoCert() error {
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
		return err
	}
	return nil
}

func (m *Server) StartTLS() error {
	fmt.Println("aaaaa")
	var err error
	config := tls.Config{}
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(
		m.CrtFilePath,
		m.KeyFilePath,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error = %s\n", err.Error())
		return err
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", m.Port),
		Handler:      m.Engine,
		TLSConfig:    &config,
		ReadTimeout:  m.RequestTimeout,
		WriteTimeout: m.RequestTimeout,
		IdleTimeout:  m.RequestTimeout,
	}
	fmt.Fprintf(os.Stderr, "--- Started At Port [:%d] ---\n", m.Port)

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error = %s\n", err.Error())
		return err
	}
	return nil
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

func (m *Server) RequestClient(url string, obj interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(obj); err != nil {
		return err
	}
	return nil
}
