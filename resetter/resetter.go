// Copyright (C) 2016-2017  Vincent Ambo <mail@tazj.in>
//
// This file is part of Kontemplate.
//
// Kontemplate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This file contains the implementation of a template function for retrieving
// IP addresses from DNS

package resetter

import (
  "bufio"
  "os"
  "regexp"
  "bytes"
  "fmt"
)

func LoadAndResetConfig(configFile *string) (string, error) {
  // Dereference configFile
  configFileValue := *configFile

  file, err := os.Open(configFileValue)
  if err != nil {
    return "", fmt.Errorf("Error while opening file %s: %v", configFile, err)
  }
  // It wil close file when main function returns
  defer file.Close()

  scanner := bufio.NewScanner(file)
  var bufferResettedLines bytes.Buffer
  for scanner.Scan() {
    resettedLine := resetLine(scanner.Text())

    _, err := bufferResettedLines.WriteString(resettedLine + "\n")
    if err != nil {
      return "", fmt.Errorf("Error while writting string in buffer %s: %v", configFile, err)
      break
    }
  }

  if err := scanner.Err(); err != nil {
    return "", fmt.Errorf("Error while reading file %s: %v", configFile, err)
  }

  return bufferResettedLines.String(), nil
}

func resetLine(line string) string {
  nonResetableKeyValues := map[string]bool{
    "global": true,
    "includes": true,
    "values": true,
    "import": true}
  isFixedLine, _ := regexp.MatchString("# FIXED", line)

  if isFixedLine {
    return line
  } else {
    var mapLineRegexp = regexp.MustCompile(`^([\s-]*\S*): (\S*)(.*)$`)
    var lineSubmatches = mapLineRegexp.FindStringSubmatch(line)

    // Get mapkey
    var mapKey string = ""
    if len(lineSubmatches) > 2 {
      mapKey = lineSubmatches[1]
    }

    if mapKey != "" && !nonResetableKeyValues[mapKey] {
      return mapKey + ": ''"
      } else {
        return line
      }
    }
}
