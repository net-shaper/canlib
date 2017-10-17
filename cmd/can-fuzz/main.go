package main

import (
    "github.com/buffersandbeer/canlib"
    "flag"
    "fmt"
    "time"
    "strconv"
    "os"
    "strings"
    "bufio"
    "math"
)

func check(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func main() {

    // Setup command line flags
    caniface := flag.String("interface", "vcan0", "The CAN interface to fuzz")
    displayPkts := flag.Int("granularity", 1000000, "The number of packets to send before printing seed")
    seed := flag.String("seed", "0x0", "Seed to begin fuzzing at")
    rateLimit := flag.Int("rate-limit", 0, "The number of miliseconds to sleep in between packets")
    targetFile := flag.String("ids-file", "NONE", "Newline delimited file of hex CAN IDs to fuzz with")

    // Customize usage display
    flag.Usage = func() {
        appNameTmp := strings.Split(os.Args[0], "/")
        appName := appNameTmp[len(appNameTmp)-1]
        fmt.Fprintf(os.Stderr, "Tool to fuzz CAN bus interfaces\n\n")
        fmt.Fprintf(os.Stderr, "Usage:\n")
        fmt.Fprintf(os.Stderr, "\t%s (-ids-file <path>) [-option <argument>]...\n", appName)
        fmt.Fprintf(os.Stderr, "\nOptions:\n")
        flag.PrintDefaults()
        fmt.Fprintf(os.Stderr, "\nExamples:\n")
        fmt.Fprintf(os.Stderr, "\t%s -interface slcan0 -ids-file ./ids.txt\t# Fuzz slcan0 with ids from ids.txt\n", appName)
    }
    
    // Parse command line arguments
    flag.Parse()

    // Check for mandatory file parameter
    if *targetFile == "NONE" {
        fmt.Fprintf(os.Stderr, "You must supply a file of CAN ids\n\n")
        flag.Usage()
        os.Exit(2)
    }

    // Split seed and verify format
    splitSeed := strings.Split(*seed, "x")
    if len(splitSeed) != 2 {
        fmt.Fprintf(os.Stderr, "Invalid seed format. Seed must match the format int!int, such as 0x0.")
        flag.Usage()
        os.Exit(2)
    }

    // Convert seed values back to int
    var dataStart, sizeStart int
    var err error
    dataStart, err = strconv.Atoi(splitSeed[1])
    check(err)
    sizeStart, err = strconv.Atoi(splitSeed[0])
    check(err)

    // Create channels
    canout := make(chan canlib.RawCanFrame, 100)
    errChan := make(chan error)
    targetIDs := readIds(*targetFile) 

    // Start send Can parallel process and fuzzing function
    go canlib.SendCanConcurrent(*caniface, canout, errChan)
    fuzzCan(targetIDs, canout, *displayPkts, uint64(dataStart), sizeStart , *rateLimit)
}

// readIds will open and parse a file for CAN ids, and  panic if the file or contents are invalid
func readIds(path string) []uint32 {

    ids := make([]uint32, 0)
    if file, err := os.Open(path); err == nil {
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            id, _ := strconv.ParseInt(scanner.Text(), 16, 32)
            ids = append(ids, uint32(id))
        }

        if err = scanner.Err(); err != nil {
            panic("Error reading file: " + err.Error())
        }

    } else {
        panic(" error opening file")
    }

    return ids
}

// fuzzCan will generate fuzz frames for the CAN bus and pass them to the CanSend goroutine via channel
func fuzzCan(targets []uint32, output chan<- canlib.RawCanFrame, iterBeforeDisplay int, dataStart uint64, sizeStart int, rateLimit int) {
    fuzzFrame := new(canlib.RawCanFrame)
    fuzzDataStart := dataStart
    fuzzSizeStart := sizeStart
    displayTracker := iterBeforeDisplay

    // For all possible lengths of data in the packets (0-8 bytes)...
    for fuzzSize := fuzzSizeStart; fuzzSize < 8; fuzzSize ++ {
        fuzzBytes := make([]byte, fuzzSize)

        // Set the maximum data value based on target buffer size
        fuzzMax := uint64(math.Pow(2, float64((8 * fuzzSize)))-1)

        // For the range of data values of current target data length...
        fuzzData := fuzzDataStart
        for fuzzData <= fuzzMax {

            if fuzzSize == 0 {
                fuzzBytes = make([]byte, 0)
            } else {
                for index, _ := range fuzzBytes {
                    fuzzBytes[index] = byte(fuzzData >> uint(8 * index))
                }
            }

            // For each ID in the ID list, craft and send the packet
            for _, target := range targets {
                canlib.CreateRawFrame(fuzzFrame, uint32(target), fuzzBytes, false, false, false)
                fuzzFrame.Timestamp = time.Now().UnixNano()
                output <- *fuzzFrame
                time.Sleep(time.Duration(rateLimit) * time.Millisecond)
                displayTracker--

                // Print seed
                if displayTracker == 0 {
                    displayTracker = iterBeforeDisplay
                    timestamp := canlib.TimestampToSeconds(time.Now().UnixNano())
                    timestr := strconv.FormatFloat(timestamp, 'f', -1, 64)
                    fmt.Println(timestr, "\t", fmt.Sprintf("%dx%d",fuzzSize,fuzzData))
                }
            }
            fuzzData++
        }
    }
}
