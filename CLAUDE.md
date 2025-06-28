# Review Master - Multi-Repository Project Overview

## Project Summary

Review Master is a comprehensive Google Reviews and Google My Business management platform designed to help businesses automate review collection, analyze customer feedback, and generate monthly reports. The system integrates with Google My Business API to manage reviews across multiple business locations and provides AI-powered insights using OpenAI's GPT models.

## Architecture Overview

The project follows a microservices architecture with separate backend services and frontend applications:

### Backend Services (Go)

1. **google_reviews** - Core review collection and SMS gateway service
2. **google_my_business** - Google My Business API integration and AI analysis
3. **google_reviews_ui** - Admin UI backend service
4. **rm_client_portal** - Client portal backend service

### Frontend Applications (Vue.js/Quasar)

1. **google_reviews_fe** - Admin interface for managing review configurations
2. **rm_client_portal_fe** - Client-facing portal for viewing reports and analytics

## Module Details

### 1. google_my_business/

**Purpose**: Core module for Google My Business integration and automated review management

**Key Features**:
- OAuth 2.0 integration with Google My Business API
- AI-powered review response generation using OpenAI
- Monthly review analysis and report generation
- PDF report generation with insights and trends
- Automated email delivery via SendGrid
- Retry mechanism for failed processing

**Technologies**:
- Go programming language
- OpenAI API for AI analysis
- Google My Business API
- SendGrid for email delivery
- MySQL database

**Entry Points**:
- `my_business.go` - Main application for review processing
- `cmd/run_monthly_analysis/main.go` - Monthly analysis runner

### 2. google_reviews/

**Purpose**: SMS gateway and review request management system

**Key Features**:
- REST API for dispatcher integrations (CAB9, Cordic)
- SMS gateway for sending review requests
- Telephone number validation and barring
- Rate limiting and scheduling
- Multi-message support
- Statistics tracking

**Technologies**:
- Go programming language
- MySQL database
- HTTPS with TLS certificates
- Integration with various SMS providers

**Entry Point**:
- `google_reviews.go` - Main HTTP server application

### 3. google_reviews_fe/

**Purpose**: Admin frontend for managing review configurations

**Key Features**:
- Client management interface
- Review configuration settings
- Statistics and analytics dashboard
- User authentication with JWT
- Real-time monitoring of review campaigns

**Technologies**:
- Vue.js 3 with Quasar Framework
- Axios for API communication
- JWT for authentication
- QR code generation

### 4. google_reviews_ui/

**Purpose**: Backend service for the admin UI

**Key Features**:
- RESTful API for frontend
- User authentication and authorization
- Client and configuration management
- Statistics aggregation
- Integration with Google My Business module

**Technologies**:
- Go programming language
- JWT authentication
- MySQL database
- HTTPS server

### 5. rm_client_portal/

**Purpose**: Client portal backend for accessing reports

**Key Features**:
- Client authentication
- Report viewing and downloading
- Google My Business data access
- Monthly report retrieval
- API for frontend portal

**Technologies**:
- Go programming language
- Google OAuth integration
- MySQL database
- JWT authentication

### 6. rm_client_portal_fe/

**Purpose**: Client-facing portal for viewing analytics and reports

**Key Features**:
- Modern dashboard with charts (ApexCharts)
- Report viewing and downloading
- Authentication system
- Responsive design

**Technologies**:
- Vue.js 3 with Quasar Framework
- TypeScript
- Pinia for state management
- ApexCharts for data visualization
- Zod for validation

## Business Flow

1. **Review Request Flow**:
   - Dispatcher systems (taxi/cab companies) trigger review requests after completed rides
   - System validates customer phone numbers and applies rate limiting
   - SMS messages are sent with Google review links
   - All interactions are logged for statistics

2. **Review Management Flow**:
   - Google My Business API fetches new reviews periodically
   - AI analyzes reviews for sentiment and generates appropriate responses
   - Responses can be automatically posted or queued for approval
   - All reviews are stored for reporting

3. **Monthly Analysis Flow**:
   - Scheduled job runs on the 1st of each month
   - Fetches all reviews from the previous month
   - AI analyzes reviews for trends, themes, and insights
   - Generates PDF reports with visualizations
   - Emails reports to clients and administrators

## Key Technologies Stack

- **Backend**: Go (Golang)
- **Frontend**: Vue.js 3, Quasar Framework, TypeScript
- **Database**: MySQL
- **AI/ML**: OpenAI GPT-4
- **APIs**: Google My Business API, SendGrid API
- **Authentication**: JWT tokens
- **Communication**: RESTful APIs, HTTPS

## Security Features

- TLS/HTTPS encryption for all communications
- JWT-based authentication
- API token validation
- Database credential encryption
- Rate limiting and abuse prevention
- Telephone number barring lists

## Deployment

Each module includes deployment scripts (`deploy.sh`) and systemd service configurations for Linux deployment. The system is designed to run on Ubuntu servers with MySQL database backend.

## Configuration

The system uses property files for configuration with support for:
- Database connections
- API credentials (Google, OpenAI, SendGrid)
- SMTP settings
- Time zones and scheduling
- Rate limiting parameters

This is a production-ready system designed for multi-tenant SaaS deployment, serving multiple business clients with isolated data and customizable configurations.