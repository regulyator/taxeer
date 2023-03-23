# Taxeer App

## Telegram bot for Georgian accounting

### **How to run**

- go to app root dir
- run `chmod +x start.sh`
- run `./start.sh <bot_api_key> <db_password>` (database schema name will be - taxeer)

### **How to update**

- checkout latest `main` branch and do same step as in previous section

### **How to use commands**

- `/currency` - print current currency rate for USD -> GEL
- `/income <value>:<currency>:<date>`
    - value - income value
    - currency - income currency
    - data (optional, default current date, _DD-MM-YYYY_) - date of income received
   
  **_example:_**
  - with date `/income 1000:USD:23-01-2023`
  - without date `/income 1000.50:USD`
- `/statistic` - print last ten incomes
- `/current` - print (all values already converted in GEL): 
  - total incomes in finance year
  - total incomes in finance month
  - taxes value needed to be pay in this month
