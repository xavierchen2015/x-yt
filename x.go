package main

import (
   "github.com/rylio/ytdl"
   "github.com/spf13/viper"

   "strings"
   "os"
   "bufio"
   "fmt"
   "log"
   "path/filepath"
   "runtime"
)

func main() {

   //// use viper get config file

   viper.SetConfigType("yaml")
   viper.SetConfigName("x")
   viper.AddConfigPath(".")

   if err := viper.ReadInConfig(); err == nil {
      fmt.Println("Using config file:", viper.ConfigFileUsed())
   }

   var dir = ""

   //// win & mac current path different, use GOOS defind OS, then get dir for each
   //// win -> exe work dir
   //// mac -> x.yaml path

   if runtime.GOOS == "windows" {
       path2, _ := os.Executable()
       dir = filepath.Dir(path2)
   } else {
       path := viper.GetString("path")
       os.Chdir(path)
       dir = path
   }

   fmt.Println(dir)

   currentPath := filepath.Join(dir, "file.txt")
   file, err := os.Open(currentPath)

   if err != nil {
      log.Fatal(err)
   }
   defer file.Close()

   //// open txt file loop get youtube url

   scanner := bufio.NewScanner(file)
   for scanner.Scan() {
      url := scanner.Text()

      //// lib api get youtube info
      videoInfo, _ := ytdl.GetVideoInfo(url)


      videotitle := strings.Replace(videoInfo.Title, "/", "", -1)
      loadfile := videotitle + ".mp4"
      downloadPath := filepath.Join(dir, loadfile)

      fmt.Println(videoInfo.Title)

      file2, _ := os.Create(downloadPath)
      defer file2.Close()

      videoInfo.Download(videoInfo.Formats.Best(ytdl.FormatAudioEncodingKey)[0], file2)
   }

   if err := scanner.Err(); err != nil {
      fmt.Println(err)
      log.Fatal(err)
   }
}
