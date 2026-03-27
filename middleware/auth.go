package middleware

import (
	"strings"

	"api-gateway/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/*
These comments are not AI generated, and based on my personal understanding of what's happening.
The Fiber Handler wants a function of type : (*fiber.Ctx) error
So we give it that, via our Auth function that we can use anywhere.
*/

func Auth(cfg *config.Config) fiber.Handler {
	return func(context *fiber.Ctx) error {
		/* All we are doing here is extracting the "Authorization" header out of the context Fiber has received.
		If the context doesnt contain any header like that, we just say : Whoopsie daisy get dafaq out.*/
		authHeader := context.Get("Authorization")
		if authHeader == "" {
			return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing auth header"})
		}
		// Considering it does have the header Authorization, lets get the bearer token out of it now.
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == authHeader { // essentially saying no bearer token
			return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing bearer token"})
		}
		// We assume we have a bearer token now, so lets parse it and confirm if its good.
		token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
			/*
				Yeah, this line borked me as well, so much is happening here. Frist the Parse function is taking another function as an argument, which is fine. Returning an empty interfcae is just saying "I don't know what will come out of the other side, straigt up bytes of something, so handle its like TS saying any". ChatGPT taaught me : we are essentially checking the what class is the signing method. Method : interface (SigningMethod) which different Signing Methods implement. So when we do this, we are doing Interface.Struct -> This interface instance belongs to this Struct right? If yes, good else bad. 
			*/
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			return context.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		return context.Next()
	}
}
