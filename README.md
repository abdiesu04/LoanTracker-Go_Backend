# LoanTracker-Go_Backend
## Loan Tracker API

### Features

#### 1. **User Management**
- **User Registration**: Register a new user with email, password, and profile details.
- **Email Verification**: Verify the user's email address via a token-based system.
- **User Login**: Authenticate users with JWT-based access and refresh tokens.
- **Password Reset**: Allows users to request a password reset link and update their password.
- **User Profile**: Retrieve authenticated user profile details.

#### 2. **Loan Management**
- **Apply for Loan**: Users can submit loan applications.
- **View Loan Status**: Users can check the status of their loan applications.
- **Loan Approval/Rejection (Admin)**: Admins can approve or reject loan applications.
- **View All Loans (Admin)**: Admins can access and filter loan applications by status.
- **Delete Loan (Admin)**: Admins can delete a specific loan application.

#### 3. **Admin Functionalities**
- **Manage Users**: View all registered users and delete specific user accounts.
- **System Logs**: View detailed logs of user login attempts, loan submissions, password resets, and loan status updates.

#### 4. **Security & Performance**
- **Secure Password Handling**: Passwords are hashed using bcrypt.
- **JWT Authentication**: Stateless access and refresh tokens for authentication.
- **Role-Based Access Control (RBAC)**: Ensure only admins have access to administrative functionalities.
- **Concurrency with Goroutines**: Optimize request handling and processing tasks with Go's concurrency model.

#### 5. **Documentation**
- **Postman API Documentation**: All API endpoints are documented with sample requests, responses, and error codes.

---


######[Postman Documentation](https://documenter.getpostman.com/view/24071191/2sAXjGdu5q#34a50d84-c91f-4145-85d7-7efb24126e64)

