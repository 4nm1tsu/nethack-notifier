package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Record struct {
	GameVersion     string
	Score           int
	DungeonID       int
	DungeonLevel    int
	MaxDungeonLevel int
	HP              int
	MaxHP           int
	Unused1         int
	EndDate         int
	StartDate       int
	Unused2         int
	Class           string
	Race            string
	Gender          string
	Alignment       string
	Name            string
	Result          string
}

type discordImg struct {
	URL string `json:"url"`
	H   int    `json:"height"`
	W   int    `json:"width"`
}
type discordAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon_url"`
}
type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type discordEmbed struct {
	Title  string         `json:"title"`
	Desc   string         `json:"description"`
	URL    string         `json:"url"`
	Color  int            `json:"color"`
	Image  discordImg     `json:"image"`
	Thum   discordImg     `json:"thumbnail"`
	Author discordAuthor  `json:"author"`
	Fields []discordField `json:"fields"`
}

type discordWebhook struct {
	UserName  string         `json:"username"`
	AvatarURL string         `json:"avatar_url"`
	Content   string         `json:"content"`
	Embeds    []discordEmbed `json:"embeds"`
	TTS       bool           `json:"tts"`
}

// No slash is needed at the end of the path
var InProgressDir string
var RecordFileName string
var WebhookURL string
var AvatarURL string
var UserName string
var ServerDomain string

func sendWebhook(payload *discordWebhook) error {
	j, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", WebhookURL, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		//log.Println("sent: ", payload)
	} else {
		return err
	}
	return nil
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func listFilesInDirectory(dirPath string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func getActiveUsers() ([]string, error) {
	pattern := `^\.nfs\d{24}$`
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var usernames []string

	fileList, err := listFilesInDirectory(InProgressDir)
	if err != nil {
		return nil, err
	}

	for _, filePath := range fileList {
		parts := strings.Split(filepath.Base(filePath), ":")
		username := ""
		if len(parts) > 0 {
			if r.MatchString(parts[0]) {
				continue
			}
			username = parts[0]
		} else {
			return nil, errors.New(fmt.Sprintf("Unexpected file name: %s", filepath.Base(filePath)))
		}

		usernames = append(usernames, username)
	}

	return usernames, nil
}

func eventLoop(watcher *fsnotify.Watcher) error {
	payload := discordWebhook{
		UserName:  UserName,
		AvatarURL: AvatarURL,
		Content:   "",
		Embeds:    nil,
		TTS:       false,
	}
	oldRecords, err := parseRecord()
	if err != nil {
		return err
	}
	oldInprogressUsers, err := getActiveUsers()
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Printf("EVENT: %+v %+v\n", event.Op, event.Name)

			if InProgressDir == filepath.Dir(event.Name) {
				parts := strings.Split(filepath.Base(event.Name), ":")
				username := ""

				pattern := `^\.nfs\d{24}$`
				r, err := regexp.Compile(pattern)
				if err != nil {
					return err
				}

				if len(parts) > 0 {
					if r.MatchString(parts[0]) {
						continue
					}
					username = parts[0]
				} else {
					return errors.New(fmt.Sprintf("Unexpected file name: %s", filepath.Base(event.Name)))
				}
				if event.Op == fsnotify.Create {
					payload.Content = fmt.Sprintf("%s started exploring.\n`$ telnet %s` to watch the game in progress!", username, ServerDomain)
					if err := sendWebhook(&payload); err != nil {
						return err
					}
					oldInprogressUsers, err = getActiveUsers()
					if err != nil {
						return err
					}
				}
				if event.Op == fsnotify.Remove {
					inprogressUsers, err := getActiveUsers()
					if err != nil {
						return err
					}
				OuterLoop:
					for _, old := range oldInprogressUsers {
						for _, u := range inprogressUsers {
							if old == u {
								continue OuterLoop
							}
						}
						payload.Content = fmt.Sprintf("%s finished exploring.", old)
						if err := sendWebhook(&payload); err != nil {
							return err
						}
					}
					oldInprogressUsers = inprogressUsers
				}
			}

			if event.Op == fsnotify.Write && event.Name == RecordFileName {
				records, err := parseRecord()
				if err != nil {
					return err
				}

				// 新しく追加されたレコードを特定
				for _, old := range oldRecords {
					for i, rec := range records {
						// 完全に等しいなら、recordsから削除
						if old == rec {
							records = append(records[:i], records[i+1:]...)
							break
						}
					}
					// 古いrecordがrecordsの中にないとき
				}
				if len(records) < 1 {
					log.Printf("No new records found.")
					continue
				}
				// oldRecordsに、recordに残ったものを加える
				for _, r := range records {
					oldRecords = append(oldRecords, r)
					payload.Content = fmt.Sprintf("%v(%v-%v-%v) %v.(SCORE: %v)", r.Name, r.Class, r.Race, r.Alignment, r.Result, r.Score)
					if err := sendWebhook(&payload); err != nil {
						return err
					}
				}
			}
		}
	}
}

func parseRecord() ([]Record, error) {
	file, err := os.Open(RecordFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []Record

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		fields := strings.Fields(parts[0])
		if len(fields) != 16 || len(parts) != 2 {
			log.Printf("Invalid record: line %d", lineNumber)
			continue
			//return nil, errors.New(fmt.Sprintf("Invalid record file format: %s", line))
		}

		record := Record{
			GameVersion:     fields[0],
			Score:           atoi(fields[1]),
			DungeonID:       atoi(fields[2]),
			DungeonLevel:    atoi(fields[3]),
			MaxDungeonLevel: atoi(fields[4]),
			HP:              atoi(fields[5]),
			MaxHP:           atoi(fields[6]),
			Unused1:         atoi(fields[7]),
			EndDate:         atoi(fields[8]),
			StartDate:       atoi(fields[9]),
			Unused2:         atoi(fields[10]),
			Class:           fields[11],
			Race:            fields[12],
			Gender:          fields[13],
			Alignment:       fields[14],
			Name:            fields[15],
			Result:          parts[1],
		}

		records = append(records, record)
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func main() {
	InProgressDir = os.Getenv("IN_PROGRESS_DIR")
	RecordFileName = os.Getenv("RECORD_FILE_NAME")
	WebhookURL = os.Getenv("WEBHOOK_URL")
	AvatarURL = os.Getenv("AVATAR_URL")
	UserName = os.Getenv("USER_NAME")
	ServerDomain = os.Getenv("SERVER_DOMAIN")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	if err := watcher.Add(InProgressDir); err != nil {
		panic(err)
	}

	if err := watcher.Add(RecordFileName); err != nil {
		panic(err)
	}

	if err := eventLoop(watcher); err != nil {
		log.Fatal(err.Error())
	}
}
