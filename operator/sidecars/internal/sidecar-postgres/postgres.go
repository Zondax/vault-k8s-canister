package sidecarPostgres

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sethvargo/go-password/password"
	v1 "github.com/zondax/tororu-operator/operator/common/v1"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type database struct {
	Name string
	User string
}

type scPostgresConfig struct {
	Databases []database
}

func parseSidecarPostgresConfig(cString string) (*scPostgresConfig, error) {
	// Parse the config field into a DatabaseConfig struct
	var databaseConfig scPostgresConfig
	err := yaml.Unmarshal([]byte(cString), &databaseConfig)
	if err != nil {
		zap.S().Errorf("Error unmarshaling databaseConfig: %v", err)
		return nil, err
	}

	return &databaseConfig, err
}

// TODO: clean if not found any use
// func getPostgresAdminSecret() (string, error) {
// 	pod, err := getOwnPod()
// 	if err != nil {
// 		zap.S().Errorf("Error getting own pod: %v", err)
// 		return "", err
// 	}

// 	postgresAdminSecret := ""
// 	for _, c := range pod.Spec.Containers {
// 		for _, e := range c.Env {
// 			if e.Name == "POSTGRES_PASSWORD" {
// 				// TODO: Finalize how to get the admin secret
// 				postgresAdminSecret = e.Value
// 				break
// 			}
// 		}
// 	}

// 	// TODO: Maybe using POSTGRES_PASSWORD_FILE is better and can be placed in shared volume dir
// 	if postgresAdminSecret == "" {
// 		return "", fmt.Errorf("error finding postgres admin secret")
// 	}

// 	return postgresAdminSecret, nil
// }

var (
	dbUser = DB_USER
	dbPort = DB_PORT
	dbHost = DB_HOST
)

func checkAndCreateDatabase(dbName string) error {
	// Run the psql command to list all databases
	cmd := exec.Command("psql", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-lqt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running psql command: %v", err)
	}

	// Convert the output to a string
	dbList := string(output)

	// Check if the database name exists in the list
	if strings.Contains(dbList, dbName) {
		fmt.Printf("Database '%s' already exists\n", dbName)
		return nil
	}

	// If the database doesn't exist, create it
	createCmd := exec.Command("createdb", "-U", dbUser, "-h", dbHost, "-p", dbPort, dbName)
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr

	err = createCmd.Run()
	if err != nil {
		return fmt.Errorf("error creating database '%s': %v", dbName, err)
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
	return nil
}

func createOrUpdateUserOnDatabase(dbName, userName, userPassword string) error {
	// Check if the user already exists
	checkUserCmd := exec.Command("psql", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-d", dbName, "-tAc", fmt.Sprintf("SELECT 1 FROM pg_roles WHERE rolname='%s'", userName))
	userExistsOutput, _ := checkUserCmd.CombinedOutput()
	fmt.Println("Output: ", string(userExistsOutput))

	if strings.TrimSpace(string(userExistsOutput)) == "1" {
		// User already exists, update the password
		updatePasswordCmd := exec.Command("psql", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-d", dbName, "-c", fmt.Sprintf("ALTER ROLE %s WITH PASSWORD '%s'", userName, userPassword)) //nolint:gosec
		updatePasswordCmd.Stdout = os.Stdout
		updatePasswordCmd.Stderr = os.Stderr

		err := updatePasswordCmd.Run()
		if err != nil {
			return fmt.Errorf("error updating password for user '%s' on database '%s': %v", userName, dbName, err)
		}

		fmt.Printf("Password updated for user '%s' on database '%s'\n", userName, dbName)
	} else {
		// User doesn't exist, create the user with the password
		createUserCmd := exec.Command("psql", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-d", dbName, "-c", fmt.Sprintf("CREATE ROLE %s WITH LOGIN PASSWORD '%s'", userName, userPassword)) //nolint:gosec
		createUserCmd.Stdout = os.Stdout
		createUserCmd.Stderr = os.Stderr

		err := createUserCmd.Run()
		if err != nil {
			return fmt.Errorf("error creating user '%s' on database '%s': %v", userName, dbName, err)
		}

		fmt.Printf("User '%s' created successfully on database '%s'\n", userName, dbName)
	}

	// Grant necessary privileges to the user on the database
	grantPrivilegesCmd := exec.Command("psql", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-d", dbName, "-c", fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, userName)) //nolint:gosec
	grantPrivilegesCmd.Stdout = os.Stdout
	grantPrivilegesCmd.Stderr = os.Stderr

	err := grantPrivilegesCmd.Run()
	if err != nil {
		return fmt.Errorf("error granting privileges to user '%s' on database '%s': %v", userName, dbName, err)
	}

	fmt.Printf("Privileges granted to user '%s' on database '%s'\n", userName, dbName)
	return nil
}

func waitForPostgreSQL(maxAttempts int, retryInterval time.Duration) error {
	// Define a timeout to limit the waiting time
	timeout := time.After(time.Minute) // Adjust the timeout as needed

	// Loop until the maximum number of attempts is reached or the timeout expires
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Execute the psql command to check if the database is available
		cmd := exec.Command("psql", "-U", dbUser, "-h", dbHost, "-p", dbPort, "-c", "SELECT 1;") //nolint:gosec
		output, err := cmd.CombinedOutput()

		if err == nil && strings.Contains(string(output), "1") {
			// Successfully connected to the database
			fmt.Println("Connected to the database")
			return nil
		}

		// Sleep for the specified retry interval before the next attempt
		time.Sleep(retryInterval)

		// Check if the timeout has expired
		select {
		case <-timeout:
			// Timeout reached, return an error
			return fmt.Errorf("timeout reached while waiting for the database to start")
		default:
			// Continue to the next attempt
		}
	}

	// Maximum number of attempts reached without success, return an error
	return fmt.Errorf("maximum number of attempts reached, unable to connect to the database")
}

// TODO: Improve error handling and rollback if possible
func applyConfig(config *scPostgresConfig) (string, error) {
	err := waitForPostgreSQL(5, time.Second*5)
	if err != nil {
		return "", err
	}

	var secret string

	for _, dbConf := range config.Databases {
		err := checkAndCreateDatabase(dbConf.Name)
		if err != nil {
			zap.S().Errorf("Error applying config for %s: %v", dbConf.Name, err)
			return "", err
		}

		password, err := password.Generate(32, 10, 0, false, true)
		if err != nil {
			zap.S().Errorf("Error creating new secret password: %v", err)
			return "", err
		}

		err = createOrUpdateUserOnDatabase(dbConf.Name, dbConf.User, password)
		if err != nil {
			zap.S().Errorf("Error create/update user on db: %v", err)
			return "", err
		}

		secret = password
	}

	if secret == "" {
		return "", fmt.Errorf("empty config")
	}

	return secret, nil
}

func (t *tResInfo) updateSecretOnPostgres(crd *v1.TororuResource) (string, error) {
	// postgresAdminSecret, err := getPostgresAdminSecret()
	// if err != nil {
	// 	zap.S().Errorf("Error getting postgres admin secret: %v", err)
	// 	return "", err
	// }

	config, err := parseSidecarPostgresConfig(crd.Spec.Config)
	if err != nil {
		zap.S().Errorf("Error getting crd config: %v", err)
		return "", err
	}

	secret, err := applyConfig(config)
	return secret, err
}
