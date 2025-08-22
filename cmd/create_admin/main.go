package main

import (
	db "charity/db/sqlc"
	"charity/internal/auth"
	"charity/internal/config"
	p "charity/internal/db"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	var (
		email    = flag.String("email", "", "Admin email address")
		password = flag.String("password", "", "Admin password")
		help     = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Validate required flags
	if *email == "" || *password == "" {
		fmt.Println("❌ Error: Both email and password are required")
		fmt.Println()
		showHelp()
		os.Exit(1)
	}

	// Validate email format
	if !isValidEmail(*email) {
		fmt.Println("❌ Error: Invalid email format")
		os.Exit(1)
	}

	// Validate password strength
	if len(*password) < 6 {
		fmt.Println("❌ Error: Password must be at least 6 characters long")
		os.Exit(1)
	}

	fmt.Println("🚀 Creating admin user...")
	fmt.Printf("📧 Email: %s\n", *email)
	fmt.Printf("🔒 Password: %s\n", strings.Repeat("*", len(*password)))

	// Load configuration
	fmt.Println("\n📋 Loading configuration...")
	c, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Configuration loaded successfully")

	// Connect to database
	fmt.Println("\n🗄️  Connecting to database...")
	dbConn, err := sql.Open(c.DBDriver, p.GetDBSource(c, c.DBName))
	if err != nil {
		fmt.Printf("❌ Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			fmt.Printf("⚠️  Warning: Failed to close database connection: %v\n", err)
		}
	}()

	// Test database connection
	if err := dbConn.Ping(); err != nil {
		fmt.Printf("❌ Failed to ping database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Database connection established")

	// Create queries instance
	queries := db.New(dbConn)

	// Create auth service
	authService := auth.NewAuthService(queries, c)

	// Check if admin already exists
	fmt.Println("\n🔍 Checking if admin already exists...")
	existingAdmin, err := queries.GetAdminByEmail(context.Background(), *email)
	if err == nil {
		fmt.Printf("⚠️  Admin with email %s already exists (ID: %d)\n", existingAdmin.Email, existingAdmin.ID)
		fmt.Println("   Use a different email or update the existing admin password manually")
		os.Exit(1)
	}
	fmt.Println("✅ No existing admin found with this email")

	// Create admin user
	fmt.Println("\n👤 Creating admin user...")
	admin, err := authService.CreateAdminUser(context.Background(), *email, *password)
	if err != nil {
		fmt.Printf("❌ Failed to create admin user: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n🎉 Admin user created successfully!")
	fmt.Printf("📊 User ID: %d\n", admin.ID)
	fmt.Printf("📧 Email: %s\n", admin.Email)
	fmt.Printf("📅 Created: %s\n", admin.CreatedAt.Time.Format("2006-01-02 15:04:05"))

	fmt.Println("\n🔗 You can now login at: http://localhost:8080/admin/login")
	fmt.Println("   Use the email and password you just created")
}

func showHelp() {
	fmt.Println("Admin User Creation Tool")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Println("This tool creates a new admin user for the Chinedu Onyeizu Foundation.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/create_admin/main.go -email=admin@example.com -password=yourpassword")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -email    Admin email address (required)")
	fmt.Println("  -password Admin password (required, min 6 characters)")
	fmt.Println("  -help     Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/create_admin/main.go -email=admin@foundation.org -password=securepass123")
	fmt.Println("  go run cmd/create_admin/main.go -email=john@example.com -password=mypassword")
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".") && len(email) > 5
}
