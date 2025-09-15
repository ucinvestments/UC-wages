# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

UC Wages Data Scraper - A Python-based tool for collecting and processing University of California employee wage data from the UC Annual Wage website (https://ucannualwage.ucop.edu/wage/).

## Purpose

This subdirectory is part of the larger UC-Investments project and focuses on:
- Scraping UC employee wage data from the official UC wage transparency portal
- Processing and organizing wage data by location and year
- Supporting analysis of UC compensation trends across different campuses

## Architecture

### Main Components
- **scrape.py**: Core scraping logic for retrieving wage data from the UC Annual Wage API
  - Defines locations (all UC campuses including ASUCLA, UCOP)
  - Supports data collection for years 2010-2024
  - API endpoint: https://ucannualwage.ucop.edu/wage/search

### Data Structure
- **Berkeley/**: Directory for Berkeley campus wage data (and similar directories for other campuses)
- Data is organized by location/campus for efficient retrieval and analysis

## API Details

The UC Annual Wage API expects POST requests with the following JSON structure:
```json
{
  "op": "search",
  "page": 1,
  "rows": 20,
  "sidx": "lastname",
  "sord": "asc",
  "count": 0,
  "year": "2024",
  "firstname": "",
  "location": "Berkeley",
  "lastname": "",
  "title": "",
  "startSal": "",
  "endSal": ""
}
```

## Commands

### Running the Scraper
```bash
python scrape.py                    # Execute wage data scraping
```

### Dependencies
```bash
pip install requests                 # HTTP library for API calls (if not already installed)
```

## Data Collection Workflow

1. **Location-based Collection**: Script iterates through all UC locations
2. **Year-based Collection**: For each location, collects data from 2010-2024
3. **Pagination Handling**: Manages API pagination (20 records per page by default)
4. **Data Storage**: Saves collected data in location-specific directories

## Available UC Locations

- ALL (system-wide data)
- ASUCLA
- Berkeley
- Davis
- UC SF Law
- Irvine
- Los Angeles
- Merced
- Riverside
- San Diego
- San Francisco
- Santa Barbara
- Santa Cruz
- UCOP (UC Office of the President)

## Development Notes

- The API uses POST requests with JSON payloads
- Results are paginated (default 20 rows per page)
- Data can be sorted by various fields (lastname, firstname, salary, etc.)
- Search supports filtering by name, title, and salary range
- Consider rate limiting to avoid overwhelming the UC server