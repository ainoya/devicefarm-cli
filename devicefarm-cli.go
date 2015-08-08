package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	"github.com/codegangsta/cli"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {

	svc := devicefarm.New(&aws.Config{Region: aws.String("us-west-2")})

	app := cli.NewApp()
	app.Name = "devicefarm-cli"
	app.Usage = "allows you to interact with AWS devicefarm from the command line"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{Name: "Patrick Debois",
			Email: "Patrick.Debois@jedi.be",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "projects",
			Usage: "manage the projects",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the projects", // of an account
					Action: func(c *cli.Context) {
						fmt.Println(c.Args())
						listProjects(svc)
					},
				},
			},
		},
		{
			Name:  "artifacts",
			Usage: "manage the artifacts",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the artifacts", // of a test
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
						cli.StringFlag{
							Name:   "job",
							EnvVar: "DF_JOB",
							Usage:  "job arn or run description",
						},
						cli.StringFlag{
							Name:   "type",
							EnvVar: "DF_ARTIFACT_TYPE",
							Usage:  "type of the artifact [LOG,FILE,SCREENSHOT]",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						jobArn := c.String("job")

						filterArn := ""
						if runArn != "" {
							filterArn = runArn
						} else {
							filterArn = jobArn
						}

						artifactType := c.String("type")
						listArtifacts(svc, filterArn, artifactType)
					},
				},
			},
		},
		{
			Name:  "devicepools",
			Usage: "manage the device pools",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the devicepools", //globally
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
					},
					Action: func(c *cli.Context) {

						projectArn := c.String("project")
						listDevicePools(svc, projectArn)
					},
				},
			},
		},
		{
			Name:  "devices",
			Usage: "manage the devices",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the devices", // globally
					Action: func(c *cli.Context) {
						listDevices(svc)
					},
				},
			},
		},
		{
			Name:  "jobs",
			Usage: "manage the jobs",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the jobs", // of a test
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")

						listJobs(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "runs",
			Usage: "manage the runs",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the runs",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
					},
					Action: func(c *cli.Context) {
						projectArn := c.String("project")
						listRuns(svc, projectArn)
					},
				},
				{
					Name:  "schedule",
					Usage: "schedule a run",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
						cli.StringFlag{
							Name:   "device-pool",
							EnvVar: "DF_DEVICE_POOL",
							Usage:  "devicepool arn or devicepool name",
						},
						cli.StringFlag{
							Name:   "name",
							EnvVar: "DF_RUN_NAME",
							Usage:  "name to give to the run that is scheduled",
						},
						cli.StringFlag{
							Name:  "test-type",
							Usage: "type of test [BUILTIN_FUZZ,BUILTIN_EXPLORER,APPIUM_JAVA_JUNIT,APPIUM_JAVA_TESTNG,CALABASH,INSTRUMENTATION,UIAUTOMATION,UIAUTOMATOR,XCTEST]",
						},
						cli.StringFlag{
							Name:   "test",
							Usage:  "arn or name of the test upload to schedule",
							EnvVar: "DF_TEST",
						},
						cli.StringFlag{
							Name:   "app",
							Usage:  "arn or name of the app upload to schedule",
							EnvVar: "DF_APP",
						},
					},
					Action: func(c *cli.Context) {
						projectArn := c.String("project")
						appUploadArn := c.String("app")
						runName := c.String("name")
						devicePoolArn := c.String("device-pool")
						testUploadArn := c.String("test")
						testType := c.String("test-type")
						scheduleRun(svc, runName, projectArn, appUploadArn, devicePoolArn, testUploadArn, testType)
					},
				},
			},
		},
		{
			Name:  "samples",
			Usage: "manage the samples",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the samples",
					Action: func(c *cli.Context) {
						// Not yet implemented
						// listSamples()
					},
				},
			},
		},
		{
			Name:  "suites",
			Usage: "manage the suites",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the suites",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
						cli.StringFlag{
							Name:   "job",
							EnvVar: "DF_JOB",
							Usage:  "job arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						jobArn := c.String("job")
						filterArn := ""
						if runArn != "" {
							filterArn = runArn
						} else {
							filterArn = jobArn
						}
						listSuites(svc, filterArn)
					},
				},
			},
		},
		{
			Name:  "tests",
			Usage: "manage the tests",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list the tests", // of a Run
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
					},
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						listTests(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "problems",
			Usage: "manage the problems",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "run",
							EnvVar: "DF_RUN",
							Usage:  "run arn or run description",
						},
					},
					Usage: "list the problems", // of Test
					Action: func(c *cli.Context) {
						runArn := c.String("run")
						listUniqueProblems(svc, runArn)
					},
				},
			},
		},
		{
			Name:  "upload",
			Usage: "manages the uploads",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "creates an upload",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "name of the upload",
						},
						cli.StringFlag{
							Name:  "type",
							Usage: "type of upload [ANDROID_APP,IOS_APP,EXTERNAL_DATA,APPIUM_JAVA_JUNIT_TEST_PACKAGE,APPIUM_JAVA_TESTNG_TEST_PACKAGE,CALABASH_TEST_PACKAGE,INSTRUMENTATION_TEST_PACKAGE,UIAUTOMATOR_TEST_PACKAGE,XCTEST_TEST_PACKAGE",
						},
					},
					Action: func(c *cli.Context) {
						uploadName := c.String("name")
						uploadType := c.String("type")
						projectArn := c.String("project")
						uploadCreate(svc, uploadName, uploadType, projectArn)
					},
				},
				{
					Name:  "file",
					Usage: "uploads an file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "project",
							EnvVar: "DF_PROJECT",
							Usage:  "project arn or project description",
						},
						cli.StringFlag{
							Name:  "file",
							Usage: "path to the file to upload",
						},
						cli.StringFlag{
							Name:  "type",
							Usage: "type of upload [ANDROID_APP,IOS_APP,EXTERNAL_DATA,APPIUM_JAVA_JUNIT_TEST_PACKAGE,APPIUM_JAVA_TESTNG_TEST_PACKAGE,CALABASH_TEST_PACKAGE,INSTRUMENTATION_TEST_PACKAGE,UIAUTOMATOR_TEST_PACKAGE,XCTEST_TEST_PACKAGE",
						},
					},
					Action: func(c *cli.Context) {
						uploadType := c.String("type")
						projectArn := c.String("project")
						uploadFilePath := c.String("file")
						uploadPut(svc, uploadFilePath, uploadType, projectArn)
					},
				},
				{
					Name:  "list",
					Usage: "lists all uploads", // of a Project
					Action: func(c *cli.Context) {
						projectArn := "arn:aws:devicefarm:us-west-2:110440800955:project:f7952cc6-5833-47f3-afef-c149fb4e7c76"
						listUploads(svc, projectArn)
					},
				},
				{
					Name:  "info",
					Usage: "info about uploads",
					Action: func(c *cli.Context) {
						uploadArn := "arn:aws:devicefarm:us-west-2:110440800955:upload:f7952cc6-5833-47f3-afef-c149fb4e7c76/1d1a6d6e-554d-48d1-b53f-21f80ef94a14"
						uploadInfo(svc, uploadArn)
					},
				},
			},
		},
	}

	app.Run(os.Args)

}

// --- internal API starts here

/* List all Projects */
func listProjects(svc *devicefarm.DeviceFarm) {

	resp, err := svc.ListProjects(nil)
	failOnErr(err, "error listing projects")

	fmt.Println(awsutil.Prettify(resp))
}

/* List all DevicePools */
func listDevicePools(svc *devicefarm.DeviceFarm, projectArn string) {
	// CURATED: A device pool that is created and managed by AWS Device Farm.
	// PRIVATE: A device pool that is created and managed by the device pool developer.

	pool := &devicefarm.ListDevicePoolsInput{
		ARN: aws.String(projectArn),
	}
	resp, err := svc.ListDevicePools(pool)

	failOnErr(err, "error listing device pools")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all Devices */
func listDevices(svc *devicefarm.DeviceFarm) {

	input := &devicefarm.ListDevicesInput{}
	resp, err := svc.ListDevices(input)

	failOnErr(err, "error listing devices")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all uploads */
func listUploads(svc *devicefarm.DeviceFarm, projectArn string) {

	listReq := &devicefarm.ListUploadsInput{
		ARN: aws.String(projectArn),
	}

	resp, err := svc.ListUploads(listReq)

	failOnErr(err, "error listing uploads")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all runs */
func listRuns(svc *devicefarm.DeviceFarm, projectArn string) {

	listReq := &devicefarm.ListRunsInput{
		ARN: aws.String(projectArn),
	}

	resp, err := svc.ListRuns(listReq)

	failOnErr(err, "error listing runs")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all tests */
func listTests(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListTestsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListTests(listReq)

	failOnErr(err, "error listing tests")
	fmt.Println(awsutil.Prettify(resp))
}

/* List all unique problems */
func listUniqueProblems(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListUniqueProblemsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListUniqueProblems(listReq)

	failOnErr(err, "error listing problems")
	fmt.Println(awsutil.Prettify(resp))
}

/* List suites */
func listSuites(svc *devicefarm.DeviceFarm, filterArn string) {

	listReq := &devicefarm.ListSuitesInput{
		ARN: aws.String(filterArn),
	}

	resp, err := svc.ListSuites(listReq)

	failOnErr(err, "error listing suites")
	fmt.Println(awsutil.Prettify(resp))
}

/* Schedule Run */
func scheduleRun(svc *devicefarm.DeviceFarm, runName string, projectArn string, appUploadArn string, devicePoolArn string, testUploadArn string, testType string) {

	runReq := &devicefarm.ScheduleRunInput{
		AppARN:        aws.String(appUploadArn),
		DevicePoolARN: aws.String(devicePoolArn),
		Name:          aws.String(runName),
		ProjectARN:    aws.String(projectArn),
		Test: &devicefarm.ScheduleRunTest{
			Type: aws.String(testType),

			//TestPackageArn: aws.String(testUploadArn)
			//Parameters: // test parameters
			//Filter: // filter to pass to tests
		},
	}

	resp, err := svc.ScheduleRun(runReq)

	failOnErr(err, "error scheduling run")
	fmt.Println(awsutil.Prettify(resp))
}

/* List Artifacts */

func listArtifacts(svc *devicefarm.DeviceFarm, filterArn string, artifactType string) {

	fmt.Println(filterArn)

	listReq := &devicefarm.ListArtifactsInput{
		ARN: aws.String(filterArn),
	}

	listReq.Type = aws.String("LOG")
	resp, err := svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")
	fmt.Println(awsutil.Prettify(resp))

	listReq.Type = aws.String("SCREENSHOT")
	resp, err = svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")

	fmt.Println(awsutil.Prettify(resp))

	listReq.Type = aws.String("FILE")
	resp, err = svc.ListArtifacts(listReq)
	failOnErr(err, "error listing artifacts")

	fmt.Println(awsutil.Prettify(resp))
}

/* List Jobs */
func listJobs(svc *devicefarm.DeviceFarm, runArn string) {

	listReq := &devicefarm.ListJobsInput{
		ARN: aws.String(runArn),
	}

	resp, err := svc.ListJobs(listReq)

	failOnErr(err, "error listing jobs")
	fmt.Println(awsutil.Prettify(resp))
}

/* Create an upload */
func uploadCreate(svc *devicefarm.DeviceFarm, uploadName string, uploadType string, projectArn string) {

	uploadReq := &devicefarm.CreateUploadInput{
		Name:       aws.String(uploadName),
		ProjectARN: aws.String(projectArn),
		Type:       aws.String(uploadType),
	}

	resp, err := svc.CreateUpload(uploadReq)

	failOnErr(err, "error creating upload")
	fmt.Println(awsutil.Prettify(resp))
}

/* Get Upload Info */
func uploadInfo(svc *devicefarm.DeviceFarm, uploadArn string) {

	uploadReq := &devicefarm.GetUploadInput{
		ARN: aws.String(uploadArn),
	}

	resp, err := svc.GetUpload(uploadReq)

	failOnErr(err, "error getting upload info")
	fmt.Println(awsutil.Prettify(resp))
}

/* Upload a file */
func uploadPut(svc *devicefarm.DeviceFarm, uploadFilePath string, uploadType string, projectArn string) {

	// Read File
	file, err := os.Open(uploadFilePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	// Get file size
	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()

	// read file content to buffer
	buffer := make([]byte, fileSize)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type

	// Prepare upload
	uploadFileBasename := path.Base(uploadFilePath)
	uploadReq := &devicefarm.CreateUploadInput{
		Name:        aws.String(uploadFileBasename),
		ProjectARN:  aws.String(projectArn),
		Type:        aws.String(uploadType),
		ContentType: aws.String("application/octet-stream"),
	}

	resp, err := svc.CreateUpload(uploadReq)
	fmt.Println(awsutil.Prettify(resp))

	uploadInfo := resp.Upload

	upload_url := *uploadInfo.URL

	fmt.Println(upload_url)

	req, err := http.NewRequest("PUT", upload_url, fileBytes)

	if err != nil {
		log.Fatal(err)
	}

	// Remove Host and split to get [0] = path & [1] = querystring
	strippedUrl := strings.Split(strings.Replace(upload_url, "https://prod-us-west-2-uploads.s3-us-west-2.amazonaws.com/", "/", -1), "?")
	req.URL.Opaque = strippedUrl[0]
	req.URL.RawQuery = strippedUrl[1]

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.FormatInt(fileSize, 10))

	// Debug Request to AWS
	debug(httputil.DumpRequestOut(req, false))

	client := &http.Client{}

	res, err := client.Do(req)

	dump, _ := httputil.DumpResponse(res, true)
	log.Printf("} -> %s\n", dump)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

}

/*
 * Helper page to exit on error with a nice message
 */
func failOnErr(err error, reason string) {
	if err != nil {
		log.Fatal("Failed: %s, %s\n\n", reason, err)
		os.Exit(-1)
	}

	return
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
