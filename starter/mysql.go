package starter

import (
	"database/sql"
	"fmt"
	"log"
	"startkit/library/times"
	"startkit/utils"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/jinzhu/gorm"
)

type Mysql struct {
	MysqlInstance
	Username                       string
	Password                       string
	Host                           string
	Port                           int
	DatabaseName                   string
	IsWebService                   bool
	Error                          error
	MaximumIdleConnection          int
	MaximumOpenConnection          int
	MaximumConnectionRetry         int
	MinimumRetryDuration           int
	MaximumConnectionKeepAliveTime int
}

type MysqlInstance struct {
	DB         *gorm.DB
	db         *sql.DB
	BasicModel MysqlModel
	ModelAddrs []interface{}
}

type MysqlModel struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at; sql:type:datetime;default:'1980-01-01'"`
	UpdatedAt time.Time  `json:"edited_at" gorm:"column:edited_at; sql:type:datetime;default:'1980-01-01'"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at; sql:type:datetime;default:NULL"`
	// CreatedTime int       `json:"created_time"`
	// EditedTime  int       `json:"edited_time"`
	// DeletedTime int       `json:"deleted_time"`
}

func TimeZero() time.Time {
	return times.Zero()
}

func (m *Mysql) Builder(c *Content) error {
	if err := m.CreateDatabase(); err != nil {
		return err
	}
	if close := m.Connector(); close != nil {
		defer close()
		m.DB.SingularTable(true)
		// m.setCreateCallback()
		// m.setUpdateCallback()
		// m.setDeleteCallback()
		m.connectionSetting()
	}
	return nil
}

func (m *Mysql) CreateDatabase() error {
	m.db, m.Error = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			m.Username,
			m.Password,
			m.Host,
			m.Port),
	)
	if m.Error != nil {
		log.Fatalln(m.Error)
		return m.Error
	}
	_, m.Error = m.db.Exec(fmt.Sprintf(
		"CREATE DATABASE "+
			"IF NOT EXISTS %s "+
			"CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;",
		m.DatabaseName))
	if m.Error != nil {
		log.Fatalln(m.Error)
		return m.Error
	}
	_, m.Error = m.db.Exec(fmt.Sprintf(
		"ALTER DATABASE %s "+
			"CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;",
		m.DatabaseName))
	if m.Error != nil {
		log.Fatalln(m.Error)
		return m.Error
	}
	return nil
}

func (m *Mysql) connectionSetting() {
	m.DB.DB().SetMaxIdleConns(m.MaximumIdleConnection)
	m.DB.DB().SetMaxOpenConns(m.MaximumOpenConnection)
	m.DB.DB().SetConnMaxLifetime(time.Duration(m.MaximumConnectionKeepAliveTime))
}

func (m *Mysql) setCreateCallback() {
	m.DB.Callback().Create().Replace("gorm:update_time_stamp", createCallBack)
	return
}

func createCallBack(scope *gorm.Scope) {
	if !scope.HasError() {
		if createTimeField, ok := scope.FieldByName("CreatedTime"); ok {
			if createTimeField.IsBlank {
				fmt.Println(time.Now().Unix())
				createTimeField.Set(time.Now().Unix())
			}
		}
		if modifyTimeField, ok := scope.FieldByName("EditedTime"); ok {
			if modifyTimeField.IsBlank {
				fmt.Println(time.Now().Unix())
				modifyTimeField.Set(time.Now().Unix())
			}
		}
	}
}

func (m *Mysql) setUpdateCallback() {
	m.DB.Callback().Update().Replace("gorm:update_time_stamp", updateCallBack)
	return
}

func updateCallBack(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		fmt.Println(time.Now().Unix())
		scope.SetColumn("EditedTime", time.Now().Unix())
	}
}

func (m *Mysql) setDeleteCallback() {
	m.DB.Callback().Delete().Replace("gorm:delete", deleteCallBack)
	return
}

func deleteCallBack(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		if deletedTimeField, hasDeletedTimeField := scope.FieldByName("DeletedTime"); !scope.Search.Unscoped && hasDeletedTimeField {
			fmt.Println(time.Now().Unix())
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedTimeField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addSpace(scope.CombinedConditionSql()),
				addSpace(extraOption),
			)).Exec()
		} else {
			fmt.Println(time.Now().Unix())
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addSpace(scope.CombinedConditionSql()),
				addSpace(extraOption),
			)).Exec()
		}
	}
}

func addSpace(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

func (m *Mysql) Connector() func() error {
	m.recursionCall(
		func() error {
			m.DB, m.Error = gorm.Open(
				"mysql",
				fmt.Sprintf(
					"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
					m.Username,
					m.Password,
					m.Host+":"+strconv.Itoa(m.Port),
					m.DatabaseName))
			return m.Error
		},
		m.MaximumConnectionRetry,
		m.MinimumRetryDuration,
		false)
	if !utils.AssertErr(m.Error) {
		return m.DB.Close
	}
	return nil
}

func (m *Mysql) recursionCall(f func() error, count, duration int, done bool) bool {
	if !done {
		m.Error = f()
		count--
	}
	if count > 0 && m.Error == nil {
		return true
	} else if count == 0 && m.Error != nil {
		return true
	} else {
		time.Sleep(time.Duration(duration) * time.Second)
	}
	return m.recursionCall(f, count, duration, false)
}

func (m *Mysql) AutoMigrateByAddr(objs ...interface{}) {
	m.ModelAddrs = append(m.ModelAddrs, objs...)
	for _, model := range m.ModelAddrs {
		defer m.Connector()()
		m.DB.Debug().AutoMigrate(model)
	}
	return
}

func (m *Mysql) New() {
	m.ModelAddrs = make([]interface{}, 0)
	return
}

func (m *Mysql) Starter(c *Content) error {
	return nil
}

func (m *Mysql) Router(s *Server) {
	return
}
