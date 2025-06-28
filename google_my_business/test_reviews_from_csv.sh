#!/bin/bash

# Simple script to test AI responses on a CSV file of reviews

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get CSV file from command line argument
CSV_FILE=${1:-"Drive Doncaster_2091990859.csv"}

# Check if the CSV file exists
if [ ! -f "$CSV_FILE" ]; then
  echo -e "${RED}Error: CSV file '$CSV_FILE' not found${NC}"
  echo "Usage: $0 [csv_file]"
  exit 1
fi

# Check if OpenAI API key is set
if [ -z "$OPENAI_API_KEY" ]; then
  echo -e "${RED}Error: OPENAI_API_KEY environment variable is not set${NC}"
  echo "Please set it with: export OPENAI_API_KEY=your_api_key"
  exit 1
fi

# Generate output filename
OUTPUT_FILE="${CSV_FILE%.*}_with_responses.csv"

echo -e "${BLUE}Testing AI Responses on CSV File${NC}"
echo "========================================"
echo -e "Input CSV: ${YELLOW}$CSV_FILE${NC}"
echo -e "Output CSV: ${YELLOW}$OUTPUT_FILE${NC}"
echo -e "OpenAI Model: ${YELLOW}${OPENAI_MODEL:-gpt-3.5-turbo}${NC}"
echo "========================================"

# Build the test tool
echo -e "${YELLOW}Building test tool...${NC}"
go build -tags=test_ai_responses -o test_ai_responses test_ai_responses.go

if [ $? -ne 0 ]; then
  echo -e "${RED}Build failed!${NC}"
  exit 1
fi

# Run the test tool with the CSV file
echo -e "${YELLOW}Processing reviews...${NC}"

# Pass the OpenAI model as an environment variable if specified
if [ ! -z "$OPENAI_MODEL" ]; then
  echo -e "Using model: ${YELLOW}$OPENAI_MODEL${NC}"
  OPENAI_MODEL=$OPENAI_MODEL ./test_ai_responses "$CSV_FILE"
else
  echo -e "Using default model: ${YELLOW}gpt-3.5-turbo${NC}"
  ./test_ai_responses "$CSV_FILE"
fi

if [ $? -ne 0 ]; then
  echo -e "${RED}Processing failed!${NC}"
  exit 1
fi

echo -e "${GREEN}Processing completed successfully!${NC}"
echo -e "Results saved to: ${YELLOW}$OUTPUT_FILE${NC}"

# Try to open the output file if possible
if command -v open &> /dev/null; then
  echo -e "${YELLOW}Opening results file...${NC}"
  open "$OUTPUT_FILE"
elif command -v xdg-open &> /dev/null; then
  echo -e "${YELLOW}Opening results file...${NC}"
  xdg-open "$OUTPUT_FILE"
fi 