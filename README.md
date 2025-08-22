# Chinedu Onyeizu Foundation

A charity foundation website with admin authentication system.

## Setup

1. Make sure you have a `.env` file with the following variables:
   ```
   DB_USERNAME=your_db_username
   DB_PASSWORD=your_db_password
   DB_HOST=localhost
   DB_PORT=5432
   DB_DRIVER=postgres
   DB_NAME=your_db_name
   SSLMODE=disable
   SERVER_PORT=8080
   JWT_SECRET=your_jwt_secret_key
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## Creating an Admin User

Since you don't have any admin users registered, you need to create one first:

```bash
go run cmd/create_admin/main.go -email=admin@example.com -password=yourpassword
```

Replace `admin@example.com` and `yourpassword` with your desired admin credentials.

## Login Flow

1. Visit `/admin/login` to access the login page
2. Enter the email and password you created above
3. Upon successful login, you'll be redirected to `/dashboard`
4. The dashboard is protected and requires authentication

## Fixed Issues

The following issues have been resolved:

1. **AuthService config missing**: The `NewAuthService` function now properly receives and stores the config
2. **Error handling**: Login errors now show proper error messages
3. **No admin users**: Created a script to easily create admin users
4. **Dashboard protection**: Added middleware to protect the dashboard route
5. **Better logging**: Added detailed logging for debugging authentication issues

## Troubleshooting

If you're still having issues:

1. Check that your database is running and accessible
2. Verify your `.env` file has all required variables
3. Make sure you've created an admin user using the script above
4. Check the server logs for any error messages
5. Ensure the JWT_SECRET is set in your environment variables 