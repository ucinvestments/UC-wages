# UC Wages Data Analysis

A comprehensive dataset of University of California employee compensation from 2010-2024, covering all UC campuses and affiliates.

## Dataset Overview

This repository contains public wage data for approximately 200,000+ UC employees across:
- 13 UC locations (10 campuses + UCOP, ASUCLA, UC SF Law)
- 15 years of historical data (2010-2024)
- Complete compensation records including base pay, overtime, and other pay

## Data Structure

```
data/
├── Berkeley/
│   ├── wages_2024.json
│   ├── wages_2023.json
│   └── ...
├── Los_Angeles/
├── San_Diego/
└── ... (all UC locations)
```

Each JSON file contains:
- Employee names and titles
- Base salary information
- Overtime and other pay
- Total compensation
- Location and year

## Quick Start

### Scraping Fresh Data
```bash
docker compose up
```

### Data Access
```python
import json
import pandas as pd

# Load single year
with open('data/Berkeley/wages_2024.json') as f:
    data = json.load(f)
    df = pd.DataFrame(data['records'])

# Load all years for analysis
dfs = []
for year in range(2010, 2025):
    with open(f'data/Berkeley/wages_{year}.json') as f:
        data = json.load(f)
        df = pd.DataFrame(data['records'])
        df['year'] = year
        dfs.append(df)
combined = pd.concat(dfs)
```

## Data Science Applications (Ai gend questions but use it as inspo)

### 1. Compensation Analysis
- **Salary Distribution Studies**: Analyze pay equity across departments, titles, and locations
- **Gender Pay Gap Analysis**: Investigate compensation differences by inferring gender from names
- **Cost of Living Adjustments**: Compare salaries across campuses relative to local housing costs
- **Top Earner Analysis**: Identify highest-paid positions and track changes over time

### 2. Temporal Trends
- **Wage Growth Patterns**: Track salary progression for specific roles over 15 years
- **Inflation Adjustment**: Compare real wage changes against CPI
- **Budget Impact Analysis**: Measure total compensation costs by department/campus
- **Hiring Trends**: Analyze workforce growth and contraction patterns

### 3. Predictive Modeling
- **Salary Prediction**: Build models to predict compensation based on title, location, and experience
- **Anomaly Detection**: Identify unusual compensation patterns or outliers
- **Career Path Analysis**: Map typical progression routes and associated salary increases
- **Retention Risk**: Model likelihood of turnover based on compensation trends

### 4. Comparative Studies
- **Inter-Campus Comparison**: Benchmark salaries across different UC locations
- **Public vs Private**: Compare UC salaries with industry standards
- **Academic vs Administrative**: Analyze compensation differences between faculty and staff
- **Department Rankings**: Rank departments by average compensation or growth

### 5. Network Analysis
- **Organizational Structure**: Infer reporting relationships from titles and compensation
- **Department Clustering**: Group similar departments based on compensation patterns
- **Title Standardization**: Use NLP to normalize job titles across campuses

### 6. Visualization Projects
- **Interactive Dashboards**: Build Tableau/PowerBI dashboards for salary exploration
- **Heat Maps**: Visualize compensation by location and department
- **Time Series Animations**: Show salary evolution over 15 years
- **Sankey Diagrams**: Display money flow across organizational units

### 7. Policy Research
- **Minimum Wage Impact**: Analyze effects of minimum wage changes on UC employees
- **Overtime Analysis**: Study overtime pay patterns and costs
- **Benefits Estimation**: Estimate total compensation including benefits
- **Budget Allocation**: Understand resource distribution across campuses

### 8. Machine Learning Applications
- **Classification**: Predict employee type (faculty/staff/student) from compensation data
- **Clustering**: Group employees with similar compensation profiles
- **Time Series Forecasting**: Predict future salary budgets
- **Natural Language Processing**: Extract insights from job titles and descriptions

## Sources

- Data source: UC Annual Wage (https://ucannualwage.ucop.edu/wage/)
- Update frequency: Annual
- Format: JSON with consistent schema
- Size: ~2-3GB total for all years and locations

## License

Data is public domain. Analysis code is GNU Affereo licensed.
