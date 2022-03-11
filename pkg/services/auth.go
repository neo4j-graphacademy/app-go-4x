package services

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/services/jwtutils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"golang.org/x/crypto/bcrypt"
)

type User map[string]interface{}

type AuthService interface {
	Save(email, plainPassword, name string) (User, error)

	FindOneByEmailAndPassword(email string, password string) (User, error)

	ExtractUserId(bearer string) (string, error)
}

type neo4jAuthService struct {
	driver     neo4j.Driver
	jwtSecret  string
	saltRounds int
}

func NewAuthService(driver neo4j.Driver, jwtSecret string, saltRounds int) AuthService {
	return &neo4jAuthService{
		driver:     driver,
		jwtSecret:  jwtSecret,
		saltRounds: saltRounds,
	}
}

// Save should create a new User node in the database with the email and name
// provided, along with an encrypted version of the password and a `userId` property
// generated by the server.
//
// The properties also be used to generate a JWT `token` which should be included
// with the returned user.
// tag::register[]
func (as *neo4jAuthService) Save(email, plainPassword, name string) (_ User, err error) {
	// TODO: Handle Unique constraints in the database
	// if email != "graphacademy@neo4j.com" {
	// 	return nil, fmt.Errorf("An account already exists with this email address")
	// }

	// user, err := fixtures.ReadObject("fixtures/user.json")
	// if err != nil {
	// 	return nil, err
	// }

	// subject := user["userId"].(string)
	// token, err := jwtutils.Sign(subject, userToClaims(user), as.jwtSecret)
	// if err != nil {
	// 	return nil, err
	// }

	// return userWithToken(user, token), nil

	encryptedPassword, err := encryptPassword(plainPassword, as.saltRounds)
	if err != nil {
		return nil, err
	}

	// Open a new Session
	// tag::catch[]
	// tag::session[]
	session := as.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()
	// end::session[]

	// tag::create[]
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			CREATE (u:User {
				  userId: randomUuid(),
				  email: $email,
				  password: $encrypted,
				  name: $name
			})
			RETURN u { .userId, .name, .email } as u`,
			map[string]interface{}{
				"email":     email,
				"encrypted": encryptedPassword,
				"name":      name,
			})
		// end::create[]

		// tag::catch[]
		// Check the error title
		if neo4jError, ok := err.(*neo4j.Neo4jError); ok && neo4jError.Title() == "ConstraintValidationFailed" {
			return nil, fmt.Errorf(fmt.Sprintf("A user already exists with email %s", email))
		}

		if err != nil {
			return nil, err
		}
		// end::catch[]

		// tag::extract[]
		// Extract safe properties from the user node (`u`) in the first row
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		user, _ := record.Get("u")
		return user, nil
		// end::extract[]
	})
	if err != nil {
		return nil, err
	}

	// tag::return[]
	user := result.(map[string]interface{})
	subject := user["userId"].(string)
	token, err := jwtutils.Sign(subject, userToClaims(user), as.jwtSecret)
	if err != nil {
		return nil, err
	}
	return userWithToken(user, token), nil
	// end::return[]
	// end::catch[]
}

// end::register[]

// tag::authenticate[]
func (as *neo4jAuthService) FindOneByEmailAndPassword(email string, password string) (_ User, err error) {
	// TODO: Authenticate the user from the database
	if email != "graphacademy@neo4j.com" {
		return nil, fmt.Errorf("Incorrect username or password")
	}

	user, err := fixtures.ReadObject("fixtures/user.json")
	if err != nil {
		return nil, err
	}

	subject := user["userId"].(string)
	token, err := jwtutils.Sign(subject, userToClaims(user), as.jwtSecret)
	if err != nil {
		return nil, err
	}

	return userWithToken(user, token), nil

	// // Open a new Session
	// // tag::catch[]
	// session := as.driver.NewSession(neo4j.SessionConfig{})
	// defer func() {
	// 	err = ioutils.DeferredClose(session, err)
	// }()

	// // tag::query[]
	// // Find the User node within a Read Transaction
	// result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
	// 	result, err := tx.Run(`
	// 		MATCH (u:User {email: $email}) RETURN u`,
	// 		map[string]interface{}{
	// 			"email": email,
	// 		})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	record, err := result.Single()
	// 	if err != nil {
	// 		// do not expose whether an account matches or not
	// 		return nil, fmt.Errorf("account not found or incorrect password")
	// 	}
	// 	user, _ := record.Get("u")
	// 	return user, nil
	// })
	// // end::query[]

	// if err != nil {
	// 	return nil, err
	// }

	// // tag::password[]
	// // Check password
	// userNode := result.(neo4j.Node)
	// user := userNode.Props
	// if !verifyPassword(password, user["password"].(string)) {
	// 	return nil, fmt.Errorf("account not found or incorrect password")
	// }
	// // end::password[]

	// // tag::authreturn[]
	// subject := userNode.Props["userId"].(string)
	// token, err := jwtutils.Sign(subject, userToClaims(user), as.jwtSecret)
	// if err != nil {
	// 	return nil, err
	// }
	// return userWithToken(user, token), nil
	// // end::authreturn[]
}

// end::authenticate[]

func (as *neo4jAuthService) ExtractUserId(bearer string) (string, error) {
	if bearer == "" {
		return "", nil
	}
	userId, err := jwtutils.ExtractToken(bearer, as.jwtSecret, func(token *jwt.Token) interface{} {
		claims := token.Claims.(jwt.MapClaims)
		return claims["sub"]
	})
	if err != nil {
		return "", err
	}
	return userId.(string), nil
}

func encryptPassword(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	return err == nil
}

func userToClaims(user User) map[string]interface{} {
	return map[string]interface{}{
		"sub":    user["userId"],
		"userId": user["userId"],
		"name":   user["name"],
	}
}

func userWithToken(user User, token string) User {
	return map[string]interface{}{
		"token":  token,
		"userId": user["userId"],
		"email":  user["email"],
		"name":   user["name"],
	}
}
