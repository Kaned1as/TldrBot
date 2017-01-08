package impl

import (
    "gopkg.in/gorp.v2"
    _ "github.com/mattn/go-sqlite3"
    "time"
    "database/sql"
    "log"
    "os"
    "github.com/mattn/go-sqlite3"
)

// ORM type for table `SCORES`
type Score struct {
    Id uint64                           `db:"id, primarykey, autoincrement"`
    // ID from Telegram
    PersonId int64
    // hit grade (4 minimum)
    Grade uint8
    // time when score was recorded
    Time time.Time
}

type ScoreStat struct {
    // ID from Telegram
    PersonId int64
    // sum of all scores
    TotalScored uint64
}

type PersistenceLayer struct {
    database *gorp.DbMap
}

func newDb(path string) *PersistenceLayer {
    pl := PersistenceLayer{}
    pl.database = initDb(path)
    return &pl
}

func (pl *PersistenceLayer) Close() {
    pl.database.Db.Close()
}

func (pl *PersistenceLayer) SaveScore(score *Score) {
    pl.database.Insert(score)
}

func (pl *PersistenceLayer) GetScores(personId int64) (total int64, highest *Score, latest *Score) {
    total, _ = pl.database.SelectInt("select sum(Grade) from scores where PersonId = ?", personId)
    max, _ := pl.database.SelectInt("select max(Grade) from scores where PersonId = ?", personId)
    maxTime, _ := pl.database.SelectStr("select max(Time) from scores where PersonId = ?", personId)

    if (max > 0) { // we have records for this user
        var maxScores, latestScores []Score
        _, err := pl.database.Select(&maxScores, "select * from scores where Grade = ?", max)
        checkErr(err, "Failed to select max scores")
        t, _ := time.Parse(sqlite3.SQLiteTimestampFormats[0], maxTime)
        _, err = pl.database.Select(&latestScores, "select * from scores where Time = ?", t)
        checkErr(err, "Failed to select latest scores")

        if len(maxScores) > 0 {
            highest = &maxScores[0]
        }

        if len(latestScores) > 0 {
            latest = &latestScores[0]
        }
    }

    return
}

func (pl *PersistenceLayer) GetTotals() (stats []ScoreStat) {
    _, err := pl.database.Select(&stats, "select PersonId, sum(Grade) as TotalScored from scores group by PersonId order by TotalScored desc")
    checkErr(err, "Failed to get total stats")
    return
}

func initDb(path string) *gorp.DbMap {
    // connect to db using standard Go database/sql API
    // use whatever database/sql driver you wish
    db, err := sql.Open("sqlite3", path)
    checkErr(err, "sql.Open failed")

    // construct a gorp DbMap
    dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    dbmap.TraceOn("[ORM]", log.New(os.Stdout, "TraceSQL:", log.Lmicroseconds))

    // add a table, setting the table name to 'posts' and
    // specifying that the Id property is an auto incrementing PK
    scoreT := dbmap.AddTableWithName(Score{}, "scores")
    scoreT.AddIndex("PersonIdx", "Btree", []string{"PersonId"})
    scoreT.AddIndex("TimeIdx", "Btree", []string{"Time"})

    // create the table. in a production system you'd generally
    // use a migration tool, or create the tables via scripts
    err = dbmap.CreateTablesIfNotExists()
    checkErr(err, "Create tables failed")
    _ = dbmap.CreateIndex() // indexes will fail

    return dbmap
}

func checkErr(err error, msg string) {
    if err != nil {
        log.Fatalln(msg, err)
    }
}