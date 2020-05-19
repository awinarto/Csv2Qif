package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

// Config is the configuration for Csv2Qif
type Config struct {
	CsvHasHeader             bool
	CsvColumnDate            int // D
	CsvColumnAmount          int // T
	CsvColumnMemo            int // M
	CsvColumnPayee           int // P
	CsvColumnCategory        int // L
	CsvColumnAddress         int // A
	CsvColumnRefNumber       int // N
	CsvColumnCleared         int // C
	CsvColumnReimburseFlag   int // F
	CsvColumnSplitCategory   int // S
	CsvColumnSplitMemo       int // E
	CsvColumnSplitAmount     int // $
	CsvColumnSplitPercentage int // %
	CsvDateFormat            string
	CsvReverseAmountSign     bool
	QifAccountType           string
	QifDateFormat            string
}

func (c *Config) init() {
	c.CsvColumnDate = -1
	c.CsvColumnAmount = -1
	c.CsvColumnMemo = -1
	c.CsvColumnPayee = -1
	c.CsvColumnCategory = -1
	c.CsvColumnAddress = -1
	c.CsvColumnRefNumber = -1
	c.CsvColumnCleared = -1
	c.CsvColumnReimburseFlag = -1
	c.CsvColumnSplitCategory = -1
	c.CsvColumnSplitMemo = -1
	c.CsvColumnSplitAmount = -1
	c.CsvColumnSplitPercentage = -1
}

// Load configuration from json file
func loadConfig(confFile string) (*Config, error) {
	cnt, err := ioutil.ReadFile(confFile)
	if err != nil {
		return nil, err
	}
	conf := Config{}
	conf.init()
	err = json.Unmarshal(cnt, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

// Convert CSV to QIF
func convertCsv2Qif(csvFile string, qifFile string, conf *Config) error {
	if conf == nil {
		return errors.New("Missing configuration")
	}
	in, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create qif file and add header
	out, err := os.Create(qifFile)
	if err != nil {
		return err
	}
	defer out.Close()
	if conf.QifAccountType != "" {
		fmt.Fprintln(out, "!Type:"+conf.QifAccountType)
	} else {
		fmt.Fprintln(out, "!Type:Bank")
	}

	reader := csv.NewReader(in)
	lineno := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		lineno++
		if conf.CsvHasHeader && lineno == 1 {
			continue
		}
		hasdata := false
		if conf.CsvColumnDate >= 0 && conf.CsvColumnDate < len(row) && row[conf.CsvColumnDate] != "" {
			tmpdatestr := row[conf.CsvColumnDate]
			if conf.CsvDateFormat != "" && conf.QifDateFormat != "" {
				tmpdate, err := time.Parse(getDateLayout(conf.CsvDateFormat), tmpdatestr)
				if err != nil {
					log.Println(err)
				} else {
					tmpdatestr = tmpdate.Format(getDateLayout(conf.QifDateFormat))
				}
			}
			fmt.Fprintln(out, "D"+tmpdatestr)
			hasdata = true
		}
		if conf.CsvColumnAmount >= 0 && conf.CsvColumnAmount < len(row) && row[conf.CsvColumnAmount] != "" {
			tmpamount := formatAmount(strings.TrimSpace(row[conf.CsvColumnAmount]))
			if conf.CsvReverseAmountSign && tmpamount != "" {
				if tmpamount[0] == '-' {
					tmpamount = tmpamount[1:]
				} else {
					tmpamount = "-" + tmpamount
				}
			}
			fmt.Fprintln(out, "T"+tmpamount)
			hasdata = true
		}
		if conf.CsvColumnMemo >= 0 && conf.CsvColumnMemo < len(row) && row[conf.CsvColumnMemo] != "" {
			fmt.Fprintln(out, "M"+row[conf.CsvColumnMemo])
			hasdata = true
		}
		if conf.CsvColumnPayee >= 0 && conf.CsvColumnPayee < len(row) && row[conf.CsvColumnPayee] != "" {
			fmt.Fprintln(out, "P"+row[conf.CsvColumnPayee])
			hasdata = true
		}
		if conf.CsvColumnCategory >= 0 && conf.CsvColumnCategory < len(row) && row[conf.CsvColumnCategory] != "" {
			fmt.Fprintln(out, "L"+row[conf.CsvColumnCategory])
			hasdata = true
		}
		if conf.CsvColumnAddress >= 0 && conf.CsvColumnAddress < len(row) && row[conf.CsvColumnAddress] != "" {
			fmt.Fprintln(out, "A"+row[conf.CsvColumnAddress])
			hasdata = true
		}
		if conf.CsvColumnRefNumber >= 0 && conf.CsvColumnRefNumber < len(row) && row[conf.CsvColumnRefNumber] != "" {
			fmt.Fprintln(out, "N"+row[conf.CsvColumnRefNumber])
			hasdata = true
		}
		if conf.CsvColumnCleared >= 0 && conf.CsvColumnCleared < len(row) && row[conf.CsvColumnCleared] != "" {
			fmt.Fprintln(out, "C"+row[conf.CsvColumnCleared])
			hasdata = true
		}
		if conf.CsvColumnReimburseFlag >= 0 && conf.CsvColumnReimburseFlag < len(row) && row[conf.CsvColumnReimburseFlag] != "" {
			fmt.Fprintln(out, "F"+row[conf.CsvColumnReimburseFlag])
			hasdata = true
		}
		if conf.CsvColumnSplitCategory >= 0 && conf.CsvColumnSplitCategory < len(row) && row[conf.CsvColumnSplitCategory] != "" {
			fmt.Fprintln(out, "S"+row[conf.CsvColumnSplitCategory])
			hasdata = true
		}
		if conf.CsvColumnSplitMemo >= 0 && conf.CsvColumnSplitMemo < len(row) && row[conf.CsvColumnSplitMemo] != "" {
			fmt.Fprintln(out, "E"+row[conf.CsvColumnSplitMemo])
			hasdata = true
		}
		if conf.CsvColumnSplitAmount >= 0 && conf.CsvColumnSplitAmount < len(row) && row[conf.CsvColumnSplitAmount] != "" {
			tmpamount := formatAmount(strings.TrimSpace(row[conf.CsvColumnSplitAmount]))
			if conf.CsvReverseAmountSign && tmpamount != "" {
				if tmpamount[0] == '-' {
					tmpamount = tmpamount[1:]
				} else {
					tmpamount = "-" + tmpamount
				}
			}
			fmt.Fprintln(out, "$"+tmpamount)
			hasdata = true
		}
		if conf.CsvColumnSplitPercentage >= 0 && conf.CsvColumnSplitPercentage < len(row) && row[conf.CsvColumnSplitPercentage] != "" {
			fmt.Fprintln(out, "%"+row[conf.CsvColumnSplitPercentage])
			hasdata = true
		}
		if hasdata {
			fmt.Fprintln(out, "^")
		}
	}
	return nil
}

// Convert date format to date layout
func getDateLayout(dtfmt string) string {
	if dtfmt == "" {
		return dtfmt
	}
	res := strings.ToUpper(dtfmt)
	res = strings.ReplaceAll(res, "DD", "02")
	res = strings.ReplaceAll(res, "D", "2")
	res = strings.ReplaceAll(res, "MMMM", "January")
	res = strings.ReplaceAll(res, "MMM", "Jan")
	res = strings.ReplaceAll(res, "MM", "01")
	res = strings.ReplaceAll(res, "M", "1")
	res = strings.ReplaceAll(res, "YYYY", "2006")
	res = strings.ReplaceAll(res, "YY", "06")
	res = strings.ReplaceAll(res, "Y", "06")
	return res
}

// Remove currency from amount
func formatAmount(str string) string {
	if str == "" {
		return str
	}
	re := regexp.MustCompile(`[-]|\d[\d,\.]*`)
	return strings.Join(re.FindAllString(str, -1), "")
}

func main() {
	csvHasHeader := flag.Bool("csvHasHeader", false, "Does the CSV file contain header?")
	csvColumnDate := flag.Int("csvColumnDate", -1, "The column index (start with 0) for Date column (D)")
	csvColumnAmount := flag.Int("csvColumnAmount", -1, "The column index (start with 0) for Amount column (T)")
	csvColumnMemo := flag.Int("csvColumnMemo", -1, "The column index (start with 0) for Memo column (M)")
	csvColumnPayee := flag.Int("csvColumnPayee", -1, "The column index (start with 0) for Payee column (P)")
	csvColumnCategory := flag.Int("csvColumnCategory", -1, "The column index (start with 0) for Category column (L)")
	csvColumnAddress := flag.Int("csvColumnAddress", -1, "The column index (start with 0) for Address column (A)")
	csvColumnRefNumber := flag.Int("csvColumnRefNumber", -1, "The column index (start with 0) for RefNumber column (N)")
	csvColumnCleared := flag.Int("csvColumnCleared", -1, "The column index (start with 0) for Cleared column (C)")
	csvColumnReimburseFlag := flag.Int("csvColumnReimburseFlag", -1, "The column index (start with 0) for ReimburseFlag column (F)")
	csvColumnSplitCategory := flag.Int("csvColumnSplitCategory", -1, "The column index (start with 0) for SplitCategory column (S)")
	csvColumnSplitMemo := flag.Int("csvColumnSplitMemo", -1, "The column index (start with 0) for SplitMemo column (E)")
	csvColumnSplitAmount := flag.Int("csvColumnSplitAmount", -1, "The column index (start with 0) for SplitAmount column ($)")
	csvColumnSplitPercentage := flag.Int("csvColumnSplitPercentage", -1, "The column index (start with 0) for SplitPercentage column (%)")
	csvDateFormat := flag.String("csvDateFormat", "", "Date format used in CSV file (e.g. YYYY-MM-DD)")
	csvReverseAmountSign := flag.Bool("csvReverseAmountSign", false, "Multiply amount with -1?")
	qifAccountType := flag.String("qifAccountType", "Bank", "The account type for QIF file")
	qifDateFormat := flag.String("qifDateFormat", "", "Date format used in QIF file (e.g. YYYY-MM-DD)")
	help := flag.Bool("help", false, "Print usage info")
	flag.Lookup("help").Hidden = true
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Println("Convert CSV file to QIF file")
		fmt.Println("")
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [--config=value ...] csvFile qifFile [configFile]")
		fmt.Println("")
		fmt.Println("csvFile      The CSV input file")
		fmt.Println("qifFile      The QIF output file")
		fmt.Println("configFile   The config file in JSON format (optional). The config can be passed as command line argument")
		fmt.Println("")
		fmt.Println("Config:")
		flag.PrintDefaults()
		fmt.Println("")
	}
	flag.Parse()

	args := flag.Args()
	if (help != nil && *help) || len(args) < 2 || len(args) > 3 {
		flag.Usage()
		if help != nil && *help {
			os.Exit(0)
		} else {
			fmt.Println("ERROR: Invalid command line argument")
			os.Exit(1)
		}
	}

	conf := &Config{}
	conf.init()
	if len(args) > 2 {
		// Has config
		tmpconf, err := loadConfig(args[2])
		if err != nil {
			log.Fatal(err)
		}
		conf = tmpconf
	} else {
		if csvHasHeader != nil {
			conf.CsvHasHeader = *csvHasHeader
		}
		if csvColumnDate != nil {
			conf.CsvColumnDate = *csvColumnDate
		}
		if csvColumnAmount != nil {
			conf.CsvColumnAmount = *csvColumnAmount
		}
		if csvColumnMemo != nil {
			conf.CsvColumnMemo = *csvColumnMemo
		}
		if csvColumnPayee != nil {
			conf.CsvColumnPayee = *csvColumnPayee
		}
		if csvColumnCategory != nil {
			conf.CsvColumnCategory = *csvColumnCategory
		}
		if csvColumnAddress != nil {
			conf.CsvColumnAddress = *csvColumnAddress
		}
		if csvColumnRefNumber != nil {
			conf.CsvColumnRefNumber = *csvColumnRefNumber
		}
		if csvColumnCleared != nil {
			conf.CsvColumnCleared = *csvColumnCleared
		}
		if csvColumnReimburseFlag != nil {
			conf.CsvColumnReimburseFlag = *csvColumnReimburseFlag
		}
		if csvColumnSplitCategory != nil {
			conf.CsvColumnSplitCategory = *csvColumnSplitCategory
		}
		if csvColumnSplitMemo != nil {
			conf.CsvColumnSplitMemo = *csvColumnSplitMemo
		}
		if csvColumnSplitAmount != nil {
			conf.CsvColumnSplitAmount = *csvColumnSplitAmount
		}
		if csvColumnSplitPercentage != nil {
			conf.CsvColumnSplitPercentage = *csvColumnSplitPercentage
		}
		if csvDateFormat != nil {
			conf.CsvDateFormat = *csvDateFormat
		}
		if csvReverseAmountSign != nil {
			conf.CsvReverseAmountSign = *csvReverseAmountSign
		}
		if qifAccountType != nil {
			conf.QifAccountType = *qifAccountType
		}
		if qifDateFormat != nil {
			conf.QifDateFormat = *qifDateFormat
		}
	}

	err := convertCsv2Qif(args[0], args[1], conf)
	if err != nil {
		log.Fatal(err)
	}
}
