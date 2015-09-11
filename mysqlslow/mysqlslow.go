package mysqlslow

import (
	"bufio"
	"encoding/base64"
	"io"
	"log"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

var (
	firstLine  = regexp.MustCompile(`/.*Version:`)
	secondLine = regexp.MustCompile(`Tcp port:`)
	legendLine = regexp.MustCompile(`Time\s+Id\s+Command\s+Argument`)
	timeLine   = regexp.MustCompile(`# Time:`)
	useLine    = regexp.MustCompile(`use `)

	userLine       = regexp.MustCompile(`# User@Host: ([^\s]*)\[(.*)\] @ ([^\s]*) \[([0-9.]*)\]\s+Id:\s([0-9]+)`)
	statisticsLine = regexp.MustCompile(`# Query_time:\s+([0-9.]+)\s+Lock_time:\s+([0-9.]+)\s+Rows_sent:\s+([0-9])\s+Rows_examined:\s([0-9]+)`)
	timestampLine  = regexp.MustCompile(`SET timestamp=([0-9]+)`)

	// SQL normalization (from mysqldumpslow)
	decDigit       = regexp.MustCompile(`\b\d+\b`)
	hexDigit       = regexp.MustCompile(`\b0x[0-9A-Fa-f]+\b`)
	emptySquote    = regexp.MustCompile(`''`)
	emptyDquote    = regexp.MustCompile(`""`)
	unescapeSquote = regexp.MustCompile(`(\\')`)
	unescapeDquote = regexp.MustCompile(`(\\")`)
	stringSquote   = regexp.MustCompile(`'[^']+'`)
	stringDquote   = regexp.MustCompile(`"[^"]+"`)

	num  = []byte{'N'}
	sstr = []byte{'\'', 'S', '\''}
	dstr = []byte{'"', 'S', '"'}
)

func Parse(r io.Reader) []SlowLog {
	reader := bufio.NewReader(r)

	result := []SlowLog{}
	record := SlowLog{
		rawsql: []byte{},
	}
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			log.Println("reach EOF")
			break
		} else if err != nil {
			log.Println(err)
			break
		}

		switch {
		case firstLine.Match(line):
			fallthrough
		case secondLine.Match(line):
			fallthrough
		case legendLine.Match(line):
			fallthrough
		case timeLine.Match(line):
			fallthrough
		case useLine.Match(line):
			// skip!
		case userLine.Match(line):
			if len(record.Sql) != 0 {
				result = append(result, record)
				//fmt.Println(string(record.NormalizedSql))
				record = SlowLog{
					rawsql: []byte{},
				}
			}
			submatches := userLine.FindSubmatch(line)
			record.User = string(submatches[1])
			record.SrcUser = string(submatches[2])
			record.Host = string(submatches[3])
			record.Address = string(submatches[4])
			//fmt.Printf("user=%s, srcUser=%s, hostname=%s, ipAddress=%s\n",
			//	user, srcUser, hostname, ipAddress)
		case statisticsLine.Match(line):
			submatches := statisticsLine.FindSubmatch(line)
			queryTime, _ := strconv.ParseFloat(string(submatches[1]), 64)
			lockTime, _ := strconv.ParseFloat(string(submatches[2]), 64)
			rowsSent, _ := strconv.ParseInt(string(submatches[3]), 10, 64)
			rowsExamined, _ := strconv.ParseInt(string(submatches[4]), 10, 64)
			record.QueryTime = queryTime
			record.LockTime = lockTime
			record.RowsSent = rowsSent
			record.RowsExamined = rowsExamined
			//fmt.Printf("QueryTime=%s, LockTime=%s RowsSent=%s RowsExamined=%s\n",
			//	queryTime, lockTime, rowsSent, rowsExamined)
		case timestampLine.Match(line):
			submatches := timestampLine.FindSubmatch(line)
			epoch, _ := strconv.ParseInt(string(submatches[1]), 10, 64)
			record.Time = time.Unix(epoch, 0)
			//fmt.Printf("timestamp=%s\n", timestamp)

		default:
			switch line[len(line)-1] {
			case ';':
				fallthrough
			case '\n':
				line = line[0 : len(line)-1]
			}
			if len(record.rawsql) != 0 {
				record.rawsql = append(record.rawsql, ' ')
			}
			record.rawsql = append(record.rawsql, line...)
			if utf8.Valid(record.rawsql) {
				record.Sql = string(record.rawsql)
			} else {
				record.Sql = base64.StdEncoding.EncodeToString(record.rawsql)
			}
			record.NormalizedSql = string(Normalize(record.rawsql))
		}

	}

	return result
}

func Normalize(sql []byte) []byte {
	sql = decDigit.ReplaceAll(sql, num)
	sql = hexDigit.ReplaceAll(sql, num)
	sql = emptySquote.ReplaceAll(sql, sstr)
	sql = emptyDquote.ReplaceAll(sql, dstr)
	sql = unescapeSquote.ReplaceAll(sql, []byte{})
	sql = unescapeDquote.ReplaceAll(sql, []byte{})
	sql = stringSquote.ReplaceAll(sql, sstr)
	sql = stringDquote.ReplaceAll(sql, dstr)
	return sql
}

type SlowLog struct {
	User          string    `json:"user"`
	SrcUser       string    `json:"src_user"`
	Host          string    `json:"host"`
	Address       string    `json:"address"`
	Time          time.Time `json:"time"`
	QueryTime     float64   `json:"query_time"`
	LockTime      float64   `json:"lock_time"`
	RowsSent      int64     `json:"rows_sent"`
	RowsExamined  int64     `json:"rows_examined"`
	Sql           string    `json:"sql"`
	NormalizedSql string    `json:"normalized_sql"`
	rawsql        []byte
}

func (s SlowLog) ToFluentLog() (time.Time, FluentLog) {
	return s.Time, FluentLog{
		User:          s.User,
		SrcUser:       s.SrcUser,
		Host:          s.Host,
		Address:       s.Address,
		QueryTime:     s.QueryTime,
		LockTime:      s.LockTime,
		RowsSent:      s.RowsSent,
		RowsExamined:  s.RowsExamined,
		Sql:           s.Sql,
		NormalizedSql: s.NormalizedSql,
	}
}

type FluentLog struct {
	User          string  `codec:"user"`
	SrcUser       string  `codec:"src_user"`
	Host          string  `codec:"host"`
	Address       string  `codec:"address"`
	QueryTime     float64 `codec:"query_time"`
	LockTime      float64 `codec:"lock_time"`
	RowsSent      int64   `codec:"rows_sent"`
	RowsExamined  int64   `codec:"rows_examined"`
	Sql           string  `codec:"sql"`
	NormalizedSql string  `codec:"normalized_sql"`
}
