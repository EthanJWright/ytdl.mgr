package main
import (
  "fmt";
  "os";
  "os/exec";
  "log";
  "bytes";
  "bufio";
  "time";
)

// create a container to track whether a download is completed, by name
type download struct {
  name string;
  completed bool;
  scanner *bufio.Scanner;
}

// track all downloads by name
var downloads = []download{}

// turn the logic in main into a function
func ytdownload(input string) {
  // create a new Command object
  cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "--output", "./downloads/%(title)s.%(ext)s", input);
  // create a buffer for stderr
  var stderr bytes.Buffer
  // attach buffer to Command
  cmd.Stderr = &stderr
  // create a buffer for stdout
  stdout, _ := cmd.StdoutPipe()
  scanner := bufio.NewScanner(stdout)
  // attach buffer to Command
  // run the command

  // add download to the list of downloads
  downloads = append(downloads, download{input, false, scanner})

  err := cmd.Start();

  for scanner.Scan() {
    fmt.Println(scanner.Text())
  }
  cmd.Wait()

  if err != nil {
    // log stderr
    log.Println(stderr.String())
    log.Fatal(err);
  }

  // set the download to completed
  for i, dl := range downloads {
    if dl.name == input {
      downloads[i].completed = true
      // print download completed
      fmt.Println("Download completed: ", input)
    }
  }
}

// create a wrapper to thread the download
func downloadWrapper(input string) {
  go ytdownload(input);
}

// create an async function that posts updates on download progress
func updateDownloads() {
  var lastReportedDownloadCount = -1
  // loop forever
  for {
    // sleep for 5 seconds
    time.Sleep(5 * time.Second)
    // count number of downloads that are false
    var count int = 0
    for _, dl := range downloads {
      if !dl.completed {
        count++
      }
    }
    // if lastReportedDownloadCount is different from current download count,
    // print it and set it
    if lastReportedDownloadCount != count {
      fmt.Println("Downloads: ", count)
      lastReportedDownloadCount = count
    }
  }
}


// create a loop for user input
func loop() {
  for {
    fmt.Println("Enter youtube url: ")
    var first string
    in := bufio.NewReader(os.Stdin)
    line, _, err := in.ReadLine()
    if err != nil {
      log.Fatal(err)
    }
    // convert line to string
    first = string(line)

    // if the user enters "quit"
    if first == "quit" {
      // exit the program
      os.Exit(0)
    }

    fmt.Println("Queueing download: ", first)
    // download the video
    downloadWrapper(first)

  }
}



// use the youtube-dl cli tool to download a youtube video
func main() {
  // every 5 seconds, update the downloads
  go updateDownloads()
  loop()
}


