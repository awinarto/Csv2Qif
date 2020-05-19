# Csv2Qif

Csv2Qif is used to convert CSV (Comma-separated values) file to QIF (Quicken Interchange Format) file.

## Installation

Install [Go](https://golang.org) and run the following command. Alternatively, download the released executable. 

```
go get github.com/awinarto/Csv2Qif
Csv2Qif 
```

## Usage

```
Usage: Csv2Qif.exe [--config=value ...] csvFile qifFile [configFile]

csvFile      The CSV input file
qifFile      The QIF output file
configFile   The config file in JSON format (optional). The config can be passed as command line argument

Config:
      --csvColumnAddress int           The column index (start with 0) for Address column (A) (default -1)
      --csvColumnAmount int            The column index (start with 0) for Amount column (T) (default -1)
      --csvColumnCategory int          The column index (start with 0) for Category column (L) (default -1)
      --csvColumnCleared int           The column index (start with 0) for Cleared column (C) (default -1)
      --csvColumnDate int              The column index (start with 0) for Date column (D) (default -1)
      --csvColumnMemo int              The column index (start with 0) for Memo column (M) (default -1)
      --csvColumnPayee int             The column index (start with 0) for Payee column (P) (default -1)
      --csvColumnRefNumber int         The column index (start with 0) for RefNumber column (N) (default -1)
      --csvColumnReimburseFlag int     The column index (start with 0) for ReimburseFlag column (F) (default -1)
      --csvColumnSplitAmount int       The column index (start with 0) for SplitAmount column ($) (default -1)
      --csvColumnSplitCategory int     The column index (start with 0) for SplitCategory column (S) (default -1)
      --csvColumnSplitMemo int         The column index (start with 0) for SplitMemo column (E) (default -1)
      --csvColumnSplitPercentage int   The column index (start with 0) for SplitPercentage column (%) (default -1)
      --csvDateFormat string           Date format used in CSV file (e.g. YYYY-MM-DD)
      --csvHasHeader                   Does the CSV file contain header?
      --csvReverseAmountSign           Multiply amount with -1?
      --qifAccountType string          The account type for QIF file (default "Bank")
      --qifDateFormat string           Date format used in QIF file (e.g. YYYY-MM-DD)
```

### Configuration

The configuration can be specified in a file or passed as command line parameter. If config file is specified, it will overwrite configuration specified in the command line parameter. The config file is in JSON format. Set the value to null or remove the config to use the default value. Below is an example of the config file.

```
{
    "csvColumnAddress": null,
    "csvColumnAmount": 4,
    "csvColumnCategory": 2,
    "csvColumnCleared": null,
    "csvColumnDate": 0,
    "csvColumnMemo": 1,
    "csvColumnPayee": null,
    "csvColumnRefNumber": null,
    "csvColumnReimburseFlag": null,
    "csvColumnSplitAmount": null,
    "csvColumnSplitCategory": null,
    "csvColumnSplitMemo": null,
    "csvColumnSplitPercentage": null,
    "csvDateFormat": "D-MMM-YY",
    "csvHasHeader": true,
    "csvReverseAmountSign": false,
    "qifAccountType": "Bank",
    "qifDateFormat": "DD/MM/YY"
}
```

The above config would convert the following CSV file

```
Date,Description,Category,Amount
16-Apr-20,Payment 1,Electricity,-$500.80
8-Apr-20,Payment 2,Water,-$150.00
2-Apr-20,Payment 3,Internet,-$39.99
26-Mar-20,Payment 4,Grocery,-$82.24
19-Mar-20,Payment 5,Homeware,-$27.30
```

to QIF file below

```
!Type:Bank
D16/04/20
T-500.80
MPayment 1
LElectricity
^
D08/04/20
T-150.00
MPayment 2
LWater
^
D02/04/20
T-39.99
MPayment 3
LInternet
^
D26/03/20
T-82.24
MPayment 4
LGrocery
^
D19/03/20
T-27.30
MPayment 5
LHomeware
^
```

### Date format

The date in the QIF file could have different format compared to the CSV file. This is done by specifying both csvDateFormat and qifDateFormat. The following are supported.

| Date Format | Description                                  |
| ----------- | -------------------------------------------- |
| D           | The day of the month (1 - 31)                |
| DD          | The day of the month (01 - 31)               |
| M           | The month (1 - 12)                           |
| MM          | The month (01 - 12)                          |
| MMM         | The abbreviated name of the month (e.g. Jan) |
| MMMM        | The full name of the month (e.g. January)    |
| YY          | The 2 digits year                            |
| YYYY        | The 4 digits year                            |

## Compile

Csv2Qif is developed using [Go](https://golang.org) language. To build the executable, install Go and run the following commands.

```
go get github.com/spf13/pflag
go build Csv2Qif.go 
```

## License

Csv2Qif is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
