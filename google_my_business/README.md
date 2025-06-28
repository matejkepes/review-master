# Google My Business Review Manager

This application helps manage and respond to Google My Business reviews using AI-powered responses and provides monthly review analysis reports.

## Configuration

The application uses a configuration file located at `config/config.properties` or in the root directory as `config.properties`.

### OpenAI API Key

You can configure the OpenAI API key in one of two ways:

1. Add it to your `config.properties` file:

   ```
   openai_api_key=your_openai_api_key
   ```

2. Set it as an environment variable:
   ```
   export OPENAI_API_KEY=your_openai_api_key
   ```

The application will first check the configuration file, and if not found, it will fall back to the environment variable.

### Sample Configuration

A sample configuration file is provided at `config/config.properties.sample`. Copy this file to `config/config.properties` and update the values as needed.

## Changelog

### 2023-05-10: Monthly Review Analysis System Added

- **New Feature**: Automated monthly review analysis and reporting
- **New Feature**: PDF report generation for client reviews
- **New Feature**: Email delivery of reports to clients
- **New Feature**: Retry mechanism for failed analysis attempts
- **New Feature**: Admin summary emails with processing statistics
- **Enhancement**: Command-line interface for manual control
- **Enhancement**: Detailed logging system

## Monthly Review Analysis System

The Monthly Review Analysis system automatically processes Google My Business reviews for all clients at the beginning of each month. It analyzes the previous month's reviews using AI to generate insights, creates PDF reports, and sends them to clients via email.

### System Architecture

The system consists of the following components:

1. **Scheduler**: A cron job that runs on the 1st of each month at 2 AM (UK time)
2. **Analysis Engine**: Processes reviews and generates insights using AI
3. **Report Generator**: Creates PDF reports from analysis results
4. **Email Service**: Sends reports to clients and summary statistics to administrators
5. **Retry Mechanism**: Handles retry attempts for failed clients

### How It Works

#### Scheduled Execution

The system runs automatically at the beginning of each month via a cron job:

```
# Run Monthly Review Analysis for Google My Business
00 02 1 * * root /home/ubuntu/Documents/code/golang/google_my_business/monthly_analysis.sh
```

This executes the `monthly_analysis.sh` script, which:

1. Determines the previous month to analyze
2. Checks if any clients need to be retried from previous failed attempts
3. Runs the analysis program with appropriate parameters
4. Logs all activity to daily log files
5. Cleans up old log files (older than 3 months)
6. Sends a summary email to administrators

#### Analysis Process

For each client with monthly review analysis enabled:

1. Retrieves all Google accounts associated with the client
2. For each account, fetches all locations configured for reporting
3. For each location, fetches all reviews from the previous month
4. Analyzes reviews using AI to identify:
   - Overall sentiment
   - Key themes (positive and negative)
   - Trends compared to previous periods
   - Specific actionable insights
5. Generates a PDF report containing the analysis
6. Sends the report to the client via email

#### Error Handling and Retry Mechanism

If processing fails for any client:

- The failure is tracked in a JSON file (`retry_tracking.json`)
- Each client can be retried up to 3 times
- Subsequent runs with the `--retry-only` flag will attempt to process only the failed clients

### Command Line Options

The analysis program supports several command-line flags:

- `--month=YYYY-MM`: Specifies the month to analyze (defaults to previous month)
- `--force-reprocess=true|false`: Whether to reprocess clients that already have reports
- `--retry-only=true|false`: Only process clients that previously failed
- `--email-summary=email@address.com`: Email address to send the processing summary

Example usage:

```
go run cmd/run_monthly_analysis/main.go --email-summary=admin@example.com
```

### Understanding Monthly Reports

Each monthly report contains:

1. **Executive Summary** - An overview of the month's review activity
2. **Location Summaries** - Detailed analysis for each business location
3. **Trend Analysis** - Comparison with previous periods
4. **Actionable Insights** - Recommended actions based on the reviews

#### Key Metrics in Reports

- **Average Rating**: The mean star rating across all reviews for the period
- **Rating Distribution**: Breakdown of reviews by star rating (1-5)
- **Sentiment Analysis**: Proportion of positive, neutral, and negative reviews
- **Thematic Analysis**: Key topics from positive and negative reviews
- **Trend Information**: How metrics have changed from previous periods

#### Best Practices for Using Reports

1. **Regular Review**: Set aside time each month to review the report
2. **Share Insights**: Distribute relevant findings to location managers and staff
3. **Create Action Plans**: Develop specific plans to address identified issues
4. **Track Progress**: Use subsequent reports to measure improvement
5. **Celebrate Success**: Recognize improvements and positive trends

## Testing AI Responses

You can test AI responses on a CSV file of reviews using the `test_reviews_from_csv.sh` script:

```
./test_reviews_from_csv.sh your_reviews.csv
```

This will process each review in the CSV file and generate AI responses, saving the results to a new CSV file with "\_with_responses" appended to the original filename.
