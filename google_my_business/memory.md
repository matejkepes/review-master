# Project Memory: AI-Powered Review Response System

## Project Overview

We're building an AI-powered system to automatically generate and validate responses to Google My Business reviews. The system uses OpenAI's GPT models to create personalized, appropriate responses to customer reviews.

## Current Architecture

### 1. AI Service Foundation (`ai_service/`)

- Base LLM provider interface
- OpenAI implementation with retry logic
- Mock provider for testing

### 2. Review Response Generator (`ai_service/review/`)

- Takes review context (text, rating, author, business name)
- Generates appropriate responses using LLM
- Configurable system prompts and parameters

### 3. Response Validator (`ai_service/review/`)

- Validates generated responses against quality criteria
- Returns structured validation results (score, reasons, suggestions)
- Configurable validation thresholds

### 4. Configuration System (`config/`)

- System-level configuration (API keys, model settings)
- Client-level configuration (response style, allowed solutions)
- Configuration validation and merging

## Implementation Details

### LLM Provider

- Interface-based design for swappable providers
- Retry mechanism for API failures
- Supports both basic generation and system/user prompt format

### Generator

- Context-aware response generation
- Customizable prompts per business
- Error handling and timeout support

### Validator

- JSON-based validation results
- Multiple validation criteria
- Provides improvement suggestions

### Configuration

- Environment variable support
- Default configurations
- Separation of system and client settings
- Centralized OpenAI API key configuration through config.properties

### Monthly Review Analysis System

- Implemented `FetchReviews` function to retrieve reviews for a specific location and time period
- Added `FetchReviewsForMonth` function to get reviews for a specific calendar month
- Created database schema for storing analysis results including client reports with proper indexing
- Implemented database access layer for saving and retrieving client reports
- Added functions to get clients with monthly review analysis enabled
- Created skeleton implementation of `AnalyzeClientReviews` orchestration function with:
  - Client, account, and location processing workflow
  - Report existence checking and handling
  - Progress tracking and detailed logging
  - Function variable aliases for easier testing
  - Database storage of analysis results
- Integrated review fetching with the analysis pipeline:
  - Added review fetching function variable for testability
  - Implemented conversion of API review format to analyzer format
  - Connected ReviewAnalyzer to process batches of reviews
  - Added proper error handling throughout the process
- Implemented PDF report generation:
  - Created HTML template for monthly reports
  - Added PDF conversion functionality
- Implemented automated monthly scheduling system:
  - Created a cron job to run on the 1st of each month
  - Added shell script to manage the process
  - Implemented command-line interface for manual control
  - Added support for specifying target month
  - Added force-reprocess option for regenerating reports
- Added email notification system:
  - Created SendPlainTextEmail function for generic emails
  - Implemented summary email generation
  - Added email delivery of statistics to administrators
- Implemented error handling and retry mechanism:
  - Created file-based retry tracking system
  - Limited retries to 3 attempts per client
  - Added retry-only mode for failed clients
  - Added detailed error logging
- Enhanced logging and monitoring:
  - Added daily log files for each run
  - Implemented log rotation (90-day retention)
  - Enhanced ProcessingSummary with detailed error information
  - Added command-line options for controlling output

## Next Steps

1. Complete Monthly Review Analysis Implementation

   - Implement report generation and delivery (HTML/PDF)
   - Add email delivery mechanism for reports using SendGrid

2. Database schema for client configurations
3. Integration of configs with generator/validator
4. Full response generation and validation flow
5. Response storage and metrics
6. Migration of ai_responses_enabled and contact_method columns from clients to google_reviews_configs table
   - Modified the GoogleReviewsConfigFromGoogleMyBusinessLocationNameAndPostalCode struct to handle NULL values for ContactMethod
   - Updated the code to handle the pointer for ContactMethod
   - Created a migration script to move the columns from clients to google_reviews_configs

## Technical Decisions

- Using Go for all components
- OpenAI GPT-4 as the default model
- JSON for structured LLM outputs
- Interface-based design for testability
- Configuration split between system and client levels
- Monthly review analysis will use the existing LLMProvider interface
- Analysis reports will be stored as JSON in the database for portal integration
- Monthly reviews are fetched using a specific date-based approach instead of a relative "months back" parameter for more flexibility
- Testing is facilitated by function variable aliases that allow easy mocking of Google API functions
- Review analysis functionality is separated from review response functionality to maintain clean separation of concerns
