package sidecarPostgres

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func getTResNames() []string {
	names := []string{}

	namesString := os.Getenv(TORORU_RESOURCE_NAMES)
	for _, n := range strings.Split(namesString, ",") {
		names = append(names, strings.TrimSpace(n))
	}

	return names
}

func getRotationDuration() time.Duration {
	secString := os.Getenv(ROTATION_SECONDS)
	secs, err := strconv.Atoi(secString)
	if err != nil {
		secs = 3600
		fmt.Println(secs)
	}

	return time.Duration(secs) * time.Second
}
