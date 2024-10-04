package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nsvirk/moneybotstds/pkg/response"
)

// AuthMiddleware creates a new authorization middleware
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return response.ErrorResponse(c, http.StatusUnauthorized, "AuthorizationException", "Missing Authorization header")
			}

			parts := strings.SplitN(auth, ":", 2)
			if len(parts) != 2 {
				return response.ErrorResponse(c, http.StatusUnauthorized, "AuthorizationException", "Invalid Authorization header format")
			}

			userID, enctoken := parts[0], parts[1]

			// Verify the enctoken
			valid, err := verifyEnctoken(enctoken)
			if err != nil || !valid {
				return response.ErrorResponse(c, http.StatusUnauthorized, "AuthorizationException", "Invalid or expired session")
			}

			// Add session data to context for use in handlers
			c.Set("userID", userID)
			c.Set("enctoken", enctoken)

			// Get from the context to verify that the data was set
			userID = c.Get("userID").(string)
			enctoken = c.Get("enctoken").(string)

			return next(c)
		}
	}
}

// verifyEnctoken verifies the enctoken
func verifyEnctoken(enctoken string) (bool, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://kite.zerodha.com/oms/user/profile", nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "enctoken "+enctoken)
	req.Header.Add("X-Kite-Version", "3")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "kite.zerodha.com")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 200, nil
}
