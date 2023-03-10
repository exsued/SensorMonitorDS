package main

import (
    "github.com/tarm/serial"
    "github.com/itchyny/timefmt-go"
    "fmt"
    "log"
    "time"
    "os"
)

func WriteToFile(str string, startTime string) {
    fileName := "dataset" + startTime + ".csv"
    //If file not exists
    if _, err := os.Stat(fileName); err != nil {
        _, err := os.Create(fileName)
        if err != nil {
            log.Fatal(err)
        }
        appendFile(fileName, "time,DHT22 humidity,DHT22 temp,BMP085 temp,BMP085 pressure,\n")
    }
    finStr := timefmt.Format(time.Now(), "%H:%M:%S") + "," + str
    appendFile(fileName, finStr)
}

func appendFile(fileName string, str string) {
    file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Println(err)
    }
    defer file.Close()
    if _, err := file.WriteString(str); err != nil {
        log.Fatal(err)
    }
}

func main() {
    c0 := &serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}
    startTime := timefmt.Format(time.Now(), "%Y%m%d%H%M%S")
    s, err := serial.OpenPort(c0)
    if err != nil {
        log.Println(err)
        return
    }
    str := ""
    buf := make([]byte, 128)
    for ; ; {
        n, err := s.Read(buf)
        if err != nil {
            log.Println(err)
            return
        }
        str += string(buf[:n])
        if len(str) == 27 {
            WriteToFile(str, startTime)
            fmt.Print(timefmt.Format(time.Now(), "%H:%M:%S") + "," + str)
            str = ""
            time.Sleep(time.Second * 25)
            s.Flush()
        }else if len(str) > 27 {
            str = ""
        }
    }
}
