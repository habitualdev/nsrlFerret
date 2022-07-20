package nsrl

import (
	"archive/zip"
	"fmt"
	progress "github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"nsrlFerret/db"
	"os"
	"path/filepath"
	"strings"
)

func GetNSrl() {
	linkUrl := "https://s3.amazonaws.com/rds.nsrl.nist.gov/RDS/current/rds_modernm.zip"
	hashUrl := "https://s3.amazonaws.com/rds.nsrl.nist.gov/RDS/current/version.txt"
	newLink := false

	resp, err := http.Get(hashUrl)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	versionData, _ := io.ReadAll(resp.Body)
	latestVersion := strings.Replace(strings.Split(strings.Split(string(versionData), "\n")[4], ",")[0], "\"", "", -1)

	_, err = os.Stat("lastHash.txt")
	if err == nil {
		lastHash, _ := os.ReadFile("lastHash.txt")
		if string(lastHash) == latestVersion {
			log.Println("No new NSRL files")
			goto UnzipStage
		} else {
			newLink = true
			log.Println("New NSRL files found")
			f, _ := os.OpenFile("lastHash.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			f.WriteString(latestVersion)
		}
	} else {
		log.Println("First Time Run")
		f, _ := os.OpenFile("lastHash.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		f.WriteString(latestVersion)
	}

	_, err = os.Stat("rds_modernm/NSRLFile.txt")
	if err != nil {
		log.Println("NSRL hash list not found, downloading...")
		resp, err := http.Get(linkUrl)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		defer resp.Body.Close()
		out, _ := os.Create("rds_modernm.zip")
		defer out.Close()

		bar := progress.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)

		io.Copy(io.MultiWriter(out, bar), resp.Body)
	}
UnzipStage:
	_, err = os.Stat("rds_modernm")
	if err != nil || newLink == true {
		if newLink {
			log.Println("Deleting old NSRL files")
			os.RemoveAll("rds_modernm")
		}
		log.Println("Unzipping NSRL...")
		unzip()
		log.Println("Starting to process new NSRL List")
		db.ProcessNSRLtxt("rds_modernm/rds_modernm/NSRLFile.txt")
	}

}

func unzip() {
	dst := "rds_modernm"
	archive, err := zip.OpenReader("rds_modernm.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		log.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			log.Fatalln("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Fatalln(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatalln(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			log.Fatalln(err)
		}

		bar := progress.DefaultBytes(
			int64(archive.File[1].UncompressedSize64),
			"unzipping",
		)

		if _, err := io.Copy(io.MultiWriter(dstFile, bar), fileInArchive); err != nil {
			log.Fatalln(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
