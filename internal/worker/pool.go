package worker

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"sentinel/internal/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	DB          *pgxpool.Pool
	Concurrency int
	JobChan     chan models.Job
	Wg          sync.WaitGroup
}

func New(db *pgxpool.Pool, concurrency int) *Pool {

	return &Pool{
		DB:          db,
		Concurrency: concurrency,
		JobChan:     make(chan models.Job, concurrency),
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.Concurrency; i++ {
		p.Wg.Add(1)
		go p.work(i)

	}
}

func (p *Pool) work(workerID int) {
	defer p.Wg.Done()

	for job := range p.JobChan {
		fmt.Printf("[Worker %d] Processing: %s\n", workerID, job.URL)

		p.processJob(job)
	}
}

func (p *Pool) processJob(job models.Job) {
	startTime := time.Now()
	resp, err := http.Get(job.URL)
	if err != nil {
		fmt.Println("Error While fetching", err)
		return
	}
	responseTime := time.Since(startTime).Milliseconds()
	defer resp.Body.Close()

	limitContent := io.LimitReader(resp.Body, 2<<20)
	body, err := io.ReadAll(limitContent)
	if err != nil {
		fmt.Println("Error While reading response", err)
		return
	}

	hash := sha256.Sum256(body)
	contentHash := hex.EncodeToString(hash[:])
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error parsing HTML", err)
		return
	}

	pageTitle := doc.Find("title").Text()
	h1 := doc.Find("h1").Text()
	metaDescription, exists := doc.Find("meta[name='description']").Attr("content")
	if !exists {
		metaDescription = ""
	}

	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}

	})

	data := models.CrawlData{
		URL:             job.URL,
		ResponseTime:    int(responseTime),
		StatusCode:      resp.StatusCode,
		ContentHash:     contentHash,
		Title:           pageTitle,
		H1:              h1,
		MetaDescription: metaDescription,
		Links:           links,
	}

	dataDb, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling data: ", err)
		return
	}

	query := "INSERT INTO results(job_id,data) VALUES ($1,$2)"

	_, err = p.DB.Exec(context.Background(), query, job.ID, dataDb)
	if err != nil {
		fmt.Println("Error adding data to DB", err)

	}

}
