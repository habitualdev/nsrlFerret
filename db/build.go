package db

import (
	"bufio"
	bloomfilter "github.com/alovn/go-bloomfilter"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/briandowns/spinner"
	progress "github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BucketCollection struct {
	Uuids        []string
	BloomFilters map[string]bloomfilter.BloomFilter
	mapMutex     sync.RWMutex
}

type NSRLDataPoint struct {
	SHA1         string
	MD5          string
	CRC32        string
	FileName     string
	FileSize     string
	ProductCode  string
	OpSystemCode string
	SpecialCode  string
}

type bucketEntry struct {
	Bucket *bloom.BloomFilter
	Type   string
	UUID   string
}

func (b *BucketCollection) CheckBloom(hash string) bool {
	for _, filter := range b.BloomFilters {
		if contains, _ := filter.MightContain([]byte(hash)); contains {
			return true
		}
	}
	return false
}

func (b *BucketCollection) AddNSRL(dbUuid string, nsrlData []string, wg *sync.WaitGroup) error {
	b.mapMutex.Lock()
	b.BloomFilters[dbUuid] = bloomfilter.NewMemoryBloomFilter(33547705)

	newBucket, _ := os.OpenFile("buckets/"+dbUuid+".bkt", os.O_RDWR|os.O_CREATE, 0666)
	for _, lineBuffer := range nsrlData {
		delimitedText := strings.Split(strings.ReplaceAll(lineBuffer, "\"", ""), ",")
		if len(delimitedText) == 1 {
			continue
		}
		newBucket.WriteString(lineBuffer + "\n")
		tempDatapoint := NSRLDataPoint{
			SHA1:         delimitedText[0],
			MD5:          delimitedText[1],
			CRC32:        delimitedText[2],
			FileName:     delimitedText[3],
			FileSize:     delimitedText[4],
			ProductCode:  delimitedText[5],
			OpSystemCode: delimitedText[6],
			SpecialCode:  delimitedText[7],
		}
		err := b.BloomFilters[dbUuid].Put([]byte(tempDatapoint.MD5))
		if err != nil {
			println(err.Error())
		}
		b.BloomFilters[dbUuid].Put([]byte(tempDatapoint.SHA1))
		b.BloomFilters[dbUuid].Put([]byte(tempDatapoint.FileName))

		if err != nil {
			return err
		}

	}
	b.Uuids = append(b.Uuids, dbUuid)
	b.mapMutex.Unlock()
	wg.Done()
	return nil
}

func ProcessNSRLtxt(filename string) BucketCollection {
	bucketNume := 0
	Buckets := BucketCollection{
		Uuids:        []string{},
		BloomFilters: map[string]bloomfilter.BloomFilter{},
	}
	var wg sync.WaitGroup
	stats, _ := os.Stat(filename)
	file, _ := os.OpenFile(filename, os.O_RDONLY, 0666)
	bar := progress.DefaultBytes(stats.Size(), "Processing NSRL File")
	scanner := bufio.NewScanner(file)
	numLines := 1000000
	nsrlBuffer := []string{}
	for scanner.Scan() {
		nsrlBuffer = append(nsrlBuffer, scanner.Text())
		bar.Add(len(scanner.Text()))
		for len(nsrlBuffer) < numLines+1 {
			nsrlBuffer = append(nsrlBuffer, scanner.Text())
			scanner.Scan()
		}
		wg.Add(1)
		go Buckets.AddNSRL(strconv.Itoa(bucketNume), nsrlBuffer, &wg)
		bucketNume++
		nsrlBuffer = []string{}

	}

	bar.Close()
	log.Println("Waiting for processing to finish")
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spinner.Start()
	wg.Wait()
	spinner.Stop()
	return Buckets
}
