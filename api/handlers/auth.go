package handlers

import (
	"context"
	"net/http"
	"time"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/middleware"
	"bryce-stabenow/grocer-me/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

// HandleSignupForm serves the signup HTML form
func HandleSignupForm(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, signupHTML)
}

// HandleSigninForm serves the signin HTML form
func HandleSigninForm(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, signinHTML)
}

// HandleSignup handles user registration
func HandleSignup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	now := time.Now()
	user := models.User{
		ID:           primitive.NewObjectID(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set JWT as HTTP-only cookie (24 hours expiration to match token)
	c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true)

	// Return response
	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User: &models.UserPublic{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
	})
}

// HandleSignin handles user login
func HandleSignin(c *gin.Context) {
	var req models.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set JWT as HTTP-only cookie (24 hours expiration to match token)
	c.SetCookie("jwt_token", token, 3600*24, "/", "", false, true)

	// Return response
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: &models.UserPublic{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
	})
}

// HandleGetMe returns the current user's information
func HandleGetMe(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userIDStr, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Find user by ID
	collection := config.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// Return user (password_hash is excluded via json:"-" tag)
	c.JSON(http.StatusOK, user)
}

// generateToken creates a JWT token for the given user ID
func generateToken(userID string) (string, error) {
	// Token expires in 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

const signupHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign Up - Grocer Me</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 40px;
            width: 100%;
            max-width: 400px;
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 28px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            color: #333;
            margin-bottom: 8px;
            font-weight: 500;
            font-size: 14px;
        }
        input {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        input:focus {
            outline: none;
            border-color: #667eea;
        }
        button {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
        }
        button:active {
            transform: translateY(0);
        }
        .message {
            margin-top: 20px;
            padding: 12px;
            border-radius: 8px;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .message.show {
            display: block;
        }
        .link {
            text-align: center;
            margin-top: 20px;
            color: #666;
            font-size: 14px;
        }
        .link a {
            color: #667eea;
            text-decoration: none;
            font-weight: 500;
        }
        .link a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Sign Up</h1>
        <p class="subtitle">Create a new account to get started</p>
        <form id="signupForm">
            <div class="form-group">
                <label for="email">Email</label>
                <input type="email" id="email" name="email" required>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required minlength="6">
            </div>
            <button type="submit">Sign Up</button>
        </form>
        <div id="message" class="message"></div>
        <div class="link">
            Already have an account? <a href="/signin">Sign In</a>
        </div>
    </div>
    <script>
        document.getElementById('signupForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const messageDiv = document.getElementById('message');
            
            try {
                const response = await fetch('/signup', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password })
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    messageDiv.className = 'message success show';
                    messageDiv.innerHTML = 'Account created successfully! Cookie set. <br><button onclick="testMe()" style="margin-top: 10px; padding: 8px 16px; background: #28a745; color: white; border: none; border-radius: 4px; cursor: pointer;">Test /me endpoint</button>';
                    localStorage.setItem('token', data.token);
                } else {
                    messageDiv.className = 'message error show';
                    messageDiv.textContent = 'Error: ' + (data.error || 'Something went wrong');
                }
            } catch (error) {
                messageDiv.className = 'message error show';
                messageDiv.textContent = 'Error: ' + error.message;
            }
        });
        
        async function testMe() {
            const messageDiv = document.getElementById('message');
            try {
                const response = await fetch('/me', {
                    method: 'GET',
                    credentials: 'include'
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    messageDiv.className = 'message success show';
                    messageDiv.innerHTML = 'Success! User info: <pre style="margin-top: 10px; background: #f5f5f5; padding: 10px; border-radius: 4px; overflow-x: auto;">' + JSON.stringify(data, null, 2) + '</pre>';
                } else {
                    messageDiv.className = 'message error show';
                    messageDiv.innerHTML = 'Error: ' + (data.error || 'Failed to fetch user info');
                }
            } catch (error) {
                messageDiv.className = 'message error show';
                messageDiv.innerHTML = 'Error: ' + error.message;
            }
        }
    </script>
</body>
</html>`

const signinHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sign In - Grocer Me</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            padding: 40px;
            width: 100%;
            max-width: 400px;
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 28px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            color: #333;
            margin-bottom: 8px;
            font-weight: 500;
            font-size: 14px;
        }
        input {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        input:focus {
            outline: none;
            border-color: #667eea;
        }
        button {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
        }
        button:active {
            transform: translateY(0);
        }
        .message {
            margin-top: 20px;
            padding: 12px;
            border-radius: 8px;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .message.show {
            display: block;
        }
        .link {
            text-align: center;
            margin-top: 20px;
            color: #666;
            font-size: 14px;
        }
        .link a {
            color: #667eea;
            text-decoration: none;
            font-weight: 500;
        }
        .link a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Sign In</h1>
        <p class="subtitle">Welcome back! Please sign in to your account</p>
        <form id="signinForm">
            <div class="form-group">
                <label for="email">Email</label>
                <input type="email" id="email" name="email" required>
            </div>
            <div class="form-group">
                <label for="password">Password</label>
                <input type="password" id="password" name="password" required>
            </div>
            <button type="submit">Sign In</button>
        </form>
        <div id="message" class="message"></div>
        <div class="link">
            Don't have an account? <a href="/signup">Sign Up</a>
        </div>
    </div>
    <script>
        document.getElementById('signinForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const messageDiv = document.getElementById('message');
            
            try {
                const response = await fetch('/signin', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password })
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    messageDiv.className = 'message success show';
                    messageDiv.innerHTML = 'Signed in successfully! Cookie set. <br><button onclick="testMe()" style="margin-top: 10px; padding: 8px 16px; background: #28a745; color: white; border: none; border-radius: 4px; cursor: pointer;">Test /me endpoint</button>';
                    localStorage.setItem('token', data.token);
                } else {
                    messageDiv.className = 'message error show';
                    messageDiv.textContent = 'Error: ' + (data.error || 'Invalid email or password');
                }
            } catch (error) {
                messageDiv.className = 'message error show';
                messageDiv.textContent = 'Error: ' + error.message;
            }
        });
        
        async function testMe() {
            const messageDiv = document.getElementById('message');
            try {
                const response = await fetch('/me', {
                    method: 'GET',
                    credentials: 'include'
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    messageDiv.className = 'message success show';
                    messageDiv.innerHTML = 'Success! User info: <pre style="margin-top: 10px; background: #f5f5f5; padding: 10px; border-radius: 4px; overflow-x: auto;">' + JSON.stringify(data, null, 2) + '</pre>';
                } else {
                    messageDiv.className = 'message error show';
                    messageDiv.innerHTML = 'Error: ' + (data.error || 'Failed to fetch user info');
                }
            } catch (error) {
                messageDiv.className = 'message error show';
                messageDiv.innerHTML = 'Error: ' + error.message;
            }
        }
    </script>
</body>
</html>`

