package shared

// MonthlyReportTemplate contains the HTML template for monthly review analysis reports.
//
// This template is designed to be self-contained and portable between Go projects.
// It uses JavaScript for all formatting and requires no external helper functions.
//
// The template expects a data structure with the following fields:
// - ClientName (string): Name of the client
// - PeriodStart (time.Time): Start date of the reporting period
// - PeriodEnd (time.Time): End date of the reporting period
// - GeneratedAt (time.Time): When the report was generated
// - LocationResults ([]LocationResult): Array of location results containing:
//   - Metadata:
//   - LocationName (string): Name of the location
//   - Analysis:
//   - OverallSummary:
//   - AverageRating (float64): Average star rating
//   - PositiveThemes ([]string): List of positive themes
//   - NegativeThemes ([]string): List of negative themes
//   - SentimentAnalysis:
//   - PositivePercentage (float64): Percentage of positive reviews
//   - NeutralPercentage (float64): Percentage of neutral reviews
//   - NegativePercentage (float64): Percentage of negative reviews
//   - TotalReviews (int): Total number of reviews analyzed
//   - SentimentTrend (string): Description of sentiment trend
//   - KeyTakeaways:
//   - Strengths ([]Insight): Each having Category, Description, and Example
//   - AreasForImprovement ([]Insight): Each having Category, Description, and Example
//   - NegativeReviewBreakdown:
//   - ImprovementRecommendations ([]string): List of suggested improvements
//   - TrainingRecommendations:
//   - ForOperators ([]string): List of training recommendations for operators
//   - ForDrivers ([]string): List of training recommendations for drivers
//
// To use this template, simply:
// 1. Import this package
// 2. Parse the template with template.New("report").Parse(shared.MonthlyReportTemplate)
// 3. Execute the template with your data structure
const MonthlyReportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Monthly Review Analysis - {{.ClientName}}</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        /* Base styles */
        body {
            font-family: 'Arial', sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        
        /* Header styles */
        .report-header {
            background-color: #04062c;
            padding: 25px 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 30px;
            flex-shrink: 0;
        }
        
        .header-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .header-title {
            flex: 1;
            text-align: left;
        }
        
        .header-meta {
            text-align: right;
            color: #ffffff;
            font-size: 14px;
        }
        
        .report-header h1 {
            color: #f5bd41;
            margin: 0 0 5px 0;
            font-size: 28px;
        }
        
        .report-header h2 {
            color: #ffffff;
            margin: 0;
            font-size: 20px;
            margin-top: 15px;
            font-weight: normal;
        }
        
        .subtitle {
            font-size: 16px;
            color: #ffffff;
            margin-top: 5px;
            font-weight: normal;
        }
        
        /* Logo styles */
        .logo-image {
            height: 24px;
            margin-top: 5px;
        }
        
        /* Executive summary styles */
        .executive-summary {
            background-color: #fff;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            margin-bottom: 40px;
            padding: 40px 30px;
            text-align: center;
        }
        
        .executive-summary h2 {
            color: #04062c;
            margin: 0 0 30px 0;
            font-size: 28px;
            font-weight: 600;
            text-align: center;
        }
        
        .metrics-dashboard {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
            gap: 25px;
            margin-top: 30px;
            max-width: 800px;
            margin-left: auto;
            margin-right: auto;
        }
        
        .metric-card {
            background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
            border-radius: 12px;
            padding: 30px 20px;
            text-align: center;
            border-left: 4px solid #04062c;
            transition: transform 0.2s ease, box-shadow 0.2s ease;
            position: relative;
            overflow: hidden;
        }
        
        .metric-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 2px;
            background: linear-gradient(90deg, #04062c, #2c5aa0);
        }
        
        .metric-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 12px rgba(0,0,0,0.15);
        }
        
        .metric-value {
            font-size: 42px;
            font-weight: 700;
            color: #04062c;
            margin-bottom: 12px;
            line-height: 1;
            letter-spacing: -1px;
        }
        
        .metric-label {
            font-size: 13px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 1px;
            font-weight: 600;
            margin-top: 8px;
        }
        
        .metric-card.positive .metric-value {
            color: #2ecc71;
        }
        
        .metric-card.positive {
            border-left-color: #2ecc71;
        }
        
        .metric-card.positive::before {
            background: linear-gradient(90deg, #2ecc71, #27ae60);
        }
        
        /* Location card styles */
        .location-cards {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
            gap: 20px;
        }
        
        .location-card {
            background-color: #fff;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        
        .location-header {
            background-color: #030d54;
            color: white;
            padding: 15px 20px;
        }
        
        .location-header h3 {
            margin: 0;
            font-size: 18px;
        }
        
        .location-address {
            color: rgba(255,255,255,0.8);
            font-size: 14px;
            margin: 5px 0 0 0;
        }
        
        .location-content {
            padding: 20px;
        }
        
        /* Rating styles */
        .rating-container {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
        }
        
        .average-rating {
            font-size: 36px;
            font-weight: bold;
            color: #f39c12;
            margin-right: 15px;
        }
        
        .rating-stars {
            font-size: 20px;
            color: #f39c12;
            font-family: Arial, sans-serif;
            letter-spacing: 2px;
        }
        
        /* Star elements for PDF compatibility */
        .star {
            display: inline-block;
            width: 20px;
            height: 20px;
            margin-right: 2px;
            background-color: #f39c12;
            clip-path: polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%);
            -webkit-clip-path: polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%);
        }
        
        .star.empty {
            background-color: transparent;
            border: 2px solid #f39c12;
        }
        
        /* Sentiment analysis styles */
        .sentiment-analysis {
            margin-bottom: 20px;
        }
        
        .sentiment-pie-container {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 15px;
            margin-top: 10px;
        }
        
        .sentiment-pie-chart {
            flex-shrink: 0;
        }
        
        .sentiment-legend {
            display: flex;
            gap: 20px;
            justify-content: center;
            flex-wrap: wrap;
        }
        
        .sentiment-legend-item {
            display: flex;
            align-items: center;
            margin-bottom: 8px;
            font-size: 14px;
        }
        
        .sentiment-legend-color {
            width: 16px;
            height: 16px;
            border-radius: 3px;
            margin-right: 10px;
            flex-shrink: 0;
        }
        
        .sentiment-legend-color.positive {
            background-color: #2ecc71;
        }
        
        .sentiment-legend-color.neutral {
            background-color: #3498db;
        }
        
        .sentiment-legend-color.negative {
            background-color: #e74c3c;
        }
        
        .sentiment-legend-label {
            flex: 1;
            display: flex;
            justify-content: space-between;
        }
        
        /* Legacy bar styles - kept for backward compatibility */
        .sentiment-bars {
            margin-top: 10px;
        }
        
        .sentiment-bar-container {
            height: 20px;
            background-color: #ecf0f1;
            border-radius: 10px;
            margin-bottom: 10px;
            overflow: hidden;
        }
        
        .sentiment-bar {
            height: 100%;
            border-radius: 10px;
        }
        
        .sentiment-label {
            display: flex;
            justify-content: space-between;
            margin-bottom: 5px;
            font-size: 14px;
        }
        
        .positive-bar {
            background-color: #2ecc71;
        }
        
        .neutral-bar {
            background-color: #3498db;
        }
        
        .negative-bar {
            background-color: #e74c3c;
        }
        
        /* Themes styles */
        .themes-container {
            margin-bottom: 20px;
        }
        
        .themes {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin-top: 10px;
        }
        
        .theme {
            background-color: #f1f1f1;
            border-radius: 15px;
            padding: 5px 12px;
            font-size: 13px;
            color: #555;
        }
        
        .theme.positive {
            background-color: rgba(46, 204, 113, 0.2);
            color: #27ae60;
        }
        
        .theme.negative {
            background-color: rgba(231, 76, 60, 0.2);
            color: #c0392b;
        }
        
        /* Insights styles */
        .insights-container {
            margin-bottom: 20px;
        }
        
        .insight {
            background-color: #f8f9fa;
            border-left: 4px solid #3498db;
            padding: 10px 15px;
            margin-bottom: 10px;
        }
        
        .insight.negative {
            border-left-color: #e74c3c;
        }
        
        .insight-category {
            font-weight: bold;
            margin-bottom: 5px;
        }
        
        .insight-example {
            font-style: italic;
            color: #7f8c8d;
            font-size: 13px;
            margin-top: 5px;
            padding-left: 10px;
            border-left: 2px solid #ddd;
        }
        
        /* Recommendations styles */
        .recommendations-container {
            margin-bottom: 20px;
        }
        
        .recommendations-list {
            padding-left: 20px;
        }
        
        .recommendations-list li {
            margin-bottom: 8px;
        }
        
        /* Training recommendations styles */
        .training-recommendations-container {
            margin-bottom: 20px;
        }
        
        .training-subsection {
            margin-bottom: 15px;
        }
        
        .training-subsection:last-child {
            margin-bottom: 0;
        }
        
        .training-subsection-title {
            color: #34495e;
            font-size: 14px;
            font-weight: 600;
            margin-bottom: 8px;
            margin-top: 0;
        }
        
        .training-list {
            padding-left: 20px;
            margin: 0;
        }
        
        .training-list li {
            margin-bottom: 6px;
            font-size: 14px;
            line-height: 1.4;
        }
        
        /* Helper classes */
        .section-title {
            color: #2c3e50;
            font-size: 16px;
            margin-bottom: 10px;
            border-bottom: 1px solid #eee;
            padding-bottom: 5px;
        }
        
        /* Screen reader only class for accessibility */
        .sr-only {
            position: absolute;
            width: 1px;
            height: 1px;
            padding: 0;
            margin: -1px;
            overflow: hidden;
            clip: rect(0, 0, 0, 0);
            white-space: nowrap;
            border-width: 0;
        }
        
        /* Responsive adjustments */
        @media (max-width: 768px) {
            .location-cards {
                grid-template-columns: 1fr;
            }
            
            .report-header {
                padding: 20px;
            }
            
            .header-content {
                flex-direction: column;
                align-items: flex-start;
            }
            
            .header-meta {
                text-align: left;
                margin-top: 15px;
            }
            
            .average-rating {
                font-size: 28px;
            }
        }
        
        /* Print styles for A4 format */
        @media print {
            @page {
                size: A4 portrait;
                margin: 1cm;
            }
            
            body {
                background-color: white;
                color: black;
            }
            
            .container {
                max-width: 100%;
                padding: 0;
                margin: 0 auto;
                text-align: center;
                padding-top: 6cm;
            }
            
            .report-header {
                background-color: #04062c !important;
                -webkit-print-color-adjust: exact !important;
                color-adjust: exact !important;
                page-break-inside: avoid;
                width: 100%;
                max-width: 21cm; /* A4 width minus margins */
                margin: 0 auto 1cm auto !important;
                text-align: left;
            }
            
            .executive-summary {
                background-color: white !important;
                border: 1px solid #ddd;
                page-break-inside: avoid !important;
                page-break-after: always !important;
                margin-bottom: 0 !important;
                display: block !important;
            }
            
            .main-content {
                display: block !important;
            }
            
            .location-cards {
                page-break-before: always !important;
                display: block !important;
                width: 100%;
                max-width: 21cm; /* A4 width minus margins */
                margin: 0 auto;
                text-align: left;
            }
            
            .location-card {
                break-inside: auto !important;
                page-break-inside: auto !important;
                display: block !important;
                margin-bottom: 20px;
                width: 100% !important;
            }
            
            .location-card + .location-card {
                break-before: page !important;
                page-break-before: always !important;
                margin-top: 1cm !important;
            }
            
            .location-cards .location-card:last-child {
                page-break-after: auto !important;
            }
            
            .location-header {
                background-color: #030d54 !important;
                -webkit-print-color-adjust: exact;
                print-color-adjust: exact;
                color-adjust: exact;
            }
            
            .sentiment-bar {
                -webkit-print-color-adjust: exact;
                print-color-adjust: exact;
                color-adjust: exact;
            }
            
            .insight, .section-title {
                break-inside: avoid !important;
                page-break-inside: avoid !important;
            }
            
            .insight {
                page-break-after: auto !important;
            }
            
            .section-title {
                page-break-after: avoid !important;
                page-break-before: auto !important;
            }
            
            .insights-container, .recommendations-container, .themes-container, .sentiment-analysis {
                break-inside: auto !important;
                page-break-inside: auto !important;
                page-break-after: auto !important;
                margin-bottom: 25px !important;
            }
            
            .training-recommendations-container {
                break-inside: avoid !important;
                page-break-inside: avoid !important;
                page-break-after: auto !important;
                margin-bottom: 25px !important;
            }
            
            .negative-breakdown-container {
                break-inside: avoid !important;
                page-break-inside: avoid !important;
                page-break-after: auto !important;
                margin-bottom: 25px !important;
            }
            
            .rating-stars {
                font-family: "DejaVu Sans", Arial, "Helvetica Neue", sans-serif;
                font-weight: normal;
            }
            
            .metrics-dashboard {
                display: grid !important;
                grid-template-columns: repeat(3, 1fr) !important;
                gap: 15px !important;
            }
            
            .metric-card {
                background-color: #f8f9fa !important;
                border: 1px solid #ddd !important;
                -webkit-print-color-adjust: exact !important;
                color-adjust: exact !important;
            }
        }
        
        /* Themes styles */
        .themes-container {
            margin-bottom: 20px;
        }
        
        .themes {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin-top: 10px;
        }
        
        .theme {
            background-color: #f1f1f1;
            border-radius: 15px;
            padding: 5px 12px;
            font-size: 13px;
            color: #555;
        }
        
        .theme.positive {
            background-color: rgba(46, 204, 113, 0.2);
            color: #27ae60;
        }
        
        .theme.negative {
            background-color: rgba(231, 76, 60, 0.2);
            color: #c0392b;
        }
        
        /* Negative review breakdown styles */
        .negative-breakdown-container {
            margin-bottom: 20px;
        }
        
        .negative-breakdown-chart-container {
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 15px;
            margin-top: 10px;
        }
        
        .negative-breakdown-pie-chart {
            flex-shrink: 0;
        }
        
        .negative-breakdown-legend {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 8px;
            width: 100%;
            max-width: 400px;
        }
        
        .negative-breakdown-legend-item {
            display: flex;
            align-items: center;
            font-size: 12px;
        }
        
        .negative-breakdown-legend-color {
            width: 14px;
            height: 14px;
            border-radius: 3px;
            margin-right: 8px;
            flex-shrink: 0;
        }
        
        .negative-breakdown-legend-label {
            flex: 1;
            display: flex;
            justify-content: space-between;
        }
        
        /* Main content wrapper */
        .main-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            justify-content: flex-start;
            align-items: stretch;
            margin-top: 20px;
        }
    </style>
    <script>
        // Formatting helper functions
        function formatDate(dateString) {
            try {
                // Handle Go time format which may include quotes and have a different format
                if (dateString) {
                    // Remove any quotes if present
                    dateString = dateString.replace(/"/g, '');
                    
                    // Try to parse with different formats
                    let date;
                    
                    // Try standard ISO format first
                    date = new Date(dateString);
                    
                    // Check if date is valid
                    if (isNaN(date.getTime())) {
                        // Try parsing Go's default format: 2006-01-02 15:04:05 -0700 MST
                        const goParts = dateString.match(/(\d{4})-(\d{2})-(\d{2})[ T](\d{2}):(\d{2}):(\d{2})/);
                        if (goParts) {
                            date = new Date(
                                parseInt(goParts[1]),   // year
                                parseInt(goParts[2])-1, // month (0-based)
                                parseInt(goParts[3]),   // day
                                parseInt(goParts[4]),   // hour
                                parseInt(goParts[5]),   // minute
                                parseInt(goParts[6])    // second
                            );
                        }
                    }
                    
                    // Format the date if it's valid
                    if (!isNaN(date.getTime())) {
                        const options = { year: 'numeric', month: 'long', day: 'numeric' };
                        return date.toLocaleDateString('en-US', options);
                    }
                }
                
                // Fallback for any other issues
                return dateString;
            } catch (e) {
                console.error("Error formatting date:", e);
                return dateString;
            }
        }
        
        function formatPercentage(value) {
            return parseFloat(value).toFixed(1) + '%';
        }
        
        function calculateExecutiveMetrics() {
            let totalReviews = 0;
            let totalRatingSum = 0;
            let totalPositiveReviews = 0;
            
            // Extract data from existing DOM elements instead of raw JSON
            const locationCards = document.querySelectorAll('.location-card');
            
            console.log("Found", locationCards.length, "location cards");
            
            locationCards.forEach(function(card) {
                // Get reviews count from sentiment analysis text
                const sentimentText = card.querySelector('.sentiment-analysis div[style*="font-size: 13px"]');
                if (sentimentText) {
                    const reviewsMatch = sentimentText.textContent.match(/Based on (\d+) reviews/);
                    if (reviewsMatch) {
                        const reviews = parseInt(reviewsMatch[1]);
                        totalReviews += reviews;
                        
                        // Get rating from the rating display
                        const ratingElement = card.querySelector('.average-rating');
                        if (ratingElement) {
                            const rating = parseFloat(ratingElement.textContent);
                            totalRatingSum += (rating * reviews);
                            
                            console.log("Location - Reviews:", reviews, "Rating:", rating);
                        }
                        
                        // Get positive count from sentiment pie chart data
                        const pieChart = card.querySelector('[data-format="sentiment-pie"]');
                        if (pieChart) {
                            const positivePercentage = parseFloat(pieChart.getAttribute('data-positive')) || 0;
                            const positiveCount = Math.round((positivePercentage / 100) * reviews);
                            totalPositiveReviews += positiveCount;
                            
                            console.log("Positive percentage:", positivePercentage, "Positive count:", positiveCount);
                        }
                    }
                }
            });
            
            // Calculate weighted average rating
            const overallRating = totalReviews > 0 ? (totalRatingSum / totalReviews) : 0;
            
            // Calculate overall positive sentiment percentage
            const overallPositivePercentage = totalReviews > 0 ? (totalPositiveReviews / totalReviews * 100) : 0;
            
            console.log("Final metrics:", {
                totalReviews: totalReviews,
                overallRating: overallRating,
                positivePercentage: overallPositivePercentage
            });
            
            return {
                totalReviews: totalReviews,
                overallRating: overallRating,
                positivePercentage: overallPositivePercentage
            };
        }
        
        function populateExecutiveMetrics() {
            const metrics = calculateExecutiveMetrics();
            
            // Update total reviews
            const totalReviewsElement = document.getElementById('total-reviews-metric');
            if (totalReviewsElement) {
                totalReviewsElement.textContent = metrics.totalReviews.toLocaleString();
            }
            
            // Update overall rating
            const overallRatingElement = document.getElementById('overall-rating-metric');
            if (overallRatingElement) {
                overallRatingElement.textContent = metrics.overallRating.toFixed(1);
            }
            
            // Update positive sentiment
            const positiveSentimentElement = document.getElementById('positive-sentiment-metric');
            if (positiveSentimentElement) {
                positiveSentimentElement.textContent = Math.round(metrics.positivePercentage) + '%';
            }
        }
        
        function displayStars(rating) {
            // Round up for ratings of X.9 or higher
            if (rating % 1 >= 0.9) {
                rating = Math.ceil(rating);
            }
            const fullStars = Math.floor(rating);
            let stars = '';
            
            // Use styled HTML elements instead of Unicode characters for better PDF compatibility
            for (let i = 0; i < 5; i++) {
                if (i < fullStars) {
                    stars += '<span class="star filled"></span>';
                } else {
                    stars += '<span class="star empty"></span>';
                }
            }
            
            // Add text alternative for accessibility and fallback
            stars += '<span class="sr-only">' + rating + ' out of 5 stars</span>';
            
            return stars;
        }

        function createPieChart(positive, neutral, negative) {
            // Use original colors
            const colors = {
                positive: '#2ecc71',  // Original green
                neutral: '#3498db',   // Keep existing blue
                negative: '#e74c3c'   // Original red
            };
            
            const size = 220;  // Increased to 240px for better visibility
            const centerX = size / 2;
            const centerY = size / 2;
            const radius = 100;  // Increased proportionally (was 65 for 160px)
            
            // Calculate percentages and round to ensure they add up to 100%
            const total = positive + neutral + negative;
            if (total === 0) return '<div>No sentiment data available</div>';
            
            let rawPositive = (positive / total) * 100;
            let rawNeutral = (neutral / total) * 100;  
            let rawNegative = (negative / total) * 100;
            
            // Round to nearest integer
            let roundedPositive = Math.round(rawPositive);
            let roundedNeutral = Math.round(rawNeutral);
            let roundedNegative = Math.round(rawNegative);
            
            // Adjust to ensure sum equals 100%
            let sum = roundedPositive + roundedNeutral + roundedNegative;
            if (sum !== 100) {
                let diff = 100 - sum;
                // Add/subtract difference to largest value
                if (roundedPositive >= roundedNeutral && roundedPositive >= roundedNegative) {
                    roundedPositive += diff;
                } else if (roundedNeutral >= roundedNegative) {
                    roundedNeutral += diff;
                } else {
                    roundedNegative += diff;
                }
            }
            
            let svgContent = '<svg width="' + size + '" height="' + size + '" viewBox="0 0 ' + size + ' ' + size + '">';
            
            // Only create slices for non-zero values
            let slices = [];
            if (roundedPositive > 0) slices.push({label: 'Positive', value: roundedPositive, color: colors.positive});
            if (roundedNeutral > 0) slices.push({label: 'Neutral', value: roundedNeutral, color: colors.neutral});
            if (roundedNegative > 0) slices.push({label: 'Negative', value: roundedNegative, color: colors.negative});
            
            if (slices.length === 0) {
                // No data to display
                svgContent += '<circle cx="' + centerX + '" cy="' + centerY + '" r="' + radius + '" fill="#f0f0f0" stroke="#ddd" stroke-width="2"/>';
                svgContent += '<text x="' + centerX + '" y="' + centerY + '" text-anchor="middle" dominant-baseline="middle" fill="#999" font-size="12">No data</text>';
            } else if (slices.length === 1) {
                // Single slice - draw full circle
                const slice = slices[0];
                svgContent += '<circle cx="' + centerX + '" cy="' + centerY + '" r="' + radius + '" fill="' + slice.color + '"/>';
                svgContent += '<text x="' + centerX + '" y="' + centerY + '" text-anchor="middle" dominant-baseline="middle" fill="white" font-size="14" font-weight="bold">' + slice.value + '%</text>';
            } else {
                // Multiple slices - draw pie chart
                let currentAngle = -Math.PI / 2; // Start at top (12 o'clock)
                
                for (let i = 0; i < slices.length; i++) {
                    const slice = slices[i];
                    const sliceAngle = (slice.value / 100) * 2 * Math.PI;
                    const endAngle = currentAngle + sliceAngle;
                    
                    const x1 = centerX + radius * Math.cos(currentAngle);
                    const y1 = centerY + radius * Math.sin(currentAngle);
                    const x2 = centerX + radius * Math.cos(endAngle);
                    const y2 = centerY + radius * Math.sin(endAngle);
                    
                    const largeArcFlag = sliceAngle > Math.PI ? 1 : 0;
                    
                    // Create path for slice
                    svgContent += '<path d="M ' + centerX + ' ' + centerY + ' L ' + x1 + ' ' + y1 + ' A ' + radius + ' ' + radius + ' 0 ' + largeArcFlag + ' 1 ' + x2 + ' ' + y2 + ' Z" fill="' + slice.color + '"/>';
                    
                    // Add percentage label inside slice if there's enough space
                    if (slice.value >= 10) {
                        const labelAngle = currentAngle + sliceAngle / 2;
                        const labelRadius = radius * 0.6;
                        const labelX = centerX + labelRadius * Math.cos(labelAngle);
                        const labelY = centerY + labelRadius * Math.sin(labelAngle);
                        
                        svgContent += '<text x="' + labelX + '" y="' + labelY + '" text-anchor="middle" dominant-baseline="middle" fill="white" font-size="12" font-weight="bold">' + slice.value + '%</text>';
                    }
                    
                    currentAngle = endAngle;
                }
            }
            
            svgContent += '</svg>';
            return svgContent;
        }

        function createNegativeBreakdownPieChart(categories) {
            // Use distinct colors for different negative review categories
            const categoryColors = [
                '#FF6B6B', // Red-pink
                '#4ECDC4', // Teal
                '#45B7D1', // Light blue
                '#96CEB4', // Mint green
                '#FFEAA7', // Light yellow
                '#DDA0DD', // Plum
                '#F4A460', // Sandy brown
                '#98D8C8'  // Mint
            ];
            
            const size = 200;  // Slightly smaller than sentiment chart
            const centerX = size / 2;
            const centerY = size / 2;
            const radius = 80;
            
            if (!categories || categories.length === 0) {
                return '<div style="text-align: center; color: #999; padding: 20px;">No negative review categories available</div>';
            }
            
            let svgContent = '<svg width="' + size + '" height="' + size + '" viewBox="0 0 ' + size + ' ' + size + '">';
            
            // Calculate total and rounded percentages similar to sentiment chart
            let totalPercentage = 0;
            categories.forEach(function(category) {
                totalPercentage += category.percentage;
            });
            
            // Create slices with proper rounding
            let slices = [];
            let runningTotal = 0;
            
            for (let i = 0; i < categories.length; i++) {
                const category = categories[i];
                if (category.percentage > 0) {
                    // Calculate the rounded percentage that will be displayed
                    let displayPercentage = Math.round(category.percentage);
                    
                    // For the last slice, adjust to ensure total equals 100%
                    if (i === categories.length - 1) {
                        displayPercentage = 100 - runningTotal;
                    }
                    
                    slices.push({
                        label: category.name,
                        displayValue: displayPercentage, // What shows in legend and chart
                        originalValue: category.percentage, // Original value for angle calculation
                        color: categoryColors[i % categoryColors.length],
                        count: category.count
                    });
                    
                    runningTotal += displayPercentage;
                }
            }
            
            if (slices.length === 0) {
                svgContent += '<circle cx="' + centerX + '" cy="' + centerY + '" r="' + radius + '" fill="#f0f0f0" stroke="#ddd" stroke-width="2"/>';
                svgContent += '<text x="' + centerX + '" y="' + centerY + '" text-anchor="middle" dominant-baseline="middle" fill="#999" font-size="12">No data</text>';
            } else if (slices.length === 1) {
                // Single slice - draw full circle
                const slice = slices[0];
                svgContent += '<circle cx="' + centerX + '" cy="' + centerY + '" r="' + radius + '" fill="' + slice.color + '"/>';
                svgContent += '<text x="' + centerX + '" y="' + centerY + '" text-anchor="middle" dominant-baseline="middle" fill="white" font-size="12" font-weight="bold">' + slice.displayValue + '%</text>';
            } else {
                // Multiple slices - draw pie chart using original values for angles
                let currentAngle = -Math.PI / 2; // Start at top (12 o'clock)
                
                for (let i = 0; i < slices.length; i++) {
                    const slice = slices[i];
                    const sliceAngle = (slice.originalValue / totalPercentage) * 2 * Math.PI;
                    const endAngle = currentAngle + sliceAngle;
                    
                    const x1 = centerX + radius * Math.cos(currentAngle);
                    const y1 = centerY + radius * Math.sin(currentAngle);
                    const x2 = centerX + radius * Math.cos(endAngle);
                    const y2 = centerY + radius * Math.sin(endAngle);
                    
                    const largeArcFlag = sliceAngle > Math.PI ? 1 : 0;
                    
                    // Create path for slice
                    svgContent += '<path d="M ' + centerX + ' ' + centerY + ' L ' + x1 + ' ' + y1 + ' A ' + radius + ' ' + radius + ' 0 ' + largeArcFlag + ' 1 ' + x2 + ' ' + y2 + ' Z" fill="' + slice.color + '"/>';
                    
                    // Add percentage label inside slice if there's enough space - use display value
                    if (slice.displayValue >= 8) {
                        const labelAngle = currentAngle + sliceAngle / 2;
                        const labelRadius = radius * 0.6;
                        const labelX = centerX + labelRadius * Math.cos(labelAngle);
                        const labelY = centerY + labelRadius * Math.sin(labelAngle);
                        
                        svgContent += '<text x="' + labelX + '" y="' + labelY + '" text-anchor="middle" dominant-baseline="middle" fill="white" font-size="10" font-weight="bold">' + slice.displayValue + '%</text>';
                    }
                    
                    currentAngle = endAngle;
                }
            }
            
            svgContent += '</svg>';
            return svgContent;
        }

        // Will run when the DOM is fully loaded
        document.addEventListener('DOMContentLoaded', function() {
            // Format all dates
            document.querySelectorAll('[data-format="date"]').forEach(function(element) {
                element.textContent = formatDate(element.getAttribute('data-value'));
            });
            
            // Format all percentages
            document.querySelectorAll('[data-format="percentage"]').forEach(function(element) {
                element.textContent = formatPercentage(element.getAttribute('data-value'));
            });
            
            // Display all star ratings - use innerHTML for proper rendering of HTML entities
            document.querySelectorAll('[data-format="stars"]').forEach(function(element) {
                element.innerHTML = displayStars(parseFloat(element.getAttribute('data-value')));
            });

            // Generate sentiment pie charts
            document.querySelectorAll('[data-format="sentiment-pie"]').forEach(function(element) {
                const positive = parseFloat(element.getAttribute('data-positive'));
                const neutral = parseFloat(element.getAttribute('data-neutral'));
                const negative = parseFloat(element.getAttribute('data-negative'));
                
                element.innerHTML = createPieChart(positive, neutral, negative);
            });

            // Generate negative breakdown pie charts
            document.querySelectorAll('[data-format="negative-breakdown-pie"]').forEach(function(element) {
                const categoriesData = element.getAttribute('data-categories');
                if (categoriesData && categoriesData !== '') {
                    try {
                        const categories = JSON.parse(categoriesData);
                        element.innerHTML = createNegativeBreakdownPieChart(categories);
                        
                        // Update legend percentages to match chart rounding
                        const container = element.closest('.negative-breakdown-container');
                        if (container && categories.length > 0) {
                            const legendItems = container.querySelectorAll('.negative-breakdown-percentage');
                            let runningTotal = 0;
                            
                            categories.forEach(function(category, index) {
                                if (index < legendItems.length) {
                                    let displayPercentage = Math.round(category.percentage);
                                    
                                    // For the last item, adjust to ensure total equals 100%
                                    if (index === categories.length - 1) {
                                        displayPercentage = 100 - runningTotal;
                                    }
                                    
                                    legendItems[index].textContent = displayPercentage + '%';
                                    runningTotal += displayPercentage;
                                }
                            });
                        }
                    } catch (e) {
                        element.innerHTML = '<div style="text-align: center; color: #999; padding: 20px;">Error loading negative review data</div>';
                    }
                } else {
                    element.innerHTML = '<div style="text-align: center; color: #999; padding: 20px;">No negative review categories available</div>';
                }
            });

            // Populate executive metrics
            populateExecutiveMetrics();
        });
    </script>
</head>
<body>
    <div class="container">
        <div class="report-header">
            <div class="header-content">
                <div class="header-title">
                    <h1>{{.ClientName}}</h1>
                    <img src="data:image/png;base64,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" alt="Google Logo" class="logo-image" >

                </div>
                <div class="header-meta">
                    <p>Period: <span data-format="date" data-value="{{.PeriodStart.Format "2006-01-02T15:04:05Z07:00"}}"></span> - <span data-format="date" data-value="{{.PeriodEnd.Format "2006-01-02T15:04:05Z07:00"}}"></span></p>
                    <p>Generated: <span data-format="date" data-value="{{.GeneratedAt.Format "2006-01-02T15:04:05Z07:00"}}"></span></p>
                </div>
            </div>
        </div>
        
        <div class="main-content">
            <div class="executive-summary">
                <h2>Executive Summary</h2>
                <div class="metrics-dashboard">
                    <div class="metric-card">
                        <div class="metric-value" id="total-reviews-metric">0</div>
                        <div class="metric-label">Total Reviews Analyzed</div>
                    </div>
                    <div class="metric-card">
                        <div class="metric-value" id="overall-rating-metric">0.0</div>
                        <div class="metric-label">Overall Average Rating</div>
                    </div>
                    <div class="metric-card positive">
                        <div class="metric-value" id="positive-sentiment-metric">0%</div>
                        <div class="metric-label">Positive Sentiment</div>
                    </div>
                </div>
            </div>
            
            <div class="location-cards">
                {{range .LocationResults}}
                <div class="location-card">
                    <div class="location-header">
                        <h3>{{.Metadata.LocationName}}</h3>
                    </div>
                    
                    <div class="location-content">
                        <!-- Rating section -->
                        <div class="rating-container">
                            <div class="average-rating">{{printf "%.1f" .Analysis.OverallSummary.AverageRating}}</div>
                            <div class="rating-stars" data-format="stars" data-value="{{.Analysis.OverallSummary.AverageRating}}"></div>
                        </div>
                        
                        <!-- Sentiment analysis section -->
                        <div class="sentiment-analysis">
                            <h4 class="section-title">Sentiment Analysis</h4>
                            <div class="sentiment-pie-container">
                                <div class="sentiment-pie-chart" data-format="sentiment-pie" data-positive="{{.Analysis.SentimentAnalysis.PositivePercentage}}" data-neutral="{{.Analysis.SentimentAnalysis.NeutralPercentage}}" data-negative="{{.Analysis.SentimentAnalysis.NegativePercentage}}"></div>
                                <div class="sentiment-legend">
                                    <div class="sentiment-legend-item">
                                        <div class="sentiment-legend-color positive"></div>
                                        <div class="sentiment-legend-label">Positive</div>
                                    </div>
                                    <div class="sentiment-legend-item">
                                        <div class="sentiment-legend-color neutral"></div>
                                        <div class="sentiment-legend-label">Neutral</div>
                                    </div>
                                    <div class="sentiment-legend-item">
                                        <div class="sentiment-legend-color negative"></div>
                                        <div class="sentiment-legend-label">Negative</div>
                                    </div>
                                </div>
                            </div>
                            <div style="font-size: 13px; margin-top: 10px; color: #555;">
                                Based on {{.Analysis.SentimentAnalysis.TotalReviews}} reviews | Trend: {{.Analysis.SentimentAnalysis.SentimentTrend}}
                            </div>
                        </div>
                        
                        <!-- Themes section -->
                        <div class="themes-container">
                            <h4 class="section-title">Key Themes</h4>
                            
                            <div class="themes">
                                {{range .Analysis.OverallSummary.PositiveThemes}}
                                <div class="theme positive">{{.}}</div>
                                {{end}}
                                
                                {{range .Analysis.OverallSummary.NegativeThemes}}
                                <div class="theme negative">{{.}}</div>
                                {{end}}
                            </div>
                        </div>
                        
                        <!-- Key insights section -->
                        <div class="insights-container">
                            <h4 class="section-title">Key Strengths</h4>
                            {{range .Analysis.KeyTakeaways.Strengths}}
                            <div class="insight">
                                <div class="insight-category">{{.Category}}</div>
                                <div>{{.Description}}</div>
                                <div class="insight-example">"{{.Example}}"</div>
                            </div>
                            {{end}}
                            
                            <h4 class="section-title">Areas for Improvement</h4>
                            {{range .Analysis.KeyTakeaways.AreasForImprovement}}
                            <div class="insight negative">
                                <div class="insight-category">{{.Category}}</div>
                                <div>{{.Description}}</div>
                                <div class="insight-example">"{{.Example}}"</div>
                            </div>
                            {{end}}
                        </div>
                        
                        <!-- Training recommendations section -->
                        {{if or .Analysis.TrainingRecommendations.ForOperators .Analysis.TrainingRecommendations.ForDrivers}}
                        <div class="training-recommendations-container">
                            <h4 class="section-title">Training Recommendations</h4>
                            
                            {{if .Analysis.TrainingRecommendations.ForOperators}}
                            <div class="training-subsection">
                                <h5 class="training-subsection-title">For Operators</h5>
                                <ul class="training-list">
                                    {{range .Analysis.TrainingRecommendations.ForOperators}}
                                    <li>{{.}}</li>
                                    {{end}}
                                </ul>
                            </div>
                            {{end}}
                            
                            {{if .Analysis.TrainingRecommendations.ForDrivers}}
                            <div class="training-subsection">
                                <h5 class="training-subsection-title">For Drivers</h5>
                                <ul class="training-list">
                                    {{range .Analysis.TrainingRecommendations.ForDrivers}}
                                    <li>{{.}}</li>
                                    {{end}}
                                </ul>
                            </div>
                            {{end}}
                        </div>
                        {{end}}
                        
                        <!-- Negative review breakdown section -->
                        {{if .Analysis.NegativeReviewBreakdown.Categories}}
                        <div class="negative-breakdown-container">
                            <h4 class="section-title">Negative Review Breakdown</h4>
                            <div class="negative-breakdown-chart-container">
                                <div class="negative-breakdown-pie-chart" data-format="negative-breakdown-pie" data-categories="[{{range $i, $cat := .Analysis.NegativeReviewBreakdown.Categories}}{{if $i}},{{end}}{&quot;name&quot;:&quot;{{$cat.Name}}&quot;,&quot;count&quot;:{{$cat.Count}},&quot;percentage&quot;:{{printf "%.1f" $cat.Percentage}}}{{end}}]"></div>
                                <div class="negative-breakdown-legend">
                                    {{range $i, $cat := .Analysis.NegativeReviewBreakdown.Categories}}
                                    <div class="negative-breakdown-legend-item">
                                        <div class="negative-breakdown-legend-color" style="background-color: {{if eq $i 0}}#FF6B6B{{else if eq $i 1}}#4ECDC4{{else if eq $i 2}}#45B7D1{{else if eq $i 3}}#96CEB4{{else if eq $i 4}}#FFEAA7{{else if eq $i 5}}#DDA0DD{{else if eq $i 6}}#F4A460{{else}}#98D8C8{{end}};"></div>
                                        <div class="negative-breakdown-legend-label">
                                            <span>{{$cat.Name}}</span>
                                            <span class="negative-breakdown-percentage">{{printf "%.0f" $cat.Percentage}}%</span>
                                        </div>
                                    </div>
                                    {{end}}
                                </div>
                            </div>
                        </div>
                        {{end}}
                        
                        <!-- Recommendations section -->
                        <div class="recommendations-container">
                            <h4 class="section-title">Suggested Improvements</h4>
                            <ul class="recommendations-list">
                                {{range .Analysis.NegativeReviewBreakdown.ImprovementRecommendations}}
                                <li>{{.}}</li>
                                {{end}}
                            </ul>
                        </div>
                    </div>
                </div>
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>
`
