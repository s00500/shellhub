package main

import (
	"net/http"
	"testing"
	//"time"

	"fmt"

	"github.com/gavv/httpexpect/v2"
	"github.com/shellhub-io/shellhub/pkg/models"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/routes"
	"github.com/shellhub-io/shellhub/api/store/mongo"
)

func testAuthUser(e *httpexpect.Expect) (user, token, tenant string) {
	type Login struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	//tries login with a wrong password
	e.POST("/api/login").WithForm(Login{"username", "<bad password>"}).
		Expect().
		Status(http.StatusUnauthorized)
	//login with  a correct password
	authUser := e.POST("/api/login").WithForm(Login{"username", "password"}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	authUser.Keys().ContainsOnly("user", "name", "tenant", "email", "token")

	token = authUser.Value("token").String().Raw()
	tenant = authUser.Value("tenant").String().Raw()
	user = authUser.Value("user").String().Raw()
	return user, token, tenant
}
func testAuthDevice(e *httpexpect.Expect, authReq models.DeviceAuthRequest, username string) (uid string) {

	authDevice := e.POST("/api/devices/auth").WithJSON(authReq).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	authDevice.Keys().ContainsOnly("name", "namespace", "token", "uid")
	authDevice.Value("name").Equal("mac")
	authDevice.Value("namespace").Equal(username)
	return authDevice.Value("uid").String().Raw()
}

func testGetDevice(e *httpexpect.Expect, uid, token string, device map[string]interface{}) {

	e.GET(fmt.Sprintf("/api/devices/%s", uid)).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsMap(device)

}

func testListDevices(e *httpexpect.Expect, device map[string]interface{}, token, tenant string) {
	listDevices := e.GET("/api/devices").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	for _, val := range listDevices.Iter() {
		val.Object().ContainsMap(device)
	}

}

func testGetToken(e *httpexpect.Expect, tenant string) {
	e.GET(fmt.Sprintf("/internal/auth/token/%s", tenant)).
		Expect().
		Status(http.StatusOK)
}
func testRenameDevice(e *httpexpect.Expect, rename []string, status []int, uid, token, username, tenant string) {

	for i, j := range rename {
		renameReq := map[string]interface{}{
			"name": j,
		}

		e.PATCH(fmt.Sprintf("/api/devices/%s", uid)).
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", tenant).
			WithHeader("X-Username", username).
			WithJSON(renameReq).
			Expect().
			Status(status[i])

	}
}
func testUpdatePendingStatus(e *httpexpect.Expect, pendingArray []string, uid, token, username, tenant string) {
	for _, j := range pendingArray {

		e.PATCH(fmt.Sprintf("/api/devices/%s/%s", uid, j)).
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", tenant).
			WithHeader("X-Username", username).
			Expect().
			Status(http.StatusOK)
	}
}
func testLookupDevice(e *httpexpect.Expect, lookup map[string]string, token, username, tenant string) {

	e.GET(fmt.Sprintf("/internal/lookup")).
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("X-Tenant-ID", tenant).
		WithHeader("X-Username", username).
		WithJSON(lookup).
		Expect().
		Status(http.StatusOK)
}

func testOfflineDevice(e *httpexpect.Expect, uid, token string) {
	e.POST(fmt.Sprintf("/internal/devices/%s/offline", uid)).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK)
}
func testCreateSession(e *httpexpect.Expect, session map[string]interface{}) {
	e.POST("/internal/sessions").WithJSON(session).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
}

func testAuthenticateSession(e *httpexpect.Expect, uid_session string, authenticated map[string]interface{}) {
	e.PATCH(fmt.Sprintf("/internal/sessions/%s", uid_session)).
		WithJSON(authenticated).
		Expect().
		Status(http.StatusOK)
}
func testGetSession(e *httpexpect.Expect, uid_session, token string, sessionAuth map[string]interface{}) {
	e.GET(fmt.Sprintf("/api/sessions/%s", uid_session)).
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().Object().ContainsMap(sessionAuth)
}

func testListSessions(e *httpexpect.Expect, token string, sessionAuth map[string]interface{}) {
	array := e.GET("/api/sessions").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).JSON().Array()

	for _, val := range array.Iter() {
		val.Object().ContainsMap(sessionAuth)
	}
}
func testFinishSession(e *httpexpect.Expect, uid_session, token, username, tenant string) {
	e.POST(fmt.Sprintf("/internal/sessions/%s/finish", uid_session)).
		WithHeader("Authorization", "Bearer"+token).
		WithHeader("X-Tenant-ID", tenant).
		WithHeader("X-Username", username).
		Expect().
		Status(http.StatusOK)
}
func testStats(e *httpexpect.Expect, token string) {
	// public tests for stats
	e.GET("/api/stats").
		WithHeader("Authorization", "Bearer "+token).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
}

func testDeleteDevice(e *httpexpect.Expect, uid, token, username, tenant string) {
	e.DELETE(fmt.Sprintf("/api/devices/%s", uid)).
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("X-Tenant-ID", tenant).
		WithHeader("X-Username", "username").
		Expect().
		Status(http.StatusOK)

}

func testUpdateUser(e *httpexpect.Expect, forms_array []interface{}, status_array []int, token, username, tenant string) {
	for i, v := range forms_array {
		e.PUT("/api/user").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("X-Tenant-ID", tenant).
			WithHeader("X-Username", username).
			WithJSON(v).
			Expect().
			Status(status_array[i])
	}

}

func testUpdateUserSecurity(e *httpexpect.Expect, status map[string]bool, token, tenant, username string) {
	e.PUT("api/user/security").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("X-Tenant-ID", tenant).
		WithHeader("X-Username", username).
		WithJSON(status).
		Expect().
		Status(http.StatusOK)

}

func GetUserSecurity(e *httpexpect.Expect, token, tenant, username string) {
	e.PUT("api/user/security").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("X-Tenant-ID", tenant).
		WithHeader("X-Username", username).
		Expect().
		Status(http.StatusOK)
}

func TestEchoClient(t *testing.T) {

	/*	e := httpexpect.WithConfig(httpexpect.Config{
		// prepend this url to all requests
		BaseURL: "http://api:8080/",

		// use http.Client with a cookie jar and timeout
		Client: &http.Client{
			Jar:     httpexpect.NewJar(),
			Timeout: time.Second * 30,
		},

		// use fatal failures
		Reporter: httpexpect.NewRequireReporter(t),

		// use verbose logging
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewDebugPrinter(t, true),
		},
	})*/

	handler := EchoHandler()

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	username, token, tenant := testAuthUser(e)

	authReq := &models.DeviceAuthRequest{
		Info: &models.DeviceInfo{
			ID:         "id",
			PrettyName: "Pretty name",
			Version:    "test",
		},
		DeviceAuth: &models.DeviceAuth{
			TenantID: tenant,
			Identity: &models.DeviceIdentity{
				MAC: "mac",
			},
			PublicKey: "key",
		},
	}

	uid := testAuthDevice(e, *authReq, username)

	device := map[string]interface{}{
		"identity": map[string]string{
			"mac": "mac",
		},
		"info": map[string]string{
			"id":          "id",
			"pretty_name": "Pretty name",
			"version":     "test",
		},
		"name":       "mac",
		"namespace":  username,
		"public_key": "key",
		"status":     "pending",
		"tenant_id":  "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	}
	testGetDevice(e, uid, token, device)
	testListDevices(e, device, token, tenant)
	rename := []string{"@#$%", "mac", "newName", "mac"}
	status := []int{http.StatusForbidden, http.StatusForbidden, http.StatusOK, http.StatusOK}
	testRenameDevice(e, rename, status, uid, token, username, tenant)
	pendingArray := []string{"unused", "pending", "rejected", "accepted"}
	testUpdatePendingStatus(e, pendingArray, uid, token, username, tenant)
	lookup := map[string]string{
		"domain":     "username",
		"name":       "mac",
		"username":   "username",
		"ip_address": "1.1.1.1",
	}
	testLookupDevice(e, lookup, token, username, tenant)
	testOfflineDevice(e, uid, token)

	session := map[string]interface{}{
		"username":      "username",
		"device_uid":    uid,
		"uid":           "uid",
		"authenticated": false,
	}
	uid_session := "uid"

	authenticated := map[string]interface{}{
		"authenticated": true,
	}

	sessionAuth := map[string]interface{}{
		"username":      "username",
		"device_uid":    uid,
		"uid":           "uid",
		"authenticated": true,
	}

	testCreateSession(e, session)
	testAuthenticateSession(e, uid_session, authenticated)
	testGetSession(e, uid_session, token, sessionAuth)
	testListSessions(e, token, sessionAuth)
	testFinishSession(e, uid_session, token, username, tenant)
	testStats(e, token)
	testDeleteDevice(e, uid, token, username, tenant)

	status_array := []int{http.StatusOK, http.StatusOK, http.StatusConflict, http.StatusForbidden}
	forms_array := []interface{}{
		map[string]interface{}{ // successfull email and username change
			"username":        "newusername",
			"email":           "new@email.com",
			"currentPassword": "",
			"newPassword":     "",
		},
		map[string]interface{}{ // successfull password change
			"username":        "",
			"email":           "",
			"currentPassword": "password",
			"newPassword":     "new_password_hash",
		},
		map[string]interface{}{ //conflict
			"username":        "username2",
			"email":           "new@email.com",
			"currentPassword": "",
			"newPassword":     "",
		},
		map[string]interface{}{ // forbidden
			"username":        "",
			"email":           "",
			"currentPassword": "wrong_password",
			"newPassword":     "new_password",
		},
	}
	testUpdateUser(e, forms_array, status_array, token, username, tenant)
	sessionRecord := map[string]bool{
		"sessionRecord": true,
	}
	testUpdateUserSecurity(e, sessionRecord, token, tenant, username)

}

func EchoHandler() http.Handler {
	{
		e := echo.New()
		e.Use(middleware.Logger())

		/*var cfg config
		if err := envconfig.Process("api", &cfg); err != nil {
			panic(err.Error())
		}*/

		db := dbtest.DBServer{}
		defer db.Stop()

		if err := mongo.ApplyMigrations(db.Client().Database("test")); err != nil {
			panic(err)
		}

		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				store := mongo.NewStore(db.Client().Database("test"))
				ctx := apicontext.NewContext(store, c)

				return next(ctx)
			}
		})

		publicAPI := e.Group("/api")
		internalAPI := e.Group("/internal")

		internalAPI.GET(routes.AuthRequestURL, apicontext.Handler(routes.AuthRequest), apicontext.Middleware(routes.AuthMiddleware))
		publicAPI.POST(routes.AuthDeviceURL, apicontext.Handler(routes.AuthDevice))
		publicAPI.POST(routes.AuthDeviceURLV2, apicontext.Handler(routes.AuthDevice))
		publicAPI.POST(routes.AuthUserURL, apicontext.Handler(routes.AuthUser))
		publicAPI.POST(routes.AuthUserURLV2, apicontext.Handler(routes.AuthUser))
		publicAPI.GET(routes.AuthUserURLV2, apicontext.Handler(routes.AuthUserInfo))
		internalAPI.GET(routes.AuthUserTokenURL, apicontext.Handler(routes.AuthGetToken))

		publicAPI.PUT(routes.UpdateUserURL, apicontext.Handler(routes.UpdateUser))
		publicAPI.PUT(routes.UserSecurityURL, apicontext.Handler(routes.UpdateUserSecurity))
		publicAPI.GET(routes.UserSecurityURL, apicontext.Handler(routes.GetUserSecurity))

		publicAPI.GET(routes.GetDeviceListURL, apicontext.Handler(routes.GetDeviceList))
		publicAPI.GET(routes.GetDeviceURL, apicontext.Handler(routes.GetDevice))
		publicAPI.DELETE(routes.DeleteDeviceURL, apicontext.Handler(routes.DeleteDevice))
		publicAPI.PATCH(routes.RenameDeviceURL, apicontext.Handler(routes.RenameDevice))
		internalAPI.POST(routes.OfflineDeviceURL, apicontext.Handler(routes.OfflineDevice))
		internalAPI.GET(routes.LookupDeviceURL, apicontext.Handler(routes.LookupDevice))
		publicAPI.PATCH(routes.UpdateStatusURL, apicontext.Handler(routes.UpdatePendingStatus))

		publicAPI.GET(routes.GetSessionsURL, apicontext.Handler(routes.GetSessionList))
		publicAPI.GET(routes.GetSessionURL, apicontext.Handler(routes.GetSession))
		internalAPI.PATCH(routes.SetSessionAuthenticatedURL, apicontext.Handler(routes.SetSessionAuthenticated))
		internalAPI.POST(routes.CreateSessionURL, apicontext.Handler(routes.CreateSession))
		internalAPI.POST(routes.FinishSessionURL, apicontext.Handler(routes.FinishSession))
		internalAPI.POST(routes.RecordSessionURL, apicontext.Handler(routes.RecordSession))
		publicAPI.GET(routes.PlaySessionURL, apicontext.Handler(routes.PlaySession))

		publicAPI.GET(routes.GetStatsURL, apicontext.Handler(routes.GetStats))

		return e
	}
}
